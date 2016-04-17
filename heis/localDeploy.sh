#!/bin/bash

directory="/home/alarm/Documents/TTK4145-Sanntid"
topDirectory="/home/alarm/Documents/"
GOPATH=$directory+"/heis"
if [ ! -d "$directory" ]; then
	cd "$topDirectory"; git clone https://github.com/lauritsriple/TTK4145-Sanntid
else
	cd "$directory"; git pull
fi

go build $directory+"src/driver"
go build $directory+"src/localqueue"
go build $directory+"src/udp"
go build $directory+"src/control"
go bulid $directory+"src"
