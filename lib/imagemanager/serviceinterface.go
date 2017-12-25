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

func (im *ImageManager)RunJob(job httpinterface.JobInfo)(error){
    log.Printf("Starting build job for: %s at: %s", job.ImageName, job.ConfigFile)

    return nil
}
