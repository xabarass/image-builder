package imagemanager;

import(
    "database/sql"

    "github.com/xabarass/image-builder/lib/images"
)

type Storage struct {
    DB *database.DB
}

func Create(dbPath string)(*Storage, error){
    //TODO: Create db and initialize tables
}

func (s *Storage)loadScionImages(originalImageName string)([]ScionImage, error){
    //TODO: FInd images and return them
}

func (s *Storage)saveScionImage(scionImage images.ScionImage)(){
    
}