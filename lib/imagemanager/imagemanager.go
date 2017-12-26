package imagemanager;

import(
    "log"    
    "os"
    "path"
    "encoding/json"
    "fmt"
    "os/exec"

    "github.com/xabarass/image-builder/lib/images"
    "github.com/xabarass/image-builder/lib/scionimagebuilder"
    "github.com/xabarass/image-builder/lib/httpinterface"
    "github.com/xabarass/image-builder/lib/imagecustomizer"
)

type ImageManager struct{
    db *images.ScionImageStorage
    images []images.OriginalImage
    imageBuilder *scionimagebuilder.ScionImageBuilder
    imageCustomizer *imagecustomizer.ImageCustomizer
    httpInterface *httpinterface.HttpInterface

    outputDir string

    mountScript string
    umountScript string
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

        mountScript: config.MountScript,
        umountScript: config.UmountScript,
    }

    log.Println("Loading image information")
    for i:=0; i<len(imgMgr.images);i++{
        img:=&imgMgr.images[i]
        log.Printf("\t%s", img.Name)
        img.ScionImages, err=imgStore.LoadReadyScionImages(img.Name)
        log.Printf("Loaded %d scion images for %s ", len(img.ScionImages), img.Name)
    }

    //TODO: Remove failed scion images from directory

    log.Printf("Creating SCION image builder from configuration %s", config.BuildConfigurationPath)
    imgMgr.imageBuilder, err=scionimagebuilder.Create(config.BuildConfigurationPath)
    if err != nil {
        return nil, err
    }

    imgMgr.httpInterface = httpinterface.CreateHttpServer(":8080", &imgMgr, map[string]bool{"milan":true}, "/tmp/downloadable_images")

    imgMgr.imageCustomizer=imagecustomizer.Create(imgMgr.httpInterface)

    log.Println("Created Image Manager!")
    return &imgMgr, err
}

func (im *ImageManager) Run(stop <-chan bool)(error){
    log.Println("Starting image manager...")

    // Create image builder
    imgBuildStop:=make(chan bool, 1)
    go im.imageBuilder.Run(imgBuildStop)

    //TODO: Create image customizer

    im.httpInterface.StartServer()
    readyImages:=make(chan *images.ScionImage, len(im.images))

    for _, img := range im.images{
        log.Printf("Starting to prepare image %s with %d scion images\n",img.Name, len(img.ScionImages))
        go im.prepareImage(&img, readyImages)
    }

    im.imageCustomizer.Run()
    
    LOOP: for{
        log.Printf("Starting to wait for requests")
        select{
           
        case readyImage:= <-readyImages:
            log.Printf("Received notification that image: %d is ready", readyImage.GetId())
            readyImage.SetMounted(true)
            readyImage.Used=false
            im.imageCustomizer.ScionImageReady(readyImage.ImageName, readyImage)
        case <-stop:
            log.Println("ImageManager >> Got request to shutdown!");
            imgBuildStop<-true
            im.httpInterface.StopServer()
            im.imageCustomizer.Stop()
            break LOOP;
        } 
    }

    log.Printf("Finishing run thread!")
    return nil
}

func (im *ImageManager)prepareImage(img *images.OriginalImage, readyImage chan<- *images.ScionImage){
    log.Printf("Starting to prepare image: %s", img.Name)

    var scimg *images.ScionImage;
    var destDirectory string

    if(len(img.ScionImages)==0){
        log.Panic("This part is not implemented fully!")    //TODO: Implement this, there is a problem need to investigate

        log.Printf("There are no generated SCION images for %s, starting build!", img.Name)

        scimg, err:=im.db.CreateScionImage(img.Name, im.outputDir)
        if(err!=nil){
            return
        }    

        destDirectory=path.Join(im.outputDir, fmt.Sprintf("img-%d",scimg.GetId()))
        log.Printf("Creating output directory at %s \n",destDirectory)

        err=os.MkdirAll(destDirectory, 0777)    //TODO: Fix permissions
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
    }

    err:=im.mountScionImage(scimg)
    if(err!=nil){
        log.Printf("Error mounting scion image! %s", err.Error())
    }

    readyImage<-scimg;
}

func (im *ImageManager)mountScionImage(scimg *images.ScionImage) (error){
    log.Printf("Mounting image %s", scimg.GetPathFor(images.ImgFile))
    cmd := exec.Command("sudo", im.mountScript, scimg.GetPathFor(images.ImgFile), 
        scimg.GetPathFor(images.Root), scimg.GetPathFor(images.Home), scimg.GetPathFor(images.Etc), "milan")   //TODO: Replace milan with env variable

    cmd.Run()

    return nil
}