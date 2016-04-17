#!/bin/bash

directory="/home/student/Documents/TTK4145-Sanntid/heis"
topDirectory="/home/student/Documents/"
GOPATH="$directory"
export GOPATH="$directory"
if [ ! -d "$directory" ]; then
	cd "$topDirectory"; git clone https://github.com/lauritsriple/TTK4145-Sanntid
else
	cd "$directory"; git pull
fi

cd $directory"/src"; go build driver
cd $directory"/src"; go build localqueue
cd $directory"/src"; go build control
cd $directory"/src"; go build udp
cd $directory"/src"; go build
chmod +x $directory"/localRunElev.sh"
