package control

import (
	"udp"
	"driver"
	"log"
	"time"
	"localqueue"
	"fsmelev"
)

var(
	myID int
	isIdle=true
	lastOrder=uint(0)
	floorOrder=make(chan uint,5)
	setLight=make(chan driver.Light,5)
	liftStatus driver.LiftStatus
	button driver.Button
	message udp.Message
	toNetwork=make(chan udp.Message,10)
	fromNetwork=make(chan udp.Message,10)
	motorStopCh=make(chan bool,3)
)

//Called by main
//Initializes all the other modules and goes to eternal for loop
func RunLift(quitCh chan bool){
	var buttonPress=make(chan driver.Button,5)
	var status=make(chan driver.LiftStatus,5)
	myID=udp.NetInit(toNetwork,fromNetwork,quitCh)
	fsmelev.Init(floorOrder,setLight,status,buttonPress,motorStopCh,quitCh)
	restoreBackup()
	liftStatus =<-status
	ticker1:=time.NewTicker(10*time.Millisecond).C
	ticker2:=time.NewTicker(5*time.Millisecond).C
	log.Println("Network UP \n Driver UP \n My id:",myID)
	for {
		select{
		case button=<-buttonPress:
			newKeyPress(button)
		case liftStatus=<-status:
			runQueue(liftStatus,floorOrder)
		case message=<-fromNetwork:
			newMessage(message)
			orderLight(message)
		case <-ticker1:
			checkTimeout()
		case <-ticker2:
			runQueue(liftStatus,floorOrder)
		case <-quitCh:
			return
		}
	}
}

//Called by RunLift
func newKeyPress(button driver.Button){
	switch button.Button{
	case driver.Up:
		log.Println("Request up button pressed:",button.Floor)
		addMessage(button.Floor,true)
		setOrderLight(button.Floor,true,true)
	case driver.Down:
		log.Println("Request down button pressed: ",button.Floor)
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

//Called by RunLift
func runQueue(liftStatus driver.LiftStatus, floorOrder chan<- uint){
	floor:=liftStatus.Floor
	if liftStatus.Running{
		if liftStatus.Direction{
			floor++
		}else{
			floor--
		}
	}

	order,direction:=localqueue.GetOrder(floor,liftStatus.Direction)
	if liftStatus.Floor == order && liftStatus.Door{
		removeFromQueue(order,direction)
		lastOrder=0
		liftStatus.Door=true
		time.Sleep(20*time.Millisecond)
	} else if order==0 && !liftStatus.Door{
		isIdle=true
	} else if order != 0{
		isIdle=false
		if lastOrder!= order && !liftStatus.Door{
			lastOrder=order
			floorOrder <- order
		}
	}
}

//Called by runQueue
func removeFromQueue(floor uint, direction bool){
	localqueue.DeleteLocalOrder(floor,direction)
	delMessage(floor,direction)
	setLight <- driver.Light{floor,driver.Command,false}
	setOrderLight(floor,direction,false)
}

//Called by RunLift
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
//Called by orderLight
func setOrderLight(floor uint, direction bool, on bool){
	if direction{
		setLight<-driver.Light{floor,driver.Up,on}
	}else{
		setLight<-driver.Light{floor,driver.Down,on}
	}
}

//Called by RunLift and restoreBackup
func addCommand(floor uint){
	localqueue.AddLocalCommand(floor)
	setLight <- driver.Light{floor,driver.Command,true}
}

//Called by RunLift
func restoreBackup(){
	for i, val := range localqueue.ReadQueueFromFile(){
		if val{
			addCommand(uint(i+1))
		}
	}
}
