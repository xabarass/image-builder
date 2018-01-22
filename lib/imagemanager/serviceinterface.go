package imagemanager;

import(
    "log"
    "fmt"
    "io/ioutil"
    "path"
    "os"

    "github.com/xabarass/image-builder/lib/httpinterface"

    "github.com/mholt/archiver"
)

func (im *ImageManager)GetAvailableImages()([]httpinterface.AvailableImage){
    output:=make(chan []httpinterface.AvailableImage, 0)

    im.syncFunctions<-func(){
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

        output<-result
    }

    return <-output
}

func (im *ImageManager)RunJob(job httpinterface.JobInfo)(error){
    log.Printf("Starting build job for: %s at: %s", job.ImageName, job.ConfigFile)

    err := archiver.TarGz.Open(job.ConfigFile, job.DestDir)
    if(err!=nil){
        return err
    }

    log.Printf("Decompress finished")

    files, err := ioutil.ReadDir(job.DestDir)
    if err != nil {
        return err
    }
    
    for _, file := range files {
        fmt.Println(file.Name())
    }

    var userDirecory string
    for _, file := range files {
        filePath:=path.Join(job.DestDir, file.Name())
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
        job.ConfigFile=userDirecory     //FIXME: Implement better way
        
        // im.imageCustomizer.AddJob(job)
    }else{
        return fmt.Errorf("Missing gen directory")
    }

    return nil

}
