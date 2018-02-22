package main

import (
	"io/ioutil"
	"log"
	"os"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func setUpLogger(debug bool) {
	if debug == true {
		Debug = log.New(os.Stdout,
			"DEBUG: ",
			log.Ldate|log.Ltime|log.Lshortfile)
		Debug.Println("Debug mode")
	} else {
		Debug = log.New(ioutil.Discard,
			"DEBUG: ",
			log.Ldate|log.Ltime|log.Lshortfile)
	}

	Info = log.New(os.Stdout,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(os.Stdout,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(os.Stderr,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
