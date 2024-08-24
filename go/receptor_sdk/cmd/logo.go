// This file is subject to the terms and conditions defined in
// file 'LICENSE.txt', which is part of this source code package.

package cmd

import (
	// "encoding/json"
	// "fmt"
	// "reflect"
	// "regexp"
	// "strings"

	"io"
	"os"

	"github.com/spf13/cobra"
)

const (
	logoUse   = "logo"
	logoShort = "Dump logo Trustero internal use"
)

type getlogo struct {
	cmd *cobra.Command
}

func (l *getlogo) getCommand() *cobra.Command {
	return l.cmd
}

func (l *getlogo) setup() {
	l.cmd = &cobra.Command{
		Use:          logoUse,
		Short:        logoShort,
		Args:         cobra.MinimumNArgs(0),
		RunE:         logo,
		SilenceUsage: true,
	}
	l.cmd.FParseErrWhitelist.UnknownFlags = true
}

// Cobra executes this function on logo command.
func logo(_ *cobra.Command, args []string) (err error) {
	svg, err := readSVGFile("logo.svg")
	println(svg)
	return
}

func readSVGFile(filePath string) (svgContent string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}

	// Return the contents as a string
	return string(data), nil
}
