#!/bin/bash

x=1000
while [[ $x -gt 0 ]]; do
  msg="Pedroz: $RANDOM"
  x=$((x-1))
  echo $msg | nc -4u -w0 localhost 1961
done
