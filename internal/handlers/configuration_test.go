package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestConfiguration(t *testing.T) {
	type test struct {
		address               commons.Address
		authorization         string
		body                  string
		expConfigurationTable commons.ConfigurationTable
	}

	jwtPublicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(strings.Replace("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1BcotH1yuIoY7jUataLT\na2Id1EpfJ8iGaX/0HFMu/byedn2MsVb1xg9Rufj3q4zc60lhhCZB260h8ERc+NtF\nwRivceBeGK3voZ6xQYYPAovo+O+GX7sxWQHWREB5hrXv4LGqLeP/ste7FXntBEoa\n+QQqn6y4fDrkWeMtUq/wfsaNP3JE+p/gsLyFAWVQ7kWmaUPpkbQb8Tv12OsY9rEm\nfVxUAAbixFN7tOxJgt57jAQHC331Br8w8IC1P860pAmqconlU7jKKv9QCoyoH9qh\nBOgCqiZ52dCmx7sf61gXwOQ4yTMlKH67oUnfg6RSGW4tjBcIwJDGUI8cnadSTEjA\n3QIDAQAB\n-----END PUBLIC KEY-----", `\n`, "\n", -1)))
	ce := ConfigurationEnvironment{
		UpdaterEnvironment:                   commons.UpdaterEnvironment{JwtPublicKey: jwtPublicKey},
		PersistentConfigurationTableFilePath: "test.json",
	}
	_, _ = os.Create(ce.PersistentConfigurationTableFilePath)
	defer func() {
		_ = os.Remove(ce.PersistentConfigurationTableFilePath)
	}()

	invalidPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEA1IYja+OFP9TtaaX7IRlowZsmuSZkWR7pLGBVR+5W9reiNd47\nqaGc9+OaINgXIT54lZhsbpZDnAph0sTR4ZFINnz6BxWdD0dfYGJregA/DS7NySaR\n1OE9ICFh5kX+hHof5PyHj0Ser8DUvXNPB+AltYpaq3kwUphpUOE3mKapAyMy5aJO\n6IAF8iXBvSzSXdckLAVyVbICiNDzTWt8Z/o9AworPSmWyt86CJaTvt0wlu8rQTJj\nQCw5zQzcnPJD3EpE5JYuJ/xNjOc61F9lu4Ydc9638GWGKkUrCEHk0njQ7Lk+570m\nUklV6R2JWs15ScAlavxtB0VS88Ah7bIL9mIiiQIDAQABAoIBAG53XtgIhk18jv8Y\nQGWfAH2J+OHKW+SbFZZ0aD+gWmGZZ95aKmbs01tiWc3ypGsqPqF7ffkpb/Ee0GQu\n2+1Eh+WSpp+iE2ZKsf+2iKj8kcl3Z43D149wmEYoM+vuu/H+TK/It2m3wEgYkjv7\nQLwWPWOUi7aPopd8E1nwBub40ecHpH/iFjQdy/yfUE9WkzAr1ndwywirnWyjMM3x\nwocKlCcjvuILsAt1g8g1RZ2fVHUMzyxSVYPTpGYdyhkabLlnxJwEo9YRoZni4G8z\nWb3uV+RS32K7c2ucikMuGulHnV0DneyRXiQoFVbpqF6MqbEoy5rWLgFOCGAvPgmq\nlchpd4ECgYEA/I5Fu0ew8TJiO4WIEZ5DGjYRXjCRUmWNPPo4fM3y/3kugRpglKmm\nMmQmDzn4nc4L9XjOe5UH5MGli/GA0/ZgXxPBC4KFflAOUjAMGJCOgLOMhG37OOpF\nXp/whNAmjGQ1SPGuH0dTs4JIkyuRBWCbF0rYkyGOOdc+86fUctI7ohECgYEA12wb\nM+CgbUfZ+ORc2j+xzwh/pQBlk+bedWyO6VIw8J3bOoJMEABezHS39P11BueCZGAE\neurHND5Bp9ZnFxYqAD2BeLAzGAwMbAbfLez5mP0csFKAI6s1ihyje9WaVvSW6BwL\niQL2OW57cqWqXbKU9Lo66oB+q4Z8EaBRzMrJgPkCgYEAhbKzi49KRoWLp0LrY5hK\noaeZzikb5WjJOkykdr14NHJsgf/6vPiKeQa3dzwTN0cK4apQdO/SO/Gft9PjhVJa\ndjq46WTxUosC6dNxs9di+RMUAk9OvTSYAJ3e8BBZN90csD7xFLHnx5Hi5bYckaIw\nuEXxHQKYjlo6gzaHqzlMRGECgYEA0EKFMewRbFSwV7x+P8igH6T5sgzmJsxleGQQ\n7WQ2SAh0LuZUnoudGlAkn8aA4sHh/yQMcMCVAN7HHnlahKk6xaywhHrjinXrdGxY\nSs/0pYDdwWCg3NriEbmKG2fvo+0mDxM926FvZSp4UefzAk5pTwbTem3rB+wl+exh\n6HiLwhECgYEA7yYhys1C5NFGKQRQMSSWr2As3N16nztdQQtIVYlEeebPc4U2/tz0\nM9/XnIUfS7nIr83klpOPaxj4bP5f+5dykggRz1AiW5cKZHH9L/VAGcZLojayH/k1\n/ZLJ3FaVKVwjx+716EaISo3qAmy+41FoDOwaJMnpBBp9hZIGuxUZrsA=\n-----END RSA PRIVATE KEY-----"))
	validPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA1BcotH1yuIoY7jUataLTa2Id1EpfJ8iGaX/0HFMu/byedn2M\nsVb1xg9Rufj3q4zc60lhhCZB260h8ERc+NtFwRivceBeGK3voZ6xQYYPAovo+O+G\nX7sxWQHWREB5hrXv4LGqLeP/ste7FXntBEoa+QQqn6y4fDrkWeMtUq/wfsaNP3JE\n+p/gsLyFAWVQ7kWmaUPpkbQb8Tv12OsY9rEmfVxUAAbixFN7tOxJgt57jAQHC331\nBr8w8IC1P860pAmqconlU7jKKv9QCoyoH9qhBOgCqiZ52dCmx7sf61gXwOQ4yTMl\nKH67oUnfg6RSGW4tjBcIwJDGUI8cnadSTEjA3QIDAQABAoIBACCQ2cg1Bvt9AOz6\nrHo9YTc43pmtPcUvDix+4C3FPA7r+mz5RDQqxRw/V+41Y1otC/L41odjFHO1tuNc\nq+XuTVyj/LSAnhIuCwCHDHKiMgXTE9e9d6WVabNgC9V5DO/5WbbnsNjRDH8ajvhy\nO0mQU2nvENhpvf4dUNjYgCoVJ9lfPbUlASf1Vm5KcluCwlCK1uOmAfiJioLigFDl\nHMxhBT35eb1/xsgEfBDsy9004+qCBr+475Rk/GUdb1CXyBhFoJ28ZrFpVliaLLar\nbeOIf1Y2Xs2m5drxp395ArqHqoOAH57ifPKiPs+K7KB5RzurVSfRWtOgNuiEENON\nOtI+o8UCgYEA+Loq/RZQ0PGUqMn+CN3km7/8CDsqhHe4O4ulDZalpHJ95K5+vOZi\noVFVB+n7M9wvWdVTtzOYvrnsn8TqesIyJK3jWeIEq3Ohg2srQJcNhstQhV11F8X0\nn1lyjrCBKla28wBU+19faysAOJb82M1/sayniPbzWHxCpQeXdQHenG8CgYEA2kq/\nuvTojhB9WL624wj/ZJqNxdlpb0ThIyeCmUp6Bx4ozWuyN9UoSLx0j9MGGpj2QCHl\nxoeN6jW3w9Fb04Loatf7/GA5yVeSlCZ8vOk7vj4GlZgNkFS+R1uMAL5t4rnP+2YO\noVLmI9IbPuQo5KdSd2rBhqGfZlTYqCKJ6i0ytXMCgYEAgZZ1SVE0H/iNzHcZDMOX\nFPKsvBkfaM77RMLX5sGDYa9pGhkz3PVnk8bNN4dXoshoPzSfHkcaoNw7hW5SAE1n\nVKboWe8hIbboApF2gntwx7bsJ9/uXsl7Tv6Kaf9Z/JYbUXXt0N3619oajmFUMRy4\nF/jtfLW1SXfMPTE8XvDva68CgYA/F6e+45ZrqTxxb/3wNOOMMleTvbkJngDsZjkL\nEUmf0Qm+Bcim9ocELuQJxdZXzaou8x00em18KjaI9HPz+Dww1WhRk4ZgOV12UIFx\nIIBSBh9lEWOqObQdb8pRVLjx9P62DNNhsVIvPET/snZXUD03orV4sjIeI5vTTNhL\nxPCJHQKBgFttjVrst2PGU6mSgRkFeoQxajWkR919+MZTNkbjR9fvlmIECh1vTUJ8\n3nBkKoGCIexIFPT1wS7idzsqtzpCjgfn+gIhnDZw76NlcJiTBXq4is71igZqwzI8\nXi/TLYh7PXZv+qCEkCdEznsm85FcQv0/CJ/BNFcULGSG/BAKJKcX\n-----END RSA PRIVATE KEY-----"))

	configurations.UpdateActiveConfigurationTable(nil)

	tests := []test{
		{
			address:               "1.2.3.4:80",
			authorization:         "",
			body:                  "",
			expConfigurationTable: nil,
		},
		{
			address: "127.0.0.1:36521",
			authorization: fmt.Sprintf("Bearer %s", func() string {
				token, _ := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{"usr": "sys.admin@zippo.apps.com"}).SignedString(invalidPrivateKey)
				return token
			}()),
			body:                  "",
			expConfigurationTable: nil,
		},
		{
			address: "34.52.1.172:5420",
			authorization: fmt.Sprintf("Bearer %s", func() string {
				token, _ := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{"usr": "sys.admin@zippo.apps.com"}).SignedString(validPrivateKey)
				return token
			}()),
			body: "[{\"package\":\"com.softinit.iquitos.mainapp\",\"country_code\":\"US\",\"percentile_min\":0,\"percentile_max\":25,\"main_sku\":\"rdm_premium_v3_020_trial_7d_monthly\"},{\"package\":\"com.softinit.iquitos.mainapp\",\"country_code\":\"US\",\"percentile_min\":25,\"percentile_max\":50,\"main_sku\":\"rdm_premium_v3_030_trial_7d_monthly\"}]",
			expConfigurationTable: commons.ConfigurationTable{
				{
					ConfigurationRequirement: commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					ConfigurationChance: commons.ConfigurationChance{
						PercentileMin: 0,
						PercentileMax: 25,
						MainSku:       "rdm_premium_v3_020_trial_7d_monthly",
					},
				},
				{
					ConfigurationRequirement: commons.ConfigurationRequirement{
						Package:     "com.softinit.iquitos.mainapp",
						CountryCode: "US",
					},
					ConfigurationChance: commons.ConfigurationChance{
						PercentileMin: 25,
						PercentileMax: 50,
						MainSku:       "rdm_premium_v3_030_trial_7d_monthly",
					},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			router := gin.Default()
			router.POST("/", ce.Configuration)
			recorder := httptest.NewRecorder()
			request, _ := http.NewRequest("POST", "/", bytes.NewReader([]byte(tt.body)))
			request.Header.Set("Authorization", tt.authorization)
			router.ServeHTTP(recorder, request)
			time.Sleep(time.Second)
			subscriberChannel := make(chan commons.ConfigurationTable)
			go configurations.AddSubscriber(tt.address, subscriberChannel)
			retConfigurationTable := <-subscriberChannel
			if !reflect.DeepEqual(tt.expConfigurationTable, retConfigurationTable) {
				t.Errorf("configuration table mismatch: expected '%v', received '%v'", tt.expConfigurationTable, retConfigurationTable)
			}
		})
	}
}
