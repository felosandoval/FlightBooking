package main

type Avion struct {
	Modelo             string `bson:"modelo" json:"modelo"`
	Numero_de_serie    string `bson:"numero_de_serie" json:"numero_de_serie"`
	Stock_de_pasajeros int    `bson:"stock_de_pasajeros" json:"stock_de_pasajeros"`
}
