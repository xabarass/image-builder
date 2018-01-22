package imagecustomizer

import(
    "log"
    "os"
    "os/exec"
    "path"
    "fmt"

    "github.com/xabarass/image-builder/lib/images"
)

func (ic *ImageCustomizer)customizeImage(image *images.ScionImage, configDirectory string, outputDir string) (string, error){
    log.Printf("Starting to customize image: %s", image.Name) 

    cmd:=exec.Command(ic.customizeScript, configDirectory, 
        image.GetPathFor(images.Home), image.GetPathFor(images.Etc), 
        image.File, outputDir)

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    if(cmd.Run()==nil){
        log.Printf("Success customizing image!")
        return path.Join(outputDir, "scion.img.bz2"), nil
    }else{
        log.Printf("There was an error customizing image!")
        return "", fmt.Errorf("Error while running customization script")
    }
}