package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"os/exec"
)

func CheckError(err error) {
	if err  != nil {
		fmt.Println("Error: " , err)
		os.Exit(0)
	}
	
}


func udp_receive_blocking() bool {
	serverAddr,err := net.ResolveUDPAddr("udp",":20012")
	CheckError(err)
	

	serverConn, err := net.ListenUDP("udp",serverAddr)
	CheckError(err)
	if err!=nil{
		serverConn.Close()
		return true
	}
	
	defer serverConn.Close()

	buf:= make([]byte,1024)
	serverConn.SetReadDeadline(time.Now().Add(5 * time.Second))
	n,addr,err := serverConn.ReadFromUDP(buf)
	if err!=nil{
		serverConn.Close()
		return true
	}
	fmt.Println("received ",string(buf[0:n])," from ",addr)
	CheckError(err)
	serverConn.Close()
	return false
}

func udp_receive_nonBlocking(recv chan<-int, alive chan<-bool){
	serverAddr,err := net.ResolveUDPAddr("udp",":20012")
	CheckError(err)
	
	serverConn, err := net.ListenUDP("udp",serverAddr)
	CheckError(err)
	
	serverConn.Close()

	buf:= make([]byte,1024)
	
	for{
		serverConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n,addr,err := serverConn.ReadFromUDP(buf)
		if (err != nil){
			alive<-false
			serverConn.Close()
		}
		
		i,err :=  strconv.Atoi(string(buf[0:n]))
		fmt.Println("received ",i," from ",addr)
		alive<-true
		recv<-i
		
	} 

}

func udp_send(count int){
	ServerAddr,err := net.ResolveUDPAddr("udp","127.0.0.1:20012")
	CheckError(err)
 
	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)
 
	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)
 
	defer Conn.Close()
	for {
		msg := strconv.Itoa(count)
		count++
		buf := []byte(msg)
		_,err := Conn.Write(buf)
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
	}
}

func isMaster()bool{
	master := udp_receive_blocking()
	if master{
		return true
	} else{
		return false
	}


}

func initSlave()int{
	master := false
	count := 0
	ch:= make(chan int)
	alive_ch := make(chan bool)
	go udp_receive_nonBlocking(ch,alive_ch)
	for master == false{
		if !<-alive_ch{
			return count
		}
		count = <- ch
	}
	return 0
}


func main(){
	count:=0
	master:=isMaster()
	for{
		if master{
			fmt.Println("Master spawned")
			cmd:= exec.Command("gnome-terminal","-x","bash","-c", "go run /home/student/TTK4145-Sanntid/Exercises/EX06/process_pair.go")
			err:= cmd.Start()
			fmt.Println(err)
			udp_send(count)
		} else{
			fmt.Println("Slave spawned")
			count=initSlave()
			//Master died
			master=true
		}
	}
}

 
	
