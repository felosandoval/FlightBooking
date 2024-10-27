package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	// Carga las variables de entorno desde el archivo .env
	//CONFIGURACIONES
	err := godotenv.Load("./var_entorno.env")
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	// Obtiene el valor de CONNECTION_STRING del archivo .env
	connectionString := os.Getenv("CONNECTION_STRING")

	if connectionString == "" {
		log.Fatal("No se encontró la variable de entorno CONNECTION_STRING")
	}

	// Utiliza el valor de connectionString en tu código de conexión de MongoDB
	clientOptions := options.Client().ApplyURI(connectionString)

	// Conectarse a la base de datos.
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Comprobar que la conexión es correcta.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conexión a la base de datos de MongoDB exitosa!")
	// TERMINO DE CONFIGURACIONES

	db := client.Database("aerolinea")
	vueloCollection := db.Collection("vuelo")
	reservaCollection := db.Collection("reserva")
	costoCollection := db.Collection("costo")
	router := gin.Default()

	router.GET("/aerolinea/vuelo", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		getListaVuelos(ctx, vueloCollection, c)
	})

	router.PUT("/aerolinea/vuelo", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		updateStock(ctx, vueloCollection, c)
	})

	router.GET("/aerolinea/reserva", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		getReserva(ctx, reservaCollection, c)
	})

	router.PUT("/aerolinea/reserva", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		updateReserva(ctx, reservaCollection, c)
	})

	router.POST("/aerolinea/reserva", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		postReserva(ctx, reservaCollection, c)
	})

	router.GET("/aerolinea/ancillares", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		getAncillariesCosto(ctx, costoCollection, c)
	})
	router.GET("/aerolinea/estadisticas", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		RutaMayor(ctx, reservaCollection, c)
		RutaMenor(ctx, reservaCollection, c)
		PromedioPasajeros(ctx, reservaCollection, c)

	})

	// Obtiene el valor de PORT del archivo .env
	puerto := os.Getenv("PORT")

	if puerto == "" {
		log.Fatal("No se encontró el puerto PORT")
	}

	router.Run(":" + puerto)
}
