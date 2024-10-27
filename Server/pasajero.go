package main

type Balance struct {
	Ancillaries_ida    int `json:"ancillaries_ida" bson:"ancillaries_ida"`
	Vuelo_ida          int `json:"vuelo_ida" bson:"vuelo_ida"`
	Ancillaries_vuelta int `json:"ancillaries_vuelta" bson:"ancillaries_vuelta"`
	Vuelo_vuelta       int `json:"vuelo_vuelta" bson:"vuelo_vuelta"`
}

type Pasajero struct {
	Nombre      string             `json:"nombre" bson:"nombre"`
	Apellido    string             `json:"apellido" bson:"apellido"`
	Edad        int                `json:"edad" bson:"edad"`
	Ancillaries []AncillariesViaje `json:"ancillaries" bson:"ancillaries"`
	Balances    Balance            `json:"balances" bson:"balances"`
}
