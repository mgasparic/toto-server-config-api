package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type EnvVars struct {
	JwtPrivateKeyPath string `envconfig:"JWT_PRIVATE_KEY_PATH"`
	UsrClaim          string `envconfig:"USR_CLAIM" default:""`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var envVars EnvVars
	err := envconfig.Process("", &envVars)
	if err != nil {
		log.Fatal(err)
	}

	keyRaw, err := os.ReadFile(envVars.JwtPrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyRaw)
	if err != nil {
		log.Fatal(err)
	}

	if len(envVars.UsrClaim) > 0 {
		token, err := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), jwt.MapClaims{"usr": envVars.UsrClaim}).SignedString(key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("TOKEN:")
		fmt.Println(token)
	} else {
		token, err := jwt.New(jwt.GetSigningMethod("RS256")).SignedString(key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("TOKEN:")
		fmt.Println(token)
	}
}
