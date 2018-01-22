package httpinterface

// Return message types

type createImageResponse struct {
    JobId string  `json:"id"`
    Image string  `json:"image"`
}

type errorResponse struct {
    Message string  `json:"message"`
    ErrorCode int32   `json:"err_code"`
}

type AvailableImage struct {
    Name string             `json:"name"`
    DisplayName string      `json:"display_name"`
    Description string      `json:"description"`
    Version string          `json:"version"`
}

// Callback interface

type ImageBuilderService interface {
    GetAvailableImages()([]AvailableImage)
    RunJob(imageName, configFile, destDir, jobId string)(error)
}
