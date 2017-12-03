package images;

import(
    "log"
    "os"
    "encoding/json"
)

type Configuration struct {
    Images []OriginalImage `json:"images"`
    DBPath string   `json:"db_path"`

    ImageStore *ScionImageStorage
}

func LoadConfiguration(configFilePath string)(*Configuration, error){

    configFile, err := os.Open(configFilePath)
    defer configFile.Close()
    if err != nil {
        return nil, err
    }

    // Load configuration from file
    var config Configuration
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)

    log.Printf("Loaded list of available images, loading SCION image information from database: %s", config.DBPath);
    imgStore,err:=Open(config.DBPath)
    if err != nil {
        return nil, err
    }

    for _, img := range config.Images{
        img.ScionImages, err=imgStore.LoadScionImages(img.Name)
    }

    config.ImageStore=imgStore

    return &config, err
}

