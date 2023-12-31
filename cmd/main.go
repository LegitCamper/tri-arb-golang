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
	// create channels and opens websockets
	platform := platforms.Handler(crypto.New(true))

	func() {
		for x := range platform.Market_conn {
			log.Println(x)
		}
		for x := range platform.User_conn {
			log.Println(x)
		}
	}()

	// pass reader to detector and writer for orders
	detector.Test()
	time.Sleep(time.Hour)

}
