package main

import (
	"api/internal/commons"
	"api/internal/configurations"
	"api/internal/handlers"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"log"
)

type EnvVars struct {
	ServicePort        int    `envconfig:"SERVICE_PORT" default:"9000"`
	ServiceEnvironment string `envconfig:"SERVICE_ENVIRONMENT" default:"local"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	// todo rewrite:
	configurations.UpdateConfigurationTable(commons.ConfigurationTable{
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 0,
				PercentileMax: 25,
				MainSku:       "rdm_premium_v3_020_trial_7d_monthly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 25,
				PercentileMax: 50,
				MainSku:       "rdm_premium_v3_030_trial_7d_monthly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 50,
				PercentileMax: 75,
				MainSku:       "rdm_premium_v3_100_trial_7d_yearly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "US",
			},
			commons.ConfigurationChance{
				PercentileMin: 75,
				PercentileMax: 100,
				MainSku:       "rdm_premium_v3_150_trial_7d_yearly",
			},
		},
		{
			commons.ConfigurationRequirement{
				Package:     "com.softinit.iquitos.mainapp",
				CountryCode: "ZZ",
			},
			commons.ConfigurationChance{
				PercentileMin: 0,
				PercentileMax: 100,
				MainSku:       "rdm_premium_v3_050_trial_7d_yearly",
			},
		},
	})

	if envVars.ServiceEnvironment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	service := router.Group("/distributor")
	service.GET("/offering/:package", handlers.Offering)

	log.Fatal(router.Run(fmt.Sprintf(":%d", envVars.ServicePort)))
}
