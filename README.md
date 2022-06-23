# LuxuryNFT

Description
The aim of this project is to create an application using GoLang to connect to an ethereum node to mint, burn and transfer ERC721 NFTs.

Pre-requisities
The basics required to get the application working on local are :
  docker
  postgres
  go
  node
  truffle
  
User-Roles
The application has 2 user roles
  Manufacturer 
      -- Can add items to the DB
      -- Can deploy the ERC721 contract
      -- Can mint ERC721 NFTs
      -- Can transfer the minted NFTs
      -- Can burn the NFTs
  Consumer
      -- Can transfer an owned NFT
      
Use-Case
This application can be used to mint and assign NFTs to luxury physical assets in order to help prevent farudlent duplicates from being sold.

Set-Up
Clone this repo using


Navigate to the database directory and start the postgres db container
```
cd database
docker-compose up

```

Navigate to the project root and run the application
```
go run .

```
You should see a message in the terminal 
```
Successfully connected to database!
Successfully connected to localhost!

```

This indicates that you can now access the application using 
```
http://localhost:8080/api/v1/{}

```

