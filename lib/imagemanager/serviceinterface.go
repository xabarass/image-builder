package imagemanager;

import(
    "log"
    "fmt"
    "io/ioutil"
    "path"
    "os"

    "github.com/xabarass/image-builder/lib/httpinterface"
    "github.com/xabarass/image-builder/lib/images"

    "github.com/mholt/archiver"
)

func (im *ImageManager)GetAvailableImages()([]httpinterface.AvailableImage){
    var result []httpinterface.AvailableImage

    for _, img := range im.images{
        if(img.IsMounted()){
            result=append(result, httpinterface.AvailableImage{ Name:img.Name, 
                                                                DisplayName:img.DisplayName, 
                                                                Description:img.Description, 
                                                                Version:img.Version,
                                                              })
        }
    }

    return result
}

func (im *ImageManager)RunJob(imageName, configFile, destDir, jobId string)(error){
    log.Printf("Starting build job for: %s at: %s", imageName, configFile)

    var scionImg *images.ScionImage
    if img, ok := im.images[imageName]; ok {
        scionImg=img
    }else{
        return fmt.Errorf("Unknown image name: %s", imageName)
    }

    err := archiver.TarGz.Open(configFile, destDir)
    if(err!=nil){
        return err
    }

    log.Printf("Decompress finished")

    files, err := ioutil.ReadDir(destDir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        fmt.Println(file.Name())
    }

    var userDirecory string
    for _, file := range files {
        filePath:=path.Join(destDir, file.Name())
        if info, _ := os.Stat(filePath); info.Mode().IsDir(){
            userDirecory=filePath
        }
    }
    if (userDirecory==""){
        return fmt.Errorf("Not a valid directory structure")
    }else{
        log.Printf("Users config directory: %s", userDirecory)
    }

    if info, _ := os.Stat(path.Join(userDirecory, "gen")); info.Mode().IsDir(){
        im.imageCustomizer.CustomizeImage(scionImg, userDirecory, destDir, jobId)
    }else{
        return fmt.Errorf("Missing gen directory")
    }

    return nil
}
