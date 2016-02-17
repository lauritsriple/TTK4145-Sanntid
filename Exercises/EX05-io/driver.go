package driver

import (
	"time"
	"log"
	."io"
)

const N_FLOORS = 4
const N_LIGHTS = 4
const N_BUTTONS = 3

//Milliseconds between each polling round
const POLLRATE = 20*time.Millisecond

type MotorDirection int
const (
	MD_up = 1
	MD_down = -1
	MD_stop = 0
)
type BtnEvent struct{
	Floor int
	Button int
}

var driverInitialized=false

driver_init(){
	if driverInitialized{
		log.fatal("ERROR, driver already initialized")
	}
	else {
		driverInitialized=true
		if io_init()==0 {
			log.fatal("ERROR, could not initialize driver")
		}
		else {
			//sucess
			return nil
		}
	}
}

//LIGHTS: FLOOR INDICATORS
var floorIndicatorChannels = [N_FLOORS] int {
	ch.LIGHT_FLOOR_IND1,
	ch.LIGHT_FLOOR_IND2,
	ch.LIGHT_FLOOR_IND3,
	ch.LIGHT_FLOOR_IND4
}

func driver_setFloorIndicator(floor int, val bool){
	if val{
		io_setBit(floorIndicatorChannels[floor])
	}
	else{
		io_clearBit(floorIndicatorChannels[floor])
	}
}

//LIGHTS: UP;DOWN;COMMAND
var lightChannels = [N_FLOORS][N_LIGHTS] int {
	{ch.LIGHT_UP1,0,ch.LIGHT_COMMAND1},
	{ch.LIGHT_UP2,ch.LIGHT_DOWN2,ch.LIGHT_COMMAND2},
	{ch.LIGHT_UP3,ch.LIGHT_DOWN3,ch.LIGHT_COMMAND3},
	{0,ch.LIGHT_DOWN4,ch.LIGHT_COMMAND4}
}

func driver_setBtnLight(floor int, btn int, val bool){
	if val{
		io_setBit(lightChannels[floor][btn])
	}
	else{
		io_clearBit(lightChannels[floor][btn])
	}
}

func driver_setStopLight(val bool){
	if val{
		io_setBit(ch.LIGHT_STOP)
	}
	else {
		io_clearBit(ch.LIGHT_STOP)
	}
}

func driver_setDoorLight(val bool){
	if val{
		io_setBit(ch.LIGHT_DOOR_OPEN)
	}
	else {
		io_clearBit(ch.LIGHT_DOOR_OPEN)
	}
}

//BUTTONS: UP;DOWN;COMMAND
var btnChannels = [N_FLOORS][N_BUTTONS] int {
	{ch.BUTTON_UP1,0,ch.BUTTON_COMMAND1},
	{ch.BUTTON_UP2,ch.BUTTON_DOWN2,ch.BUTTON_COMMAND2},
	{ch.BUTTON_UP3,ch.BUTTON_DOWN3,ch.BUTTON_COMMAND3},
	{0,ch.BUTTON_DOWN4,ch.BUTTON_COMMAND4}
}

func driver_btnPoller(recv chan <- BtnEvent){}
	var prev [N_FLOORS][N_BUTTONS] int

	for {
		time.Sleep(POLLRATE)
		for f:=0; f<N_FLOORS; f++{
			for b:=0; b<N_BUTTONS; b++{
				curr:=io_readBit(btnChannels[floor][btn])
				if (curr != 0 && curr != prev[f][b]):
					recv <- BtnEvent{f,b}
				prev[f][b]=curr
			}
		}
	}
}

func driver_btnStopPoller(recv chan <- int){
	var prev int
	for{
		time.Sleep(POLLRATE)
		curr:=io_readBit(ch.STOP)
		if (curr!=0 && curr!=prev):
			recv <- curr
			prev=curr
	}
}

func driver_obstructionPoller(recv chan <- int){
	var prev int
	for{
		time.Sleep(POLLRATE)
		curr:=io_readBit(ch.OBSTRUCTION)
		if (curr!=0 && curr!=prev):
			recv <- curr
			prev=curr
	}
}

//FLOORSENSORS
var floorSensorChannels = [N_FLOORS] int {
	ch.SENSOR_FLOOR1,
	ch.SENSOR_FLOOR2,
	ch.SENSOR_FLOOR3,
	ch.SENSOR_FLOOR4
}

func driver_floorSensorPoller(recv chan <- int){
	var prev int
	for{
		time.Sleep(POLLRATE)
		for f:=0; f<N_FLOORS; f++{
			curr:=io_readBit(floorSensorChannels[f])
			if (curr!=0 && f!=prev){
				receiver <- f
				prev=f
			}
		}
	}
}

func driver_setMotorDir(dir MotorDirection){
	switch dir{
	case MD_stop:
		io_writeAnalog(ch.MOTOR,0)
	case MD_up:
		io_clearBit(io.MOTORDIR)
		io_writeAnalog(ch.MOTOR,2800)
	case MD_stop:
		io_setBit(io.MOTORDIR)
		io_writeAnalog(ch.MOTOR,2800)
	}
}