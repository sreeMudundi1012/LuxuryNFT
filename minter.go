package main

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
)

func (c *Contract) MintToken(client Client, tokenURI string, tokenID *big.Int) (*types.Transaction, error) {
	
	tx, err := c.Instance.MintNFT(client.Auth, tokenURI, tokenID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create signed mint transaction")
	}

	return tx, nil
}
