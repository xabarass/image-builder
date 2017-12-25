package httpinterface

type AvailableImage struct {
    Device string   `json:"device"`
    Name string     `json:"name"`
}

type ImageBuilderService interface {
    GetAvailableImages()([]AvailableImage)
    CreateBuildJob(jobId string, configPath string, imageName string)(error)
}

