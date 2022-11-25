package main

import (
	"api/internal/commons"
	"api/internal/configurations"
	"api/internal/handlers"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"time"
)

type EnvVars struct {
	ServicePort        int    `envconfig:"SERVICE_PORT" default:"9000"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	ConfigTableUrl     string `envconfig:"CONFIG_TABLE_URL"`
	AuthorizationJwt   string `envconfig:"AUTHORIZATION_JWT"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})
	go func() {
		doneSignaled := false
		for {
			connection, _, err := websocket.DefaultDialer.Dial(envVars.ConfigTableUrl, http.Header{"Authorization": []string{fmt.Sprintf("Bearer %s", envVars.AuthorizationJwt)}})
			if err != nil {
				log.Print(err)
				time.Sleep(time.Minute)
				continue
			}
			for {
				err := connection.SetReadDeadline(time.Time{})
				if err != nil {
					log.Print(err)
					break
				}
				messageType, message, err := connection.ReadMessage()
				if err != nil {
					log.Print(err)
					break
				}
				if messageType == websocket.TextMessage {
					var configurationTable commons.ConfigurationTable
					err = json.Unmarshal(message, &configurationTable)
					if err != nil {
						log.Fatal(err)
					}
					configurations.UpdateConfigurationTable(configurationTable)
					if !doneSignaled {
						doneSignaled = true
						close(done)
					}
				}
			}
			_ = connection.Close()
		}
	}()
	<-done

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/offering/:package/:country", handlers.Offering)
	router.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}
