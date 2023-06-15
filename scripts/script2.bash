#!/bin/bash


pc(){
  curl -X POST http://localhost:17000 -d "$1" 
}

updMv(){
for ((i = 0; i < $3; i++))
do
  pc "move $1 $2"
  pc "update"
  sleep 0.05
done
}

pc "reset"
pc "white"
pc "figure 0 0"

updMv 20 0 20
updMv 0 20 20
updMv -10 -10 40
updMv 40 80 5
updMv 30 0 10
updMv -20 -20 10

pc "update"
