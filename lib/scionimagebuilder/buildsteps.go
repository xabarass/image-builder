package scionimagebuilder;

import(
    "log"
    "os/exec"
    "io/ioutil"
    "io"
    "bufio"
    "strings"
)

func (job *buildJob)copyImage(ib *ScionImageBuilder)(error){
    log.Println("Copying file!")

    cmd := exec.Command("cp", job.inputImage, job.outputImage)
    return ib.logCommandOutput("copy", cmd, ioutil.Discard)
}

func (job *buildJob)resizeAndMount(ib *ScionImageBuilder)(string, error){
    log.Println("Resizing image!")

    cmd := exec.Command("sudo", ib.ResizeScript, job.outputImage)

    r, w := io.Pipe()
    reader:=bufio.NewReader(r)
    defer r.Close()
    defer w.Close()

    lastLine:=""
    // We want to take last line from stdout
    go func(){
        log.Println("Starting to read from pipe!")
        var err error
        for ;err==nil;{
            lastLine,err= reader.ReadString('\n')    
        }
        
        log.Println("Exiting pipe read gorutine")
    }()

    err:=ib.logCommandOutput("resize", cmd, w)

    lastLine = strings.TrimRight(lastLine, "+")
    lastLine = strings.TrimLeft(lastLine, "+")
    return lastLine, err
}

func (job *buildJob)installScion(ib *ScionImageBuilder, loopDevice string)(error){
    log.Println("Installing SCION!")

    cmd := exec.Command("sudo", ib.BuildScript, loopDevice)
    return ib.logCommandOutput("make", cmd, ioutil.Discard)
}

func (job *buildJob)closeLoopDevice(ib *ScionImageBuilder)(error){
    log.Println("Closing loop devices!")

    cmd := exec.Command("sudo", "/sbin/kpartx", "-d", "-v", job.outputImage)
    return ib.logCommandOutput("close_loop_dev", cmd, ioutil.Discard)
}
