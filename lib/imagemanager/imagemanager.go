package imagemanager;

import(
    "log"    
    "os"
    "path"
    "encoding/json"
    "fmt"
    // "os/exec"

    "github.com/xabarass/image-builder/lib/images"
    "github.com/xabarass/image-builder/lib/scionimagebuilder"
    "github.com/xabarass/image-builder/lib/httpinterface"
)

type ImageManager struct{
    db *images.ScionImageStorage
    images []images.OriginalImage
    imageBuilder *scionimagebuilder.ScionImageBuilder

    outputDir string
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
        outputDir:config.BuildOutputDirectory, 
    }

    log.Println("Loading image information")
    for _, img := range imgMgr.images{
        log.Printf("\t%s", img.Name)
        img.ScionImages, err=imgStore.LoadScionImages(img.Name)
    }

    log.Printf("Creating SCION image builder from configuration %s", config.BuildConfigurationPath)
    imgMgr.imageBuilder, err=scionimagebuilder.Create(config.BuildConfigurationPath)
    if err != nil {
        return nil, err
    }

    log.Println("Created Image Manager!")
    return &imgMgr, err
}

func (im *ImageManager) Run(stop <-chan bool)(error){
    log.Println("Starting image manager...")

    // Create image builder
    imgBuildStop:=make(chan bool, 1)
    go im.imageBuilder.Run(imgBuildStop)

    //TODO: Create image customizer

    srv := httpinterface.CreateHttpServer(":8080", im, map[string]bool{"milan":true}, "/tmp/downloadable_images")
    // readyImages:=make(chan images.ScionImage, len(im.images))

    // for _, img := range im.images{
    //     go im.prepareImage(&img, readyImages)
    // }
    
    LOOP: for{
        log.Printf("Starting to wait for requests")
        select{
           
        case <-stop:
            log.Println("ImageManager >> Got request to shutdown!");
            imgBuildStop<-true
            srv.Shutdown(nil)
            break LOOP;
        } 
    }

    log.Printf("Finishing run thread!")
    return nil
}

func (im *ImageManager)prepareImage(img *images.OriginalImage, readyImage chan<- images.ScionImage){
    log.Printf("Starting to prepare image: %s", img.Name)

    var scimg images.ScionImage;
    var destDirectory string

    if(len(img.ScionImages)==0){
        log.Printf("There are no generated SCION images for %s, starting build!", img.Name)

        scimg, err:=im.db.CreateScionImage(img.Name)
        if(err!=nil){
            return
        }    

        destDirectory=path.Join(im.outputDir, fmt.Sprintf("img-%d",scimg.GetId()))
        log.Printf("Creating output directory at %s \n",destDirectory)

        err=os.MkdirAll(destDirectory, 0777)
        if(err!=nil){
            log.Panic(err.Error())
        }

        outputFile:=path.Join(destDirectory, img.Name+".img");

        resultChan:=make(chan scionimagebuilder.ImageBuildResult, 1)
        err=im.imageBuilder.StartBuildJob(img.Path, outputFile, resultChan)
        if err!=nil{
            log.Panic(err.Error())   
        }

        res:=<-resultChan

        if(!res.Success){
            log.Printf("Error building image: %s, %s", img.Name, res.Error.Error())
            //TODO: Send error message to main thread
            return
        }
    }else{
        log.Printf("There is stored record of the images %s, loading data", img.Name)

        //TODO: For now we assume there is only one SCION image, later we will extend
        scimg=img.ScionImages[0]
        destDirectory=path.Join(im.outputDir, fmt.Sprintf("img-%d",scimg.GetId()))
    }

    // etcDir:=path.Join(destDirectory, "etc");
    // homeDir:=path.Join(destDirectory, "home");

    // if(exec.Command("mountpoint", etcDir).Run()!=0){
    //     //We need to mount this
    //     log.Printf("etc is not mounted! Mounting it")
    // }

    // if(exec.Command("mountpoint", homeDir).Run()!=0){
    //     //We need to mount this
    //     log.Printf("home is not mounted! Mounting it")
    // }
    
    // err=scimg.Ready("1.0", outputFile, etcDir, homeDir)
    // if err!=nil{
    //     log.Printf("Error making scion image ready!")
    // }

    readyImage<-scimg;
}