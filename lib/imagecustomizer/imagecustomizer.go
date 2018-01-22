package imagecustomizer

import(
    "log"

    "github.com/xabarass/image-builder/lib/images"
)

type ImageCustomizer struct{
    customizeScript string
    
    jobQueue chan *customizeJob
    jobRequester CustomizeJobRequester

    stop chan bool
}

const MAX_QUEUE_SIZE int = 10

func Create(customizeScript string, jobRequester CustomizeJobRequester)(*ImageCustomizer){
    ic:=new(ImageCustomizer)

    ic.jobRequester=jobRequester
    ic.customizeScript=customizeScript

    ic.jobQueue=make(chan *customizeJob, MAX_QUEUE_SIZE)
    ic.stop=make(chan bool, 1)

    return ic
}

func (ic *ImageCustomizer)CustomizeImage(image *images.ScionImage, configDirectory string, destinationDir string, jobId string){
    log.Printf("Creating new build job for %s", image.Name)
    ic.jobQueue<-&customizeJob{image:image, configDirectory:configDirectory, destinationDir:destinationDir, jobId:jobId}
}

func (ic *ImageCustomizer)Run(){
    go func(){
        LOOP: for{
            log.Printf("ImageCustomizer: Starting to wait for requests")
            select{
            
            case newJob:=<-ic.jobQueue:
                log.Println("Got request to customize image");       
                createdFile, err:=ic.customizeImage(newJob.image, newJob.configDirectory, newJob.destinationDir)
                if(err==nil){
                    ic.jobRequester.OnCustomizeJobSuccess(newJob.image, newJob.jobId, createdFile)
                }else{
                    ic.jobRequester.OnCustomizeJobError(newJob.image, newJob.jobId, err)
                }

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
