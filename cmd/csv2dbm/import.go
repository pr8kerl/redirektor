package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/jsimonetti/berkeleydb"
	"github.com/mitchellh/cli"
	"io"
	"os"
)

type DbmImportCommand struct {
	Db  string
	Csv string
	Ui  cli.Ui
}

func dbmImportCmdFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	return &DbmImportCommand{
		Db:  "",
		Csv: "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorGreen,
		},
	}, nil
}

func (c *DbmImportCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("import", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Db, "database", "", "the name of the dbm file to use")
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
		flag.Usage()
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

	db, err := berkeleydb.NewDB()
	if err != nil {
		fmt.Printf("error creating db connection : %s\n", err)
		return 1
	}

	//err = db.Open(c.Db, berkeleydb.DbBtree, berkeleydb.DbCreate)
	err = db.Open(c.Db, berkeleydb.DbHash, berkeleydb.DbCreate)
	if err != nil {
		fmt.Printf("error opening db %s : %s\n", c.Db, err)
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

		dberr := db.Put(record[0], record[1])
		if dberr != nil {
			fmt.Printf("dbm insert error: %s\n", dberr)
		}
	}

	return 0

}

func (c *DbmImportCommand) Help() string {
	return fmt.Sprintf("csv2dbm : import csv data to a dbm database file\n\n\t\t--database <filename>\tthe dbm filename to use\n\t\t--csv <filename>\tthe csv filename to import data from\n\n")
}

func (c *DbmImportCommand) Synopsis() string {
	return "import csv to a dbm file"
}
