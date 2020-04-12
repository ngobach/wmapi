package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngobach/wmapi/config"
	"github.com/ngobach/wmapi/wm"
)

func StartServer() error {
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	g := gin.Default()
	g.GET("/:country", func(context *gin.Context) {
		country := wm.CountryFrom(context.Param("country"))
		defer func() {
			if recover() != nil {
				context.String(200, "Failed you")
			}
		}()
		resp, err := wm.GetStatistics(&country)
		if err != nil {
			_ = context.Error(err)
			return
		}
		context.JSON(200, resp)
	})
	return g.Run(addr)
}
