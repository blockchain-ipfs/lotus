package main

import (
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"golang.org/x/xerrors"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/extern/sector-storage/fr32"
)

var fr32Cmd = &cli.Command{
	Name:        "fr32",
	Description: "fr32 tools",
	Subcommands: []*cli.Command{
		fr32WriteCmd,
		fr32ReadCmd,
	},
}

var fr32WriteCmd = &cli.Command{
	Name:  "write",
	Usage: "write fr32",
	Action: func(context *cli.Context) error {
		w := fr32.NewPadWriter(os.Stdout)
		if _, err := io.Copy(w, os.Stdin); err != nil {
			return err
		}
		return w.Close()
	},
}

var fr32ReadCmd = &cli.Command{
	Name:  "read",
	Usage: "read fr32",
	Action: func(context *cli.Context) error {
		st, err := os.Stdin.Stat()
		if err != nil {
			return err
		}

		pps := abi.PaddedPieceSize(st.Size())
		if pps == 0 {
			return xerrors.Errorf("zero size input")
		}

		if err := pps.Validate(); err != nil {
			return err
		}

		r, err := fr32.NewUnpadReader(os.Stdin, pps)
		if err != nil {
			return err
		}
		if _, err := io.Copy(os.Stdout, r); err != nil {
			return err
		}
		return nil
	},
}
