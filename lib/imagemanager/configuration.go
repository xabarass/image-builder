package imagemanager;

import(
    "github.com/xabarass/image-builder/lib/images"
)

type Configuration struct {
    MountScript string              `json:"mount_script"`
    UmountScript string             `json:"umount_script"`
    CustomizeScript string          `json:"customize_script"`

    Images []*images.ScionImage            `json:"available_images"`

    BindAddress string              `json:"bind_address"`
    OutputDirectory string          `json:"output_directory"`
}
