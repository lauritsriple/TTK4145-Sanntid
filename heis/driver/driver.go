package driver

import (
	"io"
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
	MD_up   = 1
	MD_down = -1
	MD_stop = 0
)

type ByttonType int //Enum for buttons
const (
	Up ButtonType = iota
	Down
	Command
	Stop
	Obstruction
	door //Not actual button, but used for door light, therefore not exported
)

type Button struct {
	Floor  int
	Button ButtonType
}

type Light struct {
	Floor  uint
	Button ButtonType
	On     bool
}

type LiftStatus struct {
	Running   bool
	Floor     uint
	Direction bool
	Door      bool
}

var floorSensorChannels = [N_FLOORS]int{
	SENSOR_FLOOR1,
	SENSOR_FLOOR2,
	SENSOR_FLOOR3,
	SENSOR_FLOOR4}

var buttons = []int{
	FLOOR_COMMAND1,
	FLOOR_COMMAND2,
	FLOOR_COMMAND3,
	FLOOR_COMMAND4,
	FLOOR_UP1,
	FLOOR_UP2,
	FLOOR_UP3,
	FLOOR_UP4,
	FLOOR_DOWN1,
	FLOOR_DOWN2,
	FLOOR_DOWN3,
	FLOOR_DOWN4,
	STOP,
	OBSTRUCTION}

var buttonsKeyType = []int{
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
		if io.Init() == false {
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
			io.SetBit(lightmap[lightKeyType[int(light.Button)]+int(light.Floor)])
		} else {
			io.ClearBit(lightmap[lightKeyType[int(light.Button)]+int(light.Floor)])
		}
	}
}

func SetFloorIndicator(floor int) {
	if (floor < 1) || (floor > N_FLOORS) {
		log.Fatal("Floororder out of range: ", floor)
	}
	switch floor {
	case 1:
		io.ClearBit(LIGHT_FLOOR_IND1)
		io.ClearBit(LIGHT_FLOOR_IND2)
	case 2:
		io.ClearBit(LIGHT_FLOOR_IND1)
		io.SetBit(LIGHT_FLOOR_IND2)
	case 3:
		io.SetBit(LIGHT_FLOOR_IND1)
		io.ClearBit(LIGHT_FLOOR_IND2)
	case 4:
		io.SetBit(LIGHT_FLOOR_IND1)
		io.SetBit(LIGHT_FLOOR_IND2)

	}
}

func readButtons(keypress chan<- Button) {
	for index, key := range buttons {
		if readButton(index, key) {
			keypress <- Button{uint(index%4 + 1), buttonsKeyType[index]}
		}
	}
}

func readButton(key int, index int) bool {
	if io.ReadBit(key) {
		if !lastPress[index] {
			lastPress[index] = true
			return true
		}
	} else if lastPress[index] {
		lastPress[index] = false
	}
}

func ReadFloorSensors(floorSeen chan<- int) {
	atFloor := false
	for f := 0; f < N_FLOORS; f++ {
		sensor := io.ReadBit(floorSensorChannels[f])
		if sensor != 0 && sensor != currentFloor {
			currentFloor = sensor
			atFloor = true
			floorSeen <- sensor
			return
		}
	}
	if !atFloor && sensor != currentFloor {
		currentFloor = sensor
		floorSeen <- uint(sensor)
	}
}

func SetMotorDir(dir MotorDirection) {
	switch dir {
	case MD_stop:
		io.WriteAnalog(MOTOR, 0)
	case MD_up:
		io.ClearBit(MOTORDIR)
		io.WriteAnalog(MOTOR, 2800)
	case MD_down:
		io.SetBit(MOTORDIR)
		io.WriteAnalog(MOTOR, 2800)
	}
}
