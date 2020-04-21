package server

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ngobach/wmapi/config"
	"github.com/ngobach/wmapi/wm"
	"github.com/olekukonko/tablewriter"
	"time"
)

func StartServer() error {
	addr := fmt.Sprintf("0.0.0.0:%d", config.Port)
	tz, err := time.LoadLocation("asia/Ho_Chi_Minh")
	if err != nil {
		tz, err = time.LoadLocation("")
	}
	if err != nil {
		panic(err)
	}
	g := gin.Default()
	g.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/viet-nam")
	})
	g.GET("/:country", func(context *gin.Context) {
		country := wm.CountryFrom(context.Param("country"))
		resp, err := wm.GetStatistics(country)
		if err != nil {
			_ = context.Error(err)
			context.String(400, "Sorry! I'm not able to process the request now.")
			return
		}
		if _, ok := context.GetQuery("json"); ok {
			context.JSON(200, resp)
		} else {
			buf := &bytes.Buffer{}
			table := tablewriter.NewWriter(buf)
			table.SetHeader([]string{"Date", "Total cases", "Active cases", "New cases"})
			for _, day := range resp.Days {
				table.Append([]string{
					day.Date.Format("02/01/2006"),
					fmt.Sprintf("%d", day.Total),
					fmt.Sprintf("%d", day.Active),
					fmt.Sprintf("%d", day.New),
				})
			}
			table.SetFooter([]string{
				"Summary",
				fmt.Sprintf("%d total", resp.Total),
				fmt.Sprintf("%d recovered", resp.Recovered),
				fmt.Sprintf("%d deaths", resp.Deaths),
			})
			table.SetCaption(true, fmt.Sprintf("Updated at %s", resp.UpdatedAt.In(tz).Format("01/02/2006 15:04")))
			table.Render()
			context.String(200, "%s", buf.String())
		}
	})
	return g.Run(addr)
}
