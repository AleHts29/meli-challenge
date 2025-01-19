package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	APIKey     string
	ServerPort string
	APIUrl     string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file: " + err.Error())
	}

	// Obtener variables de entorno o carga de parametros por default
	config := &Config{
		APIKey:     getEnvironment("API_KEY", ""),
		ServerPort: getEnvironment("SERVER_PORT", "9090"),
		APIUrl:     getEnvironment("API_URL", ""),
	}
	return config, nil
}

func getEnvironment(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value // Si la variable existe, devuelve su valor
	}
	return defaultValue // Si no existe, devuelve el valor por defecto
}
