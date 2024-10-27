package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AncillariesViaje struct {
	Ida    []AncillariesCliente `bson:"ida" json:"ida,omitempty"`
	Vuelta []AncillariesCliente `bson:"vuelta" json:"vuelta,omitempty"`
}

type AncillariesCliente struct {
	SSR      string `bson:"ssr" json:"ssr"`
	Cantidad int    `bson:"cantidad" json:"cantidad"`
}

type Ancillarie struct {
	Nombre string `bson:"nombre" json:"nombre"`
	Stock  int    `bson:"stock" json:"stock"`
	SSR    string `bson:"ssr" json:"ssr"`
}

type AncillariesCosto struct {
	SSR   string `bson:"ssr" json:"ssr"`
	Costo int    `bson:"costo" json:"costo"`
}

func getAncillariesCosto(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Obtenemos el SSR de la consulta
	ssr := c.Query("ssr")

	// Creamos el filtro para la búsqueda
	filter := bson.M{
		"ssr": ssr,
	}

	// Realizamos la búsqueda en la base de datos
	var costo AncillariesCosto

	if err := collection.FindOne(ctx, filter).Decode(&costo); err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Costo de ancillary no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(costo)
	// Retornamos el costo encontrado
	c.JSON(http.StatusOK, gin.H{"ssr": costo.SSR, "costo": costo.Costo})
}
