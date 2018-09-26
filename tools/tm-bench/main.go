package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/kit/log/term"

	"github.com/Alex-duzhichao/tendermint/libs/log"
	tmrpc "github.com/Alex-duzhichao/tendermint/rpc/client"
)

var logger = log.NewNopLogger()

func main() {
	var durationInt, txsRate, connections, txSize int
	var verbose bool
	var outputFormat, broadcastTxMethod string

	flagSet := flag.NewFlagSet("tm-bench", flag.ExitOnError)
	flagSet.IntVar(&connections, "c", 1, "Connections to keep open per endpoint")
	flagSet.IntVar(&durationInt, "T", 10, "Exit after the specified amount of time in seconds")
	flagSet.IntVar(&txsRate, "r", 1000, "Txs per second to send in a connection")
	flagSet.IntVar(&txSize, "s", 250, "The size of a transaction in bytes, must be greater than or equal to 40.")
	flagSet.StringVar(&outputFormat, "output-format", "plain", "Output format: plain or json")
	flagSet.StringVar(&broadcastTxMethod, "broadcast-tx-method", "async", "Broadcast method: async (no guarantees; fastest), sync (ensures tx is checked) or commit (ensures tx is checked and committed; slowest)")
	flagSet.BoolVar(&verbose, "v", false, "Verbose output")

	flagSet.Usage = func() {
		fmt.Println(`Tendermint blockchain benchmarking tool.

Usage:
	tm-bench [-c 1] [-T 10] [-r 1000] [-s 250] [endpoints] [-output-format <plain|json> [-broadcast-tx-method <async|sync|commit>]]

Examples:
	tm-bench localhost:26657`)
		fmt.Println("Flags:")
		flagSet.PrintDefaults()
	}

	flagSet.Parse(os.Args[1:])

	if flagSet.NArg() == 0 {
		flagSet.Usage()
		os.Exit(1)
	}

	if verbose {
		if outputFormat == "json" {
			fmt.Fprintln(os.Stderr, "Verbose mode not supported with json output.")
			os.Exit(1)
		}
		// Color errors red
		colorFn := func(keyvals ...interface{}) term.FgBgColor {
			for i := 1; i < len(keyvals); i += 2 {
				if _, ok := keyvals[i].(error); ok {
					return term.FgBgColor{Fg: term.White, Bg: term.Red}
				}
			}
			return term.FgBgColor{}
		}
		logger = log.NewTMLoggerWithColorFn(log.NewSyncWriter(os.Stdout), colorFn)

		fmt.Printf("Running %ds test @ %s\n", durationInt, flagSet.Arg(0))
	}

	if txSize < 40 {
		fmt.Fprintln(
			os.Stderr,
			"The size of a transaction must be greater than or equal to 40.",
		)
		os.Exit(1)
	}

	if broadcastTxMethod != "async" &&
		broadcastTxMethod != "sync" &&
		broadcastTxMethod != "commit" {
		fmt.Fprintln(
			os.Stderr,
			"broadcast-tx-method should be either 'sync', 'async' or 'commit'.",
		)
		os.Exit(1)
	}

	var (
		endpoints     = strings.Split(flagSet.Arg(0), ",")
		client        = tmrpc.NewHTTP(endpoints[0], "/websocket")
		initialHeight = latestBlockHeight(client)
	)
	logger.Info("Latest block height", "h", initialHeight)

	transacters := startTransacters(
		endpoints,
		connections,
		txsRate,
		txSize,
		"broadcast_tx_"+broadcastTxMethod,
	)

	// Wait until transacters have begun until we get the start time
	timeStart := time.Now()
	logger.Info("Time last transacter started", "t", timeStart)

	duration := time.Duration(durationInt) * time.Second

	timeEnd := timeStart.Add(duration)
	logger.Info("End time for calculation", "t", timeEnd)

	<-time.After(duration)
	for i, t := range transacters {
		t.Stop()
		numCrashes := countCrashes(t.connsBroken)
		if numCrashes != 0 {
			fmt.Printf("%d connections crashed on transacter #%d\n", numCrashes, i)
		}
	}

	logger.Debug("Time all transacters stopped", "t", time.Now())

	stats, err := calculateStatistics(
		client,
		initialHeight,
		timeStart,
		durationInt,
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	printStatistics(stats, outputFormat)
}

func latestBlockHeight(client tmrpc.Client) int64 {
	status, err := client.Status()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return status.SyncInfo.LatestBlockHeight
}

func countCrashes(crashes []bool) int {
	count := 0
	for i := 0; i < len(crashes); i++ {
		if crashes[i] {
			count++
		}
	}
	return count
}

func startTransacters(
	endpoints []string,
	connections,
	txsRate int,
	txSize int,
	broadcastTxMethod string,
) []*transacter {
	transacters := make([]*transacter, len(endpoints))

	wg := sync.WaitGroup{}
	wg.Add(len(endpoints))
	for i, e := range endpoints {
		t := newTransacter(e, connections, txsRate, txSize, broadcastTxMethod)
		t.SetLogger(logger)
		go func(i int) {
			defer wg.Done()
			if err := t.Start(); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			transacters[i] = t
		}(i)
	}
	wg.Wait()

	return transacters
}
