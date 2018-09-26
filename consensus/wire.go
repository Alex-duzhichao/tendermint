package consensus

import (
	"github.com/tendermint/go-amino"
	"github.com/Alex-duzhichao/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterConsensusMessages(cdc)
	RegisterWALMessages(cdc)
	types.RegisterBlockAmino(cdc)
}
