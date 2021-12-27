package main

import (
	"log"
	"os"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"

	"rocky-list/db"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/RockList", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/RockList/Profile", func(c *gin.Context) {
		c.HTML(http.StatusOK, "profile.html", nil)
	})

	router.GET("/RockList/Disc", func(c *gin.Context) {
		released, _ := strconv.Atoi(c.Query("released"))
		title := c.Query("title")
		list := c.Query("list")
		artist := c.Query("artist")
		disc := db.FindBy(list, title, artist, released)

		c.HTML(http.StatusOK, "disc.html", disc)
	})

	router.GET("/RockList/List", func(c *gin.Context) {
		list := c.Query("list")
		discs := db.FindByList(list)

		c.HTML(http.StatusOK, "result.html", gin.H{
			"discs": discs,
			"list":  list,
		})
	})

	router.Run(":" + port)
}
