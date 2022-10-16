#!/bin/bash

if [[ $1 == "" ]]; then
   echo "You must enter the IP address of a Vector running WireOS."
   exit 0
fi

echo "Waking vector up..."

curl http://$1:8080/api/initSDK
curl -d "priority=high" http://$1:8080/api/assume_behavior_control

sleep 3
echo
echo
echo "Vector should now be ready to control! W for forward, A for left, R for right, S for backwards, Q for quit. Do not use two keys at a time."

isMovingForward=false
isMovingLeft=false
isMovingRight=false
isMovingBack=false
movementStopped=false

while true;
do
read -rsn1 -t 0.2 letter
if [[ ${letter} == "" ]] && [[ ${movementStopped} == false ]]; then
   isMovingForward=false
   isMovingLeft=false
   isMovingRight=false
   isMovingBack=false
   curl -d "lw=0&rw=0" http://$1:8080/api/move_wheels &
   movementStopped=true
fi
if [[ ${letter} == "w" ]] && [[ ${isMovingForward} == false ]]; then
   movementStopped=false
   curl -d "lw=150&rw=150" http://$1:8080/api/move_wheels &
   isMovingForward=true
   sleep 0.3
fi
if [[ ${letter} == "a" ]] && [[ ${isMovingLeft} == false ]]; then
   movementStopped=false
   curl -d "lw=-60&rw=120" http://$1:8080/api/move_wheels &
   isMovingLeft=true
   sleep 0.3
fi
if [[ ${letter} == "d" ]] && [[ ${isMovingRight} == false ]]; then
   movementStopped=false
   curl -d "lw=120&rw=-60" http://$1:8080/api/move_wheels &
   isMovingRight=true
   sleep 0.3
fi
if [[ ${letter} == "s" ]] && [[ ${isMovingBack} == false ]]; then
   movementStopped=false
   curl -d "lw=-150&rw=-150" http://$1:8080/api/move_wheels &
   isMovingBack=true
   sleep 0.3
fi
if [[ ${letter} == "q" ]]; then
   echo "Quitting"
   curl http://$1:8080/api/release_behavior_control
   exit 0
fi
done
