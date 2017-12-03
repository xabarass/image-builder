package imagemanager;

import(
    "github.com/xabarass/image-builder/lib/images"
)

type Configuration struct {
    Images []OriginalImage  `json:"images"`
    DBPath string           `json:"db_path"`
}
