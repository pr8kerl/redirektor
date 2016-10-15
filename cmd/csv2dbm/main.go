package main

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/

func main() {

	c := cli.NewCLI("csv2dbm", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"import": dbmImportCmdFactory,
		"export": dbmExportCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)

}
