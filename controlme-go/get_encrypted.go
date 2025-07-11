package main

import (
	"fmt"
	"github.com/thecontrolapp/controlme-go/internal/auth"
	"github.com/thecontrolapp/controlme-go/internal/config"
)

func main() {
	cfg, _ := config.Load()
	crypto := auth.NewLegacyCrypto(cfg.Legacy.CryptoKey)
	encrypted, _ := crypto.Encrypt("testpass1")
	fmt.Println("Encrypted password for testpass1:", encrypted)
}
