package main

import (
	"flag"
	"os"

	"github.com/Alex-duzhichao/tendermint/crypto/ed25519"
	cmn "github.com/Alex-duzhichao/tendermint/libs/common"
	"github.com/Alex-duzhichao/tendermint/libs/log"

	"github.com/Alex-duzhichao/tendermint/privval"
)

func main() {
	var (
		addr        = flag.String("addr", ":26659", "Address of client to connect to")
		chainID     = flag.String("chain-id", "mychain", "chain id")
		privValPath = flag.String("priv", "", "priv val file path")

		logger = log.NewTMLogger(
			log.NewSyncWriter(os.Stdout),
		).With("module", "priv_val")
	)
	flag.Parse()

	logger.Info(
		"Starting private validator",
		"addr", *addr,
		"chainID", *chainID,
		"privPath", *privValPath,
	)

	pv := privval.LoadFilePV(*privValPath)

	rs := privval.NewRemoteSigner(
		logger,
		*chainID,
		*addr,
		pv,
		ed25519.GenPrivKey(),
	)
	err := rs.Start()
	if err != nil {
		panic(err)
	}

	cmn.TrapSignal(func() {
		err := rs.Stop()
		if err != nil {
			panic(err)
		}
	})
}
