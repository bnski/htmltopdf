package main

import (
	"html2pdf/htmltopdf"
	"log"
	"os"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // Provide a default port if none is specified
	}

	p := htmltopdf.New()
	g := gin.Default()
	g.POST("/page-to-base64", p.ToBase64)
	g.POST("/page-to-s3", p.ToS3)
	log.Fatal(g.Run(":" + port))
}
