package main

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/cmd/server/handler"
	"github.com/AleHts29/meli-challenge/internal/config"
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/AleHts29/meli-challenge/pkg/api"
	"github.com/AleHts29/meli-challenge/pkg/store"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Cargar el .env
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	ipStore, err := store.NewIpStore("./IP-COUNTRY-DB26.BIN")
	if err != nil {
		panic(err)
		return
	}

	// creacion de instancias
	apiCountries := api.NewCountries(cfg.APIKey, cfg.APIUrl)
	apiCurrencies := api.NewCurrencies(cfg.APIKey, cfg.APIUrl)

	repository := ipinfo.NewRepository(apiCountries, apiCurrencies, ipStore)
	service := ipinfo.NewService(repository)

	// Test para IP
	results, err := service.GetCountryByIP("45.5.164.0")
	if err != nil {
		return
	}
	fmt.Printf("country_short: %s\n", results.CountryCode)
	fmt.Printf("country_long: %s\n", results.CountryName)

	newHandler := handler.NewHandler(service)

	router := gin.Default()
	//router.SetTrustedProxies([]string{"192.168.1.0/24"})// Si se usan proxies, especifica los rangos de IP confiables

	ip := router.Group("/api/ip")
	{
		ip.GET("/:ip", newHandler.GetCountryByIP())
		ip.POST("/block", newHandler.BlockIPs())
	}

	// Iniciar el servidor
	log.Printf("Servidor escuchando en el puerto %s...\n", cfg.ServerPort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
