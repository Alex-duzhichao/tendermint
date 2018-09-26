package blockchain

import (
	"github.com/tendermint/go-amino"
	"github.com/Alex-duzhichao/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterBlockchainMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
