#!/bin/bash

if [ "$#" -lt 1  ]; then
    echo "Usage: $0 LOOP_DEVICE"
    exit 1
fi

loop_dev=$1

echo "loop device is: $loop_dev"

cd /home/milan/projects/scion/rpi-img-builder

echo "My directory is"
pwd

make LOOP_DEV_NAME="$loop_dev"

# Unmount and clean stuff (Doesn't remove output image)
make clean