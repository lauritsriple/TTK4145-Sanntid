package driver
/*
#cgo CFLAGS: -std=c11
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
#include "channels.h"
*/
import "C"
import "log"

//Private functions
func checkError(err error){
	if err != nil {
		log.Fatal("Error interfacing the c-driver ",err)
	}
}

//Public functions
func io_Init() bool {
	n, err := C.io_init()
	checkError(err)
	return citob(n)
}


func SetBit(channel int){
	_,err := C.io_set_bit(C.int(channel))
	checkError(err)
}

func ClearBit(channel int){
	_,err := C.io_clear_bit(C.int(channel))
	checkError(err)
}

func WriteAnalog(channel,value int){
	_,err := C.io_write_analog(C.int(channel),C.int(value))
	checkError(err)
}

func ReadBit(channel int) bool {
	n,err := C.io_read_bit(C.int(channel))
	checkError(err)
	return citob(n)
}

func ReadAnalog(channel int) int{
	n,err := C.io_read_analog(C.int(channel))
	checkError(err)
	return int(n)
}

func citob(i C.int) bool{
	if i==0{
		return false
	} else {
		return true
	}
}

