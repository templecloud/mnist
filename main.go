package main

import (
	"log"
	"os"
)

import (
	"github.com/templecloud/mnist/idxfile"
)

func main() {
	pwd, _ := os.Getwd()
	log.Printf("Current: %s\n", pwd)

	err := idxfile.DownloadMNIST()
	if err != nil {
		log.Println(err)
	}
}
