package udp_test

//How to test
//run: go test -v udp_test.go

import (
	"fmt"
	"testing"
	"time"
	"udp"
)

func TestUdpModule(t *testing.T) {
	var id int
	toNetwork := make(chan udp.Message, 10)
	fromNetwork := make(chan udp.Message, 10)
	quit := make(chan bool)
	id = udp.NetInit(toNetwork, fromNetwork, quit)
	fmt.Println("id is: ", id)
	count := 0
	for {
		message := udp.Message{
			LiftId:    111,
			Floor:     111,
			Direction: true,
			Status:    udp.New,
			Weight:    5,
			TimeRecv:  time.Now()}
		toNetwork <- message

		time.Sleep(time.Second * 1)
		select {
		case recv := <-fromNetwork:
			fmt.Println("received message", recv.LiftId)
		default:
			fmt.Println("no message")
		}
		count++
		if count > 5 {
			break
		}
	}
}
