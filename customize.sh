#!/bin/bash

if [ "$#" -ne 5 ]; then
    echo "Usage: $0 CONFIG_DIR HOME ETC IMG_FILE DEST_DIR"
    exit 1
fi

set -e

config_dir=$1
home_dir=$2
etc_dir=$3
img_file=$4
dest_dir=$5

#img_file="${img_file}1"

echo "Copy gen folder"
rm -rf  "${home_dir}/go/src/github.com/scionproto/scion/gen"
cp -r "${config_dir}/gen" "${home_dir}/go/src/github.com/scionproto/scion"

rm -rf "${etc_dir}/openvpn/client.conf"
if [ -f "${config_dir}/client.conf" ]; then
    echo "Copy client OpenVPN configuration"
    cp "${config_dir}/client.conf" "${etc_dir}/openvpn"
fi

sync

lbzip2 -zk --fast $img_file

mv "${img_file}.bz2" "${dest_dir}/scion.img.bz2"

echo "Done"
