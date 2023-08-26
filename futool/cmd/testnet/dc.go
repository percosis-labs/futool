package testnet

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func DcCmd() *cobra.Command {
	dcCmd := &cobra.Command{
		Use:   "dc",
		Short: "A convenience command that runs `docker-compose <...args>` on the generated config.",
		Example: `Follow logs of running chain ("--" necessary to correctly interpret docker-compose flags):
$ futool testnet dc -- logs -f

Get a shell in the fury node container:
$ futool testnet dc exec furynode /bin/bash

Run some fury cli commands:
$ futool testnet dc exec furynode fury keys add magic-account
$ futool testnet dc exec furynode -- fury tx bank send whale <address> 1000000000ufury --gas-prices 1000ufury -y

If the testnet was started with --ibc, you can IBC transfer coins ot Fury from the ibc chain.
Transfer 1 ATOM from the ibc chain to random account on Fury:
$ futool t dc exec ibcnode -- fury tx ibc-transfer transfer \
  transfer channel-0 fury1dwktgf8jcuusc885myae3hjk63jrc6tsz69muu 1000000uatom \
	--from whale2 --gas auto --gas-adjustment 1.2 --gas-prices 0.01uatom

Checking the balance on the Fury Chain:
$ futool t dc exec furynode -- fury q bank balances fury1dwktgf8jcuusc885myae3hjk63jrc6tsz69muu
balances:
- amount: "1000000"
  denom: ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2
`,
		Args: cobra.ArbitraryArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			cmd := []string{"docker-compose", "--file", generatedPath("docker-compose.yaml")}
			cmd = append(cmd, args...)
			fmt.Println("running:", strings.Join(cmd, " "))
			if err := replaceCurrentProcess(cmd...); err != nil {
				return fmt.Errorf("could not run command: %v", err)
			}
			return nil
		},
	}

	return dcCmd
}
