package imagemanager;

import(
    "log"
    "os"
    "encoding/json"
    "fmt"
    "os/exec"

    "github.com/xabarass/image-builder/lib/images"
    "github.com/xabarass/image-builder/lib/httpinterface"
    "github.com/xabarass/image-builder/lib/imagecustomizer"
)

type SyncFunc func()

type ImageManager struct{
    images map[string]*images.ScionImage
    imageCustomizer *imagecustomizer.ImageCustomizer
    httpInterface *httpinterface.HttpInterface

    readyImages chan *images.ScionImage
    syncFunctions chan SyncFunc
    buildJobs chan *httpinterface.JobInfo

    config Configuration
}

func Create(configFilePath string)(*ImageManager, error){
    configFile, err := os.Open(configFilePath)
    if err != nil {
        return nil, err
    }
    defer configFile.Close()

    var imgMgr ImageManager
    imgMgr.images=make(map[string]*images.ScionImage)
    imgMgr.readyImages=make(chan *images.ScionImage, 10)  //TODO: remove magic constant
    imgMgr.syncFunctions=make(chan SyncFunc, 10)
    
    // Load configuration from file
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&imgMgr.config)

    for _, img := range imgMgr.config.Images{
        if _, exists := imgMgr.images[img.Name]; exists {
            return nil, fmt.Errorf("Error! Image names are not unique!")
        }

        imgMgr.images[img.Name]=img
    }
   
    //TODO: Fix loading secret token
    // imgMgr.httpInterface = httpinterface.CreateHttpServer(imgMgr.config.BindAddress, &imgMgr, map[string]bool{"milan":true}, imgMgr.config.OutputDirectory)
    imgMgr.imageCustomizer=imagecustomizer.Create(imgMgr.config.CustomizeScript, &imgMgr)

    log.Println("Created Image Manager!")
    return &imgMgr, err
}

func (im *ImageManager) Run(stop <-chan bool)(error){
    log.Println("Starting image manager...")

    // im.httpInterface.StartServer()


    log.Printf("Mounting available images")
    for _, img := range im.images{
        im.mountScionImage(img)
        img.SetMounted(true)
    }

    im.imageCustomizer.Run()
    
    imgToIdMap := make(map[string]string)

    LOOP: for{
        log.Printf("Starting to wait for requests")
        select{
           
        case readyImage:= <-im.readyImages:
            delete(imgToIdMap, readyImage.Name)           
            readyImage.SetUsed(false)

        // Execute all functions from worker thread, avoiding mutexes
        case f := <-im.syncFunctions:
            f()

        case bj := <- im.buildJobs:


        case <-stop:
            log.Println("ImageManager >> Got request to shutdown!");
            im.httpInterface.StopServer()
            im.imageCustomizer.Stop()
            break LOOP;
        } 
    }

    log.Printf("Finishing run thread!")
    return nil
}

func (im *ImageManager)mountScionImage(scimg *images.ScionImage) (error){
    log.Printf("Mounting image %s for user: %s", scimg.File, os.Getenv("USER"))

    cmd := exec.Command("sudo", im.config.MountScript, scimg.File, 
        scimg.GetPathFor(images.Root), scimg.GetPathFor(images.Home), scimg.GetPathFor(images.Etc), os.Getenv("USER"))

    cmd.Run()

    return nil
}

func (im *ImageManager)OnCustomizeJobSuccess(image *images.ScionImage, generatedFile string){
    // TODO: Notify http module

    im.readyImages<-image
}

func (im *ImageManager)OnCustomizeJobError(image *images.ScionImage, err error){
    // TODO: Notify http module

    im.readyImages<-image
}