package images;

import (
    "database/sql"
    "log"

    _ "github.com/mattn/go-sqlite3"
)

type ScionImageStorage struct {
    db *sql.DB
    dbPath string

    preparedInsert *sql.Stmt
}

func Open(dbPath string)(*ScionImageStorage, error){
    log.Println("Opening database file...")
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }

    createStmt:=`
        CREATE TABLE IF NOT EXISTS scion_images (id INTEGER PRIMARY KEY,
                                                 version TEXT DEFAULT '', 
                                                 imgDir TEXT DEFAULT '',
                                                 fromImage TEXT DEFAULT '',
                                                 ready INTEGER DEFAULT 0
                                                );
    `

    _, err = db.Exec(createStmt)
    if err != nil {
        return nil, err
    }

    log.Println("Creating prepared statements...")
    inserStmt:="INSERT INTO scion_images (fromImage, imgDir) VALUES (?, ?)"
    preparedInsert, err:=db.Prepare(inserStmt)
    if(err != nil){
        return nil, err   
    }

    storage := &ScionImageStorage{db:db, dbPath:dbPath, preparedInsert:preparedInsert}

    return storage, nil
}

func (st *ScionImageStorage)Close(){
    st.preparedInsert.Close()
    st.db.Close()
}

func (st *ScionImageStorage)LoadReadyScionImages(name string)([]*ScionImage, error){
    log.Printf("Loading ready SCION images for %s", name)
    rows, err := st.db.Query("SELECT id, version, imgDir, fromImage FROM scion_images WHERE fromImage=? AND ready=1", name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    resultImages:=make([]*ScionImage,0);

    for rows.Next() {
        var si ScionImage

        err = rows.Scan(&si.id, &si.Version, &si.imgDir, &si.ImageName)
        if err != nil {
            return nil, err
        }
        si.storage=st
        si.initializePaths()

        resultImages=append(resultImages, &si)
        log.Printf("We have new image with id %d", si.id)
    }

    return resultImages, nil
}

func (st *ScionImageStorage)CreateScionImage(fromImage, imgDir string)(*ScionImage, error){
    result, err := st.preparedInsert.Exec(fromImage, imgDir)
    if err != nil {
        log.Printf("Error creating new scion image record")
        return nil, err
    }

    lastID, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }
    si:= &ScionImage{
        id:lastID,
        Used:false,
        storage:st,
        imgDir:imgDir,
        ready:false,
        ImageName:fromImage,
    }
    si.initializePaths()

    return si, nil
}

