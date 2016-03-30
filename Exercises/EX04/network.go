package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
)

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}	
}

func tcp_receive() {
	laddr, err := net.ResolveTCPAddr("tcp",":30000")
	CheckError(err)

	ln, err := net.ListenTCP("tcp",laddr)
	CheckError(err)

	conn,err:=ln.Accept()
	CheckError(err)
	fmt.Println("here")

	buf:= make([]byte,1024)
	
	read_len,err := conn.Read(buf[:])
	if read_len > 0{
		fmt.Println("received tcp:",string(buf[0:read_len]))
		conn.Write([]byte("backloop"))
	}
}

func tcp_talk(){
	rAddr, err := net.ResolveTCPAddr("tcp","129.241.187.23:33546")
	CheckError(err)
	conn,err := net.DialTCP("tcp",nil,rAddr)
	CheckError(err)
	
	defer conn.Close()

	conn.Write([]byte("HelloWorld!\000"))

	laddr, err := net.ResolveTCPAddr("tcp4",":30000")
	CheckError(err)
	listen,err := net.ListenTCP("tcp",laddr)
	CheckError(err)

	listen.SetDeadline(time.Now().Add(2*time.Second))
	defer listen.Close()

	//conn.Write([]byte("Connect to: 129.241.187.144:30000\000"))
	//conn.Write([]byte("Connect to: 129.241.187.144:30000\000"))

	deadline := time.Now().Add(2*time.Second)
	for time.Now().Before(deadline) {
		conn2, err := listen.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Socket acceptance error: " + err.Error())
			continue
		}
		defer conn2.Close()
	}
}

func tcp_send(){
	rAddr, err := net.ResolveTCPAddr("tcp","129.241.187.23:33546")
	CheckError(err)
	
	conn, err := net.DialTCP("tcp",nil, rAddr)
	CheckError(err)
	
	var test = "hello world\x00"
	_,err=conn.Write([]byte(test))
}

func udp_receive() {
	serverAddr,err := net.ResolveUDPAddr("udp",":20012")
	CheckError(err)
	
	serverConn, err := net.ListenUDP("udp",serverAddr)
	CheckError(err)
	
	defer serverConn.Close()

	buf:= make([]byte,1024)
	
	for{
		n,addr,err := serverConn.ReadFromUDP(buf)
		CheckError(err)
		fmt.Println("received ",string(buf[0:n])," from ",addr)
	} 
}

func udp_send(){
	ServerAddr,err := net.ResolveUDPAddr("udp","129.241.187.23:20012")
	CheckError(err)
 
	LocalAddr, err := net.ResolveUDPAddr("udp", "129.241.187.144:0")
	CheckError(err)
 
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
 
	defer Conn.Close()
	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_,err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func udp_broadcast(){
	ServerAddr,err := net.ResolveUDPAddr("udp","129.241.187.255:20012")
	CheckError(err)
 
	Conn, err := net.DialUDP("udp", nil, ServerAddr)
	CheckError(err)
 
	defer Conn.Close()
	i := 0
	for {
		msg := strconv.Itoa(i)
		i++
		buf := []byte(msg)
		_,err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func main(){
	//go sendUDPBroadcast()
	//go receiveUDP()
	//go tcp_receive()
	//go tcp_send()
	tcp_talk()
	for {

	}
}
 
	
