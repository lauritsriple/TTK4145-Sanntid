package driver

import (
	"log"
	"time"
)

const N_FLOORS = 4
const N_LIGHTS = 4
const N_BUTTONS = 3

//Milliseconds between each polling round
const POLLRATE = 20 * time.Millisecond

type MotorDirection int

const (
	MD_up MotorDirection = iota
	MD_down
	MD_stop
)

type ButtonType int //Enum for buttons

const (
	Up ButtonType = iota
	Down
	Command
	Stop
	Obstruction
	door //Not actual button, but used for door light, therefore not exported
)

type Button struct {
	Floor  uint
	Button ButtonType
}

type Light struct {
	Floor  uint
	Button ButtonType
	On     bool
}

type LiftStatus struct {
	Running   bool
	Floor     uint // why?  	(driver.currentFloor)
	Direction MotorDirection // better?	(Direction MotorDirection)
	Door      bool
}

var floorSensorChannels = [N_FLOORS]int{
	SENSOR_FLOOR1,
	SENSOR_FLOOR2,
	SENSOR_FLOOR3,
	SENSOR_FLOOR4}

var buttons = []int{
	BUTTON_COMMAND1,
	BUTTON_COMMAND2,
	BUTTON_COMMAND3,
	BUTTON_COMMAND4,
	BUTTON_UP1,
	BUTTON_UP2,
	BUTTON_UP3,
	BUTTON_UP4,
	BUTTON_DOWN1,
	BUTTON_DOWN2,
	BUTTON_DOWN3,
	BUTTON_DOWN4,
	STOP,
	OBSTRUCTION}

var buttonsKeyType = []ButtonType{
	Command,
	Command,
	Command,
	Command,
	Up,
	Up,
	Up,
	Up,
	Down,
	Down,
	Down,
	Down,
	Stop,
	Obstruction}

var lightmap = []int{
	LIGHT_COMMAND1,
	LIGHT_COMMAND2,
	LIGHT_COMMAND3,
	LIGHT_COMMAND4,
	LIGHT_UP1,
	LIGHT_UP2,
	LIGHT_UP3,
	LIGHT_UP4,
	LIGHT_DOWN1,
	LIGHT_DOWN2,
	LIGHT_DOWN3,
	LIGHT_DOWN4,
	LIGHT_STOP,
	LIGHT_DOOR_OPEN}

var lightKeyType = []int{
	Command: -1,
	Up:      3,
	Down:    7,
	Stop:    12,
	door:    13}

var (
	currentFloor      = -1
	driverInitialized = false
	lastPress         [14]bool //Remembers last state of buttons
)

func Init() bool {
	if driverInitialized {
		log.Fatal("ERROR, driver already initialized")
	} else {
		driverInitialized = true
		if Io_Init() == false {
			log.Fatal("ERROR, could not initialize driver")
		} else {
			//sucess
			return true
		}
	}
	return false
}

func setLight(lightch chan Light) {
	select {
	default:
		return
	case light := <-lightch:
		if light.On {
			Io_SetBit(lightmap[lightKeyType[int(light.Button)]+int(light.Floor)])
		} else {
			Io_ClearBit(lightmap[lightKeyType[int(light.Button)]+int(light.Floor)])
		}
	}
}

func setFloorIndicator(floor int) {
	if (floor < 1) || (floor > N_FLOORS) {
		log.Fatal("Floororder out of range: ", floor)
	}
	switch floor {
	case 1:
		Io_ClearBit(LIGHT_FLOOR_IND1)
		Io_ClearBit(LIGHT_FLOOR_IND2)
	case 2:
		Io_ClearBit(LIGHT_FLOOR_IND1)
		Io_SetBit(LIGHT_FLOOR_IND2)
	case 3:
		Io_SetBit(LIGHT_FLOOR_IND1)
		Io_ClearBit(LIGHT_FLOOR_IND2)
	case 4:
		Io_SetBit(LIGHT_FLOOR_IND1)
		Io_SetBit(LIGHT_FLOOR_IND2)

	}
}

func ReadButtons(keypress chan<- Button) {
	for index, key := range buttons {
		if readButton(index, key) {
			keypress <- Button{uint(index%4 + 1), buttonsKeyType[index]}
		}
	}
}

func readButton(key int, index int) bool {
	if Io_ReadBit(key) {
		if !lastPress[index] {
			lastPress[index] = true
			return true
		}
	} else if lastPress[index] {
		lastPress[index] = false
	}
	return false
}

func ReadFloorSensors(floorSeen chan<- uint) {
	atFloor := false
	for f := 0; f < N_FLOORS; f++ {
		sensor := Io_ReadBit(floorSensorChannels[f])
		if sensor && (f+1) != currentFloor {
			currentFloor = f+1
			floorSeen <- uint(currentFloor)
			setFloorIndicator(f+1)
			atFloor = true
			return
		}
	}
	if !atFloor && currentFloor !=0{
		currentFloor = 0
		floorSeen <- uint(0)
	}
}

func SetMotorDir(dir MotorDirection) {
	switch dir {
	case MD_stop:
		Io_WriteAnalog(MOTOR, 0)
	case MD_up:
		Io_ClearBit(MOTORDIR)
		Io_WriteAnalog(MOTOR, 2800)
	case MD_down:
		Io_SetBit(MOTORDIR)
		Io_WriteAnalog(MOTOR, 2800)
	}
}
