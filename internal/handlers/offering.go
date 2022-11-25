package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Offering(ctx *gin.Context) {
	mainSku, ok := configurations.RandomMainSku(commons.ConfigurationRequirement{
		Package:     commons.Package(ctx.Param("package")),
		CountryCode: commons.CountryCode(ctx.Param("country")),
	})
	if !ok {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	} else {
		ctx.JSON(http.StatusOK, commons.ApiResponseParameters{MainSku: mainSku})
	}
}
