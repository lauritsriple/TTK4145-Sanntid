package udp

import (
	"fmt"
	"net"
	"os"	
)

//How to make libary
//Put libary in src/udp
//go build udp
//go install udp

const PORT string = "20012"
const BROADCASTADDR string = "129.241.187.255"

//Checks the error message and kills the program
func checkError(err error) {
	if err != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}	
}
	
func Listen(msgCh chan<- string ) { //Writing to chan, ch should be buffered
	serverAddr,err := net.ResolveUDPAddr("udp",":"+PORT)
	checkError(err)
	
	serverConn, err := net.ListenUDP("udp",serverAddr)
	checkError(err)
	
	defer serverConn.Close() //Closes connection when program stops

	buf:= make([]byte,1024)
	for {
		n,_,err := serverConn.ReadFromUDP(buf) //blocking
		checkError(err)

		msgCh <- string(buf[0:n]) //Sends the message to the buffer
	}
	 
}

func Broadcast(msgCh <-chan string){ //reading from chan
	ServerAddr,err := net.ResolveUDPAddr("udp",BROADCASTADDR+":"+PORT)
	checkError(err)
 
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	checkError(err)
 
	defer Conn.Close() //Closes connection when program stops
	for {
		buf := []byte(<-msgCh) //Blocking, reads from channel
		_,err := Conn.Write(buf)
		checkError(err)
		//time.Sleep(time.Second * 1)
	}
}


