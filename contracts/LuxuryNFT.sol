// SPDX-License-Identifier: MIT
pragma solidity >=0.4.20 <0.9.0;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";

contract LuxuryNFT is ERC721URIStorage, Ownable{

    constructor() ERC721("Luxury NFT", "KLD"){
    }

    function mintNFT(string memory tokenURI, uint256 tokenID) public onlyOwner returns (uint256){
        _mint(msg.sender, tokenID);
        _setTokenURI(tokenID, tokenURI);
        return tokenID;
    }

    function burnNFT(uint256 tokenID) public {
        _burn(tokenID);
    }

    function transferNFT(uint256 tokenID, address from, address to )public onlyOwner{
        transferFrom(from, to, tokenID);
    }
}