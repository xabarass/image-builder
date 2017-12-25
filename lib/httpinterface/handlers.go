package httpinterface

import (
    "net/http"
    "log"
    "github.com/gorilla/mux"
    "encoding/json"
    "os"
    "path"
    "io"
    "strings"

    "github.com/xabarass/image-builder/lib/utils"
)

func TokenAuth(authTokens map[string]bool, h func(http.ResponseWriter, *http.Request)) (func(http.ResponseWriter, *http.Request)) {
    return func(w http.ResponseWriter, r *http.Request) {
        token:=r.Header.Get("Auth")
        if _, ok := authTokens[token]; ok {
            log.Printf("Handling request")
            h(w, r)
        }else{
            log.Printf("Request not authenticated")
            http.Error(w, "Not authenticated", 403)
        }
    }
}

type HttpInterface struct {
    imgMgr ImageBuilderService
    validTokens map[string]bool
    rootDir string
}

func (i *HttpInterface)IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Welcome to SCION Lab image builder service!\n"))
}

func (i *HttpInterface)CreateImageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    imageName:=vars["image_name"]

    log.Printf("Create image request %s", imageName)

    jobId:=utils.GenerateRandomString(64);

    log.Printf("Creating destination dir")
    destDir:=path.Join(i.rootDir, jobId)
    os.MkdirAll(destDir, os.ModePerm);
    confFileDest:=path.Join(destDir,"config.zip")

    file, header, err := r.FormFile("config_file")
    if err != nil {
        panic(err)
    }
    defer file.Close()
    name := strings.Split(header.Filename, ".")
    log.Printf("Uploading file name %s\n", name[0])
    
    dest, err := os.OpenFile(confFileDest, os.O_WRONLY|os.O_CREATE, 0666)
    defer dest.Close()
    io.Copy(dest, file)

    err=i.imgMgr.CreateBuildJob(jobId, confFileDest, imageName) 
    if(err!=nil){
        log.Printf(err.Error())
        w.Write([]byte(err.Error()))
    }else{
        w.Write([]byte("Done!\n"))       
    }
}

func (i *HttpInterface)GetImages(w http.ResponseWriter, r *http.Request) {
    availableImages := i.imgMgr.GetAvailableImages()

    json.NewEncoder(w).Encode(availableImages)
}

func createHandler(imgMgr ImageBuilderService, authorizedTokens map[string]bool, rootDir string)(*mux.Router){
    hi:=&HttpInterface{
        imgMgr:imgMgr,
        rootDir:rootDir,
    }

    r := mux.NewRouter()
    r.HandleFunc("/", hi.IndexHandler)
    r.HandleFunc("/create/{image_name}", TokenAuth(authorizedTokens, hi.CreateImageHandler)).Methods("POST")
    r.HandleFunc("/get-images/", TokenAuth(authorizedTokens, hi.GetImages)).Methods("GET")

    return r
}

