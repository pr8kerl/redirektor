package main

import (
	"bufio"
	"bytes"
	"fmt"
	"gopkg.in/redis.v4"
	"os"
	"path/filepath"
)

var (
	basename string
	cfgfile  string
)

func init() {
	basename = filepath.Base(os.Args[0])
}

func main() {

	cfgfile = os.Getenv("REDIREKTOR_CFG")
	// config defaults
	// use basename as default key prefix
	var cfg = &ConfigSection{
		DBAddr:     "localhost:6379",
		DBID:       0,
		DBPassword: "",
		KeyPrefix:  basename,
	}
	if cfgfile != "" {
		// load custom settings from configfile
		var err error
		cfg, err = initialiseConfig(cfgfile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error : failed to initialise config : %s", err)
			os.Exit(1)
		}
	}

	var prefix bytes.Buffer
	prefix.WriteString(cfg.KeyPrefix)
	prefix.WriteString(":")

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.DBAddr,
		Password: cfg.DBPassword,
		DB:       cfg.DBID,
		//		ReadOnly: true,
	})

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		inkey := scanner.Text()
		key := prefix
		key.WriteString(inkey)

		// lookup key
		val, err := client.Get(key.String()).Result()
		if err == redis.Nil {
			// key does not exist
			os.Stdout.Write([]byte("NULL\n"))
		} else if err != nil {
			// key lookup error
			os.Stdout.Write([]byte("NULL\n"))
			fmt.Fprintf(os.Stderr, "redirektor error looking up key %s : %s\n", inkey, err)
		} else {
			// key found
			os.Stdout.WriteString(val)
			os.Stdout.WriteString("\n")
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "redirektor scanner error: %s\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "redirektor shutdown\n")

}
