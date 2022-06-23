package main

import (
	"github.com/golang-jwt/jwt"
)

type Response struct {
	Code    uint64 `json:"code"`
	Message string `json:"message"`
}

type LogInUserDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type DBUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Passhash string `json:"passhash"`
	Role     string `json:"role"`
}

type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type JWTClaim struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

type DeployContractDetails struct {
	Nonce     int64 `json:"nonce"`
	GasPrice  int64 `json:"gasprice"`
	GasLimit  int64 `json:"gaslimit"`
	FundValue int64 `json:"fundvalue"`
}

type TransactionDetails struct {
	ID          string `json:"id"`
	Hash        string `json:"hash"`
	FromAddress string `json:"fromaddress"`
	ToAddress   string `json:"toaddress"`
	TxType      string `json:"txtype"`
}

type Item struct {
	ID       string `json:"id"`
	Brand    string `json:"brand"`
	Price    int    `json:"price"`
	OwnerID  string `json:"ownerID"`
	TokenURI string `json:"tokenURI"`
}

type MintNFTDetails struct {
	TokenID  int64  `json:"tokenID"`
	TokenURI string `json:"tokenURI"`
}
