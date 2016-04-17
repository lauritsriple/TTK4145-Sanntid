package fsmelev

import (
	"driver"
	"log"
	"time"
)

const motorStopTimeoutBase=5 //seconds
var timeGoToFloor= time.Now()
var motorStopTimerRunning=false

func Init(orderedFloorCh <-chan uint,lightCh chan driver.Light, statusCh chan driver.LiftStatus, buttonCh chan<- driver.Button,motorStopCh chan<- bool ,quitCh <-chan bool ) bool{
	if !driver.Init(){
		log.Fatal("could not initialize driver")
		return false
	}
	// clear all lights and stop motor
	driver.ClearAll()
	log.Println("cleared all lights and stopped motor.")
	var floorSensorCh = make(chan uint,5)
	var doorTimerCh = make(chan bool,2)
	var motorDirectionCh = make(chan driver.MotorDirection, 5)
	go driverLoop(lightCh, buttonCh, floorSensorCh, motorDirectionCh,quitCh)
	go executeOrder(orderedFloorCh, lightCh, statusCh, floorSensorCh, doorTimerCh, motorDirectionCh,motorStopCh, quitCh)
	log.Println("fsmElev initialized")
	return true
}

func driverLoop(lightCh <-chan driver.Light, buttonCh chan<- driver.Button, floorSensorCh chan<- uint, motorDirectionCh <-chan driver.MotorDirection ,quitCh <-chan bool){
	for{
		select{
		case <-quitCh:
			driver.ClearAll()
			return
		default:
			driver.ReadButtons(buttonCh)
			driver.ReadFloorSensors(floorSensorCh)
			driver.RunMotor(motorDirectionCh)
			driver.SetLight(lightCh)
			time.Sleep(5 * time.Millisecond)
		}
	}
}

func executeOrder(orderedFloorCh <-chan uint, lightCh chan<- driver.Light, statusCh chan<- driver.LiftStatus, floorSensorCh <-chan uint, doorTimerCh chan bool, motorDirectionCh chan<- driver.MotorDirection,motorStopCh chan<- bool ,quitCh <-chan bool){
	var (
		currentFloor uint
		stopFloor uint
		status driver.LiftStatus
		)
	status.Direction = true

	// not in state, go up until floor
	motorDirectionCh <-driver.MD_up
	for{
		currentFloor = <-floorSensorCh
		if currentFloor != 0{
			log.Println("found floor")
			break
		}
	}
	motorDirectionCh <-driver.MD_stop
	status.Floor=currentFloor
	log.Println(currentFloor)
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
			lightCh<-driver.Light{0,driver.Door, false}
			status.Door = false
			statusCh <-status
		case currentFloor = <-floorSensorCh:
			updateStatus(currentFloor, &status,motorDirectionCh,statusCh)
		default:
			time.Sleep(5*time.Millisecond)
			if stopFloor != 0{
				stopAtFloor(currentFloor, &status, &stopFloor, motorDirectionCh,lightCh,statusCh,doorTimerCh,motorStopCh)
				goToFloor(currentFloor, &status, &stopFloor,motorDirectionCh,statusCh)// goes to floor if not door is open
			}
		}
	}
}

func stopAtFloor(currentFloor uint, status *driver.LiftStatus, stopFloor *uint, motorDirectionCh chan<- driver.MotorDirection, lightCh chan<- driver.Light,statusCh chan<- driver.LiftStatus,doorTimerCh chan<- bool,motorStopCh chan<- bool){
	if status.Floor == *stopFloor{
		motorDirectionCh <-driver.MD_stop
		status.Running = false
		status.Door = true
		lightCh <-driver.Light{0, driver.Door, true}
		go func(){
			time.Sleep(3* time.Second)
			doorTimerCh<- true
		}()
		*stopFloor = 0
		motorStopCh<-false
		statusCh <-*status
	} else if time.Now().Sub(timeGoToFloor) > time.Duration(motorStopTimeoutBase)*time.Second && motorStopTimerRunning{
		//motorstop detected
		motorStopTimerRunning=false
		motorStopCh<-true
	}
}


func goToFloor(currentFloor uint, status *driver.LiftStatus, stopFloor *uint, motorDirectionCh chan<- driver.MotorDirection,statusCh chan<- driver.LiftStatus){
	if !status.Door && !status.Running{
		if currentFloor < *stopFloor{
			motorDirectionCh <- driver.MD_up
			status.Direction = true
		} else {
			motorDirectionCh <- driver.MD_down
			status.Direction = false
		}
		motorStopTimerRunning=true
		timeGoToFloor=time.Now()
		status.Running = true
		statusCh <- *status
	}
}

func updateStatus(currentFloor uint, status *driver.LiftStatus, motorDirectionCh chan<- driver.MotorDirection,statusCh chan<- driver.LiftStatus){
	switch currentFloor{
		case 0:
			if status.Door{
				log.Fatal("lift should not be moving, door is open")
			}
			if !status.Running{
				log.Fatal("lift should not be moving, motor is off")
			}
			return
		case 1,4:
			motorDirectionCh <-driver.MD_stop
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









