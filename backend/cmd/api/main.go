package main

import (
	"context"
	"log"
	"os"
	"time"
    "strconv"

	"github.com/joho/godotenv"

	"github.com/SHILOP0P/Yardly/backend/internal/app/httpserver"
	bookingpg "github.com/SHILOP0P/Yardly/backend/internal/booking/pgrepo"
    itempg "github.com/SHILOP0P/Yardly/backend/internal/item/pgrepo"
    userpg "github.com/SHILOP0P/Yardly/backend/internal/user/pgrepo"
	"github.com/SHILOP0P/Yardly/backend/internal/db"
    "github.com/SHILOP0P/Yardly/backend/internal/auth"
)

func main() {
    _ = godotenv.Load()

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
        log.Println("APP_PORT not set, using default 8080")
    }

    ttlMinutesStr := os.Getenv("JWT_TTL_MINUTES")
    ttlMinute := 60
    if ttlMinutesStr !=""{
        if v, err := strconv.Atoi(ttlMinutesStr); err ==nil && v>0{
            ttlMinute = v
        }
    }
    jwtTTL := time.Duration(ttlMinute) * time.Minute

    jwtSvc := auth.NewJWT(
    os.Getenv("JWT_SECRET"),
    jwtTTL,
    )


    databaseURL := os.Getenv("DATABASE_URL")
    if databaseURL == "" {
        log.Fatal("DATABASE_URL isn’t set")
    }

    pool, err := db.NewPool(databaseURL)
    if err != nil {
        log.Fatalf("db connection failed: %v", err)
    }
    defer pool.Close()

    // 🔹 СОЗДАЁМ REPO ЗДЕСЬ (composition root)
    bookingRepo := bookingpg.New(pool)
    itemRepo:= itempg.New(pool)
    userRepo :=userpg.New(pool)


    srv := httpserver.New(port, pool, itemRepo, bookingRepo, userRepo, jwtSvc)

    jobCtx, jobCancel := context.WithCancel(context.Background())

    srv.RegisterOnShutdown(func(){
        jobCancel()
    })

    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()

        runOnce:=func(){
            ctx, cancel := context.WithTimeout(jobCtx, 5*time.Second)
            defer cancel()

            n, err:= bookingRepo.ExpireOverdueHandovers(ctx, time.Now().UTC())
            if err!= nil{
                 log.Println("expire overdue handovers error:", err)
                return
            }
            if n > 0{
                log.Println("expired overdue handovers:", n)
            }
        }

        runOnce()

        for range ticker.C {
           select{
           case<-jobCtx.Done():
            log.Println("expire job stopped")
            return
           case <- ticker.C:
            runOnce()
           }
        }
    }()

    log.Printf("Starting HTTP server on :%s\n", port)
    if err := srv.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}
