package main

import (
	"log"

	"github.com/benmanns/goworker"
)

func init() {
	goworker.Register("Archive", archiveWorker)
}

func archiveWorker(queue string, args ...interface{}) error {
	log.Printf("Archive on %s, with params %s", queue, args)
	return nil
}

func main() {
	if err := goworker.Work(); err != nil {
		log.Println(err)
	}
}
