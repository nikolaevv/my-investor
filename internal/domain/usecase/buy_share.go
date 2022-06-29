package usecase

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/nikolaevv/my-investor/internal/domain/entity"
	"github.com/nikolaevv/my-investor/internal/transport/dto/request"
	"github.com/nikolaevv/my-investor/pkg/tinkoff/investapi"
)

func (h *handler) BuyShare(c *gin.Context) {
	signingKey := h.Config.Auth.JWTSecret
	claims, err := h.Auth.AuthorizateUser(c.Request.Header, signingKey)
	if err != nil {
		c.String(http.StatusForbidden, fmt.Sprintf("error: %s", err))
		return
	}

	var req request.BuyingShare
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	if err := req.Validate(); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// 1. getting info about share
	instrumentClient := investapi.CreateInstrumentsServiceClient(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token)

	gettingShareReq := investapi.InstrumentRequest{
		IdType:    2,
		ClassCode: req.ClassCode,
		Id:        req.Id,
	}

	protoShareMsg, err := instrumentClient.ShareBy(ctx, &gettingShareReq)
	if err != nil {
		c.String(http.StatusNotFound, "error: not found")
		return
	}

	share := &entity.Share{}
	copier.Copy(share, req)
	share.UserID = uint(claims.Id)
	share.Code = req.Id

	shareId, err := h.Repo.Share.CreateShare(share)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("error: %s", err))
		return
	}

	user, err := h.Repo.User.GetUserByID(claims.Id)
	if err != nil {
		c.String(http.StatusNotFound, "error: not found")
		return
	}

	// 2. buy share by FIGI id, gotten in previous step
	sandboxClient := investapi.CreateSandboxServiceClient(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token)
	buyingShareReq := investapi.PostOrderRequest{
		Quantity:  int64(req.Quantity),
		Figi:      protoShareMsg.Instrument.Figi,
		Direction: 1,
		AccountId: user.AccountID,
		OrderType: 2,
		OrderId:   fmt.Sprint(shareId),
	}
	fmt.Println(protoShareMsg.Instrument.Figi)

	protoOrderMsg, err := sandboxClient.PostSandboxOrder(ctx, &buyingShareReq)
	if err != nil {
		c.String(http.StatusBadGateway, fmt.Sprintf("error: %s", err))
		return
	}

	c.JSON(200, protoOrderMsg)
}
