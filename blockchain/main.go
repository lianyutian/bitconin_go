package main

import (
	"github.com/boltdb/bolt"
	"log"
)

func main() {
	bc := NewBlockChain()
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			log.Panic(err)
		}
	}(bc.db)

	cli := CLI{bc}
	cli.Run()
}
