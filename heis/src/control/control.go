package control

import (
	"../udp"
	"../driver"
	"../localQueue"
	"log"
	"time"
)

var(
	myID int
	isIdle=true
	lastOrder=uint(0)
	floorOrder=make(chan uint,5)
	setLight=make(chan driver.Light,5)
	liftStatus=driver.Status
	maxFloor=driver.N_FLOORS
	button driver.Button
	message udp.Message
)

func RunLift(quit *chan bool){
	quit *chan bool
	buttonPress:=make(chan driver.Button,5)
	status:=make(chan driver.Status,5)
	toNetwork=make(chan udp.Message,10)
	fromNetwork=make(chan udp.Message,10)
	myID=udp.NetInit(toNetwork,fromNetwork,quit)
	fsm.Init(floorOrder,setLight,status,buttonPress,quit)
	restoreBackup()
	liftStatus <- status
	ticker1:=time.NewTicker(10*time.Millisecond).C
	ticker2:=time.NewTicker(5*time.Millisecond).C
	log.Println("Network UP \n Driver UP \n My id:",myID)
	for {
		select{
		case button=<-buttonPress:
			newKeypress(button)
		case liftStatus=<-status:
			runQueue()
		case message=<-fromNetwork:
			newMessage(message)
			orderLight(message)
		case <-ticker1:
			checkTimeout()
		case <-ticker2:
			runQueue()
		case <-*quit:
			return
		}
	}
}

func newKeyPress(button driver.Button){
	switch driver.Button{
	case driver.Up:
		log.Println("Request up button pressed:",button.Floor)
		addMessage(button.Floor,true)
		setOrderLight(button.Floor,true,true)
	case driver.Down:
		log.Println("Request down button pressed: ",button.Floor
		addMessage(button.Floor,false)
		setOrderLight(button.Floor,false,true)
	case driver.Command:
		log.Println("Command button pressed: ",button.Floor)
		addCommand(button.Floor)
	case driver.Stop:
		log.Println("Stop button pressed")
		//Action is not needed for this project
	case driver.Obstruction:
		log.Println("Obstruction button pressed")
		//Action is not needed for this project
	}
}

// Called by RunLift
func runQueue(liftStatus driver.Status, floorOrder chan<- uint){
	floor:=liftStatus.Floor
	if liftStatus.Running{
		if liftStatus.Direction{
			floor++
		}
		else{
			floor--
		}
	}

	// Get order from localQueue
	order,direction:=localQueue.GetOrder(floor,liftStatus.Direction)
	// Reported status is the ordered floor and door open
	if liftStatus.Floor == order && liftStatus.Door{
		removeFromQueue(order,direction)
		lastOrder=0
		liftStatus.Door=true
		time.sleep(20*time.Millisecond)
	} else if order==0 && !liftStatus.Door{ // No order and door closed, idle elevator
		isIdle=true
	} else if order != 0{ // We have an order
		isIdle=false
		if lastOrder!= order && !liftStatus.Door{
			// Check that the new order is not the same as last order, and door closed
			lastOrder=order
			floorOrder <- order
		}
	}
}

// Called by runQueue
func removeFromQueue(floor uint, direction bool,setLight chan<- driver.Light){
	localQueue.DeleteLocalOrder(floor,direction)
	messageparser.DelMessage(floor,direction)
	setLight <- driver.Light{floor,driver.Command,false}
	setOrderLight(floor,direction,false)
}

// Called by RunLift
func orderLight(message udp.Message){
	switch message.Status{
	case udp.Done:
		setOrderLight(message.Floor,message.Direction,false)
	case udp.New:
		setOrderLight(message.Floor,message.Direction,true)
	case udp.Accepted:
		setOrderLight(message.Floor,message.Direction,true)
	}
}

// Called by RunLift and ReadQueueFromFile
func addCommand(floor uint,setLight chan<- driver.Light){
	localQueue.AddLocalCommand(floor)
	setLight <- driver.Light{floor,driver.Command,true}
}

// Called by RunLift
func restoreBackup(setLight chan<- driver.Light){
	for i,val:=range localQueue.ReadQueueFromFile(){
		if val{
			addCommand(uint(i+1),setLight)
		}
	}
}
