#!/bin/bash

for i in "$@"
do
gnome-terminal --window-with-profile=hold -e "ssh $i '/home/student/Documents/TTK4145-Sanntid/heis/localRunElev.sh'"
done
wait
