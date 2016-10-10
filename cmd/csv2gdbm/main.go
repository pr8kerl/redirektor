package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/cfdrake/go-gdbm"
	"github.com/mitchellh/cli"
	"io"
	"os"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/

func main() {

	c := cli.NewCLI("csv2gdbm", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"import": gdbmCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)

}

type GdbmCommand struct {
	Db  string
	Csv string
	Ui  cli.Ui
}

func gdbmCmdFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	return &GdbmCommand{
		Db:  "",
		Csv: "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorGreen,
		},
	}, nil
}

func (c *GdbmCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("import", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Db, "database", "", "the name of the gdbm file to use")
	cmdFlags.StringVar(&c.Csv, "csv", "", "the name of the csv to import")
	if err := cmdFlags.Parse(args); err != nil {
		fmt.Printf("parse error : %s\n", err)
		return 1
	}
	var requiredFlags = 0

	if c.Db != "" {
		requiredFlags++
		fmt.Printf("using db file : %s\n", c.Db)
	}
	if c.Csv != "" {
		requiredFlags++
		fmt.Printf("using csv file name : %s\n", c.Csv)
	}

	if requiredFlags < 2 {
		cmdFlags.Usage()
		return 1
	}

	csvfile, err := os.Open(c.Csv)
	if err != nil {
		fmt.Printf("error opening csv file %s : %s\n", c.Csv, err)
		return 1
	}
	// automatically call Close() at the end of current method
	defer csvfile.Close()
	//
	reader := csv.NewReader(csvfile)
	//reader.Comma = '\t'
	db, err := gdbm.Open(c.Db, "c")
	if err != nil {
		fmt.Printf("error opening gdbm file %s : %s\n", c.Db, err)
		return 1
	}
	defer db.Close()

	for {
		// read just one record, but we could ReadAll() as well
		record, err := reader.Read()
		// end-of-file is fitted into err
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("error:", err)
			return 1
		}
		fmt.Printf("inurl: %s\touturl: %s\n", record[0], record[1])

		dberr := db.Insert(record[0], record[1])
		if dberr != nil {
			fmt.Printf("gdbm insert error: %s\n", dberr)
		}
	}

	return 0

}

func (c *GdbmCommand) Help() string {
	return fmt.Sprintf("csv2gdbm : import csv data to a gdbm database file\n\n\t\t--database <filename>\tthe gdbm filename to use\n\t\t--csv <filename>\tthe csv filename to import data from\n\n")
}

func (c *GdbmCommand) Synopsis() string {
	return "import csv to a gdbm file"
}
