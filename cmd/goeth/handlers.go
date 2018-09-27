package main

import (
	"github.com/gin-gonic/gin"

	"github.com/dneprix/goeth"
	"github.com/dneprix/goeth/proto"
)

func SetHandlers(r *gin.Engine, svc goeth.Service) {
	r.POST("/api/v1/SendEth", SendEthHandler(svc))
	r.GET("/api/v1/GetLast", GetLastHandler(svc))
}

func SendEthHandler(svc goeth.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req proto.SendEthRequest
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(400, proto.NewErrorResponse(err))
			return
		}

		results, err := svc.SendEth() //req.IDs
		if err != nil {
			ctx.JSON(500, proto.NewErrorResponse(err))
			return
		}
		_ = results
		ctx.JSON(200, proto.SendEthResponse{
			//Users: results,
		})
		return
	}
}

func GetLastHandler(svc goeth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {

		results, err := svc.GetLast()
		_ = results
		if err != nil {
			c.JSON(500, proto.NewErrorResponse(err))
			return
		}
		c.JSON(200, proto.GetLastResponse{
			//Users: results,
		})
		return
	}
}
