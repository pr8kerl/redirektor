package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/jsimonetti/berkeleydb"
	"github.com/mitchellh/cli"
	"os"
)

type DbmExportCommand struct {
	Db  string
	Csv string
	Ui  cli.Ui
}

func dbmExportCmdFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	return &DbmExportCommand{
		Db:  "",
		Csv: "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorGreen,
		},
	}, nil
}

func (c *DbmExportCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("export", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&c.Db, "database", "", "the name of the dbm file to export")
	cmdFlags.StringVar(&c.Csv, "csv", "", "the name of the csv to output to")
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

	csvfile, err := os.Create(c.Csv)
	if err != nil {
		fmt.Printf("error opening csv file %s : %s\n", c.Csv, err)
		return 1
	}
	// automatically call Close() at the end of current method
	defer csvfile.Close()
	//
	w := csv.NewWriter(csvfile)
	defer w.Flush()

	db, err := berkeleydb.NewDB()
	if err != nil {
		fmt.Printf("error opening db connection : %s\n", err)
		return 1
	}

	//err = db.Open(c.Db, berkeleydb.DbBtree, berkeleydb.DbCreate)
	err = db.Open(c.Db, berkeleydb.DbHash, berkeleydb.DbRdOnly)
	if err != nil {
		fmt.Printf("error opening db %s : %s\n", c.Db, err)
		return 1
	}
	defer db.Close()

	val, err := db.Get("/Essentials")
	if err != nil {
		fmt.Printf("Unexpected error in Get: %s\n", err)
		return 1
	}
	fmt.Printf("value: %s\n", val)

	curs, err := db.Cursor()
	if err != nil {
		fmt.Printf("error obtaining db cursor: %s\n", err)
		return 1
	}

	var key, value string
	var dberr error
	key, value, dberr = curs.GetFirst()
	if dberr != nil {
		fmt.Printf("dbm read error: %s\n", dberr)
	}

	for dberr != nil {
		key, value, dberr = curs.GetNext()
		if dberr != nil {
			fmt.Printf("dbm read error: %s\n", dberr)
			break
		}

		record := []string{
			key,
			value,
		}

		if err := w.Write(record); err != nil {
			fmt.Printf("error writing record to csv: %s\n", err)
		}
		fmt.Printf("%s,%s\n", key, value)

	}

	w.Flush()
	if err := w.Error(); err != nil {
		fmt.Printf("error writing to csv: %s\n", err)
		return 1
	}

	return 0

}

func (c *DbmExportCommand) Help() string {
	return fmt.Sprintf("csv2dbm : export csv data from a dbm database file\n\n\t\t--database <filename>\tthe dbm filename to use\n\t\t--csv <filename>\tthe csv filename to export data to\n\n")
}

func (c *DbmExportCommand) Synopsis() string {
	return "export csv to a dbm file"
}
