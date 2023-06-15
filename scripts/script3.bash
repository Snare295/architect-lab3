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
pc "figure 150 150"
pc "figure 150 650"
pc "figure 650 150"
pc "figure 650 650"
pc "update"

for i in {0..1000}
do
  let x=i%15-9 
  let y=i%10-5 
  updMv $x $y 1
done
