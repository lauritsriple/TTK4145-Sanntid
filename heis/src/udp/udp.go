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

const maddr string = "239.0.0.49:2000"

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

//Called by NetInit
//Sets up the network using multicast
func multicastInit(send <-chan Message, recv chan<- Message, iface *net.Interface, quitCh <-chan bool) {
	group, err := net.ResolveUDPAddr("udp", maddr)
	checkError(err)
	conn, err := net.ListenMulticastUDP("udp", iface, group)
	checkError(err)
	defer conn.Close()
	go multicastListen(recv, conn)
	go multicastSend(send, conn, group, quitCh)
	fmt.Println("Network running")
	<-quitCh //w8 for channel to be true. Defer will be called to close the connection
}
//Called by MulticastInit
//Channel can be used to send messages to the network
func multicastSend(send <-chan Message, conn *net.UDPConn, addr *net.UDPAddr, quitCh <-chan bool) {
	for {
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

//Called by MulticastInit
//All received messsages piped to buffered channel
func multicastListen(recv chan<- Message, conn *net.UDPConn) {
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

//Called by RunLift
//Sets up the network and returns the id of the elevator
func NetInit(send <-chan Message, recv chan<- Message, quitch <-chan bool) int {
	addr, iface, err := findIP()
	if err != nil {
		fmt.Println("Error findig the interface", err)
		return 0
	}
	go multicastInit(send, recv, iface, quitch)
	return findID(addr)
}
