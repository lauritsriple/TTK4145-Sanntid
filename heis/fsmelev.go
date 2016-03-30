package fsmelev

import (
	"driver"
	"log"
	"fmt"
	// messagePasser
)

type elevFSM struct{
	stopIsPressed bool
	obstackle bool
	floor int
	direction int
	destinatedFloor int
}

func (fsmData elevFSM) LoopIO(){
	floor := make(chan int,1)
	go driver.Driver_floorSensorPoller(floor)
	select{
	case fsmData.floor <- floor:
		pass
	default:
		pass
	}

	stopButtonPressed := make(chan bool,1)
	go driver.Driver_btnStopPoller(stopButtonPressed)
	select{
	case fsmData.stopIsPressed <- stopButtonPressed:
		pass
	default:
		pass
	}

	obstackleChan := make(chan bool,1)
	go driver.DriverObstructionPoller(obstacleChan)
	select{
	case fsmData.obstacle <- obstacleChan:
		pass
	default:
		pass
	}

	// orderedFloor := make(chan int)
	// go messagePasser.getfloor(orderedFloor)
	// select{
	// case fsmData.destinatedFloor <- orderedFloor:
	// 	pass
	// default:
	// 	pass
	//}

	// btnLightChan := make(chan int)
	// go messagePasser.getlight(btnLightChan)
	// select{
	// case driver.Driver_setBtnLight(<-btnLightChan):
	// 	pass
	// default:
	// 	pass
	// }

	//btnChannel := make(chan int)
	//go driver.Driver_btnPoller(btnChannel)
	//select{
	//case messagePasser.sendTCP(<- btnChannel)
	// 	pass
	// default:
	// 	pass
	// }
}	

func Initialize(){
	driver.DriverInit()
	// go up to next floor if not in any floor
	// send message to master - is initialized

}

func (fsmData elevFSM)FSM(){
	switch destinatedFloor{
		case 0:
			driver.Driver_setMotorDir(driver.MD_stop)
		case 1:
			if fsmData.floor != 1{
				driver.Driver_setMotorDir(driver.MD_down)
			} else {
				// send message to master that the elevator has arrived at foor 1
				driver.Driver_setMotorDir(driver.MD_stop)


