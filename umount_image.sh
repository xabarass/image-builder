#!/bin/bash

if [ "$#" -ne 4 ]; then
    echo "Usage: $0 IMAGE ROOT HOME ETC"
    exit 1
fi

img_file=$1
root_dir=$2
home_dir=$3
etc_dir=$4

lostup_out=$(losetup -j "$img_file")
if [[ ! -z "$lostup_out" ]]
then
    # WARNING! This thing is really image and platform dependent
    loop_dev_name=$(echo "$lostup_out" | awk -F'[/:]' '{print $(3)}' )    
    loop_dev_name="/dev/mapper/${loop_dev_name}p2"

    mountpoint "$home_dir"
    if [ $? -eq 0 ]; then
        echo "Unmounting $home_dir"
        umount "$home_dir"
    fi

    mountpoint "$etc_dir"
    if [ $? -eq 0 ]; then
        echo "Unmounting $etc_dir"
        umount "$etc_dir"
    fi

    mountpoint "$root_dir"
    if [ $? -eq 0 ]; then
        echo "Unmounting $root_dir"
        umount "$root_dir"
    fi

    kparted_out=$(kpartx -d -v -s "$img_file")

else
    echo "Image is not mounted, skipping umount part..."
fi

echo "done"
