package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"triple-s/handlers"
	"triple-s/internal/utils"
)

var (
	help = flag.Bool("help", false, "Show help screen")
	port = flag.String("port", "3000", "Port number")
	dir  = flag.String("dir", "data", "Path to the directory")
)

func main() {
	flag.Parse()

	if *help {
		message := `Simple Storage Service.
**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help
		
**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory`
		fmt.Println(message)
		return
	}

	createStorage()
	// Creatin new server
	mux := http.NewServeMux()

	mux.HandleFunc("/*", handlers.BucketHandler(*dir))
	mux.HandleFunc("//", handlers.ObjectHnadler)

	log.Printf("starting server on :%v\n", *port)
	err := http.ListenAndServe(":"+*port, mux)
	log.Fatal(err)
}

func createStorage() {
	err := os.MkdirAll(*dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	if !isExist(*dir + "/buckets.csv") {
		utils.CreateCSV(*dir, "buckets", []string{"Name", "CreationTime", "LastModifiedTime", "Status"})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
