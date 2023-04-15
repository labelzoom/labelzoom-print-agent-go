package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type printjob struct {
	PrinterHostname string `json:"printerHostname"`
	Text            string `json:"text"`
}

func main() {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://labelzoom.net", "https://www.labelzoom.net", "http://local.labelzoom.net", "http://localhost", "http://localhost:3000"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Origin"},
		// ExposeHeaders:    []string{"Content-Length"},
		// AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		// MaxAge: 12 * time.Hour,
	}))
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/print", func(c *gin.Context) {
		var newPrintJob printjob
		if err := c.BindJSON(&newPrintJob); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("error parsing json: %v", err))
			return
		}

		// send data to printer
		connection, err := net.Dial("tcp", newPrintJob.PrinterHostname+":9100")
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error connecting to printer: %v", err))
			return
		}
		defer connection.Close()
		_, err = connection.Write([]byte(newPrintJob.Text))
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error sending ZPL to printer: %v", err))
			return
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
