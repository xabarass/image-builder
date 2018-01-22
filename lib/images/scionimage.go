package images;

import (
    "path"
)

const (
    Root=iota
    Etc=iota
    Home=iota
)

var mountPoints=[...]string{"mnt_root","mnt_etc","mnt_home"}

type ScionImage struct{
    File string             `json:"path"`
    Directory string        `json:"mount_directory"`
    Name string             `json:"name"`
    DisplayName string      `json:"display_name"`
    Description string      `json:"description"`
    Version string          `json:"version"`

    used bool
    mounted bool
}

func (si *ScionImage)GetPathFor(what int64) string{
    return path.Join(si.Directory, mountPoints[what]);
}

func (si *ScionImage) IsMounted()(bool) {
    return si.mounted;
}

func (si *ScionImage) IsUsed()(bool) {
    return si.used;
}

func (si *ScionImage)SetMounted(mounted bool){
    si.mounted=mounted
}

func (si *ScionImage)SetUsed(used bool){
    si.used=used
}