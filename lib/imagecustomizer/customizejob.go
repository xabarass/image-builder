package imagecustomizer

import(
    "github.com/xabarass/image-builder/lib/images"
)

type customizeJob struct{
    image *images.ScionImage    
    configDirectory string
    destinationDir string
}

type CustomizeJobRequester interface {
    OnCustomizeJobSuccess(image *images.ScionImage, generatedFile string)
    OnCustomizeJobError(image *images.ScionImage, err error)
}