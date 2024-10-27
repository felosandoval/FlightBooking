package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Reserva struct {
	PNR       string     `bson:"pnr" json:"pnr"`
	Vuelos    []Vuelo    `bson:"vuelos" json:"vuelos"`
	Pasajeros []Pasajero `bson:"pasajeros" json:"pasajeros"`
}

func getReserva(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Obtenemos los parámetros de la consulta
	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	// Creamos el filtro para la búsqueda
	filter := bson.M{
		"pnr":                pnr,
		"pasajeros.apellido": apellido,
	}

	// Realizamos la búsqueda en la base de datos
	var reserva Reserva

	if err := collection.FindOne(ctx, filter).Decode(&reserva); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retornamos la reserva encontrada
	c.JSON(http.StatusOK, gin.H{"reserva": reserva})

}

func updateReserva(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Obtenemos los parámetros de la consulta
	pnr := c.Query("pnr")
	apellido := c.Query("apellido")

	// Creamos el filtro para la búsqueda
	filter := bson.M{
		"pnr":                pnr,
		"pasajeros.apellido": apellido,
	}

	// Creamos la estructura para el cuerpo de la solicitud
	var reserva Reserva
	if err := c.BindJSON(&reserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Creamos el objeto de actualización
	update := bson.M{
		"$set": reserva,
	}
	fmt.Println(reserva)
	// Realizamos la actualización en la base de datos
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verificamos si se actualizó alguna reserva
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
		return
	}

	// Retornamos un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"message": "Reserva actualizada"})
}

func postReserva(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Creamos la estructura para el cuerpo de la solicitud
	var reserva Reserva
	if err := c.BindJSON(&reserva); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Insertamos la reserva en la base de datos
	_, err := collection.InsertOne(ctx, reserva)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Obtenemos el ID de la reserva insertada
	id := reserva.PNR

	// Retornamos un mensaje de éxito
	c.JSON(http.StatusOK, gin.H{"PNR": id})
}
