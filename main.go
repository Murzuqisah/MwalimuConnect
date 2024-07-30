package main

import (
	"net/http"
	"time"

	"mwalimu/functions"
)

func main() {
	http.HandleFunc("/login", functions.LoginHandler)
	http.HandleFunc("/logout", functions.LogoutHandler)
	http.HandleFunc("/protected", functions.AuthMiddleware(functions.ProtectedHandler))

	blockchain := functions.Blockchain{}
	genesisBlock := functions.Block{
		Timestamp:    time.Now().Unix(),
		Data:         []byte("Genesis Block"),
		PreviousHash: []byte{},
	}
	genesisBlock.Hash = functions.CalculateHash(genesisBlock)
	blockchain.Blocks = append(blockchain.Blocks, genesisBlock)

	block1Data := []byte("")
	http.ListenAndServe(":8080", nil)
}
