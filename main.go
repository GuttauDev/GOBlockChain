package main

import (
	"fmt"
	"log"

	wallet "./Wallet"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	w := wallet.NewWallet()
	//fmt.Println(w.PrivateKeyStr())
	//fmt.Println(w.PublicKeyStr())
	//fmt.Println(w.BlockchainAddress())

	t := wallet.NewTransaction(w.PrivateKey(), w.PublicKey(), w.BlockchainAddress(), "B", 1.0)
	fmt.Printf("Signature %s\n", t.GenerateSignature())

}
