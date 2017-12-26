package imagecustomizer

import(
    "log"

    "github.com/xabarass/image-builder/lib/images"
    "github.com/xabarass/image-builder/lib/httpinterface"
)

type newScionImage struct{
    image *images.ScionImage    
    name string
}

type ImageCustomizer struct{
    scionImages map[string][]*images.ScionImage
    jobs []httpinterface.JobInfo
    jobRequester httpinterface.JobRequester

    newScionImages chan newScionImage
    newJobs chan httpinterface.JobInfo
    finishedJobs chan *images.ScionImage

    stop chan bool
}

func Create(jobRequester httpinterface.JobRequester)(*ImageCustomizer){
    ic:=new(ImageCustomizer)

    ic.jobRequester=jobRequester
    ic.scionImages=make(map[string][]*images.ScionImage)
    ic.jobs=make([]httpinterface.JobInfo,0)

    ic.newScionImages=make(chan newScionImage, 10)
    ic.newJobs=make(chan httpinterface.JobInfo, 10)
    ic.finishedJobs=make(chan *images.ScionImage, 10)

    ic.stop=make(chan bool, 1)

    return ic
}

func (ic *ImageCustomizer)ScionImageReady(imageName string, scionImage *images.ScionImage){
    log.Printf("Scion image is ready! Adding it to queue")
    ic.newScionImages<-newScionImage{image:scionImage, name:imageName}
} 

func (ic *ImageCustomizer)AddJob(newJob httpinterface.JobInfo){
    log.Printf("Adding new job to queue")
    ic.newJobs<-newJob
}

func (ic *ImageCustomizer)findAvailableScionImage(name string)(*images.ScionImage){
    log.Printf("Looking for image %s, there are %d scion images in list", name, len(ic.scionImages[name]))

    for i:=0; i<len(ic.scionImages[name]); i++ {

        if (!ic.scionImages[name][i].Used){
            return ic.scionImages[name][i]
            log.Printf("Image is NOT used")
        }else{
            log.Printf("Image is used")
        }
    }

    return nil
}

func (ic *ImageCustomizer)schedule(){
    log.Printf("Scheduling")
    for i:=0; i<len(ic.jobs); i++ {
        if readyImage:=ic.findAvailableScionImage(ic.jobs[i].ImageName); readyImage!=nil{
            // We found available image
            readyImage.Used=true
            job:=ic.jobs[i]

            ic.jobs=append(ic.jobs[:i], ic.jobs[i+1:]...)   //delete job from list

            log.Printf("Found ready image! Schedule complete")
            go customizeImage(readyImage, job, ic.finishedJobs, ic.jobRequester)

            return
        }
    }

    log.Printf("Impossible to schedule job, not available images")
}

func (ic *ImageCustomizer)Run(){
    go func(){
        LOOP: for{
            log.Printf("ImageCustomizer: Starting to wait for requests")
            select{
            
            case newJob:=<-ic.newJobs:
                ic.jobs=append(ic.jobs, newJob)
                ic.schedule()

            case newImage:=<-ic.newScionImages:
                log.Printf("Adding new image to queue of images %s", newImage.name)
                ic.scionImages[newImage.name]=append(ic.scionImages[newImage.name], newImage.image)
                ic.schedule()

            case availableImage:=<-ic.finishedJobs:
                availableImage.Used=false
                ic.schedule()            

            case <-ic.stop:
                log.Println("ImageCustomizer >> Got request to shutdown!");
                break LOOP;
            } 
        }
    }()
}

func (ic *ImageCustomizer)Stop(){
    log.Printf("Received stop command!")
    ic.stop<-true
}
