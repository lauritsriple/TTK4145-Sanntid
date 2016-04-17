package localqueue

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

const N_FLOORS = 4

type OrderQueue struct {
	Up      [N_FLOORS]bool //Request
	Down    [N_FLOORS]bool //Request
	Command [N_FLOORS]bool //Commands
}

var backupFile = "backupQueue.txt"
var localQueue = OrderQueue{}

func writeQueueToFile() {
	commandQueue, err := json.Marshal(localQueue.Command)
	if err != nil {
		log.Println(err)
	}
	err = ioutil.WriteFile(backupFile, commandQueue, 0600)
	if err != nil {
		log.Println("Error writing to file: ", err)
	}
}

func ReadQueueFromFile() []bool {
	byt, err := ioutil.ReadFile(backupFile)
	if err != nil {
		log.Println("Error opening backupFile", err)
	}
	var cmd []bool
	if err := json.Unmarshal(byt, &cmd); err != nil {
		log.Println("JSON:", err)
		log.Println("Got: ", cmd)
	}
	return cmd
}

func AddLocalCommand(floor uint) {
	localQueue.Command[floor-1] = true
	writeQueueToFile()
}

func AddLocalRequest(floor uint, Direction bool) {
	if Direction {
		localQueue.Up[floor-1] = true
	} else {
		localQueue.Down[floor-1] = true
	}
}

func DeleteLocalRequest(floor uint, Direction bool) {
	if Direction {
		localQueue.Up[floor-1] = false
	} else {
		localQueue.Down[floor-1] = false
	}
}

func DeleteLocalOrder(floor uint, Direction bool) {
	localQueue.Command[floor-1] = false
	writeQueueToFile()
	if Direction {
		localQueue.Up[floor-1] = false
	} else {
		localQueue.Down[floor-1] = false
	}
}

func GetOrder(currentFloor uint, direction bool) (uint, bool) {
	if direction {
		if nextStop := checkUp(currentFloor, N_FLOORS); nextStop > 0 {
			return nextStop, true
		} else if nextStop := checkDown(N_FLOORS, 1); nextStop > 0 {
			return nextStop, false
		} else {
			return checkUp(1, N_FLOORS), true
		}
	} else {
		if nextStop := checkDown(currentFloor, 1); nextStop > 0 {
			return nextStop, false
		} else if nextStop := checkUp(1, N_FLOORS); nextStop > 0 {
			return nextStop, true
		} else {
			return checkDown(N_FLOORS, 1), false
		}
	}
}

func checkUp(start uint, stop uint) uint {
	for i := int(start) - 1; i <= int(stop)-1; i++ {
		if i > N_FLOORS-1 || i < 0 {
			log.Println("In localqueue, checkUp out of bounds.Stop: ", stop, " start: ", start, "i: ", i)
			return 0
		} else if localQueue.Up[i] || localQueue.Command[i] {
			return uint(i + 1)
		}
	}
	return 0
}

func checkDown(start uint, stop uint) uint {
	for i := int(start) - 1; i >= int(stop)-1; i-- {
		if i > N_FLOORS-1 || i < 0 {
			log.Println("In localqueue, checkDown out of bounds.Stop: ", stop, " start: ", start, "i: ", i)
			return 0
		} else if localQueue.Down[i] || localQueue.Command[i] {
			return uint(i + 1)
		}
	}
	return 0
}
