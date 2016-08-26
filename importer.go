package main

import (
	"encoding/csv"
	"fmt"
	"github.com/boltdb/bolt"
	"io"
	"log"
	"os"
)

// documentation for csv is at http://golang.org/pkg/encoding/csv/

func main() {

	rdb := os.Args[1]
	bucket := os.Args[2]
	fn := os.Args[3]

	file, err := os.Open(fn)
	if err != nil {
		fmt.Printf("error opening file %s : %s\n", fn, err)
		return
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	//reader.Comma = '\t'
	db, err := bolt.Open(rdb, 0600, nil)
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
			return
		}
		fmt.Printf("inurl: %s\touturl: %s\n", record[0], record[1])

		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return err
			}

			b.Put([]byte(record[0]), []byte(record[1]))

			return nil
		})
	}

}
