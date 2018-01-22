package httpinterface

import(
    "time"
    "log"
    "os"
)

// Image is valid 10 minutes
const IMAGE_VALIDITY_MINUTES float64 = 10

func isExpired(job *jobInfo)(bool){
    if(job.finished){
        duration := time.Since(job.timestamp)
        if(duration.Minutes()>IMAGE_VALIDITY_MINUTES){
            return true;
        }
    }

    return false
}

// TODO: Make thread safe!
// Periodically checks for old jobs and removes them
func (hi *HttpInterface)startCleanupService(){
    go func(){
        LOOP: for{
            select{
            
            case <-time.After(time.Minute):
                itemsToDelete := make([]string, 0)
                for k, j := range hi.activeJobs { 
                    log.Printf("Running cleanup")
                    if(isExpired(j)){
                        log.Printf("Job %s scheduled for cleanup", k)
                        itemsToDelete=append(itemsToDelete, k)
                    }
                }

                for _,k:=range itemsToDelete{
                    if(os.RemoveAll(hi.activeJobs[k].DestDir)==nil){
                        delete(hi.activeJobs, k)    
                    }
                }

            case <-hi.stop:
                log.Println("CleanupService >> Got request to shutdown!");
                break LOOP;
            } 
        }
    }()
}

func (hi *HttpInterface)stopCleanupService(){
    hi.stop<-true
}