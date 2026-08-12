package main

import (
	"fmt"
	"log"
	"net/http"

	"domainhub/internal/config"
	"domainhub/internal/database"
	"domainhub/internal/handlers"
	"domainhub/internal/middleware"
	"domainhub/internal/repository"
	"domainhub/internal/routes"
	"domainhub/internal/service"

	"github.com/go-chi/chi/v5"
)

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

	fmt.Printf("Server running on http://localhost:%s\n", cfg.Port)

	err = http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal("Server Error:", err)
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
