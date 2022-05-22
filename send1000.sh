#!/bin/bash

read -r -d '' msg << EOM
* * * FIRE WARN 20210320T221735Z003
HOUSE ON FIRE                 13-41
Fire started on 2nd floor by kids!!
END
EOM

x=1000
while [[ $x -gt 0 ]]; do
  x=$((x-1))
  echo "$msg" | nc -4u -w0 localhost 1961
done
