package imagecustomizer

import(
    "log"
    "os/exec"

    "github.com/xabarass/image-builder/lib/images"
    "github.com/xabarass/image-builder/lib/httpinterface"
)

func customizeImage(readyImage *images.ScionImage, job httpinterface.JobInfo, finishedImage chan *images.ScionImage, jobRequester httpinterface.JobRequester){
    log.Printf("Starting to customize image: %s for job id: %s", job.ImageName, job.JobId)    
    if(exec.Command("./customize.sh", job.ConfigFile, readyImage.GetPathFor(images.Home), readyImage.GetPathFor(images.Etc), readyImage.GetPathFor(images.ImgFile)).Run()==nil){
        log.Printf("Success customizing image!")
    }else{
        log.Printf("There was an error customizing image!")
    }

    jobRequester.JobFinished(job.JobId)

    finishedImage<-readyImage
}