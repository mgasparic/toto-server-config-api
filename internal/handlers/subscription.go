package handlers

import (
	"api/internal/commons"
	"api/internal/configurations"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
)

type SubscriptionEnvironment commons.UpdaterEnvironment

var upgrader = websocket.Upgrader{HandshakeTimeout: 5 * time.Second, EnableCompression: true}

func (se SubscriptionEnvironment) Subscription(ctx *gin.Context) {
	auth := ctx.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	_, err := jwt.Parse(auth[7:], func(token *jwt.Token) (interface{}, error) { return se.JwtPublicKey, nil })
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	connection, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Print(err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if connection == nil {
		log.Print("connection is nil")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer func() {
		_ = connection.Close()
	}()

	address := commons.Address(ctx.Request.RemoteAddr)
	channel := make(chan commons.ConfigurationTable)
	go configurations.AddSubscriber(address, channel)
	defer configurations.DelSubscriber(address)
	for {
		select {
		case newConfigurationTable, ok := <-channel:
			if !ok {
				return
			}
			err := connection.WriteJSON(&newConfigurationTable)
			if err != nil {
				log.Print(err)
				return
			}
		case <-time.After(time.Minute):
			err := connection.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Print(err)
				return
			}
		}
		err := connection.SetWriteDeadline(time.Now().Add(2 * time.Minute))
		if err != nil {
			log.Print(err)
			return
		}
	}
}
