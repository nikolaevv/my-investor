package usecase

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/transport/dto/request"
)

func prepareGetShareReq(c *gin.Context) (*request.GettingShare, error) {
	var params request.GettingShare
	if err := c.Bind(&params); err != nil {
		return nil, err
	}
	if err := params.Validate(); err != nil {
		return nil, err
	}
	return &params, nil
}

func (h *handler) GetShare(c *gin.Context) {
	reqParams, err := prepareGetShareReq(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	shareData, err := getShareInfoByTicker(h, reqParams.ClassCode, reqParams.Id, c)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, shareData)
}
