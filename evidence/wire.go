package evidence

import (
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/Alex-duzhichao/tendermint/crypto/encoding/amino"
	"github.com/Alex-duzhichao/tendermint/types"
)

var cdc = amino.NewCodec()

func init() {
	RegisterEvidenceMessages(cdc)
	cryptoAmino.RegisterAmino(cdc)
	types.RegisterEvidences(cdc)
}
