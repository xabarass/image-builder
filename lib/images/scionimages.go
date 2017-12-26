package images;

import (
    "path"
    "strconv"
)

const (
    ImgFile=iota
    Root=iota
    Etc=iota
    Home=iota

    NodeCount=iota
)

type ScionImage struct{
    id int64
    Version string
    imgDir string
    nodes [NodeCount]string
    Used bool
    ImageName string
    ready bool
    mounted bool

    storage *ScionImageStorage
}

func (si *ScionImage)GetPathFor(what int64) string{
    return si.nodes[what];
}

func (si *ScionImage) IsReady()(bool){
    return si.ready;
}

func (si *ScionImage) GetId()(int64){
    return si.id;
}

func (si *ScionImage) IsMounted()(bool){
    return si.mounted;
}

func (si *ScionImage)Ready(version string)(error){
    stmt, err := si.storage.db.Prepare("UPDATE scion_images SET version=?, ready=1")
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(version)
    if err != nil {
        return err
    }

    si.Version=version
    si.Used=false
    si.ready=true

    return nil
}

func (si *ScionImage)SetMounted(mounted bool){
    si.mounted=mounted
}

func (si *ScionImage)initializePaths(){
    si.mounted=false
    si.nodes[ImgFile]=path.Join(si.imgDir, strconv.FormatInt(si.id, 10), "scion.img")
    si.nodes[Root]=path.Join(si.imgDir, strconv.FormatInt(si.id, 10), "mnt_root")
    si.nodes[Etc]=path.Join(si.imgDir, strconv.FormatInt(si.id, 10), "mnt_etc")
    si.nodes[Home]=path.Join(si.imgDir, strconv.FormatInt(si.id, 10), "mnt_home")
}