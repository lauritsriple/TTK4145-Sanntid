package udp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

//How to make libary
//Put libary in src/udp
//go build udp
//go install udp

const PORT string = "20012"
const BROADCASTADDR string = "192.168.43.255"
const braddr string = "239.0.0.49:2000"

type Orderstatus int

const (
	New Orderstatus = iota
	Accepted
	Done
	Reassign
)

type Message struct {
	LiftId    int
	ReassId   int
	Floor     uint
	Direction bool
	Status    Orderstatus
	Weight    int
	TimeRecv  time.Time
}

//Checks the error message and kills the program
func checkError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(0)
	}
}

//Called by NetInit
//Returns IPv4 address for lift
func findIP() (string, *net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", nil, err
	}
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if strings.Contains(a.String(), "129.") {
				return a.String(), &iface, nil
			}
		}
	}
	return "", nil, errors.New("Unable to find IPv4 adress")
}

//Called by NetInit
//Returns 3 last digits from IPv4 adress
func findID(a string) int {
	fmt.Println(a)
	id, err := strconv.Atoi(strings.Split(a, ".")[3][:3])
	if err != nil {
		fmt.Println("Error converting IP to ID", err)
	}
	return id
}

//Sets up the broadcast
//Called by NetInit
func BroadcastInit(send <-chan Message, recv chan<- Message, iface *net.Interface, quitCh <-chan bool) {
	group, err := net.ResolveUDPAddr("udp", braddr)
	checkError(err)
	conn, err := net.ListenMulticastUDP("udp", iface, group)
	checkError(err)
	defer conn.Close()
	go broadcastListen(recv, conn)
	go broadcastSend(send, conn, group, quitCh)
	fmt.Println("Network running")
	<-quitCh //w8 for channel to be true. Defer will be called
}

//Workerthread. Called by BroadcastInit
func broadcastSend(send <-chan Message, conn *net.UDPConn, addr *net.UDPAddr, quitCh <-chan bool) {
	for {
		fmt.Println("trying to send")
		select {
		case m := <-send:
			buf, err := json.Marshal(m)
			if err != nil {
				fmt.Println("JSON encoding: ", err)
			} else {
				_, err := conn.WriteToUDP(buf, addr)
				if err != nil {
					fmt.Println("NET: ", err)
				}
			}
		case <-quitCh:
			return
		}
	}
}

//Workerthread. Called by BroadcastInit
func broadcastListen(recv chan<- Message, conn *net.UDPConn) {
	for {
		buf := make([]byte, 512)
		l, _, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Println("NET: ", err)
		}
		var m Message
		err = json.Unmarshal(buf[:l], &m)
		if err != nil {
			fmt.Println("JSON unpacking: ", err)
		} else {
			m.TimeRecv = time.Now()
			recv <- m
		}
	}
}

//Called by
//Sets up the network and returns the id of the elevator
func NetInit(send <-chan Message, recv chan<- Message, quitch <-chan bool) int {
	addr, iface, err := findIP()
	if err != nil {
		fmt.Println("Error findig the interface", err)
		return 0
	}
	go BroadcastInit(send, recv, iface, quitch)
	return findID(addr)
}
