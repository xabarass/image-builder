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
                                                 path TEXT DEFAULT '',
                                                 etcPath TEXT DEFAULT '',
                                                 homePath TEXT DEFAULT '',
                                                 mounted INTEGER DEFAULT 0,
                                                 fromImage TEXT DEFAULT '',
                                                 ready INTEGER DEFAULT 0
                                                );
    `

    _, err = db.Exec(createStmt)
    if err != nil {
        return nil, err
    }

    log.Println("Creating prepared statements...")
    inserStmt:="INSERT INTO scion_images (fromImage) VALUES (?)"
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

func (st *ScionImageStorage)LoadScionImages(name string)([]ScionImage, error){
    rows, err := st.db.Query("SELECT id, version, path, etcPath, homePath, mounted FROM scion_images WHERE fromImage=? AND ready=1", name)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    resultImages:=make([]ScionImage, 10);

    for rows.Next() {
        var si ScionImage

        err = rows.Scan(&si.id, &si.Version, &si.Path, &si.MountPoints[Etc], &si.MountPoints[Home], &si.Mounted)
        if err != nil {
            return nil, err
        }

        resultImages=append(resultImages, si)
    }

    return resultImages, nil
}

func (st *ScionImageStorage)CreateScionImage(fromImage string)(*ScionImage, error){
    result, err := st.preparedInsert.Exec(fromImage)
    if err != nil {
        log.Printf("Error creating new scion image record")
        return nil, err
    }

    lastID, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &ScionImage{
        id:lastID,
        Mounted:false,
        Used:false,
        storage:st,
    }, nil
}

