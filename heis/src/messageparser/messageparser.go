package messageparser

import (
	"../udp"
	"../localQueue"
	"log"
	"time"
)

const acceptedTimeoutBase=4 //seconds
const newTimeoutBase = 500 //milliseconds

var globalQueue=make(map[uint]udp.Message)

func generateKey(floor uint, direction bool) uint{
	if direction{
		floor+=10
	}
	return floor
}

func addMessage(floor uint, direction bool,toNetwork chan<- udp.Message){
	key:= generateKey(floor, direction)
	message:= udp.Message{
		LiftId:myID,
		Floor:floor,
		Direction,direction,
		Status: udp.New,
		Weight: cost(floor,direction),
		TimeRecv:time.Now()}

	if _,inQueue:=globalQueue[key];inQueue{
		return
	}
	globalQueue[key]=message
	toNetwork<-message
}

func delMessage(floor uint, direction bool){
	key:=generateKey(message.Floor,message.Direction)
	val,inQueue:=globalQueue[ke]
	if inQueue{
		switch message.Status{
		case udp.Done:
			delete(globalQueue,key)
		case udp.Accepted:
			globalQueue[key]=message
		case udp.New:
			if val.Weight <= message.Weight{
				globalQueue[key]=message
			}
		case udp.Reassign:
			if message.ReassId!=myID{
				if val.Weight <= message.Weight{
					globalQueue[key]=message
				}
			}
		default:
			log.Println("Unknown status: ", message.Status, "Ignoring message")
		}
	}else{
		switch message.Status{
		case udp.Done:
			//Promptly ignore
		case udp.Accepted:
			if val.Status==udp.Reassign && val.ReassId==myID{
				localQueue.DeleteLocalRequest(message.Floor, message.Direction)
			}
			globalQueue[key]=message
		case udp.Reassign:
			fs:=cost(message.Floor,messageDirection)
			if fs > message.Weight{
				message.Weight=fs
				message.LiftId=myID
				globalQueue[key]=message
				toNetwork<-message
				log.Println("Reassign from lift ", message.reassId," to ",myID)
			}else{
				globalQueue[key]=message
			}
		case udp.New:
			fs:=cost(message.Floor,message.Direction)
			if fs > message.Weight{
				message.Weight=fs
				message.LiftId=myID
				globalQueue[key]=message
				toNetwork<-message
			}else{
				globalQueue[key]=message
			}
		default:
			log.Println("Unknown status: ",message.Status, "Ignoring message")
		}
	}
}


func checkTimeout(){
	newTimeout:=time.Duration(newTimeoutBase)
	acceptTimeout:=time.Duration(acceptedTimeoutBase)
	for key,val:=range globalQueue{
		if val.Status==udp.New || val.Status == udp.Reassign{
			timediff:= time.Now().Sub(val.TimeRecv)
			if timediff >((3*newTimeout)*time.Milisecond){
				newOrderTimeout(key,3)
			}else if timediff >((2*newTimeout)*time.Millisecond){
				newOrderTimeout(key,2)
			}else if timediff >((1* newTimeout)*time.Millisecond){
				newOrderTimeout(key,1)
			}
		} else if val.Status == udp.New && val.LiftId != myID{
			timediff:=time.Now().Sub(val.TimeRecv)
			if timediff > ((4*acceptedTimeout)*time.Second){
				acceptedOrderTimeout(key,3)
			} else if timediff >((3*acceptedTimeout)*time.Second){
				acceptedOrderTimeout(key,2)
			} else if timediff >((2*acceptedTimeout=*time.Second){
				acceptedOrderTimeout(key,1)
			}
		} else if val.Status == udp.Accepted && val.LiftIf==myID {
			timediff:=time.Now().Sub(val.TimeRecv)
			if timediff > (acceptedTimeout * time.Second){
				val.Weight=cost(val.Floor,val.Direction)
				val.TimeRecv=time.Now()
				globalQueue[key]=val
				toNetwork<-globalQueue[key]
			}
		}
	}
}

func newOrderTimeout(key,critical uint){
	switch critical{
	case 3:
		takeOrder(key)
	case 2:
		if isIdle{
			takeOrder(key)
		} else if cost(globalQueue[key].Floor,globalQueue[key].Direction) > globalQueue[key].Weight{
			takeOrder(key)
		}
	case 1:
		if globalQueue[key]==myID{
			takeOrder(key)
		}
	}
}

func acceptedTimeout(key uint, critical uint){
	switch critical{
	case 3:
		log.Println("ERROR! Reassigning orders failed. FALLBACK")
		takeOrder(key)
	case 2:
		takeOrder(key)
	case 1:
		reassignOrder(key)
	}
}

func takeOrder(key uint){
	if val,inQueue:=globalQueue[key];!inQueue{
		log.Println("Trying to accept order not in queue")
	} else{
		log.Println("Accepted order",globalQueue[key])
		val.LiftId=myID
		val.Status=udp.Accepted
		val.TimeRecv=time.Now()
		localQueue.AddLocalRequest(val.Floor,val.Direction)
		globalQueue[key]=val
		toNetwork<-globalQueue[key]
	}
}

func reassignOrder(key uint){
	if val,inQueue:=globalQueue[key];!inQueue{
		log.Println("Trying to reassign order not in queue")
	} else {
		log.Println("Reassigning order",globalQueue[key])
		val.Status=udp.Reassign
		val.ReassId=myID
		val.Weight=cost(val.Floor,val.Direction)
		val.TimeRecv=time.Now()
		globalQueue[key]=val
		toNetwork<-globalQueue[key]
	}
}

func cost(reqFloor uint, reqDir bool) int{
	statusFloor:=liftStatus.Floor
	statusDir:=liftStatus.Direction
	if isIdle{
		if reqFloor==statusFloor{
			return 6
		} else{
			return N_FLOORS+1-diff(reqFloor,statusFloor)
		}
	}else if reqDir == statusDir{
		if (statusDir && reqFloor > statusFloor) || (!statusDir && reqFloor < statusFloor){
			return N_FLOORS + 1 - diff(reqFloor,statusFloor)
		}
	} else {
		if (statusDir && reqFloor > statusFloor) || (!statusDir && reqFloor < statusFloor){
			return N_FLOORS - diff(reqFloor,statusFloor)
		}
	}
	return 1
}

diff(a uint, b uint){
	x:=int(a)
	y:=int(b)
	c:=x-y
	if c < 0{
		return c*-1
	} else{
		return c
	}
}
