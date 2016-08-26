package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

var (
	bname []byte = []byte("unknown")
)

func init() {
	basename := filepath.Base(os.Args[0])
	bname = []byte(basename)
}

func main() {

	fmt.Fprintf(os.Stderr, "redirektor bucket name: %s\n", bname)

	db, err := bolt.Open("redirektor.db", 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		fmt.Fprintf(os.Stderr, "redirektor error: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		key := scanner.Text()
		// lookup key
		err = db.View(func(tx *bolt.Tx) error {

			bucket := tx.Bucket(bname)

			if bucket == nil {
				return fmt.Errorf("bucket %q not found!", bname)
			}

			val := bucket.Get([]byte(key))
			if val != nil {
				val = append(val, "\n"...)
				os.Stdout.Write(val)
				return nil
			}
			// apache rewritemap expects NULL if no match is found
			os.Stdout.Write([]byte("NULL"))

			return nil

		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "redirektor read %s error: %s\n", bname, err)
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "redirektor scanner error: %s\n", err)
		os.Exit(1)
	}

}
