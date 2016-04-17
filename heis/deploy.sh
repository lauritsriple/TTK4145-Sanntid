#!/bin/bash

for i in "$@"
do
ssh "student@129.241.187.$i" < localDeploy.sh
done
