package main;

import (
    "log"

    "github.com/xabarass/image-builder/lib/imagemanager"
)

func main() {
    imgManager, err:=imagemanager.Create("imgconfig.json")
    if(err!=nil){
        log.Panic(err.Error())
    }

    
}