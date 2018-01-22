package httpinterface

type AvailableImage struct {
    Name string             `json:"name"`
    DisplayName string      `json:"display_name"`
    Description string      `json:"description"`
    Version string          `json:"version"`
}

type ImageBuilderService interface {
    GetAvailableImages()([]AvailableImage)
    RunJob(job JobInfo)(error)
}
