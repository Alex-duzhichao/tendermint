package types

import (
	"github.com/tendermint/go-amino"
	"github.com/Alex-duzhichao/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	types.RegisterBlockAmino(cdc)
}
