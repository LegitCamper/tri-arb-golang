package main

import (
	"tri-arb/internal/platforms"
	"tri-arb/internal/platforms/crypto"
)

func main() {
	platforms.Create(crypto.Crypto)
}
