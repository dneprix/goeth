package main

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/dneprix/goeth"
	"github.com/dneprix/goeth/model"
	"github.com/dneprix/goeth/proto"
)

func setHandlers(r *gin.Engine, svc goeth.Service) {
	r.POST("/api/v1/SendEth", SendEthHandler(svc))
	r.GET("/api/v1/GetLast", GetLastHandler(svc))
}

func SendEthHandler(svc goeth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req proto.ApiSendEthRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, proto.NewErrorResponse(err))
			return
		}
		tx, err := svc.SendWei(req.From, req.To, model.WeiFromEth(req.Amount))
		if err != nil {
			ctx.JSON(500, proto.NewErrorResponse(err))
			return
		}
		ctx.JSON(200, proto.ApiSendEthResponse{
			TxHash: tx.Hash,
		})
		return
	}
}

func GetLastHandler(svc goeth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		txs := svc.GetLast()
		res := make([]proto.ApiGetLastItemResponse, 0, len(txs))
		c.JSON(200, txs)
		return

		for _, tx := range txs {
			res = append(res, proto.ApiGetLastItemResponse{
				Date:          tx.CreatedAt.Format(time.UnixDate),
				Address:       tx.AccountTo,
				Amount:        tx.Amount.Ether(),
				Confirmations: tx.Confirmations,
			})
		}
		c.JSON(200, res)
		return
	}
}
