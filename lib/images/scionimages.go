package images;

import (
)

const (
    Etc=iota
    Home=iota

    EndpointCount=iota
)

type ScionImage struct{
    id int64
    Version string
    Path string
    MountPoints [EndpointCount]string
    Mounted bool
    Used bool
    ready bool

    storage *ScionImageStorage
}

func (si *ScionImage) IsReady()(bool){
    return si.ready;
}

func (si *ScionImage) GetId()(int64){
    return si.id;
}

func (si *ScionImage)Ready(version, path string, etcPath, homePath string)(error){
    stmt, err := si.storage.db.Prepare("UPDATE scion_images SET version=?, path=?, etcPath=?, homePath=?, mounted=1, ready=1")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(version, path, etcPath, homePath)
    if err != nil {
        return err
    }

    si.Version=version
    si.Path=path
    si.MountPoints[Etc]=etcPath
    si.MountPoints[Home]=homePath
    si.Used=false
    si.ready=true;

    return nil
}

