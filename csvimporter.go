package main

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/mitchellh/cli"
	"gopkg.in/redis.v4"
	"io"
	"log"
	"os"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/

func main() {

	c := cli.NewCLI("csvimporter", "0.0.1")
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"bolt":  boltCmdFactory,
		"redis": redisCmdFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitStatus)

}

type BoltCommand struct {
	Db     string
	Bucket string
	Csv    string
	Ui     cli.Ui
}

func boltCmdFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	return &BoltCommand{
		Db:     "",
		Bucket: "",
		Csv:    "",
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorGreen,
		},
	}, nil
}

func (c *BoltCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("bolt", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Db, "database", "", "the name of the Bolt DB file to work against")
	cmdFlags.StringVar(&c.Bucket, "bucket", "", "the name of the Bolt DB bucket name to work against")
	cmdFlags.StringVar(&c.Csv, "csv", "", "the name of the csv to import")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	var requiredFlags = 0

	if c.Db != "" {
		requiredFlags++
		fmt.Printf("using db file : %s\n", c.Db)
	}
	if c.Bucket != "" {
		requiredFlags++
		fmt.Printf("using db bucket name : %s\n", c.Bucket)
	}
	if c.Csv != "" {
		requiredFlags++
		fmt.Printf("using csv file name : %s\n", c.Csv)
	}

	if requiredFlags < 3 {
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
	db, err := bolt.Open(c.Db, 0600, nil)
	if err != nil {
		log.Fatal(err)
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

		dberr := db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(c.Bucket))
			if err != nil {
				return err
			}

			b.Put([]byte(record[0]), []byte(record[1]))

			return nil
		})
		if dberr != nil {
			fmt.Printf("bolt db update error: %s\n", err)
		}
	}

	return 0

}

func (c *BoltCommand) Help() string {
	return fmt.Sprintf("csvimporter bolt: import csv data to a BoltDB file\n\n\t\t--database <filename>\tthe BoltDB filename to use\n\t\t--bucket <value>\tthe BoltDB bucket to import data into\n\t\t--csv <filename>\tthe csv filename to import data from\n\n")
}

func (c *BoltCommand) Synopsis() string {
	return "import csv to a BoltDB file"
}

type RedisCommand struct {
	Addr     string
	Prefix   string
	Password string
	Csv      string
	DbId     int
	Ui       cli.Ui
}

func redisCmdFactory() (cli.Command, error) {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	return &RedisCommand{
		Addr:     "",
		Prefix:   "",
		Password: "",
		Csv:      "",
		DbId:     0,
		Ui: &cli.ColoredUi{
			Ui:          ui,
			OutputColor: cli.UiColorGreen,
		},
	}, nil
}

func (c *RedisCommand) Run(args []string) int {

	cmdFlags := flag.NewFlagSet("redis", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }

	cmdFlags.StringVar(&c.Addr, "database", "localhost:6379", "the redis DB connection string")
	cmdFlags.StringVar(&c.Prefix, "prefix", "", "the prefix of the redis keys to work against")
	cmdFlags.StringVar(&c.Csv, "csv", "", "the name of the csv to import")
	cmdFlags.IntVar(&c.DbId, "db", 0, "the redis database number to use - default is 0")
	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}
	var requiredFlags = 0

	if c.Prefix != "" {
		requiredFlags++
		fmt.Printf("using db key prefix : %s\n", c.Prefix)
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

	client := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Password, // no password set
		DB:       c.DbId,     // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	if err != nil {
		fmt.Printf("error connecting to redis db %s : %s\n", c.Addr, err)
		return 1
	}

	var prefix bytes.Buffer
	prefix.WriteString(c.Prefix)
	prefix.WriteString(":")
	fmt.Printf("redis prefix set to : %s\n", prefix.String())

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
		pbuf := prefix
		pbuf.WriteString(record[0])
		fmt.Printf("inurl: %s\touturl: %s\n", pbuf.String(), record[1])

		// Last argument is expiration. Zero means the key has no
		// expiration time.
		dberr := client.Set(pbuf.String(), record[1], 0).Err()
		if dberr != nil {
			fmt.Printf("error setting key: %s, %s\n", pbuf.String(), dberr)
		}

	}

	return 0
}

func (c *RedisCommand) Help() string {
	return fmt.Sprintf("csvimporter redis: import csv data to a Redis DB\n\n\t\t--database <string>\tthe redis DB connection string - default is localhost:6379\n\t\t--prefix <value>\tthe key prefix to use for import data\n\t\t--csv <filename>\tthe csv filename to import data from\n\t\t--db <int>\tthe redis DB number to use - defaults to 0\n\n")
}

func (c *RedisCommand) Synopsis() string {
	return "import csv to a RedisDB"
}
