package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nikolaevv/my-investor/internal/handlers/requests"
	"github.com/nikolaevv/my-investor/pkg/tinkoff/investapi"
)

func (h *Handler) GetShare(c *gin.Context) {
	var params requests.GettingStock
	if err := c.Bind(&params); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}
	if err := params.Validate(); err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("error: %s", err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := investapi.CreateInstrumentsServiceClient(h.Config.Tinkoff.URL, h.Config.Tinkoff.Token)

	req := investapi.InstrumentRequest{
		IdType:    2,
		ClassCode: params.ClassCode,
		Id:        params.Id,
	}

	protoMessage, err := client.ShareBy(ctx, &req)
	if err != nil {
		c.String(http.StatusNotFound, "error: not found")
		return
	}

	c.JSON(200, protoMessage)
}
