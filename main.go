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

// Define flags for command-line options
var (
	help = flag.Bool("help", false, "Show help screen")        // Flag for displaying help message
	port = flag.String("port", "3000", "Port number")          // Flag for specifying the port number
	dir  = flag.String("dir", "data", "Path to the directory") // Flag for specifying the directory path
)

func main() {
	// Parse command-line flags
	flag.Parse()

	// Check if help flag is set; if true, display help message and exit
	if *help {
		printHelp()
	}

	// Initialize storage by creating necessary directories and files in the specified path
	err := utils.CreateStorage(*dir)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new HTTP multiplexer for routing requests
	mux := http.NewServeMux()

	// Define routes
	mux.HandleFunc("/", handlers.BucketHandler(*dir))
	mux.HandleFunc("/{BucketName}/{ObjectName}", handlers.ObjectHnadler(*dir))

	// Start the HTTP server on the specified port
	log.Printf("starting server on :%v\n", *port)
	err = http.ListenAndServe(":"+*port, mux)
	log.Fatal(err)
}

// printHelp displays the help message for using the command-line options
func printHelp() {
	message := `Simple Storage Service.

**Usage:**
	triple-s [-port <N>] [-dir <S>]  
	triple-s --help
			
**Options:**
	- --help     Show this screen.
	- --port N   Port number
	- --dir S    Path to the directory`
	fmt.Println(message)
	os.Exit(0) 
}
