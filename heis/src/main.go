package main

import (
	"log"
	"control"
)

var quitCh = make(chan bool)

func main() {
	log.Println("Starting elevator. Send SIGQUIT to shutown (CTRL+\\)")
	control.RunLift(quitCh)
}
