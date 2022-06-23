package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"net/mail"
	"strings"
	"time"
	db "LuxuryNFT/database"
	"github.com/aidarkhanov/nanoid/v2"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type ClientAPIHandler struct {
	Client
}

var deployedContract Contract

func randomIDGenerator() string {
	randomID, err := nanoid.New()
	if err != nil {
		log.Fatalln(err)
	}
	return randomID
}

func generatehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func validateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false
	}
	return true
}

func generateJWT(email, role string) (string, error) {
	var mySigningKey = []byte("secretkey")
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func checkPasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func validateToken(loginJWTToken string) (claims jwt.MapClaims, valid bool) {

	var mySigningKey = []byte("secretkey")
	token, err := jwt.Parse(loginJWTToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error in parsing")
		}
		return mySigningKey, nil
	})

	if err != nil {
		fmt.Errorf("Your token has expired: %s", err.Error())
		return nil, false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		log.Printf("Invalid JWT Token")
		return nil, false
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	var signUpUser SignUpUser
	var dbUser DBUser
	var err error

	err = json.NewDecoder(r.Body).Decode(&signUpUser)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Malformed request",
		})
		return
	}

	isEmailValid := validateEmail(signUpUser.Email)
	if !isEmailValid {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Email not in the correct format",
		})
		return
	}

	row := db.DB.QueryRow("SELECT * FROM users where email= $1", signUpUser.Email)
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Passhash, &dbUser.Role)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error querying DB for users", err)
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "Error querying DB for users",
		})
		return
	}

	//checks if email is already registered
	if dbUser.Email != "" {
		fmt.Println("Email already in use")
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Email already in use",
		})
		return
	}

	//check to see the role
	if strings.ToLower(signUpUser.Role) != ("manufacturer") && strings.ToLower(signUpUser.Role) != ("consumer") {
		fmt.Println("Role is not valid")
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Role is not valid",
		})
		return
	}

	var newUser = new(DBUser)
	//create newUser details
	newUser.Passhash, err = generatehashPassword(signUpUser.Password)
	if err != nil {
		log.Fatalln("error in password hash creation")
	}

	//register new user
	fmt.Println("New User Registration")
	newUser.ID, err = nanoid.New()
	if err != nil {
		fmt.Println("Error generating UUID", err)
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "Error generating UUID",
		})
		return
	}

	_, err = db.DB.Exec("INSERT INTO users(id,username, email, passhash, role) VALUES ($1,$2, $3, $4, $5)", newUser.ID, signUpUser.Username, signUpUser.Email, newUser.Passhash, signUpUser.Role)
	if err != nil {
		fmt.Println("Error inserting new user to DB", err)
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "Error inserting new user to DB",
		})
		return
	}
	fmt.Println("New User Registration Successful")
	json.NewEncoder(w).Encode(Response{
		Code:    200,
		Message: "New User Registration Successful",
	})
	return
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var loginDetails LogInUserDetails
	var dbUser DBUser

	err := json.NewDecoder(r.Body).Decode(&loginDetails)
	if err != nil {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Malformed request",
		})
		return
	}
	isEmailValid := validateEmail(loginDetails.Email)
	if !isEmailValid {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Email not in the correct format",
		})
		return
	}

	row := db.DB.QueryRow("SELECT * FROM users where email= $1", loginDetails.Email)
	err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Passhash, &dbUser.Role)

	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Error querying DB for users", err)
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "Error querying DB for users",
		})
		return
	}

	//checks if email is registered
	if dbUser.Email == "" {
		fmt.Println(err)
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Email not registered. Please sign-up",
		})
		return
	}

	check := checkPasswordHash(loginDetails.Password, dbUser.Passhash)

	if !check {
		fmt.Println("Username or Password is Incorrect")
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "Username or Password is Incorrect",
		})
		return
	}

	validToken, err := generateJWT(dbUser.Email, dbUser.Role)
	if err != nil {
		fmt.Println("Error generating JWT token", err)
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "Error generating JWT token",
		})
		return
	}

	var token Token
	token.Email = dbUser.Email
	token.Role = dbUser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (c ClientAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	module := vars["module"]

	// Set our  Response header
	w.Header().Set("Content-Type", "application/json")

	//ensure the request has a JWT Token
	JWTToken := r.Header.Get("Authorization")
	if JWTToken == "" {
		fmt.Println("No JWT Token found")
		json.NewEncoder(w).Encode(Response{
			Code:    500,
			Message: "No JWT Token found",
		})
		return
	}
	TokenArray := strings.Split(JWTToken, " ")
	//validate and decode the JWT Token
	claims, valid := validateToken(TokenArray[1])
	if !(valid) {
		fmt.Println("JWT Token error")
		json.NewEncoder(w).Encode(Response{
			Code:    400,
			Message: "JWT Token err",
		})
		return
	}

	// Handle each request using the module parameter:
	switch module {
	case "items":
		var item Item
		var dbUser DBUser
		var err error

		if claims["role"] != "manufacturer" {
			fmt.Println("User role err")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "User role err. Only manufacturer can perform this action",
			})
			return
		}

		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}

		//retrieve the ownerID of the logged in user
		row := db.DB.QueryRow("SELECT * FROM users where email= $1", claims["email"])
		err = row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Passhash, &dbUser.Role)

		item.ID = randomIDGenerator()

		//insert item to DB
		_, err = db.DB.Exec("INSERT INTO luxury_items(id,brand, price, ownerID, tokenURI) VALUES ($1,$2, $3, $4, $5)", item.ID, item.Brand, item.Price, &dbUser.ID, item.TokenURI)
		if err != nil {
			fmt.Println("Error inserting new item to DB", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error inserting new item to DB",
			})
			return
		}

		fmt.Println("New Item saved to DB Successfully")
		json.NewEncoder(w).Encode(Response{
			Code:    200,
			Message: "New Item saved to DB Successfully",
		})
		return

	case "deploy":
		var deployContractDetails DeployContractDetails
		var transactionDetails TransactionDetails
		var err error

		if claims["role"] != "manufacturer" {
			fmt.Println("User role err")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "User role err. Only manufacturer can perform this action",
			})
			return
		}

		err = json.NewDecoder(r.Body).Decode(&deployContractDetails)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}

		c.SetNonce(big.NewInt(deployContractDetails.Nonce))
		c.SetFundValue(big.NewInt(deployContractDetails.FundValue))
		c.SetGasLimit(uint64(deployContractDetails.GasLimit))
		c.SetGasPrice(big.NewInt(deployContractDetails.GasPrice))

		//deploying the smart contract
		contract, err := c.DeployContract()
		if err != nil {
			fmt.Println("Error deploying contract blockchain network", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error deploying contract blockchain network",
			})
			return
		}
		fmt.Println("Contract address: ", contract.Address.Hex(), contract)

		transactionDetails.FromAddress = contract.Address.Hex()
		transactionDetails.TxType = "DEPLOY"
		transactionDetails.Hash = ""
		transactionDetails.ToAddress = ""

		transactionDetails.ID = randomIDGenerator()
		if err != nil {
			fmt.Println("Error creating a random ID for transaction", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error creating a random ID for transaction",
			})
			return
		}

		_, err = db.DB.Exec("INSERT INTO transactions(id, fromaddress, txtype) VALUES ($1,$2, $3)", transactionDetails.ID, transactionDetails.FromAddress, transactionDetails.TxType)
		if err != nil {
			fmt.Println("Error inserting new user to DB", err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Error inserting transaction to DB",
			})
			return
		}
		fmt.Println("Contract deployed sucessfully at address:", contract.Address.Hex())
		json.NewEncoder(w).Encode(Response{
			Code:    200,
			Message: "Contract deployed sucessfully",
		})
		return

	case "mintNFT":
		var mintNFTDetails MintNFTDetails
		var tx TransactionDetails

		if claims["role"] != "manufacturer" {
			fmt.Println("User role err")
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "User role err. Only manufacturer can perform this action",
			})
			return
		}

		err := json.NewDecoder(r.Body).Decode(&mintNFTDetails)
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    400,
				Message: "Malformed request",
			})
			return
		}

		txType := "DEPLOY"
		row := db.DB.QueryRow("SELECT * FROM transactions where txtype= $1", txType)
		err = row.Scan(&tx.ID, &tx.Hash, &tx.FromAddress, &tx.ToAddress, &tx.TxType)

		fmt.Println("************", tx)

		// var cl Contract
		// var main Main
		// var cl = Contract{"", &main}

		// tx1 , err := deployedContract.MintToken(c, "", big.NewInt(9))
		if err != nil {
			fmt.Println(err)
			json.NewEncoder(w).Encode(Response{
				Code:    500,
				Message: "Internal server error",
			})
			return
		}

		// fmt.Println("^^^^^^^", tx1)

		json.NewEncoder(w).Encode(Response{
			Code:    200,
			Message: "Successfully minted NFT",
		})

		// case "burnNFT":
		// JWTToken := r.Header.Get("Authorization")
		// TokenArray := strings.Split(JWTToken, " ")
		// if TokenArray[1] == "" {
		// 	fmt.Println("No JWT Token found")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    500,
		// 		Message: "No JWT Token found",
		// 	})
		// 	return
		// }
		// claims, valid := validateToken(TokenArray[1])
		// if !(valid) {
		// 	fmt.Println("JWT Token error")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    400,
		// 		Message: "JWT Token err",
		// 	})
		// 	return
		// }

		// if claims["role"] != "manufacturer" {
		// 	fmt.Println("User role err")
		// 	json.NewEncoder(w).Encode(  Response{
		// 		Code:    400,
		// 		Message: "User role err. Only manufacturer can perform this action",
		// 	})
		// 	return
		// }
	}
	return
}
