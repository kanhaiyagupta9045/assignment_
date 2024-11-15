package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kanhaiyagupta9045/car_management/databases"
	"github.com/kanhaiyagupta9045/car_management/middleware"
	"github.com/kanhaiyagupta9045/car_management/routes"
)

func init() {

	databases.DatabaseConnect()
}

func main() {
	router := gin.Default()
	middleware.Middleware(router)
	routes.UserRoutes(router)
	routes.CarRoutes(router)
	srv := http.Server{
		Handler: router.Handler(),
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	defer databases.DB.Close()
	select {
	case <-ctx.Done():
		log.Println("timeout of 2 seconds.")
	}
	log.Println("Server exiting")
}
