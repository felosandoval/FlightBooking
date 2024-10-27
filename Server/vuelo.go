package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Vuelo struct {
	Numero_vuelo string       `bson:"numero_vuelo" json:"numero_vuelo"`
	Origen       string       `bson:"origen" json:"origen"`
	Destino      string       `bson:"destino" json:"destino"`
	Hora_salida  string       `bson:"hora_salida" json:"hora_salida"`
	Hora_llegada string       `bson:"hora_llegada" json:"hora_llegada"`
	Fecha        string       `bson:"fecha" json:"fecha"`
	Avion_v      Avion        `bson:"avion" json:"avion"`
	Ancillaries  []Ancillarie `bson:"ancillaries" json:"ancillaries"`
}

func getListaVuelos(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Obtenemos los parámetros de la consulta
	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	// Creamos un slice para almacenar los vuelos
	var vuelos []Vuelo

	// Creamos el filtro para la búsqueda
	filter := bson.M{
		"origen":  origen,
		"destino": destino,
		"fecha":   fecha,
	}
	// Realizamos la búsqueda en la base de datos
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	// Iteramos sobre los resultados y agregamos los vuelos al slice
	for cursor.Next(ctx) {
		var vuelo Vuelo
		if err := cursor.Decode(&vuelo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		vuelos = append(vuelos, vuelo)
	}

	// Verificamos si hubo algún error durante la iteración
	if err := cursor.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Retornamos los vuelos encontrados
	c.JSON(http.StatusOK, gin.H{"vuelos": vuelos})
}

func updateStock(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Obtenemos los parámetros de la petición
	origen := c.Query("origen")
	destino := c.Query("destino")
	fecha := c.Query("fecha")

	var stock struct {
		StockDePasajeros int `json:"stock_de_pasajeros"`
	}

	if err := c.ShouldBindJSON(&stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El cuerpo de la solicitud no es válido"})
		return
	}

	// Creamos el filtro para la actualización
	filter := bson.M{
		"origen":  origen,
		"destino": destino,
		"fecha":   fecha,
	}

	// Creamos el update para modificar el stock de pasajeros
	update := bson.M{
		"$set": bson.M{"avion.stock_de_pasajeros": stock.StockDePasajeros},
	}

	// Realizamos la actualización en la base de datos
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Verificamos si se actualizó algún documento
	if result.ModifiedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No se encontró ningún vuelo que coincida con los parámetros especificados"})
		return
	}

	// Retornamos una respuesta exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Stock de pasajeros actualizado correctamente"})
}
