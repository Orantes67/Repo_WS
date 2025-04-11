package main

import (
	"log"
	"os"
	"sockets-go/infrastructure"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	// Configurar Gin
	r := gin.Default()

	// Configurar rutas
	infrastructure.Routes(r)

	// Obtener puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Iniciar servidor
	log.Printf("Servidor iniciado en el puerto %s", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}