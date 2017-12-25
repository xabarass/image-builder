package httpinterface

import (
    "net/http"
    "log"
)

func CreateHttpServer(bindAddress string, imgMgr ImageBuilderService, authorizedTokens map[string]bool, rootDir string) *http.Server {
    srv := &http.Server{
        Addr: bindAddress, 
        Handler: createHandler(imgMgr, authorizedTokens, rootDir),
    }

    go func() {
        if err := srv.ListenAndServe(); err != nil {
            log.Printf("Httpserver: ListenAndServe() error: %s", err)
        }
    }()

    return srv
}