#!/bin/bash

for i in "$@"
do
ssh $i < localDeploy.sh
done