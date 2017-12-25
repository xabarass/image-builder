package imagemanager;

import(
    "github.com/xabarass/image-builder/lib/images"
)

type Configuration struct {
    Images []images.OriginalImage  `json:"images"`
    DBPath string           `json:"db_path"`
    BuildConfigurationPath string  `json:"build_configuration"`

    BuildOutputDirectory string     `json:"output_dir"`
}
