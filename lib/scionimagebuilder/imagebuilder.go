package scionimagebuilder;

import (
    "time"
    "fmt"
    "log"
    "os/exec"
    "os"
    "io"

)

func (ib *ScionImageBuilder)buildImage(job *buildJob)(error){
    fmt.Println("Starting to build image")

    err:=job.copyImage(ib)
    if err!=nil{
        log.Println("Failed copying image!")
        return err
    }

    loopDevice, err:=job.resizeAndMount(ib)
    if err!=nil{
        log.Println("Failed resizing image and mounting loop device!")
        return err
    }
    defer job.closeLoopDevice(ib);

    log.Printf("Loop device is %s", loopDevice)
    err=job.installScion(ib, "loop1p2") //FIXME! FIgure out whats wrong with this???
    if err!=nil{
        log.Println("Failed installing SCION!")
        return err
    }

    return nil
}

func (ib *ScionImageBuilder) Run(stop <-chan bool)(error){
    
    fmt.Println("Running SCION image builder!...")
    
    go func(){

        LOOP: for{
            select{
            case job:=<-ib.buildJobs:
                fmt.Printf("Got new build job for %s ", job.inputImage)
                err:=ib.buildImage(job)
                res:=ImageBuildResult{Success:true, FromImage:job.inputImage, OutputImage:job.outputImage}
                if err!=nil{
                    res.Success=false
                    res.Error=err;
                }
                job.result<-res
            case <-stop:
                log.Println("ImageBuilder >> Got request to shutdown!");
                break LOOP;
            }    
        }
        
        log.Println("Exiting image builder thread!");
    }()

    return nil
}

func (ib *ScionImageBuilder) StartBuildJob(inImage string, outputImagePath string, result chan<- ImageBuildResult)(error){
    newJob:=&buildJob{id:1, inputImage:inImage, outputImage:outputImagePath, result:result}

    fmt.Println("Adding new job to queue!")

    ib.buildJobs<-newJob

    return nil
}

func generateTimestamp()(string){
    t := time.Now()
    return t.Format("010215040500")
}

func (ib *ScionImageBuilder)logCommandOutput(logName string, cmd *exec.Cmd, output io.Writer)(error){

    outfile, err := os.Create(fmt.Sprintf("%s/log-%s-%s", ib.BuildLogDirectory, logName, generateTimestamp()))
    if err != nil {
        return err
    }
    defer outfile.Close()

    mwriter:=io.MultiWriter(outfile, os.Stdout, output)

    cmd.Stderr = mwriter
    cmd.Stdout = mwriter

    err=cmd.Run()
    return err
}