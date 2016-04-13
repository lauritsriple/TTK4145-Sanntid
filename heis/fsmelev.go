package fsmelev

import (
	"driver"
	"log"
	"fmt"
	"io"
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
	FloorChan := make(chan int,1)
	StopButtonChan := make(chan bool,1)
	ObstackleChan := make(chan bool,1)
	go driver.BtnStopPoller(stopButtonPressed)
	go driver.FloorSensorPoller(floor)
	go driver.ObstructionPoller(obstacleChan)
	select{
	case fsmData.floor <- FloorChan:
		pass
	case fsmData.stopIsPressed <- FtopButtonChan:
		pass
	case fsmData.obstacle <- ObstacleChan:
	}

	// orderedFloor := make(chan int)
	// go messagePasser.getlight(btnLightChan)
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

func Init(orderedFloorsCh <- chan uint,lightCh <- chan driver.Light, statusCh * chan driver.LiftStatus, buttonCh chan<- driver.Button, quitCh <- chan bool ) bool{
	if !driver.init(){
		log.Fatal("could not initialize driver")
		return false
	}
	// clear all lights and stop motor
	driver.ClearAll()
	log.Println("cleared all lights and stopped motor.")
	floorSensorCh = make(chan uint,5)
	doorTimerCh = make(chan bool,2)
	motorDirectionCh = make(chan driver.MotorDirection, 5)
	go driverLoop(lightCh, buttonCh, floorSensorCh, motorDirectionCh,quitCh)
	go executeOrder(orderedFloorCh, lightCh, statusCh, floorSensorCh, doorTimerCh, motorDirectionCh, quitCh)
	return true
}

func driverLoop(lightCh <- chan driver.Light, buttonCh <- chan driver.Button, floorSensorCh <- chan uint, motorDirectionCh <-chan driver.MotorDirection ,quitCh <- chan bool){
	for{
		select{
		case <-quitCh:
			driver.ClearAll()
			return false
		default:
			driver.ReadButtons(buttonCh)
			driver.ReadFloorSensor(floorSensorCh)
			driver.RunMotor(motorDirectionCh)
			driver.SetLight(lightCh)
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func executeOrder(orderedFloorCh <- chan uint, lightCh chan <- driver.Light, statusCh chan <- driver.LiftStatus, floorSensorCh <-chan uint, doorTimerCh chan bool, motorDirectionCh chan <- driver.MotorDirection, quitCh <-chan bool){
	var (
		currentFloor uint
		stopFloor uint
		status driver.LiftStatus
		)
	status.Direction = false

	// not in state, go up until floor
	motorDirectionCh <-driver.MotorDirection{driver.MD_up}
	for{
		currentFloor = <-floorSensorCh
		if currentFloor != 0{
			log.Println("found floor")
			break
		}
	}
	motorDirectionCh <-driver.MotorDirection{driver.MD_stop}
	status.Floor = currentFloor
	status.Running = false
	status.Door = false
	statusCh <-status
	for{
		select{
		case <-quitCh:
			return
		case stopFloor = <-orderedFloorCh: // fix orderes outside range?
			// got new order
		case <-doorTimerCh:
			lightCh<-driver.Ligth{0,driver.door, false}
			status.Door = false
			statusCh <-status
		case currentFloor = <-floorSensorCh:
			updateStatus(currentFloor, &status)
		default:
			time.Sleep(5*time.Millisecond)
			if stopFloor != 0{
				stopAtFloor(currentFloor, motorDirectionCh, &status, &stopFloor)
				goToFloor(currentFloor, &status, &stopFloor)
			}
		}
	}
}

func stopAtFloor(currentFloor uint, status *driver.LiftStatus, stopFloor *uint){//check if input is correct
	if status.Floor == stopFloor{
		motorDirectionCh <-driver.MotorDirection{driver.MD_stop}
		status.Running = false
		status.Door = true
		lightCh <-driver.Light{0, driver.door, true}
		go func(){
			time.Sleep(3* time.Second)
			doorTimerCh<- true
		}
		*stopFloor = 0
		statusCh <-*status
	}
}


func goToFloor(currentFloor uint, status *driver.LiftStatus, stopFloor *uint){// check if input is correct
	if !status.Door && !status.Running{
		if currenFloor < stopFloor{
			motorDirectionCh <- driver.MotorDirection{driver.MD_up}
			status.Direction = driver.MD_up
		} else {
			motorDirectionCh <- driver.MotorDirection{driver.MD_down}
			status.Direction = driver.MD_down
		}
		status.Running = true
		statusCh <- *status
	}
}

func updateStatus(currentFloor uint, status *driver.LiftStatus){
	switch currentFloor{
		case 0:
			if status.Door{
				log.Fatal("lift should not be mooving, door is open")
			}
			if !status.Running{
				log.Fatal("lift should not be mooving, motor is off")
			}
		case 1,4:
			motorDirectionCh <-driver.MotorDirection{driver.MD_stop}
			status.Floor = currentFloor
			status.Running = false
			statusCh <-*status
		case 2,3:
			if currentFloor != status.Floor{
				status.Floor = currentFloor
				statusCh <- *status
			}
		default:
			log.Println("found unknown floor", currentFloor)
	}
}









