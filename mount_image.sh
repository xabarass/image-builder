#!/bin/bash

if [ "$#" -ne 5 ]; then
    echo "Usage: $0 IMAGE ROOT HOME ETC USER"
    exit 1
fi

img_file=$1
root_dir=$2
home_dir=$3
etc_dir=$4
regular_user=$5

lostup_out=$(losetup -j "$img_file")
if [[ ! -z "$lostup_out" ]]
then
    # WARNING! This thing is really image and platform dependent
    loop_dev_name=$(echo "$lostup_out" | awk -F'[/:]' '{print $(3)}' )    
    loop_dev_name="/dev/mapper/${loop_dev_name}p2"
else
    kparted_out=$(kpartx -a -v -s "$img_file")
    loop_dev_name=$(echo "$kparted_out" | awk 'END {print $(3)}' )
    loop_dev_name="/dev/mapper/${loop_dev_name}"
fi

# Unmount previously mounted partitions, if mounted
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

set -e

echo "Creating directories"
mkdir -p "$root_dir"
mkdir -p "$etc_dir"
mkdir -p "$home_dir"

# Mount root partition
echo "Mounting $root_dir"
mount "$loop_dev_name" "$root_dir"

# Mount 
echo "Mounting bindfs partitions"
sudo bindfs "--map=root/${regular_user}:@root/@${regular_user}" "${root_dir}/etc" "$etc_dir"
sudo bindfs "--map=scion/${regular_user}:@scion/@${regular_user}" "${root_dir}/home/scion" "$home_dir"

echo "done"
