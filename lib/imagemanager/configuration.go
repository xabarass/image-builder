package imagemanager;

import(
    "log"
    "database/sql"
    "fmt"
    "log"
    "os"

    "github.com/xabarass/image-builder/lib/images"

    _ "github.com/mattn/go-sqlite3"
)

type AvailableImages struct {
    Images []images.OriginalImage `json:"images"`
    Count int `json:"img_count"`
}

func loadConfigFile(configFilePath string)(*AvailableImages, error){

    configFile, err := os.Open(configFilePath)
    defer configFile.Close()
    if err != nil {
        return nil, err
    }

    // Load configuration from file
    var imgs AvailableImages
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&imgs)

    log.Println("Loaded list of available images");

    return &imgs, nil
}

func LoadConfiguration(configFilePath string, databasePath string){
    
}