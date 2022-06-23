# LuxuryNFT

Description<br />
The aim of this project is to create an application using GoLang to connect to an ethereum node to mint, burn and transfer ERC721 NFTs.<br />

Pre-requisities<br />
The basics required to get the application working on local are :<br />
  --docker<br />
  --postgres<br />
  --go<br />
  --node<br />
  --truffle<br />
  <br />
  
User-Roles<br />
The application has 2 user roles<br />
Manufacturer<br />
  --Can add items to the DB<br />
  --Can deploy the ERC721 contract<br />
  --Can mint ERC721 NFTs<br />
  --Can transfer the minted NFTs<br />
  --Can burn the NFTs<br />
Consumer<br />
  -- Can transfer an owned NFT<br />
      
Use-Case<br />
This application can be used to mint and assign NFTs to luxury physical assets in order to help prevent farudlent duplicates from being sold.<br />
<br />

Set-Up<br />
Clone this repo<br />
Navigate to the database directory and start the postgres db container<br />
```
cd database
docker-compose up

```

Navigate to the project root. In the main.go file add the ```endpoint```, ```privateKey``` and ```chainId``` details of your ethereum node and run the application<br />
```
go run .

```
You should see a message in the terminal<br />
```
Successfully connected to database!
Successfully connected to localhost!

```

This indicates that you can now access the application using<br />
```
http://localhost:8080/api/v1/{}

```

