package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Avion struct {
	Modelo             string `bson:"modelo" json:"modelo"`
	Numero_de_serie    string `bson:"numero_de_serie" json:"numero_de_serie"`
	Stock_de_pasajeros int    `bson:"stock_de_pasajeros" json:"stock_de_pasajeros"`
}

type Ancillarie struct {
	Nombre string
	Stock  int
	Ssr    string
}

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

type Reserva struct {
	PNR       string
	Vuelos    []Vuelo
	Pasajeros []Pasajero
}

type AncillariesViaje struct {
	Ida    []AncillariesCliente `bson:"ida" json:"ida"`
	Vuelta []AncillariesCliente `bson:"vuelta" json:"vuelta"`
}

type AncillariesCliente struct {
	SSR      string `bson:"ssr" json:"ssr"`
	Cantidad int    `bson:"cantidad" json:"cantidad"`
}

type AncillariesCosto struct {
	SSR   string `bson:"ssr" json:"ssr"`
	Costo int    `bson:"costo" json:"costo"`
}

type Pasajero struct {
	Nombre      string             `json:"nombre"`
	Apellido    string             `json:"apellido"`
	Edad        int                `json:"edad"`
	Ancillaries []AncillariesViaje `json:"ancillaries"`
	Balances    Balance            `json:"balances"`
}

type Balance struct {
	Ancillaries_ida    int `json:"ancillaries_ida" bson:"ancillaries_ida"`
	Vuelo_ida          int `json:"vuelo_ida" bson:"vuelo_ida"`
	Ancillaries_vuelta int `json:"ancillaries_vuelta" bson:"ancillaries_vuelta "`
	Vuelo_vuelta       int `json:"vuelo_vuelta" bson:"vuelo_vuelta"`
}

type VuelosResponse struct {
	Vuelos []Vuelo `json:"vuelos"`
}

type PasajerosResponse struct {
	Pasajeros []Pasajero
}

type ReservaResponse struct {
	Reserva Reserva
}

type AncillariesCostosResponse struct {
	AncillariesCosto AncillariesCosto
}

func listaDeVuelosPP(origen, destino, fecha string) VuelosResponse {
	// Escapamos los valores de los parámetros para formar la URL
	// Formamos la URL con los parámetros
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/vuelo?origen=" + origen + "&destino=" + destino + "&fecha=" + fecha)

	// Hacemos la solicitud GET a la API
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}

	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)

	result := VuelosResponse{}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
	//for _, rec := range result.Vuelos {
	//	fmt.Printf("%s;%s;%s;%s;\n", rec.Numero_vuelo, rec.Origen, rec.Destino, rec.Fecha)
	//}
	return result
}

func traerReserva(pnr, apellido string) Reserva {
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/reserva?pnr=" + pnr + "&apellido=" + apellido)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	result := ReservaResponse{}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result.Reserva
}

func stockVuelo(origen, destino, fecha string, nuevoStock int) {
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/vuelo?origen=" + origen + "&destino=" + destino + "&fecha=" + fecha)

	payload := fmt.Sprintf(`{"stock_de_pasajeros": %d}`, nuevoStock)

	req, err := http.NewRequest(http.MethodPut, url, strings.NewReader(payload))
	if err != nil {
		// Manejar el error adecuadamente
		fmt.Println("Error al crear la solicitud HTTP:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error adecuadamente
		fmt.Println("Error al enviar la solicitud HTTP:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Printf(" ")
	} else {
		// Leer el cuerpo de la respuesta para obtener el mensaje de error
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("No se pudo actualizar el stock de pasajeros del vuelo")
			return
		}

		// Manejar el error adecuadamente
		fmt.Printf("No se pudo actualizar el stock de pasajeros del vuelo: %s\n", body)
	}
}

func insertReserva(reserva Reserva) {
	// Convertir la estructura reserva en un objeto JSON
	jsonData, err := json.Marshal(reserva)
	if err != nil {
		// Manejar el error de conversión
	}

	// Crear una solicitud POST con el objeto JSON como cuerpo de la solicitud
	url := "http://127.0.0.1:5000/aerolinea/reserva"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		// Manejar el error de creación de la solicitud
	}
	req.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud POST
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		// Manejar el error de envío de la solicitud
	}
	defer resp.Body.Close()

	// Leer la respuesta de la solicitud
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Manejar el error de lectura de la respuesta
	}

	// Manejar la respuesta
	fmt.Println(string(body))
}

