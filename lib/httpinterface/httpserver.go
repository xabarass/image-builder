package httpinterface

import (
    "net/http"
    "log"
    "time"
    "path"
    "os"
    "io"
    "fmt"

    "mime/multipart"
    "github.com/xabarass/image-builder/lib/utils"
)

type JobInfo struct {
    JobId string
    DestDir string
    ConfigFile string
    ImageName string
    CreatedImage string
    timestamp time.Time
    finished bool
}

type HttpInterface struct {
    imgMgr ImageBuilderService
    rootDir string

    activeJobs map[string]*JobInfo

    server *http.Server
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

    newBuildJob:=JobInfo{
        JobId:jobId,
        DestDir:destDir,
        ConfigFile:confFileDest,
        ImageName:imageName,
        timestamp:time.Now(),
        finished:false,
    }

    hi.activeJobs[jobId]=&newBuildJob
    err = hi.imgMgr.RunJob(newBuildJob)

    if(err!=nil){
        //TODO: Cleanup (delete directory and uploaded files)
        delete(hi.activeJobs, jobId)
        return jobId, err
    }

    return jobId, nil
}

func (hi *HttpInterface)JobFinished(jobId string){
    if job, ok:=hi.activeJobs[jobId]; ok{
        job.timestamp=time.Now()
        job.finished=true
    }else{
        log.Printf("Requested jobId doesn't exist...")
    }
}

func (hi *HttpInterface)getImageForJob(jobId string)(bool, string, error){
    if job, ok:=hi.activeJobs[jobId]; ok{
        return job.finished, job.CreatedImage, nil
    }else{
        return false, "", fmt.Errorf("Provided ID doesn't match any active job")
    }
}

func CreateHttpServer(bindAddress string, imgMgr ImageBuilderService, authorizedTokens map[string]bool, rootDir string) *HttpInterface {
    hi:=&HttpInterface{
        imgMgr:imgMgr,
        rootDir:rootDir,
        activeJobs:make(map[string]*JobInfo),
    }

    srv := &http.Server{
        Addr: bindAddress, 
        Handler: createHandler(hi, authorizedTokens),
    }

    hi.server=srv

    return hi
}

func (hi *HttpInterface)StartServer(){
    go func() {
        if err := hi.server.ListenAndServe(); err != nil {
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }
    }()
}

func (hi *HttpInterface)StopServer(){
    hi.server.Shutdown(nil)
}