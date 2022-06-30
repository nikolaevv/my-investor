package usecase

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"github.com/nikolaevv/my-investor/internal/transport/dto/request"
	"github.com/nikolaevv/my-investor/pkg/gen/proto/tinkoff/investapi"
)

func getUser(h *handler, headers http.Header) (*entity.User, error) {
	signingKey := h.Config.Auth.JWTSecret
	claims, err := h.Auth.AuthorizateUser(headers, signingKey)
	if err != nil {
		return nil, err
	}
	user, err := h.Repo.User.GetUserByID(claims.Id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func prepareBuyShareReq(c *gin.Context) (*request.BuyingShare, error) {
	var req request.BuyingShare
	if err := c.BindJSON(&req); err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	return &req, nil
}

func getShareInfoByTicker(h *handler, сlassCode string, id string, ctx context.Context) (*investapi.ShareResponse, error) {
	instrumentClient := investapi.CreateInstrumentsServiceClient(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token)

	gettingShareReq := investapi.InstrumentRequest{
		IdType:    investapi.InstrumentIdType_INSTRUMENT_ID_TYPE_TICKER,
		ClassCode: сlassCode,
		Id:        id,
	}

	return instrumentClient.ShareBy(ctx, &gettingShareReq)
}

func addShareToPortfolio(h *handler, reqData *request.BuyingShare, userId uint) (uint, error) {
	share := &entity.Share{}

	if err := copier.Copy(share, reqData); err != nil {
		return 0, err
	}

	share.UserID = userId
	share.Code = reqData.Id

	return h.Repo.Share.CreateShare(share)
}

func prepareResponse(h *handler, reqData *request.BuyingShare, shareData *investapi.ShareResponse, user *entity.User, shareId uint, ctx context.Context) (*investapi.PostOrderResponse, error) {
	sandboxClient := investapi.CreateSandboxServiceClient(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token)
	return sandboxClient.PostSandboxOrder(ctx, &investapi.PostOrderRequest{
		Quantity:  int64(reqData.Quantity),
		Figi:      shareData.Instrument.Figi,
		Direction: investapi.OrderDirection_ORDER_DIRECTION_BUY,
		AccountId: user.AccountID,
		OrderType: investapi.OrderType_ORDER_TYPE_MARKET,
		OrderId:   fmt.Sprint(shareId),
	})
}

func (h *handler) BuyShare(c *gin.Context) {
	user, err := getUser(h, c.Request.Header)
	if err != nil {
		c.String(http.StatusForbidden, err.Error())
		return
	}

	reqData, err := prepareBuyShareReq(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	shareData, err := getShareInfoByTicker(h, reqData.ClassCode, reqData.Id, c)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	shareId, err := addShareToPortfolio(h, reqData, user.ID)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	resp, err := prepareResponse(h, reqData, shareData, user, shareId, c)
	if err != nil {
		c.String(http.StatusBadGateway, err.Error())
		return
	}

	c.JSON(http.StatusOK, resp)
}
