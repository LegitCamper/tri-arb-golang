package main

import (
	"log"
	"time"

	"tri-arb/internal/detector"
	"tri-arb/internal/platforms"
	"tri-arb/internal/platforms/crypto"
)

func main() {
	log.SetFlags(0)
	start()
}

func start() {
	// this platform call can change
	platform := crypto.New(true)
	platform.DownloadSymbols() // TODO: remove this
	time.Sleep(time.Hour)

	// create channels and opens websockets
	platformData := platforms.Handler(platform)

	// Creating go routine to listen to market data via websocket
	// channel and push data to PlatformApi
	go platform.ProccessMarketData(&platformData)

	go func() {
		for {
			log.Println(platformData.MarketData)
			time.Sleep(time.Second * 5)
		}
	}()

	// PlatformData to detector with marketData and writer to create orders
	detector.Test()
	time.Sleep(time.Hour)
}
