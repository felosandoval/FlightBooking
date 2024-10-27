package main

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Ruta struct {
	Origen   string
	Destino  string
	Ganancia int
}

func RutaMayor(ctx context.Context, collection *mongo.Collection, c *gin.Context) {

	var reservas []Reserva
	// Realizar la consulta a la base de datos y almacenar el resultado en una variable
	filter := bson.M{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &reservas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rutas := make(map[string]int)
	for _, res := range reservas {
		for _, vuelo := range res.Vuelos {
			ruta := vuelo.Origen + "-" + vuelo.Destino
			horaSalida, err := time.Parse("15:04", vuelo.Hora_salida)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			horaLlegada, err := time.Parse("15:04", vuelo.Hora_llegada)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			costo := int(horaLlegada.Sub(horaSalida).Minutes() * 590)
			rutas[ruta] += costo
		}
	}
	// encontrar la ruta con la mayor ganancia
	var rutaMayorGanancia string
	var gananciaMayor int
	for ruta, ganancia := range rutas {
		if ganancia > gananciaMayor {
			rutaMayorGanancia = ruta
			gananciaMayor = ganancia
		}
	}
	c.JSON(http.StatusOK, gin.H{"ruta_mayor_ganancia": rutaMayorGanancia})
}

func RutaMenor(ctx context.Context, collection *mongo.Collection, c *gin.Context) {

	var reservas []Reserva
	// Realizar la consulta a la base de datos y almacenar el resultado en una variable
	filter := bson.M{}
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	if err := cur.All(ctx, &reservas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rutas := make(map[string]int)
	for _, res := range reservas {
		for _, vuelo := range res.Vuelos {
			ruta := vuelo.Origen + "-" + vuelo.Destino
			horaSalida, err := time.Parse("15:04", vuelo.Hora_salida)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			horaLlegada, err := time.Parse("15:04", vuelo.Hora_llegada)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			costo := int(horaLlegada.Sub(horaSalida).Minutes() * 590)
			rutas[ruta] += costo
		}
	}
	// encontrar la ruta con la menor ganancia
	var rutaMenorGanancia string
	var gananciaMenor int
	first := true
	for ruta, ganancia := range rutas {
		if first {
			rutaMenorGanancia = ruta
			gananciaMenor = ganancia
			first = false
		} else {
			if ganancia < gananciaMenor {
				rutaMenorGanancia = ruta
				gananciaMenor = ganancia
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"ruta_menor_ganancia": rutaMenorGanancia})
}

func PromedioPasajeros(ctx context.Context, collection *mongo.Collection, c *gin.Context) {
	// Crear un slice de strings con los nombres de los meses en español
	meses := []string{"Enero", "Febrero", "Marzo", "Abril", "Mayo", "Junio", "Julio", "Agosto", "Septiembre", "Octubre", "Noviembre", "Diciembre"}

	// Crear un map para almacenar la cantidad de pasajeros por mes
	pasajerosPorMes := make(map[string]int)

	// Realizar la consulta a la base de datos y almacenar el resultado en una variable ordenada por fecha
	filter := bson.M{}
	opts := options.Find().SetSort(bson.M{"fecha": 1})
	cur, err := collection.Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reserva no encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)
	var reservas []Reserva
	if err := cur.All(ctx, &reservas); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recorrer las reservas y acumular la cantidad de pasajeros por mes
	for _, reserva := range reservas {
		for _, vuelo := range reserva.Vuelos {
			fecha, err := time.Parse("02/01/2006", vuelo.Fecha)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			mes := meses[fecha.Month()-1]
			pasajerosPorMes[mes] += len(reserva.Pasajeros)
		}
	}

	// Calcular el promedio de pasajeros por mes
	promedioPorMes := make(map[string]float64)
	for mes, pasajeros := range pasajerosPorMes {
		promedioPorMes[mes] = float64(pasajeros) / 30
	}

	c.JSON(http.StatusOK, gin.H{"promedio_pasajeros": promedioPorMes})
}
