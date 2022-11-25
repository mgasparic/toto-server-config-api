package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestSubscription(t *testing.T) {
	type test struct {
		authorization         string
		expIsHandshakeErr     bool
		expConfigurationTable commons.ConfigurationTable
	}

	jwtPublicKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(strings.Replace("-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAxmovukuj5QtrJADJe5uy\noNZ8/TQ9Ng066oH8gdyrQoGse+pNcxresJ6mdDMJOuzxiITfpvc5tmOFGa8CzsHb\nDS5wZHMcwSbpbHJ/1H6bTdswk9iw9YShhQT4twW618Q1hqIkx/Exj9QD/txcHwjF\nVzyKOsNAXBvLBOS/ehQ8Z7EModhwZWUZ4TohFbI96JRr/GuQGD8pinPA1dALeC4U\nwBSJqPV9jyt6Yobz3xHfWb4scxZ1IY3NVyds0GNyKbz/CHp+pphXLlj0OkjeklUM\nrygVYZjHryP1BWIt3tRNoNddy5NWQdmG4B+xHUp22rXr0iWzgJmceXqhMxL+6M1j\n3QIDAQAB\n-----END PUBLIC KEY-----", `\n`, "\n", -1)))
	se := SubscriptionEnvironment{JwtPublicKey: jwtPublicKey}

	invalidPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAsrWxna5lwLHuqbw4HZpMx3sGhG23hs7BRd/QraC/nr6n/ZaX\nSg2c0t+04FrBYeR5saq2ISgpVusR8EQ3S+Pdc07AuIptRAJCWXuv+0kMoIf2uCGM\n9+x3aaFs9TnBTeOaQLPOKmvU0qKXLthb2Llh/57Wlpqt7mMVIB5SDi8jTwNLtsaH\n6A9JkjCSa1JvAxWAnyKAH10R7tYK3rW4uv1qXkiNwjTaxG27IBKzamiDiIqhFsk3\nGmkapvaB6GWkLmDFzbgl1Z08jBmuIvcerX0qWUr7ZwvOdB/MeWIvg3arhToiixy7\nNA6SasQz8JgMXj/sgo/2xM084wmqwq+XCoWLUQIDAQABAoIBAGtQA9pF8UAnGMPA\nDIpqL5TR4XYAVGaVHSYSYkMA9lAi+MmkdjC3v+Y5A2p91QYtpo3zju6WKUzSV7FJ\npVLjAAXP0pZ+OWLPYHxPc7uXgAed3Z9wNjBiRMqfbshK6DMXa3dTAYgjvGbHz7UH\now+fMqPHA+Dn/W3a/cKNsoRl9fPpa6TTOe2LvwNNmg599k8uBwZ5Jq8ObTVdfyHW\n3DzvSMgNlGbV21fqJt7REvvOy1McX11ic5VTzRDro+NA11h1d/lWqKtH6WnM7NPr\npZwWY4P+M0Sj2MYSmhpnrQxiFHCkge9dVieDfUDBV02n1OgeCGYoUtHpoaP/3H6i\naYbBZWECgYEA3q3Uz8TzWIhHqrMEhM8mWLzJ6ZYpH2l0oLbhFCSf5NotXbL7ofr3\nc2kOJr9j7ATaLtt9x74/IGiwpQGIhIh5sBPLqvVcqwZ6O9nn2JJLU2jNsrlFCczd\nsyRUmOSy+TGR8eqrQdaylHtula2X1k0e6sLpNr1zUb5CgsOm83gLzL0CgYEAzXOH\nyTo3nEVgAP97iT57s+VPqyfXTPSmxLKBZ1mfzAJXSp+YVmzspBs/f+Y9zY8dNSnq\nEvNbxxUyBGhf7GSMLi4WQWJ6zBUVSbbf/DKPlH6JIL7Mpw5NZ1JB4Srw6hVbwHRR\nVEfDzsI2JXgnG5d3GroQuWEG6rrQVpG3hRT2BCUCgYEAxy/ZgCzdvGDTMpdFwIMH\n+zKMrpy6ljWftK9B8OG+AVlEYV9cBY3X0W200eY7vICupmGCRq0gDJ54/HC9rEfk\nCiP6+CbYyMdXibgm0qwyIx9JeMiPP+/4lGk4Hzfb/FpiAXL31EH3pigvSEZq3rBH\nqpnFQUkAIau+FPhsm7bTX0UCgYBM+vjp4TJtD2GFfjTfm9Xl9gn/65G0eAb0tJ3g\nB6jkXAwxHdOKro2Mf9kvJyoJF5KuFOJr71t4IOz40dL/VD1iIEqefTPdvBiZ4jfS\n228JCNCAwH0WKzm9eQjOQbS7QP18AxlmJu2rTwHS0E52/C3YidcyXrSkxxLBHSBS\ncSQtMQKBgD1Sws3uyS7TRnkew9qvUuN8sQWrbGqrEf3cgKKkpnATqhZFOxBZGD83\nCJaMG5UQR1w8fYTeIBTw+ZCLVauFMDhGZve0ZjPs/KEibpyzqam8fUbyRlQmDEEl\noSMBSKcw7mAETw9+Kp/u8chvDVGJOxOCaBn4dM/FSP82AWpfYZtT\n-----END RSA PRIVATE KEY-----"))
	validPrivateKey, _ := jwt.ParseRSAPrivateKeyFromPEM([]byte("-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAxmovukuj5QtrJADJe5uyoNZ8/TQ9Ng066oH8gdyrQoGse+pN\ncxresJ6mdDMJOuzxiITfpvc5tmOFGa8CzsHbDS5wZHMcwSbpbHJ/1H6bTdswk9iw\n9YShhQT4twW618Q1hqIkx/Exj9QD/txcHwjFVzyKOsNAXBvLBOS/ehQ8Z7EModhw\nZWUZ4TohFbI96JRr/GuQGD8pinPA1dALeC4UwBSJqPV9jyt6Yobz3xHfWb4scxZ1\nIY3NVyds0GNyKbz/CHp+pphXLlj0OkjeklUMrygVYZjHryP1BWIt3tRNoNddy5NW\nQdmG4B+xHUp22rXr0iWzgJmceXqhMxL+6M1j3QIDAQABAoIBAQCeuhvsYM6AUR4o\n8yg75laELJJYzQ2azAKxz8L88FdhIPOnPc0vo/M6P/DRTHK53QtsNz/kBir5Kaw0\n27jXRmXCqb/n297I9iQOSZrRl2cOZab634LRJoVAMLX3VHIgYiqfxd6+xMjtUqLP\no+FeQln8a8X5NHGsPd+vzn8PXljkaWnrFEdjTJdGjDXCZK+eVyIOZVmpKhBf7W4H\n7wMjMzvovnQmBsSchMb49HfWNe0vppMw4rdhJu8ru1hY+VfiwOdQOgAiqWoOVpUs\n9mhe2suSBKJOi2rLFtZjO8eEN9puVVas1ttEwO22oQJt0xM2u1XI79NYnPj8fLKk\nNisRz9IBAoGBAOmDhAi1msvKD/edrTYTDkzFkzqA006plu0stCJKKeHCIDiXRAt3\nF1HJJxjcglL7Q/+BOEyaIBSwwzDfAR/QwBx2UVpjoXxdXdzmsm/aMycHDYfGDmGP\nhck2497hm8X/eVySEpASv9YGZYtU9Xv2DtbE4hMhGC4gJ/pcVAhzASq9AoGBANmF\nbQ+y+VgOR8UzX7SsaoWcVXkCap3zIqR8OG5zxYHcSjtKAgUci3CzfagjA3iTba5Y\nDCxgm61tA913nbuBRw5G6zSRMzknqKGcsj0pXGHakn6rgk5rZG28KYc4cYyKpzQ/\n+uKGssldYEIb5m0p1yqYAifDfGxUfYcYfDUZET+hAoGBALDmYO/4I/yeZto/KSj9\n5qdiRdbcIThGYX/rjcssQ+4zEhXNAk9tOM1MhcHfyxryHuFBE1V0rTj/b99mEYP4\nsJDfUWIYeA2u7ZybaDI0Kuw1+5oQAHUINWHpo1cFsuycTWRDhKyAh7OrxOF4yz/N\neBBbtqinOZo34hFYQJDmDsxBAoGBAKFHnKWXPelcLTq57Kw5aoHGeFlQwVx7eaQb\ntnuuuzKd8ywio3zGvVzCuNiBnYE5TomGHwCIYOUlf0gl+H2eTOD5FEvVnPUzwoSR\nelZ+5FBpj1T8NZGPbtcuPFxWLVrXM0I6bfqnhMhow4ZAyYkHDNI4AuEYwJhRzQDt\n1qpH+9IBAoGAF1L3YsNhgdnrLXE1pN8yYTfrlBw8i++/bCrJxM8Df9YUAjCLfj9X\naE/3O49Ap7LComPYzHRwG4MN5xdCxBmBLJo9nXK8DTpMhJ9A4neUxfEptP/psNws\nsbOAUvzut9BTqxwbHy5LIpjhcTvf5KSgsdjijzM+BgoqUGScl7vscNo=\n-----END RSA PRIVATE KEY-----"))

	configurations.UpdateActiveConfigurationTable(commons.ConfigurationTable{
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
	})

	tests := []test{
		{
			authorization:         "",
			expIsHandshakeErr:     true,
			expConfigurationTable: nil,
		},
		{
			authorization: fmt.Sprintf("Bearer %s", func() string {
				token, _ := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{"usr": "sys.admin@zippo.apps.com"}).SignedString(invalidPrivateKey)
				return token
			}()),
			expIsHandshakeErr:     true,
			expConfigurationTable: nil,
		},
		{
			authorization: fmt.Sprintf("Bearer %s", func() string {
				token, _ := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{"usr": "sys.admin@zippo.apps.com"}).SignedString(validPrivateKey)
				return token
			}()),
			expIsHandshakeErr: false,
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

	router := gin.Default()
	router.GET("/", se.Subscription)
	go func() {
		_ = router.Run("127.0.0.1:8080")
	}()

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			connection, _, retHandshakeErr := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080", http.Header{"Authorization": []string{tt.authorization}})
			retIsHandshakeErr := retHandshakeErr != nil
			if tt.expIsHandshakeErr != retIsHandshakeErr {
				t.Errorf("handshake error mismatch: expected '%v', received '%v'", tt.expIsHandshakeErr, retIsHandshakeErr)
			}
			if !retIsHandshakeErr {
				_ = connection.SetReadDeadline(time.Now().Add(time.Second))
				_, message, _ := connection.ReadMessage()
				var retConfigurationTable commons.ConfigurationTable
				_ = json.Unmarshal(message, &retConfigurationTable)
				if !reflect.DeepEqual(tt.expConfigurationTable, retConfigurationTable) {
					t.Errorf("configuration table mismatch: expected '%v', received '%v'", tt.expConfigurationTable, retConfigurationTable)
				}
			}
		})
	}
}
