package main

import (
    "log"
    "net/http"
    "os"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/api"
)

func main() {
    r := api.SetupRouter()

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: r,
    }

    log.Printf("CyberKube backend started on :%s", port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("server error: %v", err)
    }
}

