package main

import (
	"fmt"
	fac "github.com/FactomProject/factom"
	"github.com/FactomProject/factomd/common/primitives"
)

func main() {
	signer := "FA33h3E6rmW5EChgYjN5b6PvpMB6aQw9qSZFeoHDX6X9H4Bg9X9h"
	data := []byte("hello Factom Project")

	sig, err := fac.SignData(signer, data)
	if err != nil {
		fmt.Printf("[ERROR] unable to sign data, err : %v\n", err)
		return
	}

	fmt.Println("[Signature] : ", string(sig.Signature))
	fmt.Println("[Pubkey] : ", string(sig.PubKey))

	err2 := primitives.VerifySignature(data, sig.PubKey, sig.Signature)
	if err2 != nil {
		fmt.Printf("[ERROR] unable to verify signature : %v \n", err2)
		return
	}
	fmt.Println("[INFO] verified ...")

}
