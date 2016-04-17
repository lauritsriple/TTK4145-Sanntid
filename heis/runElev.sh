#!/bin/bash

for i in "$@"
do
ssh $i < localRunElev.sh
done
