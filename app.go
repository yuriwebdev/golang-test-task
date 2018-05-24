package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var BIND_ADDR = flag.String("port", ":8008", "bind addr")

func init() {
	fmt.Println("Application initialized")
	flag.Parse()

	f, err := os.OpenFile("./log.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

	if err != nil {
		panic(err)
	}

	log.SetOutput(f)

}
