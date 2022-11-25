package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

type ConfigurationEnvironment struct {
	commons.UpdaterEnvironment
	PersistentConfigurationTableFilePath string
}

func (ce ConfigurationEnvironment) Configuration(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	claims := jwt.MapClaims{"usr": ""}
	_, err := jwt.ParseWithClaims(auth[7:], &claims, func(token *jwt.Token) (interface{}, error) { return ce.JwtPublicKey, nil })
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var configurationTable commons.ConfigurationTable
	err = ctx.BindJSON(&configurationTable)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	configurations.UpdateActiveConfigurationTable(configurationTable)
	log.Printf("configuration table updated by %v", claims["usr"])
	configurationTableRaw, err := json.MarshalIndent(configurationTable, "", "\t")
	if err != nil {
		log.Print(err)
	}
	err = os.WriteFile(ce.PersistentConfigurationTableFilePath, configurationTableRaw, 0644)
	if err != nil {
		log.Print(err)
	}
}
