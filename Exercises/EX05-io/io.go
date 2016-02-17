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
func checkError(err){
	if err != nil {
		log.Fatal("Error interfacing the c-driver ",err)
}

func cIntToBool(i C.int) bool{
	if i == 0{
		return false
	}
	else{
		return true
	}
}


//Public functions
func io_init() bool {
	n, err := C.io_init()
	checkError(err)
}

func io_setBit(channel int){
	_,err := C.io_set_bit(C.int(channel))
	checkError(err)
}

func io_clearBit(channel int){
	_,err := C.io_clear_bit(C.int(channel))
	checkError(err)
}

func io_writeAnalog(channel,value int){
	_,err := C.io_write_analog(C.int(channel,C.int(value))
	checkError(err)
}

func io_readBit(channel int) bool {
	n,err := C.io_read_bit(C.int(channel))
	checkError(err)
	return cIntToBool(n)
}

func io_readAnalog(channel int) int{
	n,err := C.io_read_analog(C.int(channel))
	checkError(err)
	return int(n)
}