func modificarReserva(pnr, apellido string, reserva Reserva) {
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/reserva?pnr=" + pnr + "&apellido=" + apellido)

	response, err := http.Get(url)
	// Enviar la solicitud HTTP PUT para actualizar la reserva
	payload, err := json.Marshal(reserva)
	if err != nil {
		// Manejar el error adecuadamente
		fmt.Println("Error al codificar la reserva en JSON:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(payload))
	if err != nil {
		// Manejar el error adecuadamente
		fmt.Println("Error al crear la solicitud HTTP:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err = client.Do(req)
	if err != nil {
		// Manejar el error adecuadamente
		fmt.Println("Error al enviar la solicitud HTTP:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		fmt.Println("La reserva se actualizó correctamente")
	} else {
		// Leer el cuerpo de la respuesta para obtener el mensaje de error
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Println("No se pudo actualizar la reserva")
			return
		}

		// Manejar el error adecuadamente
		fmt.Printf("No se pudo actualizar la reserva: %s\n", body)
	}
}

func traerCosto(ssr string) AncillariesCosto {
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/ancillares?ssr=" + ssr)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error en solicitud HTTP:", err.Error())
		return AncillariesCosto{} // Devolver una instancia vacía de AncillariesCosto en caso de error
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println("Respuesta HTTP:", string(body)) // Imprimir la respuesta recibida

	result := AncillariesCosto{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error al deserializar JSON:", err.Error())
		return AncillariesCosto{} // Devolver una instancia vacía de AncillariesCosto en caso de error
	}
	return result
}

func mostrarEstadisticas() {
	url := fmt.Sprintf("http://127.0.0.1:5000/aerolinea/estadisticas")
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error en solicitud HTTP:", err.Error())
		return // Devolver una instancia vacía de AncillariesCosto en caso de error
	}
	defer response.Body.Close()

	// Leer el contenido de la respuesta HTTP
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err.Error())
		return // Devolver una instancia vacía de AncillariesCosto en caso de error
	}

	// Imprimir el contenido de la respuesta HTTP
	fmt.Println(string(body))
}

