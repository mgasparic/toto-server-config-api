package main

import (
	"api/internal/commons"
	"api/internal/configurations"
	"api/internal/handlers"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
	"net/http"
	"os"
)

type EnvVars struct {
	ServicePort                          int    `envconfig:"SERVICE_PORT" default:"8080"`
	ServiceEnvironment                   string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
	PersistentConfigurationTableFilePath string `envconfig:"CONFIGURATION_TABLE_PATH"`
	ConfigurationJwtPublicKeyPath        string `envconfig:"CONFIGURATION_JWT_PUBLIC_KEY_PATH"`
	SubscriptionJwtPublicKeyPath         string `envconfig:"SUBSCRIPTION_JWT_PUBLIC_KEY_PATH"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	persistentConfigurationTable, err := os.ReadFile(envVars.PersistentConfigurationTableFilePath)
	if err != nil {
		log.Fatal(err)
	}

	var configurationTable commons.ConfigurationTable
	err = json.Unmarshal(persistentConfigurationTable, &configurationTable)
	if err != nil {
		log.Fatal(err)
	}
	configurations.UpdateActiveConfigurationTable(configurationTable)

	ce := handlers.ConfigurationEnvironment{
		UpdaterEnvironment:                   commons.UpdaterEnvironment{JwtPublicKey: mustParsePublicKey(envVars.ConfigurationJwtPublicKeyPath)},
		PersistentConfigurationTableFilePath: envVars.PersistentConfigurationTableFilePath,
	}
	se := handlers.SubscriptionEnvironment{JwtPublicKey: mustParsePublicKey(envVars.SubscriptionJwtPublicKeyPath)}

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.POST("/configuration", ce.Configuration)
	router.GET("/subscription", se.Subscription)
	router.GET("/", func(ctx *gin.Context) {
		ctx.AbortWithStatus(http.StatusOK)
	})

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}

func mustParsePublicKey(path string) *rsa.PublicKey {
	keyRaw, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyRaw)
	if err != nil {
		log.Fatal(err)
	}
	return key
}
