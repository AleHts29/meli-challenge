package main

import (
	"fmt"
	"github.com/AleHts29/meli-challenge/cmd/server/handler"
	"github.com/AleHts29/meli-challenge/internal/config"
	"github.com/AleHts29/meli-challenge/internal/ipinfo"
	"github.com/AleHts29/meli-challenge/pkg/api"
	"github.com/AleHts29/meli-challenge/pkg/store"
	"github.com/gin-contrib/cors"
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

	ipStore, err := store.NewIpStore(cfg.IPStorePath)
	if err != nil {
		panic(err)
		return
	}

	// creacion de instancias
	apiCountries := api.NewCountries(cfg.APIKey, cfg.APIUrl)
	apiCurrencies := api.NewCurrencies(cfg.APIKey, cfg.APIUrl)
	repository := ipinfo.NewRepository(apiCountries, apiCurrencies, ipStore)
	service := ipinfo.NewService(repository, cfg.BlockedIPsFilePath)
	newHandler := handler.NewHandler(service)

	router := gin.Default()

	// Habilitar CORS
	router.Use(cors.Default())

	// Configurar rutas y servidor
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	ip := router.Group("/api/ip")
	{
		ip.GET("/:ip", newHandler.GetCountryByIP())      // Obtener informacion de paises mediante una IP
		ip.POST("/block", newHandler.BlockIPs())         // Bloquear una o varias IPs
		ip.GET("/events", newHandler.NotifyBlockedIPs()) // Emitir eventos de bloqueo
	}

	// Iniciar el servidor
	log.Printf("Servidor escuchando en el puerto %s...\n", cfg.ServerPort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
