package imagecustomizer

import(
    "github.com/xabarass/image-builder/lib/images"
)

type customizeJob struct{
    image *images.ScionImage    
    configDirectory string
    destinationDir string

    jobId string
}

type CustomizeJobRequester interface {
    OnCustomizeJobSuccess(image *images.ScionImage, jobId string, generatedFile string)
    OnCustomizeJobError(image *images.ScionImage, jobId string, err error)
}