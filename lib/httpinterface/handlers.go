package httpinterface

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "encoding/json"
)

func sendError(w http.ResponseWriter, message string, errorCode int32){
    resp:=errorResponse{
        Message:message,
        ErrorCode:errorCode,
    }

    json.NewEncoder(w).Encode(resp)
}

func TokenAuth(authTokens map[string]bool, h func(http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request)) {
    return func(w http.ResponseWriter, r *http.Request) {
        r.ParseMultipartForm(0)
        log.Println(r.Form)
        token := r.PostFormValue("token")

        if _, ok := authTokens[token]; ok {
            log.Printf("Handling request")
            h(w, r)
        }else{
            log.Printf("Request not authenticated")
            sendError(w, "Not authenticated", 403)
        }
    }
}

func (i *HttpInterface)IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to SCION Lab image builder service!\n"))
}

func (i *HttpInterface)CreateImageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    imageName:=vars["image_name"]
    log.Printf("Create image request for %s", imageName)

    file, header, err := r.FormFile("config_file")
    if err != nil {
        sendError(w, err.Error(), 400)
        return
    }
    defer file.Close()

    log.Printf("Uploading file name %s\n", header.Filename)
    
    jobId, err:= i.createNewJob(imageName, file)
    if(err!=nil){
        sendError(w, err.Error(), 400)  //TODO: Fix error codes
        return
    }

    resp:=createImageResponse{
        JobId:jobId,
        Image:imageName,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (i *HttpInterface)DownloadImage(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    jobId:=vars["job_id"]
    log.Printf("Trying to download image from job %s", jobId)

    finished, imageFile, err:=i.getImageForJob(jobId)
    if(err!=nil){
        sendError(w, err.Error(), 400)
        return
    }

    log.Printf("Job id is valid, checking if job is finished")

    // TODO: Fix this mess, make normal responses
    if(finished){
        log.Printf("Sending file %s ", imageFile)
        w.Header().Set("Content-Disposition", "attachment; filename=scion.img.bz2")
        http.ServeFile(w, r, imageFile)
    }else{
        sendError(w, "File not ready", 404)
    }
}

func (i *HttpInterface)GetImageStatus(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    jobId:=vars["job_id"]
    log.Printf("Getting status for image %s", jobId)

    finished, _, err:=i.getImageForJob(jobId)
    status := ImageBuildStatus{Id:jobId, Finished:finished}
    if(err!=nil){
        status.Exists=false
    }else{
        status.Exists=true
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(status)
}

func (i *HttpInterface)GetAllImagesStatus(w http.ResponseWriter, r *http.Request) {
    log.Printf("Get all images")

    allJobs:=i.getAllJobsStatus()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(allJobs)
}

func (i *HttpInterface)GetImages(w http.ResponseWriter, r *http.Request) {
    availableImages := i.imgMgr.GetAvailableImages()

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(availableImages)
}

func createHandler(hi *HttpInterface, authorizedTokens map[string]bool)(*mux.Router){
    r := mux.NewRouter()

    r.HandleFunc("/create/{image_name}", TokenAuth(authorizedTokens, hi.CreateImageHandler)).Methods("POST")
    r.HandleFunc("/get-images", hi.GetImages).Methods("GET")
    r.HandleFunc("/download/{job_id}", hi.DownloadImage).Methods("GET")
    r.HandleFunc("/status/{job_id}", hi.GetImageStatus).Methods("GET")
    r.HandleFunc("/status", hi.GetAllImagesStatus).Methods("GET")

    r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./web/"))))

    return r
}

