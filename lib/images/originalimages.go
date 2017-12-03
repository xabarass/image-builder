package images;

const (
    RaspberryPi2=iota
    Odroid=iota
)

type OriginalImage struct{
    Type    int `json:"type"`
    Path string `json:"path"`
    // This name should be unique across all images, identifier
    Name string `json:"name"`

    ScionImages []ScionImage
}

