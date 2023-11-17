package Views

import (
	"encoding/json"
	"net/http"

	"github.com/evanyip05/Cloud/config"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)


func Web() {
	router := gin.Default()

	// Serve the frontend
	router.StaticFile("/", "Frontend/build/index.html")
	router.Static("/static", "Frontend/build/static")

	// Define a WebSocket route
	router.GET("/mongo/get", func(c *gin.Context) {
		//TelemetryHandler(c.Writer, c.Request)
		r, err := http.Get("localhost:8080/get")
		if err!= nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		MarshalAndSend(c, r)
	})

	router.POST("/mongo/put", func(c *gin.Context) {
		var locationData config.MongoPutRequest

		c.BindJSON(&locationData)

	})

	

	// log.Info().Str("url", Network.GetOutboundIP(":" + Config.ServerPort)).Msg("Listening on port")

	// go Browser.OpenURL(Config.ServerPort)

	// go PillbotIO.Launch()

	log.Info().Msg("listening on https://localhost:3000")
	
	err := router.Run(":3000")
	if err != nil {
		log.Error().Err(err).Msg("Error listening and serving")
	}
	
}

func MarshalAndSend(c *gin.Context, data any) {
    json, err := json.Marshal(data)
    if err != nil {c.Writer.WriteHeader(http.StatusInternalServerError); return}
    c.Writer.Write(json)
}
