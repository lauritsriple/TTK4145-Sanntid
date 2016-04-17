#!/bin/bash

for i in "$@"
do
gnome-terminal -x bash -c "ssh $i < localRunElev.sh" &
done
