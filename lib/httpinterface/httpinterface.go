package httpinterface

import (
    "net/http"
    "log"
    "sync"
    "os"
)

type HttpInterface struct {
    imgMgr ImageBuilderService
    rootDir string

    mapLock *sync.Mutex
    activeJobs map[string]*jobInfo

    stop chan bool

    server *http.Server
}

func CreateHttpServer(bindAddress string, imgMgr ImageBuilderService, authorizedTokens map[string]bool, rootDir string) *HttpInterface {
    hi:=&HttpInterface{
        imgMgr:imgMgr,
        rootDir:rootDir,
        activeJobs:make(map[string]*jobInfo),
        stop:make(chan bool, 1),
        mapLock:&sync.Mutex{},
    }

    err:=os.RemoveAll(rootDir)
    if(err!=nil){
        log.Println(err)
    }

    srv := &http.Server{
        Addr: bindAddress, 
        Handler: createHandler(hi, authorizedTokens),
    }

    hi.server=srv

    return hi
}

func (hi *HttpInterface)StartServer(){
    hi.startCleanupService()
    go func() {
        if err := hi.server.ListenAndServe(); err != nil {
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }
    }()
}

func (hi *HttpInterface)StopServer(){
    hi.server.Shutdown(nil)
    hi.stopCleanupService()
}

func (hi *HttpInterface)addJob(job *jobInfo){
    hi.mapLock.Lock()
    defer hi.mapLock.Unlock()

    hi.activeJobs[job.JobId]=job
}

func (hi *HttpInterface)getJob(jobId string)(*jobInfo, bool){
    hi.mapLock.Lock()
    defer hi.mapLock.Unlock()

    job, exists := hi.activeJobs[jobId]

    return job, exists
}

func (hi *HttpInterface)removeJob(jobId string){
    hi.mapLock.Lock()
    defer hi.mapLock.Unlock()

    delete(hi.activeJobs, jobId)
}

func (hi *HttpInterface)getAllJobsStatus()([]*ImageBuildStatus){
    hi.mapLock.Lock()
    defer hi.mapLock.Unlock()

    allJobs := make([]*ImageBuildStatus, 0, len(hi.activeJobs))

    for _, job := range hi.activeJobs {
        allJobs=append(allJobs, &ImageBuildStatus{Id:job.JobId, Exists:true, Finished: job.finished})
    }

    return allJobs
}