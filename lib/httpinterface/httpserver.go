package httpinterface

import (
    "log"
    "time"
    "path"
    "os"
    "io"
    "fmt"

    "mime/multipart"
    "github.com/xabarass/image-builder/lib/utils"
)

type jobInfo struct {
    JobId string
    DestDir string
    ConfigFile string
    ImageName string
    CreatedImage string
    timestamp time.Time
    finished bool
}

func (hi *HttpInterface)createNewJob(imageName string, configFile multipart.File)(string, error){
    log.Printf("Creating new job for image %s", imageName)

    jobId:=utils.GenerateRandomString(64);

    destDir:=path.Join(hi.rootDir, jobId)
    os.MkdirAll(destDir, os.ModePerm);
    confFileDest:=path.Join(destDir, "config.tar.gz")

    dest, err := os.OpenFile(confFileDest, os.O_WRONLY|os.O_CREATE, 0666)
    defer dest.Close()
    io.Copy(dest, configFile)

    newBuildJob:=jobInfo{
        JobId:jobId,
        DestDir:destDir,
        ConfigFile:confFileDest,
        ImageName:imageName,
        timestamp:time.Now(),
        finished:false,
    }

    hi.addJob(&newBuildJob)
    err = hi.imgMgr.RunJob(imageName, confFileDest, destDir, jobId)
    if(err!=nil){
        log.Printf("Error starting build job")
        //TODO: Cleanup (delete directory and uploaded files)
        hi.removeJob(jobId)
        return jobId, err
    }

    return jobId, nil
}

func (hi *HttpInterface)JobFinished(jobId string, createdFile string){
    log.Printf("Marking job %s as finished", jobId)
    if job, ok := hi.getJob(jobId); ok{
        log.Printf("Job has been marked as finished, file is: %s", createdFile)
        job.timestamp=time.Now()
        job.finished=true
        job.CreatedImage=createdFile
    }else{
        log.Printf("Requested jobId doesn't exist...")
    }
}

func (hi *HttpInterface)getImageForJob(jobId string)(bool, string, error){
    if job, ok := hi.getJob(jobId); ok{
        return job.finished, job.CreatedImage, nil
    }else{
        return false, "", fmt.Errorf("Provided ID doesn't match any active job")
    }
}
