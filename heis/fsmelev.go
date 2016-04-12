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

func Initialize() bool{
	//driver.init()
	DirectionChan := make(chan int,1)
	LightChan := make(chan driver.Ligth,5)// namechange
	LiftStatusChan := make(chan int,1)
	DestinatedFloorChan := make(chan int,1)// messages from bypassed trough driver from queue
	FloorSensorChan := make(chan int,1)
	if <-FloorSensorChan{
		driver.SetMotorDir(driver.MD_up)//
		for (!<-FloorSensorChan){
			// elevating
		}
		driver.SetMotorDir(driver.MD_stop)
	}
	return 1;

}

func (fsmData elevFSM)GoToFloor(){
	for !Initialize(){
	// wait
	}
	switch <-DestinatedFloor{
		case 0:
			driver.SetMotorDir(driver.MD_stop)
			LightChan<- driver.Light{0,driver.Stop,true}
			LiftStatusChan<-driver.LiftStatus{false,driver.currentFloor,false,false}
			//wait - reset queue - continue
		case 1:
			driver.SetFloorIndicator(<-driver.currentFloor)// probably not best way 
			if fsmData.floor != 1{
				driver.SetMotorDir(driver.MD_down)
				LightChan <-driver.Light{0,driver.Down,true}
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?

			} else {
				driver.SetMotorDir(driver.MD_stop)
				LightChan <-driver.Light{1,driver.Command,false}// maybe make this code more readable?
				LiftStatusChan<-driver.LiftStatus{false,driver.currentFloor,false,true}//Direction?
				//LightChan <-driver.Light{0,driver.door,true}
				//wait
				//LightChan <-driver.Light{0,driver.door,false}
				}
		case 2:
			driver.SetFloorIndicator(<-driver.currentFloor))
			if fsmData.floor < 2{
				driver.SetMotorDir(driver.MD_up)
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?
			} else if fsmData.floor != 2 {
				driver.SetMotorDir(driver.MD_down)
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?
			} else {
				driver.SetMotorDir(driver.MD_stop)
				LightChan <-driver.Light{2,driver.Command,false}
				LiftStatusChan<-driver.LiftStatus{false,driver.currentFloor,false,true}//Direction?
				//LightChan <-driver.Light{0,driver.door,true}
				//wait
				//LightChan <-driver.Light{0,driver.door,false}

			}
		case 3:
			driver.SetFloorIndicator(<-driver.currentFloor))
			if fsmData.floor > 3{
				driver.SetMotorDir(driver.MD_down)
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?
			} else if fsmData.floor != 3{
				driver.SetMotorDir(driver.MD_up)
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?
			} else {
				driver.SetMotorDir(driver.MD_stop)
				LightChan <-driver.Light{3,driver.Command,false}
				LiftStatusChan<-driver.LiftStatus{false,driver.currentFloor,false,true}//Direction?
				//LightChan <-driver.Light{0,driver.door,true}
				//wait
				//LightChan <-driver.Light{0,driver.door,false}


			}
		case 4:
			driver.SetFloorIndicator(<-driver.currentFloor))
			if fsmData.floor != 4{
				driver.SetMotorDir(driver.MD_up)
				LiftStatusChan<-driver.LiftStatus{true,driver.currentFloor,false,false}//Direction?
			} else {
				driver.SetMotorDir(driver.MD_stop)
				LightChan <-driver.Light{4,driver.Command,false}
				LiftStatusChan<-driver.LiftStatus{false,driver.currentFloor,false,true}//Direction?
				//LightChan <-driver.Light{0,driver.door,true}
				//wait
				//LightChan <-driver.Light{0,driver.door,false}


			}
		}
	}