func main() {
	// Iniciamos un ciclo infinito
	for {
		var opcion int
		fmt.Println("Menu")
		fmt.Println("1. Gestionar reserva")
		fmt.Println("2. Obtener estadisticas")
		fmt.Println("3. Salir")
		fmt.Printf("Ingrese una opcion: ")

		fmt.Scan(&opcion)

		switch opcion {
		case 1: // GESTIONAR RESERVA
			// Iniciamos un ciclo infinito
			for {
				var subOpcion int
				fmt.Println("Submenu: ")
				fmt.Println("1. Crear reserva")
				fmt.Println("2. Obtener reserva")
				fmt.Println("3. Modificar reserva")
				fmt.Println("4. Salir")
				fmt.Printf("Ingrese una opcion: ")

				fmt.Scan(&subOpcion)

				switch subOpcion {
				case 1: // CREAR RESERVA
					var fecha_ida string
					var fecha_regreso string
					var origen string
					var destino string
					var cantidad_pasajeros int
					var ida int
					var vuelta int
					var ancillaries string
					const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

					fmt.Printf("Ingrese fecha de ida: ")
					fmt.Scan(&fecha_ida)
					fmt.Printf("Ingrese fecha de regreso(si no regresa, ingrese 'no'): ")
					fmt.Scan(&fecha_regreso)
					fmt.Printf("Ingrese origen: ")
					fmt.Scan(&origen)
					fmt.Printf("Ingrese destino: ")
					fmt.Scan(&destino)
					fmt.Printf("Ingrese cantidad de pasajeros: ")
					fmt.Scan(&cantidad_pasajeros)

					//----------FALTA VERIFICACION DE SI HAY STOCK EN EL AVION

					//---------- TRAER VUELOS DISPONIBLES EN LA FECHA DE IDA
					get_vuelos := listaDeVuelosPP(origen, destino, fecha_ida)
					fmt.Println("Vuelos disponibles: ")

					valores_ida := []int{}
					valores_vuelta := []int{}
					fmt.Println("Ida: ")
					for indice, vuelo := range get_vuelos.Vuelos {
						hora1, _ := time.Parse("15:04", vuelo.Hora_salida)
						hora2, _ := time.Parse("15:04", vuelo.Hora_llegada)
						duracion := hora2.Sub(hora1)
						minutos := int(duracion.Minutes())
						costo := minutos * 590
						valores_ida = append(valores_ida, costo)
						fmt.Println(indice+1, ".", vuelo.Numero_vuelo, vuelo.Hora_salida, "-", vuelo.Hora_llegada, "$", costo)

					}
					//---------- HACER ELEGIR AL USUARIO CUAL VUELO DISPONIBLE QUIERE
					fmt.Printf("Ingrese una opcion: ")
					fmt.Scan(&ida)
					vuelosResponse := VuelosResponse{}
					//--------- SI HAY ESPACIO CONTINUA

					if get_vuelos.Vuelos[ida-1].Avion_v.Stock_de_pasajeros >= cantidad_pasajeros {
						cambio := get_vuelos.Vuelos[ida-1].Avion_v.Stock_de_pasajeros - cantidad_pasajeros
						//----------ACTUALIZAMOS STOCK DE VUELO
						stockVuelo(origen, destino, fecha_ida, cambio)
						get_vuelos_actualizado := listaDeVuelosPP(origen, destino, fecha_ida)
						//AGREGAMOS A LA LISTA DE VUELOS
						vuelosResponse.Vuelos = append(vuelosResponse.Vuelos, get_vuelos_actualizado.Vuelos[ida-1])
					} else {
						fmt.Println("No quedan espacios")
						break
					}
					//---------- TRAER VUELOS DISPONIBLES EN LA FECHA DE VUELTA
					if fecha_regreso != "no" {
						//----------TRAE VUELOS
						get_vuelos := listaDeVuelosPP(origen, destino, fecha_regreso)
						fmt.Println("Vuelta: ")
						for indice, vuelo := range get_vuelos.Vuelos {
							hora1, _ := time.Parse("15:04", vuelo.Hora_salida)
							hora2, _ := time.Parse("15:04", vuelo.Hora_llegada)
							duracion := hora2.Sub(hora1)
							minutos := int(duracion.Minutes())
							costo := minutos * 590
							valores_vuelta = append(valores_vuelta, costo)
							fmt.Println(indice+1, ".", vuelo.Numero_vuelo, vuelo.Hora_salida, "-", vuelo.Hora_llegada, "$", costo)
						}
						//INGRESA OPCION
						fmt.Printf("Ingrese una opcion: ")
						fmt.Scan(&vuelta)
						if get_vuelos.Vuelos[vuelta-1].Avion_v.Stock_de_pasajeros >= cantidad_pasajeros {
							cambio := get_vuelos.Vuelos[vuelta-1].Avion_v.Stock_de_pasajeros - cantidad_pasajeros
							//------ACTUALIZAMOS STOCK DE VUELO DE VUELTA
							stockVuelo(origen, destino, fecha_regreso, cambio)
							get_vuelos_actualizado := listaDeVuelosPP(origen, destino, fecha_regreso)
							//AGREGARLO A LA LISTA DE VUELOS
							vuelosResponse.Vuelos = append(vuelosResponse.Vuelos, get_vuelos_actualizado.Vuelos[vuelta-1])
						}
					}
					//comprobamos actualizacion
					// PEDIMOS, CREAMOS Y GUARDAMOS A LOS PASAJEROS

					pasajerosResponse := PasajerosResponse{}
					for i := 1; i <= cantidad_pasajeros; i++ {
						var nombre, apellido string
						var edad int

						fmt.Printf("Pasajero %d:\n", i)
						fmt.Printf("Ingrese nombre: ")
						fmt.Scan(&nombre)
						fmt.Printf("Ingrese apellido: ")
						fmt.Scan(&apellido)
						fmt.Printf("Ingrese edad: ")
						fmt.Scan(&edad)
						apellido = strings.ToLower(apellido)

						//---------- TRAER ANCILLARIES DE IDA
						fmt.Println("Ancillaries Ida: ")
						var ssr AncillariesCosto
						for indice, anciller := range vuelosResponse.Vuelos[ida-1].Ancillaries {
							ssr = traerCosto(anciller.Ssr)
							fmt.Println(indice+1, ".", anciller.Nombre, ":", ssr.Costo)
						}

						//PREGUNTA QUE ANCILLARIES QUIERE DE IDA
						fmt.Printf("Ingrese los ancillaries (separados por comas): ")
						fmt.Scan(&ancillaries)
						ancLista := strings.Split(ancillaries, ",")

						//LOS VA A AGREGANDO A LA LISTA DEL PASAJERO
						ancillariesViaje := AncillariesViaje{}
						ancillariesIda := []AncillariesCliente{}
						ancillariesVuelta := []AncillariesCliente{}
						balances := Balance{}
						valor_anciller_ida := 0
						// iteramos sobre los valores de ancLista
						for _, valor := range ancLista {
							num, err := strconv.Atoi(valor)
							if err != nil {
								fmt.Println("Error al convertir valor a número:", err)
								continue
							}

							// verificamos si el SSR ya está en la lista de AncillariesCliente
							found := false
							for i, ancillariesCliente := range ancillariesIda {
								if ancillariesCliente.SSR == vuelosResponse.Vuelos[ida-1].Ancillaries[num-1].Ssr {
									// si el SSR ya está en la lista, incrementamos la cantidad
									ancillariesIda[i].Cantidad++
									valor_anciller_ida += traerCosto(ancillariesCliente.SSR).Costo
									vuelosResponse.Vuelos[ida-1].Ancillaries[num-1].Stock -= 1
									found = true
									break
								}
							}

							// si el SSR no está en la lista, lo agregamos con una cantidad de 1
							if !found {
								ancillariesCliente := AncillariesCliente{
									SSR:      vuelosResponse.Vuelos[ida-1].Ancillaries[num-1].Ssr,
									Cantidad: 1,
								}
								valor_anciller_ida += traerCosto(ancillariesCliente.SSR).Costo
								ancillariesIda = append(ancillariesIda, ancillariesCliente)
								vuelosResponse.Vuelos[ida-1].Ancillaries[num-1].Stock -= 1
							}
						}

						//OBJETO DE IDA
						fmt.Println("El valor de los ancilleries de ida es: ", valor_anciller_ida)
						ancillariesViaje.Ida = ancillariesIda
						balances.Ancillaries_ida = valor_anciller_ida

						//---------- TRAER ANCILLARIES DE VUELTA
						if fecha_regreso != "no" {
							valor_anciller_vuelta := 0
							fmt.Println("Ancillaries Vuelta: ")
							var ssr int
							for indice, anciller := range vuelosResponse.Vuelos[ida-1].Ancillaries {
								ssr = traerCosto(anciller.Ssr).Costo
								fmt.Println(indice+1, ".", anciller.Nombre, ":", ssr)
							}

							//PREGUNTA QUE ANCILLARIES QUIERE DE VU
							fmt.Printf("Ingrese los ancillaries (separados por comas): ")
							fmt.Scan(&ancillaries)
							ancLista := strings.Split(ancillaries, ",")

							// iteramos sobre los valores de ancLista
							for _, valor := range ancLista {
								num, err := strconv.Atoi(valor)
								if err != nil {
									fmt.Println("Error al convertir valor a número:", err)
									continue
								}
								// verificamos si el SSR ya está en la lista de AncillariesCliente
								found := false
								for i, ancillariesCliente := range ancillariesVuelta {
									if ancillariesCliente.SSR == vuelosResponse.Vuelos[vuelta-1].Ancillaries[num-1].Ssr {
										// si el SSR ya está en la lista, incrementamos la cantidad
										ancillariesVuelta[i].Cantidad++
										valor_anciller_vuelta += traerCosto(ancillariesCliente.SSR).Costo
										vuelosResponse.Vuelos[vuelta-1].Ancillaries[num-1].Stock -= 1
										found = true
										break
									}
								}

								// si el SSR no está en la lista, lo agregamos con una cantidad de 1
								if !found {
									ancillariesCliente := AncillariesCliente{
										SSR:      vuelosResponse.Vuelos[vuelta-1].Ancillaries[num-1].Ssr,
										Cantidad: 1,
									}
									valor_anciller_vuelta += traerCosto(ancillariesCliente.SSR).Costo
									ancillariesVuelta = append(ancillariesVuelta, ancillariesCliente)
									vuelosResponse.Vuelos[vuelta-1].Ancillaries[num-1].Stock -= 1
								}
							}
							fmt.Println("El costo de los ancilleries de vuelta son: ", valor_anciller_vuelta)
							//OBJETO DE VUELTA
							ancillariesViaje.Vuelta = ancillariesVuelta
							balances.Ancillaries_vuelta = valor_anciller_vuelta
						}
						//SE ARMA LA IDA Y LA VUELTA DENTRO DE UN OBJETO
						nuevo_ancillar_ida := AncillariesViaje{
							Ida:    ancillariesIda,
							Vuelta: nil,
						}
						nuevo_ancillar_vuelta := AncillariesViaje{
							Ida:    nil,
							Vuelta: ancillariesVuelta,
						}
						//OBJETOS DE BALANCES PARA
						balances.Vuelo_ida = valores_ida[ida-1]

						if fecha_regreso != "no" {
							balances.Vuelo_vuelta = valores_vuelta[vuelta-1]
						}

						//SE ASIGNA
						ancillar := []AncillariesViaje{nuevo_ancillar_ida, nuevo_ancillar_vuelta}
						//SE LE INGRESA A LA PERSONA
						persona := Pasajero{
							Nombre:      nombre,
							Apellido:    apellido,
							Edad:        edad,
							Ancillaries: ancillar,
							Balances:    balances,
						}

						//AGREGAR A LA LISTA DE PASAJEROS
						pasajerosResponse.Pasajeros = append(pasajerosResponse.Pasajeros, persona)
						fmt.Println("El total de el viaje es: ", balances.Ancillaries_ida+balances.Ancillaries_vuelta+balances.Vuelo_ida+balances.Vuelo_vuelta)
					}
					//Generar PNR
					rand.Seed(time.Now().UnixNano())
					pnr := make([]byte, 6)
					for i := range pnr {
						pnr[i] = charset[rand.Intn(len(charset))]
					}

					//---------- HACER ELEGIR AL USUARIO
					reserva := Reserva{
						PNR:       string(pnr),
						Vuelos:    vuelosResponse.Vuelos,
						Pasajeros: pasajerosResponse.Pasajeros,
					}
					//INSERTAMOS LA RESERVA
					insertReserva(reserva)

				case 2: // OBTENER RESERVA
					var pnr string
					var apellido string

					fmt.Printf("Escriba el PNR: ")
					fmt.Scan(&pnr)
					fmt.Printf("Escriba el apellido: ")
					fmt.Scan(&apellido)
					apellido = strings.ToLower(apellido)
					reserva := traerReserva(pnr, apellido)

					//FORMATO DE IMPRESION
					if len(reserva.Vuelos) > 1 {
						fmt.Println("Ida: ")
						fmt.Println(reserva.Vuelos[0].Numero_vuelo, " ", reserva.Vuelos[0].Hora_salida, " - ", reserva.Vuelos[0].Hora_llegada)
						fmt.Println("Vuelta: ")
						fmt.Println(reserva.Vuelos[1].Numero_vuelo, " ", reserva.Vuelos[1].Hora_salida, " - ", reserva.Vuelos[1].Hora_llegada)

						fmt.Println("Pasajeros: ")
						for _, pasajeros := range reserva.Pasajeros {
							var ancListaIda []string
							var ancListaVuelta []string

							fmt.Println(pasajeros.Nombre, " ", pasajeros.Apellido, " ", pasajeros.Edad)
							for _, ancillaries := range pasajeros.Ancillaries[0].Ida {
								ancListaIda = append(ancListaIda, ancillaries.SSR)
							}
							fmt.Printf("Ancillaries ida: ")
							for indice, s := range ancListaIda {
								if indice < len(ancListaIda)-1 {
									fmt.Print(s, ",")
								} else {
									fmt.Println(s)
								}
							}

							for _, ancillaries := range pasajeros.Ancillaries[1].Vuelta {
								ancListaVuelta = append(ancListaVuelta, ancillaries.SSR)
							}
							fmt.Printf("Ancillaries vuelta: ")
							for indice, s := range ancListaVuelta {
								if indice < len(ancListaVuelta)-1 {
									fmt.Print(s, ",")
								} else {
									fmt.Println(s)
								}
							}
						}

					} else {
						fmt.Printf("Ida: ")
						fmt.Println(reserva.Vuelos[0].Numero_vuelo, " ", reserva.Vuelos[0].Hora_salida, " - ", reserva.Vuelos[0].Hora_llegada)
						var ancLista []string
						for _, pasajeros := range reserva.Pasajeros {
							fmt.Println("Pasajeros: ")
							fmt.Println(pasajeros.Nombre, " ", pasajeros.Edad)
							for _, ancillaries := range pasajeros.Ancillaries[0].Ida {
								ancLista = append(ancLista, ancillaries.SSR)
							}
						}

					}
				case 3: // MODIFICAR RESERVA
					var opcion3 int
					var pnr string
					var apellido string

					fmt.Printf("Escriba el PNR: ")
					fmt.Scan(&pnr)
					fmt.Printf("Escriba el apellido: ")
					fmt.Scan(&apellido)
					fmt.Println("Opciones: ")
					fmt.Println("1. Cambiar fecha de vuelo")
					fmt.Println("2. Adicionar ancillaries")
					fmt.Println("3. Salir")
					fmt.Printf("Ingrese una opcion: ")
					apellido = strings.ToLower(apellido)
					fmt.Scan(&opcion3)

					switch opcion3 {
					case 1: // CAMBIAR FECHA DE VUELO
						resultado := traerReserva(pnr, apellido)
						var opcion4 int
						var nueva_fecha string
						if len(resultado.Vuelos) > 1 {
							fmt.Println("Vuelos: ")
							fmt.Println("1. Ida: ", resultado.Vuelos[0].Numero_vuelo, resultado.Vuelos[0].Hora_salida, "-", resultado.Vuelos[0].Hora_llegada)
							fmt.Println("2. Vuelta: ", resultado.Vuelos[1].Numero_vuelo, resultado.Vuelos[1].Hora_salida, "-", resultado.Vuelos[1].Hora_llegada)
							fmt.Printf("Ingrese una opcion: ")
							fmt.Scan(&opcion4)
							fmt.Printf("Ingrese nueva fecha: ")
							fmt.Scan(&nueva_fecha)
							// GET DE VUELOS DISPONIBLES EN ESA NUEVA
							get_vuelos := listaDeVuelosPP(resultado.Vuelos[opcion4-1].Origen, resultado.Vuelos[opcion4-1].Destino, nueva_fecha)

							//CASO 1 DONDE LA FECHA DE IDA NUEVA NO PUEDE IR DESPUES DE LA VUELTA
							if opcion4 == 1 {
								// Definimos las fechas a comparar
								fecha_vuelta := resultado.Vuelos[1].Fecha

								// Convertimos las fechas a objetos time.Time
								layout := "02/01/2006" // dd/mm/aaaa
								fecha_vueltaTime, _ := time.Parse(layout, fecha_vuelta)
								fecha_nuevaTime, _ := time.Parse(layout, nueva_fecha)

								// Comparamos las fechas
								if fecha_nuevaTime.Before(fecha_vueltaTime) {
									fmt.Println("Vuelos disponibles: ")
									i := 1
									for _, vuelo := range get_vuelos.Vuelos {
										fmt.Println(i, ".", vuelo.Numero_vuelo, vuelo.Hora_salida, "-", vuelo.Hora_llegada)
									}

									var opcion5 int
									fmt.Printf("Ingrese una opcion: ")
									fmt.Scan(&opcion5)
									get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros -= len(resultado.Pasajeros)
									vuelos := VuelosResponse{Vuelos: []Vuelo{get_vuelos.Vuelos[opcion5-1], resultado.Vuelos[1]}}
									reserva_nueva := Reserva{
										PNR:       pnr,
										Vuelos:    vuelos.Vuelos,
										Pasajeros: resultado.Pasajeros,
									}

									// MODIFICAMOS LA RESERVA
									modificarReserva(pnr, apellido, reserva_nueva)

									//CAMBIO DE STOCK PARA EL VUELO NUEVO (DISMINUIMOS STOCK)
									cambio := get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros
									stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, nueva_fecha, cambio)

									//CAMBIO DE STOCK PARA EL VUELO ANTIGUO (AUMENTAMOS STOCK)
									cambio_antiguo := resultado.Vuelos[0].Avion_v.Stock_de_pasajeros + len(resultado.Pasajeros)
									stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, resultado.Vuelos[0].Fecha, cambio_antiguo)

								} else if fecha_nuevaTime.After(fecha_vueltaTime) {
									fmt.Println("Fecha nueva de IDA no puede superar la fecha de vuelta")
									break
								}

								//CASO 2 DONDE LA FECHA DE VUELTA NUEVA NO PUEDE IR ANTES DE LA IDA
							} else if opcion4 == 2 {
								// Definimos las fechas a comparar
								fecha_ida := resultado.Vuelos[0].Fecha

								// Convertimos las fechas a objetos time.Time
								layout := "02/01/2006" // dd/mm/aaaa
								fecha_idaTime, _ := time.Parse(layout, fecha_ida)
								fecha_nuevaTime, _ := time.Parse(layout, nueva_fecha)

								// Comparamos las fechas
								if fecha_nuevaTime.Before(fecha_idaTime) {
									fmt.Println("Fecha nueva de IDA no puede superar la fecha de vuelta")
									break

								} else if fecha_nuevaTime.After(fecha_idaTime) {
									fmt.Println("Vuelos disponibles: ")
									i := 1
									for _, vuelo := range get_vuelos.Vuelos {
										fmt.Println(i, ".", vuelo.Numero_vuelo, vuelo.Hora_salida, "-", vuelo.Hora_llegada)
									}

									var opcion5 int
									fmt.Printf("Ingrese una opcion: ")
									fmt.Scan(&opcion5)

									get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros -= len(resultado.Pasajeros)
									vuelos := VuelosResponse{Vuelos: []Vuelo{resultado.Vuelos[0], get_vuelos.Vuelos[opcion5-1]}}
									reserva_nueva := Reserva{
										PNR:       pnr,
										Vuelos:    vuelos.Vuelos,
										Pasajeros: resultado.Pasajeros,
									}

									// MODIFICAMOS LA RESERVA
									modificarReserva(pnr, apellido, reserva_nueva)

									//CAMBIO DE STOCK PARA EL VUELO NUEVO (DISMINUIMOS STOCK)
									cambio := get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros
									stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, nueva_fecha, cambio)

									//CAMBIO DE STOCK PARA EL VUELO ANTIGUO (AUMENTAMOS STOCK)
									cambio_antiguo := resultado.Vuelos[0].Avion_v.Stock_de_pasajeros + len(resultado.Pasajeros)
									stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, resultado.Vuelos[0].Fecha, cambio_antiguo)

								}
							}

						} else {
							fmt.Println("1. Ida: ", resultado.Vuelos[0].Numero_vuelo, resultado.Vuelos[0].Hora_salida, "-", resultado.Vuelos[0].Hora_llegada)
							fmt.Printf("Ingrese nueva fecha: ")
							fmt.Scan(&nueva_fecha)
							get_vuelos := listaDeVuelosPP(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, nueva_fecha)

							fmt.Println("Vuelos disponibles: ")
							i := 1
							for _, vuelo := range get_vuelos.Vuelos {
								fmt.Println(i, ".", vuelo.Numero_vuelo, vuelo.Hora_salida, "-", vuelo.Hora_llegada)
							}

							var opcion5 int
							fmt.Printf("Ingrese una opcion: ")
							fmt.Scan(&opcion5)
							get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros -= len(resultado.Pasajeros)
							vuelos := VuelosResponse{Vuelos: []Vuelo{get_vuelos.Vuelos[opcion5-1]}}
							reserva_nueva := Reserva{
								PNR:       pnr,
								Vuelos:    vuelos.Vuelos,
								Pasajeros: resultado.Pasajeros,
							}

							// MODIFICAMOS LA RESERVA
							modificarReserva(pnr, apellido, reserva_nueva)

							//CAMBIO DE STOCK PARA EL VUELO NUEVO (DISMINUIMOS STOCK)
							cambio := get_vuelos.Vuelos[opcion5-1].Avion_v.Stock_de_pasajeros
							stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, nueva_fecha, cambio)

							//CAMBIO DE STOCK PARA EL VUELO ANTIGUO (AUMENTAMOS STOCK)
							cambio_antiguo := resultado.Vuelos[0].Avion_v.Stock_de_pasajeros + len(resultado.Pasajeros)
							stockVuelo(resultado.Vuelos[0].Origen, resultado.Vuelos[0].Destino, resultado.Vuelos[0].Fecha, cambio_antiguo)

						}

					case 2: // AGREGAR ANCILLARIES
						var ancillaries string
						get_reserva := traerReserva(pnr, apellido)

						if len(get_reserva.Vuelos) > 1 {
							fmt.Println("1. Ida: ", get_reserva.Vuelos[0].Numero_vuelo, get_reserva.Vuelos[0].Hora_salida, "-", get_reserva.Vuelos[0].Hora_salida)
							fmt.Println("2. Vuelta: ", get_reserva.Vuelos[1].Numero_vuelo, get_reserva.Vuelos[1].Hora_salida, "-", get_reserva.Vuelos[1].Hora_llegada)
							fmt.Printf("Ingrese una opcion: ")
							fmt.Scan(&opcion)

							if opcion == 1 {
								contador := 1
								fmt.Println("Ancilleres disponibles: ")
								valor_anciller := 0
								for _, ancilleres := range get_reserva.Vuelos[0].Ancillaries {
									if ancilleres.Stock >= 1 {
										costo := traerCosto(ancilleres.Ssr)
										fmt.Println(contador, ".", ancilleres.Nombre, ":", costo.Costo)
										contador++
									}
								}
								fmt.Printf("Ingrese los ancillaries (separados por comas): ")
								fmt.Scan(&ancillaries)
								ancLista := strings.Split(ancillaries, ",")
								for indice, pasajero := range get_reserva.Pasajeros {
									if strings.ToLower(pasajero.Apellido) == strings.ToLower(apellido) {
										for _, valor := range ancLista {
											num, err := strconv.Atoi(valor)
											if err != nil {
												fmt.Println("Error al convertir valor a número:", err)
												continue
											}
											// verificamos si el SSR ya está en la lista de AncillariesCliente
											found := false

											for i, ancillariesCliente := range pasajero.Ancillaries[0].Ida {
												if ancillariesCliente.SSR == get_reserva.Vuelos[0].Ancillaries[num-1].Ssr {
													// si el SSR ya está en la lista, incrementamos la cantidad
													pasajero.Ancillaries[0].Ida[i].Cantidad++
													valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
													get_reserva.Vuelos[0].Ancillaries[num-1].Stock -= 1
													found = true
													break
												}
											}

											// si el SSR no está en la lista, lo agregamos con una cantidad de 1
											if !found {
												ancillariesCliente := AncillariesCliente{
													SSR:      get_reserva.Vuelos[0].Ancillaries[num-1].Ssr,
													Cantidad: 1,
												}
												valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
												pasajero.Ancillaries[0].Ida = append(pasajero.Ancillaries[0].Ida, ancillariesCliente)
												get_reserva.Vuelos[0].Ancillaries[num-1].Stock -= 1
											}

										}

										get_reserva.Pasajeros[indice].Balances.Ancillaries_ida += valor_anciller

									}
								}

								reserva := Reserva{
									PNR:       pnr,
									Vuelos:    get_reserva.Vuelos,
									Pasajeros: get_reserva.Pasajeros,
								}

								modificarReserva(pnr, apellido, reserva)
								fmt.Println("Total ancillaries: ", valor_anciller)

							} else if opcion == 2 {
								contador := 1
								fmt.Println("Ancilleres disponibles: ")
								valor_anciller := 0
								for _, ancilleres := range get_reserva.Vuelos[0].Ancillaries {
									if ancilleres.Stock >= 1 {
										costo := traerCosto(ancilleres.Ssr)
										fmt.Println(contador, ".", ancilleres.Nombre, ":", costo.Costo)
										contador++
									}
								}
								fmt.Printf("Ingrese los ancillaries (separados por comas): ")
								fmt.Scan(&ancillaries)
								ancLista := strings.Split(ancillaries, ",")
								for indice, pasajero := range get_reserva.Pasajeros {
									if strings.ToLower(pasajero.Apellido) == strings.ToLower(apellido) {
										for _, valor := range ancLista {
											num, err := strconv.Atoi(valor)
											if err != nil {
												fmt.Println("Error al convertir valor a número:", err)
												continue
											}
											// verificamos si el SSR ya está en la lista de AncillariesCliente
											found := false

											for i, ancillariesCliente := range pasajero.Ancillaries[1].Vuelta {
												if ancillariesCliente.SSR == get_reserva.Vuelos[1].Ancillaries[num-1].Ssr {
													// si el SSR ya está en la lista, incrementamos la cantidad
													pasajero.Ancillaries[1].Vuelta[i].Cantidad++
													valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
													get_reserva.Vuelos[1].Ancillaries[num-1].Stock -= 1
													found = true
													break
												}
											}

											// si el SSR no está en la lista, lo agregamos con una cantidad de 1
											if !found {
												ancillariesCliente := AncillariesCliente{
													SSR:      get_reserva.Vuelos[1].Ancillaries[num-1].Ssr,
													Cantidad: 1,
												}
												valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
												pasajero.Ancillaries[1].Vuelta = append(pasajero.Ancillaries[1].Vuelta, ancillariesCliente)
												get_reserva.Vuelos[1].Ancillaries[num-1].Stock -= 1
											}
										}
										get_reserva.Pasajeros[indice].Balances.Ancillaries_vuelta += valor_anciller

									}
								}

								reserva := Reserva{
									PNR:       pnr,
									Vuelos:    get_reserva.Vuelos,
									Pasajeros: get_reserva.Pasajeros,
								}

								modificarReserva(pnr, apellido, reserva)
								fmt.Println("Total ancillaries: ", valor_anciller)
							}
						} else {
							contador := 1
							fmt.Println("Ancilleres disponibles: ")
							valor_anciller := 0
							for _, ancilleres := range get_reserva.Vuelos[0].Ancillaries {
								if ancilleres.Stock >= 1 {
									costo := traerCosto(ancilleres.Ssr)
									fmt.Println(contador, ".", ancilleres.Nombre, ":", costo.Costo)
									contador++
								}
							}
							fmt.Printf("Ingrese los ancillaries (separados por comas): ")
							fmt.Scan(&ancillaries)
							ancLista := strings.Split(ancillaries, ",")
							for indice, pasajero := range get_reserva.Pasajeros {
								if strings.ToLower(pasajero.Apellido) == strings.ToLower(apellido) {
									for _, valor := range ancLista {
										num, err := strconv.Atoi(valor)
										if err != nil {
											fmt.Println("Error al convertir valor a número:", err)
											continue
										}
										// verificamos si el SSR ya está en la lista de AncillariesCliente
										found := false

										for i, ancillariesCliente := range pasajero.Ancillaries[0].Ida {
											if ancillariesCliente.SSR == get_reserva.Vuelos[0].Ancillaries[num-1].Ssr {
												// si el SSR ya está en la lista, incrementamos la cantidad
												pasajero.Ancillaries[0].Ida[i].Cantidad++
												valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
												found = true
												break
											}
										}

										// si el SSR no está en la lista, lo agregamos con una cantidad de 1
										if !found {
											ancillariesCliente := AncillariesCliente{
												SSR:      get_reserva.Vuelos[0].Ancillaries[num-1].Ssr,
												Cantidad: 1,
											}
											valor_anciller += traerCosto(ancillariesCliente.SSR).Costo
											pasajero.Ancillaries[0].Ida = append(pasajero.Ancillaries[0].Ida, ancillariesCliente)
										}
									}
									get_reserva.Pasajeros[indice].Balances.Ancillaries_ida += valor_anciller

								}
							}

							reserva := Reserva{
								PNR:       pnr,
								Vuelos:    get_reserva.Vuelos,
								Pasajeros: get_reserva.Pasajeros,
							}

							modificarReserva(pnr, apellido, reserva)
							fmt.Println("Total ancillaries: ", valor_anciller)

						}

					case 3: // SALIR
						break

					default:
						fmt.Println("Ingrese una opcion valida")
					}

					// Si la subOpción es 4, salimos del ciclo del submenú 2
					if opcion3 == 3 {
						break
					}

				case 4: // SALIR
					// Rompemos el ciclo del submenú al ingresar la opción 4
					break

				default:
					fmt.Println("Ingrese una opcion valida")
				}

				// Si la subOpción es 4, salimos del ciclo del submenú
				if subOpcion == 4 {
					break
				}
			}

		case 2: // OBTENER ESTADISTICAS
			mostrarEstadisticas()

		case 3: // SALIR
			fmt.Println("Programa finalizado!")
			// Rompemos el ciclo al ingresar la opción 3
			break
		default:
			fmt.Println("Ingrese una opcion valida")
		}
		// Si la opción es 3, salimos del ciclo principal
		if opcion == 3 {
			break
		}

	}
}
