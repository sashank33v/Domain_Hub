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

	"domainhub/internal/config"
	"domainhub/internal/database"
	"domainhub/internal/handlers"
	"domainhub/internal/middleware"
	"domainhub/internal/repository"
	"domainhub/internal/routes"
	"domainhub/internal/service"

	"github.com/go-chi/chi/v5"
)

// @title DomainHub API
// @version 1.0
// @description Domain management REST API.
// @host localhost:8000
// @BasePath /
func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	fmt.Printf("Starting %s...\n", cfg.AppName)

	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}
	defer db.Close()

	fmt.Println("Database Connected Successfully!")

	domainRepo := repository.NewDomainRepository(db)
	domainService := service.NewDomainService(domainRepo)
	domainHandler := handlers.NewDomainHandler(domainService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.CORS)
	routes.RegisterRoutes(r, domainHandler)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	fmt.Printf("Server running on http://localhost:%s\n", cfg.Port)

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatal("Server Error:", err)
		}

	case <-stop:
		fmt.Println("shutting down server....")

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			log.Fatal("Server shutdown Error:", err)
		}

		fmt.Println("server stopped successfully")
	}
}

// package main

// import (
// 	"domainhub/internal/config"
// 	"domainhub/internal/database"
// 	"domainhub/internal/routes"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	"github.com/go-chi/chi/v5"
// )

// func main(){
// 	cfg, err := config.LoadConfig()
// 	if err != nil{
// 		log.Fatal(err)
// 	}
// 	fmt.Println(cfg.AppName)

// 	db, err := database.ConnectDB(cfg)
// 	if err != nil{
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	fmt.Println(db)

// 	r := chi.NewRouter()

// 	routes.RegisterRoutes(r)

// 	err = http.ListenAndServe(":"+cfg.Port,r)
// 	if err!= nil{
// 		log.Fatal(err)
// 	}

// }
