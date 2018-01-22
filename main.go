package main;

import (
    "log"
    "os"
    "os/signal"

    "github.com/xabarass/image-builder/lib/imagemanager"
    "github.com/xabarass/image-builder/lib/utils"
)

func signalHandler(stopApplication chan<- bool){
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt)

    <-stop

    log.Println("Got Ctrl + C, stopping!")

    stopApplication <- true
}

func main() {
    utils.InitializeRandomSeed()

    imgManager, err:=imagemanager.Create("config.json")
    if(err!=nil){
        log.Panic(err.Error())
    }

    imgMgrStop:=make(chan bool, 1)
    go signalHandler(imgMgrStop)

    err=imgManager.Run(imgMgrStop)
    if(err!=nil){
        log.Panic(err.Error())
    }
}

