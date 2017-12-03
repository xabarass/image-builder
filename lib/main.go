package main;

import (
    "fmt"

    "github.com/xabarass/image-builder/lib/images"
)

func main() {
    config, err:=images.LoadConfiguration("imgconfig.json")
    if(err!=nil){
        fmt.Println(err.Error())
    }
    defer config.ImageStore.Close()


    fmt.Println("Creating new image entry")
    newImage, err:= config.ImageStore.CreateScionImage("milan")
    if(err!=nil){
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("Created SCION image version %d \n", newImage.GetId())

    err = newImage.Ready("1.0","path1", "etcpath", "homepath");
    if(err!=nil){
        fmt.Println(err.Error())
    }    


    fmt.Println("Done!")
}