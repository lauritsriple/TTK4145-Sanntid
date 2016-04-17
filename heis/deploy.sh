#!/bin/bash

for i in "$@"
do
ssh $i < test.sh
done
