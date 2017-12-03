package scionimagebuilder;

import(
    "encoding/json"
    "fmt"
    "os"
)

type ImageBuildResult struct{
    Success bool
    FromImage string
    OutputImage string

    // If there was an error
    Error error
}

type buildJob struct {
    id int
    inputImage string
    outputImage string

    result chan<- ImageBuildResult
}

type ScionImageBuilder struct{
    ResizeScript string `json:"img_resize_script"`
    BuildScript string `json:"img_build_script"`

    BuildLogDirectory string `json:"log_dir"`

    buildJobs chan *buildJob
}

func Create(configFilePath string)(*ScionImageBuilder, error){
    configFile, err := os.Open(configFilePath)
    defer configFile.Close()
    if err != nil {
        return nil, err
    }

    // Load configuration from file
    var ib ScionImageBuilder
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&ib)

    fmt.Printf("Creating SCION image builder [resize: %s, build: %s] \n", ib.ResizeScript, ib.BuildScript);

    ib.buildJobs=make(chan *buildJob, 5) // I guess 5 jobs in queue are enough

    return &ib, nil
}