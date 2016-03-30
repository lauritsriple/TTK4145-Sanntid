package udp_test

//How to test
//run: go test -v udp_test.go

import ("fmt"
	"udp"
	"time"
	"strconv"
	"testing"
)

func TestUdpModule(t *testing.T){
	udpListenChan:=make(chan string)
	udpBroadcastChan:=make(chan string)
	go udp.Listen(udpListenChan)
	go udp.Broadcast(udpBroadcastChan)
	count:=0
	for {
		
		msg:=strconv.Itoa(count)
		fmt.Println("sending message",msg)
		udpBroadcastChan<-msg
		time.Sleep(time.Second * 1)
		select{
		case recv:=<-udpListenChan:
			fmt.Println("received message",recv)
		default:
			fmt.Println("no message")
		}
		count++
		if count>5{
			break
		}
	}
}
