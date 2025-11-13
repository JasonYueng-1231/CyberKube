package main

import (
    "log"
    "net/http"
    "os"

    "github.com/JasonYueng-1231/CyberKube/backend/internal/api"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/database"
    "github.com/JasonYueng-1231/CyberKube/backend/internal/model"
)

func main() {
    if err := database.InitMySQL(); err != nil {
        log.Fatalf("mysql init error: %v", err)
    }
    if err := database.AutoMigrate(
        &model.User{},
        &model.Cluster{},
    ); err != nil {
        log.Fatalf("auto migrate error: %v", err)
    }

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
