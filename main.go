package main

import (
	_ "embed"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed resources/labelzoom_logo.txt
var logo string

type printjob struct {
	PrinterHostname string `json:"printerHostname"`
	Text            string `json:"text"`
}

func handlePing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func handlePrint(c *gin.Context) {
	var newPrintJob printjob
	if err := c.BindJSON(&newPrintJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("error parsing JSON: %v", err),
		})
		return
	}

	// send data to printer
	connection, err := net.Dial("tcp", newPrintJob.PrinterHostname+":9100")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error connecting to printer: %v", err),
		})
		return
	}
	defer connection.Close()
	_, err = connection.Write([]byte(newPrintJob.Text))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("error sending ZPL to printer: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "OK",
	})
}

func main() {
	fmt.Println(logo)
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1", "::1"})
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://labelzoom.net", "https://www.labelzoom.net", "http://local.labelzoom.net", "http://localhost", "http://localhost:3000"},
		AllowMethods: []string{"POST"},
		AllowHeaders: []string{"Origin", "Content-Type"}, // TODO: Acess to fetch at 'http://localhost:8080/print' from origin 'http://localhost:3000' has been blocked by CORS policy: Request header field content-type is not allowed by Access-Control-Allow-Headers in preflight response.
		// ExposeHeaders:    []string{"Content-Length"},
		// AllowCredentials: true,
		// AllowOriginFunc: func(origin string) bool {
		// 	return origin == "https://github.com"
		// },
		// MaxAge: 12 * time.Hour,
	}))
	r.GET("/ping", handlePing)
	r.POST("/print", handlePrint)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
