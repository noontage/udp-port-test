package main

import (
	"breaker"
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"time"
)

var (
	peerMgr *breaker.PeerManager
)

func main() {
	go udpServer()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/udp-peers", hUdpPeers)
	e.Logger.Fatal(e.Start(":8081"))
}

func hUdpPeers(c echo.Context) error {
	return c.JSON(http.StatusOK, peerMgr.List())
}

func udpServer() {
	us := breaker.NewUDPServer("127.0.0.1", 39998, onUDP)

	go func() {
		peerMgr = breaker.NewPeerManager()
		for {
			peerMgr.Clean()
			time.Sleep(60 * time.Second)
		}
	}()

	if err := us.Listen(); err != nil {
		log.Fatalln(err)
	}
}

func onUDP(ctx context.Context, uq breaker.Request, data []byte) {
	var req struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(data, &req); err != nil {
		log.Printf(err.Error())
		return
	}

	log.Printf("receve - from: %s, key: %s", uq.FromAddr.String(), req.Key)
	peerMgr.Set(req.Key, uq.FromAddr)

	if err := uq.Conn.SetWriteDeadline(time.Now().Add(time.Minute)); err != nil {
		log.Printf(err.Error())
		return
	}

	_, err := uq.Conn.WriteTo([]byte(""), uq.FromAddr)
	if err != nil {
		log.Printf(err.Error())
		return
	}
}
