package httpinterface

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "encoding/json"

)

type createImageResponse struct {
    JobId string  `json:"id"`
    Image string  `json:"image"`
}

type errorResponse struct {
    Message string  `json:"message"`
    ErrorCode int32   `json:"err_code"`
}

func sendError(w http.ResponseWriter, message string, errorCode int32){
    resp:=errorResponse{
        Message:message,
        ErrorCode:errorCode,
    }

    json.NewEncoder(w).Encode(resp)
}

func TokenAuth(authTokens map[string]bool, h func(http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request)) {
    return func(w http.ResponseWriter, r *http.Request) {
        token:=r.Header.Get("Auth")
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
        //TODO: Handle error!
        panic(err)

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

    // TODO: Fix this mess, make normal responses
    if(finished){
        http.ServeFile(w, r, imageFile)
    }else{
        sendError(w, "File not ready", 404)
    }
}

func (i *HttpInterface)GetImages(w http.ResponseWriter, r *http.Request) {
    availableImages := i.imgMgr.GetAvailableImages()

    json.NewEncoder(w).Encode(availableImages)
}

func createHandler(hi *HttpInterface, authorizedTokens map[string]bool)(*mux.Router){
    r := mux.NewRouter()
    r.HandleFunc("/", hi.IndexHandler)
    r.HandleFunc("/create/{image_name}", TokenAuth(authorizedTokens, hi.CreateImageHandler)).Methods("POST")
    r.HandleFunc("/get-images", TokenAuth(authorizedTokens, hi.GetImages)).Methods("GET")
    r.HandleFunc("/download/{job_id}", hi.DownloadImage).Methods("GET")

    return r
}

