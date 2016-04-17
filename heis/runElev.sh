#!/bin/bash

for i in "$@"
do
gnome-terminal --window-with-profile=hold -e "ssh student@129.241.187.$i '/home/student/Documents/TTK4145-Sanntid/heis/localRunElev.sh'"
done
wait
