package main

import (
	// "database/sql"
	db "LuxuryNFT/database"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	//Create an instance to connect to the postgres DB
	if err := db.ConnectDB(); err != nil {
		panic(err)
	}
	defer db.Close()
	err := db.DB.Ping()
	if err != nil {
		fmt.Println("Error connecting to database", err)
		panic(err)
	}
	fmt.Println("Successfully connected to database!")
	if err != nil {
		fmt.Println(err)
	}


	//Create a connection the ethereum node
	endpoint := "{end point of your local ganache or remote alchemy/infura app}"
	privateKey := "{private key of the account where the contract is to be deployed}"
	//ID of your testnet 
	chainId := big.NewInt(4)
	client, err := NewClient(endpoint, privateKey, chainId)
	if err != nil {
		fmt.Println("Error connecting to blockchain network", err)
		panic(err)
	}

	// Create a mux router
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/sign_up", SignUp)
	r.HandleFunc("/api/v1/sign_in", SignIn)
	r.Handle("/api/v1/luxury_nft/{module}", ClientAPIHandler{client})
	fmt.Println("Successfully connected to localhost!")
	log.Fatal(http.ListenAndServe(":8080", r))
}
