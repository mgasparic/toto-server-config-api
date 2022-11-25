package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"io"
	"log"
	"net/http"
)

type EnvVars struct {
	Port               string `envconfig:"PORT"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	HostUrl            string `envconfig:"HOST_URL"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/offering/:package", func(ctx *gin.Context) {
		request, err := http.Get(fmt.Sprintf("%s/offering/%s/%s", envVars.HostUrl, ctx.Param("package"), ctx.GetHeader("X-Appengine-Country")))
		if err != nil {
			log.Print(err)
			ctx.AbortWithStatus(http.StatusServiceUnavailable)
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			log.Print(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		_, err = ctx.Writer.Write(body)
		if err != nil {
			log.Print(err)
		}
	})
	router.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})
	log.Fatal(router.Run(fmt.Sprintf(":%s", envVars.Port)))
}
