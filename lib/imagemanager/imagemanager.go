package imagemanager;

import(
    "log"    

    "github.com/xabarass/image-builder/lib/images"
)

type ImageManager struct{
    db *ScionImageStorage
    images []OriginalImage
}

func Create(configFilePath string)(*ImageManager, error){
    configFile, err := os.Open(configFilePath)
    if err != nil {
        return nil, err
    }
    defer configFile.Close()
    
    // Load configuration from file
    var config Configuration
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&config)


    log.Printf("Loaded configuration file, database path is: %s, loading image information", config.DBPath);
    imgStore, err:=images.Open(config.DBPath)
    if err != nil {
        return nil, err
    }

    imgMgr:=ImageManager{
        db:imgStore,
        images:config.Images,       
    }

    for _, img := range imgMgr.images{
        img.ScionImages, err=imgStore.LoadScionImages(img.Name)
    }

    log.Println("Created Image Manager!")
    return &imgMgr, err
}