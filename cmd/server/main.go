package main

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/cmd/server/handler"
	"github.com/AleHts29/meli-challenge/internal/config"
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/AleHts29/meli-challenge/pkg/api"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Cargar el .env
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// creacion de instancias
	apiCountries := api.NewCountries(cfg.APIKey, cfg.APIUrl)
	repository := ipinfo.NewRepository(apiCountries)
	service := ipinfo.NewService(repository)

	newHandler := handler.NewHandler(service)

	router := gin.Default()
	//router.SetTrustedProxies([]string{"192.168.1.0/24"})// Si se usan proxies, especifica los rangos de IP confiables

	countries := router.Group("/countries")
	{
		countries.GET("", newHandler.GetCountries())
		//countries.GET("/:countryId", contriesHandler.GetCountryById())
	}

	// Iniciar el servidor
	log.Printf("Servidor escuchando en el puerto %s...\n", cfg.ServerPort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}

}
