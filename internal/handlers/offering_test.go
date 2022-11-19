package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOffering(t *testing.T) {
	type test struct {
		packageName string
		countryCode string
		expBody     string
		expCode     int
	}

	configurations.UpdateConfigurationTable(commons.ConfigurationTable{
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

	tests := []test{
		{packageName: "com.softinit.iquitos.mainapp",
			countryCode: "SI",
			expBody:     "{\"main_sku\":\"rdm_premium_v3_050_trial_7d_yearly\"}",
			expCode:     http.StatusOK},
		{packageName: "com.softinit.iquitos.newapp",
			countryCode: "US",
			expBody:     "",
			expCode:     http.StatusNoContent},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			router := gin.Default()
			router.GET("/:package", Offering)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("GET", fmt.Sprintf("/%s", tt.packageName), nil)
			router.ServeHTTP(recorder, request)
			retBody, _ := io.ReadAll(recorder.Body)
			retCode := recorder.Code
			if tt.expBody != string(retBody) {
				t.Errorf("body mismatch: expected '%s', received '%s'", tt.expBody, string(retBody))
			}
			if tt.expCode != retCode {
				t.Errorf("code mismatch: expected '%d', received '%d'", tt.expCode, retCode)
			}
		})
	}
}
