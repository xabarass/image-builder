package main;

import(
    "fmt"
    "time"

    "github.com/xabarass/image-builder/lib/scionimagebuilder"
)


func main(){
    fmt.Println("Starting image building service!")

    imageBuilder, err:= scionimagebuilder.Create("imgbuilder.json");
    if err!=nil{
        fmt.Println(err.Error())
    }

    stop:=make(chan bool, 1)
    imageBuilder.Run(stop)

    result:=make(chan scionimagebuilder.ImageBuildResult, 1)
    imageBuilder.StartBuildJob("/home/milan/Downloads/original_images/ubuntu-16.04.3-preinstalled-server-armhf+raspi2.img", "/tmp/raspberrypi2.img", result)

    output:=<-result
    if(output.Success){
        fmt.Println("Exiting with success!")    
    }else{
        fmt.Println(output.Error.Error())
    }

    time.Sleep(500 * time.Millisecond)
    stop<-true
}