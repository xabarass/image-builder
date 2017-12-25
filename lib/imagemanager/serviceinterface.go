package imagemanager;

import(
    "log"

    "github.com/xabarass/image-builder/lib/httpinterface"
)

func (im *ImageManager)GetAvailableImages()([]httpinterface.AvailableImage){
    return []httpinterface.AvailableImage{
        {Device:"rpi2",Name:"ubuntu"},
        {Device:"odroid",Name:"ubuntu"},
    }
}

func (im *ImageManager)CreateBuildJob(jobId string, configPath string, imageName string)(error){
    log.Printf("Starting build job for: %s at: %s", imageName, configPath)

    return nil
}