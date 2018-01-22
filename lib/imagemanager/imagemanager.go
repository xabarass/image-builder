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

type ImageManager struct{
    images map[string]*images.ScionImage
    imageCustomizer *imagecustomizer.ImageCustomizer
    httpInterface *httpinterface.HttpInterface

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
    
    // Load configuration from file
    jsonParser := json.NewDecoder(configFile)
    jsonParser.Decode(&imgMgr.config)

    for _, img := range imgMgr.config.Images{
        if _, exists := imgMgr.images[img.Name]; exists {
            return nil, fmt.Errorf("Error! Image names are not unique!")
        }

        imgMgr.images[img.Name]=img
    }
   
    token := os.Getenv("ACCESS_TOKEN")
    if(token==""){
        return nil, fmt.Errorf("Env variable 'ACCESS_TOKEN' not specified!")
    }

    imgMgr.httpInterface = httpinterface.CreateHttpServer(imgMgr.config.BindAddress, &imgMgr, map[string]bool{token:true}, imgMgr.config.OutputDirectory)
    imgMgr.imageCustomizer=imagecustomizer.Create(imgMgr.config.CustomizeScript, &imgMgr)

    log.Println("Created Image Manager!")
    return &imgMgr, err
}

func (im *ImageManager) Run(stop <-chan bool)(error){
    log.Println("Starting image manager...")

    log.Printf("Mounting available images")
    for _, img := range im.images{
        im.mountScionImage(img)
        img.SetMounted(true)
    }

    im.httpInterface.StartServer()
    im.imageCustomizer.Run()
    
    <-stop
    // Wait for stop
    log.Println("ImageManager >> Got request to shutdown!");
    im.httpInterface.StopServer()
    im.imageCustomizer.Stop()

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

func (im *ImageManager)OnCustomizeJobSuccess(image *images.ScionImage, jobId string, generatedFile string){
    log.Printf("Customization JOB finished!")
    im.httpInterface.JobFinished(jobId, generatedFile)
}

func (im *ImageManager)OnCustomizeJobError(image *images.ScionImage, jobId string, err error){
    //TODO: Implement!
    log.Printf("Error customizing job")
}