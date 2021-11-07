package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"time"

	pb "example.com/go-jugador-grpc/jugador"
	"google.golang.org/grpc"
)

var cantidadJugadores int = 0
var etapa int = 0 ////Etapa donde se parte: 0 = primer juego, 1= segundo juego, 2= tercer juego
var eleccion_lider int = 0
var jugadores []int
var jugadoresEtapa1 []int
var sumaetapa1 []int
var jugadoresVivos int = 16
var jugadoresMuertos int = 0
var jugadoresActivos int = 16
var jugadoresRespuestas int = 0
var ronda int = 1
var esperaEtapa = true //Para no enviar la etapa si el
var jugarEtapa = true  //Para no ejecutar el menu si se esta jugando una etapa
var jugadorMuerto int = 0
var jugadoresGanadores int = 0
var suma_primero int = 0
var suma_segundo int = 0
var primer_equipo []int
var segundo_equipo []int
var sumaetapa2 = true
var TodoNada []int
var ListoTodoNada int = 0
var ListoTodoNada2 int = 0
var toychato3 = true
var toychato1 = true
var jugadoresEnviar int = 0
var jugadoresVivosEtapa2 int
var terminar = false

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedJugadoresServer
}

func random_range(min int, max int) int {
	number := (rand.Intn(max-min) + min)
	return number
}

func etapa1() int {
	rand.Seed(time.Now().UnixNano())
	return random_range(6, 10)
}

func etapa2() int {
	rand.Seed(time.Now().UnixNano())
	return random_range(1, 4)
}

func etapa3() int {
	rand.Seed(time.Now().UnixNano())
	return random_range(1, 10)
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

func findAndDelete(s []int, item int) []int {
	index := 0
	for _, i := range s {
		if i != item {
			s[index] = i
			index++
		}
	}
	return s[:index]
}

func (s *server) EnviarSolicitud(ctx context.Context, in *pb.Peticion) (*pb.Respuesta, error) {

	cantidadJugadores++
	fmt.Println(">> El jugador ", cantidadJugadores, " se ha unido")

	for cantidadJugadores != 16 {
	}

	return &pb.Respuesta{Mensaje: "El juego ha comenzado"}, nil
	/*if cantidadJugadores == 16 {
		return &pb.Respuesta{Mensaje: "El juego ha comenzado"}, nil
	}
	return &pb.Respuesta{Mensaje: "El juego ha comenzado"}, nil*/
}

func (s *server) EnviarInicio(ctx context.Context, in *pb.Inicio) (*pb.Juego, error) {
	//mensaje := in.GetInicio()
	for esperaEtapa {
	}
	if etapa == 1 {
		return &pb.Juego{Juego: "La Etapa 1 comenzara ahora"}, nil
	} else if etapa == 2 {
		return &pb.Juego{Juego: "La Etapa 2 comenzara ahora"}, nil
	} else {
		return &pb.Juego{Juego: "La Etapa 3 comenzara ahora"}, nil
	}

}

func (s *server) EnviarJugada(ctx context.Context, in *pb.Numero) (*pb.Estado, error) {
	mensaje := in.GetNumero()  //Numero escogido por el jugador
	jugador := in.GetJugador() //Jugador que envio la jugada
	ganar := false
	estado := ""
	equipo := 1
	cantidad := time.Duration(100 * int(jugador))
	ronda_local := ronda

	time.Sleep(250 * time.Millisecond)

	for i := 0; i < len(primer_equipo); i++ {
		if int(jugador) == primer_equipo[i] {
			equipo = 1
			//suma_primero += int(mensaje)
		} else if int(jugador) == segundo_equipo[i] {
			//suma_segundo += int(mensaje)
			equipo = 2
		}
	}

	time.Sleep(cantidad * time.Millisecond)
	jugadoresRespuestas++
	time.Sleep(cantidad * time.Millisecond)

	//fmt.Println(jugadoresRespuestas)
	time.Sleep(100 * time.Millisecond)
	if etapa == 1 {
		index := jugador - 1

		/*for i := 0; i < len(jugadores); i++ {
			if jugadores[i] == int(jugador) {
				index = i
				break
			}
		}*/

		if jugadoresRespuestas == jugadoresActivos {
			toychato1 = false
		}

		for toychato1 {

		} //Todos han enviado su respuesta
		time.Sleep(cantidad * time.Millisecond)
		ronda_local++

		esperaEtapa = true
		fmt.Println(">> El jugador ", jugador, " ha escogido ", mensaje)
		if int(mensaje) >= eleccion_lider {
			//RIP
			jugadoresMuertos++
			//Quitar jugador que murio de la lista

			jugadoresEtapa1[index] = -1
			sumaetapa1[index] = -1

			fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto")
			estado = "Morir"

			//return &pb.Estado{Estado: "Morir", Ronda: int32(ronda)}, nil
		} else {
			//Sumar a la lista de suma de etapa 1(???)
			sumaetapa1[index] = sumaetapa1[index] + int(mensaje)
			estado = "Vivir"
			//Si su suma es mayor a 21
			if sumaetapa1[index] >= 15 {
				estado = "Ronda"
				fmt.Println(">> Jugador ", jugador, " ha ganado la etapa")
				jugadoresGanadores++
				ganar = true //Si ha ganado la etapa
				//return &pb.Estado{Estado: "Vivir", Ronda: int32(-1)}, nil
			} else if sumaetapa1[index] < 15 && ronda == 4 { //En la ultima ronda no llegaron a 21   ///**********************//////////////Se modifico para testeos AAAAAAHHHHHHHHH
				//RIP
				jugadoresMuertos++
				//Quitar jugador que murio de la lista

				jugadoresEtapa1[index] = -1
				sumaetapa1[index] = -1

				fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto por no llegar a 21")
				estado = "Morir"
			}

			//return &pb.Estado{Estado: "Vivir", Ronda: int32(ronda)}, nil
		}

		//fmt.Println("Jugadores activos: ", jugadoresActivos)
		if jugadoresEnviar == jugadoresActivos-1 { //Si todos los jugadores han enviado su respuesta
			ronda++
			eleccion_lider = etapa1() //Cambiar el random?
			jugadoresRespuestas = 0   //Inicializar jugadoresRespuesta
			jugadoresVivos = jugadoresVivos - jugadoresMuertos
			jugadoresActivos = jugadoresActivos - jugadoresMuertos - jugadoresGanadores
			jugadoresMuertos = 0
			jugadoresGanadores = 0

			//Para que no puedan volver seguir con otra etapa del juego mientras se ejecuta esta
			toychato1 = true
			if jugadoresVivos == 1 && jugadoresEtapa1[index] != -1 {
				fmt.Println(">> El jugador", jugador, " ha ganado el juego")
				jugadoresActivos--
				terminar = true
				return &pb.Estado{Estado: "Ganador", Ronda: int32(-1)}, nil

			}

			if ronda < 5 {
				fmt.Println(">> La ronda ", ronda, " ha comenzado")
				//fmt.Println(">>> El lider ha escogido ", eleccion_lider)

			} else { //Ronda 5: Si siguen vivos y su suma es menor a 21, mueren
				//RIP
				//RIP
				/*jugadoresVivos--
				jugadoresActivos--
				//Quitar jugador que murio de la lista

				jugadores = remove(jugadores, index)
				sumaetapa1 = remove(sumaetapa1, index)

				estado = "Morir"*/
				esperaEtapa = true
				fmt.Println(">> Se acabaron las rondas")
				//Jugadores activos = jugadoresvivos para resetear
				jugarEtapa = false
			}
			jugadoresEnviar = 0
			//fmt.Println(">> Quedan ", jugadoresVivos, " jugadores vivos y ", jugadoresActivos, " jugadores activos.")
			return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
		}
	} else if etapa == 2 {
		index := 0
		for i := 0; i < len(jugadores); i++ {
			if jugadores[i] == int(jugador) {
				index = i
				break
			}
		}

		if jugadorMuerto == int(jugador) {
			estado = "Morir"
			fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto, este es el random")
			return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
		}

		//fmt.Println(">> El jugador ", jugador, " ha escogido ", mensaje)

		if equipo == 1 {
			suma_primero += int(mensaje)
		} else {
			suma_segundo += int(mensaje)
		}

		if jugadoresRespuestas == jugadoresVivos {
			sumaetapa2 = false
			esperaEtapa = true
		}
		for sumaetapa2 { //Esperar que todos envien sus numeros
		}

		time.Sleep(cantidad * time.Millisecond)

		jugadoresRespuestas = 0

		//fmt.Println("suma del equipo 1: ", suma_primero)
		//fmt.Println("suma del equipo 2: ", suma_segundo)
		paridad_lider := eleccion_lider % 2

		if suma_primero%2 != paridad_lider && suma_segundo%2 != paridad_lider {
			numerito := random_range(1, 2)
			if equipo == numerito {
				fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto, empate, era del equipo ------ ", equipo)
				sumaetapa2 = true
				jugadoresVivos--
				estado = "Morir"
				jugadores = remove(jugadores, index)
				return &pb.Estado{Estado: "Morir", Ronda: int32(ronda)}, nil
			}
		} else {
			if suma_primero%2 == paridad_lider {
				if equipo == 1 {
					sumaetapa2 = true
					estado = "Vivir"
					//return &pb.Estado{Estado: "Vivir", Ronda: int32(-1)}, nil
				}
			} else if suma_primero%2 != paridad_lider {
				if equipo == 1 {
					fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto, era del equipo ------ ", equipo)
					sumaetapa2 = true
					jugadoresVivos--
					estado = "Morir"
					jugadores = remove(jugadores, index)
					//return &pb.Estado{Estado: "Morir", Ronda: int32(ronda)}, nil
				}
			}
			if suma_segundo%2 == paridad_lider {
				if equipo == 2 {
					sumaetapa2 = true
					estado = "Vivir"
					//return &pb.Estado{Estado: "Vivir", Ronda: int32(-1)}, nil
				}
			} else if suma_segundo%2 != paridad_lider {
				if equipo == 2 {
					fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto, era del equipo ------ ", equipo)
					jugadoresVivos--
					sumaetapa2 = true
					estado = "Morir"
					jugadores = remove(jugadores, index)
					//return &pb.Estado{Estado: "Morir", Ronda: int32(ronda)}, nil
				}
			}
		}

		jugadoresEnviar++
		/*
			for jugadoresEnviar != jugadoresVivosEtapa2 {

			}*/
		if jugadoresVivos == 1 && int(jugador) == jugadores[0] {
			fmt.Println(">> El jugador", jugador, " ha ganado el juego")
			terminar = true
			return &pb.Estado{Estado: "Ganador", Ronda: int32(-1)}, nil
		}

		//print(jugador)

		jugarEtapa = false

		return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil

	} else if etapa == 3 {
		if jugadorMuerto == int(jugador) {
			estado = "Morir"
			fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto, este es el random")
			return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
		}
		fmt.Println(">> El jugador ", jugador, " ha escogido ", mensaje)

		index := 0
		for i := 0; i < len(jugadores); i++ {
			if jugadores[i] == int(jugador) {
				index = i
				break
			}
		}

		TodoNada[index] = int(mensaje)

		if jugadoresRespuestas == jugadoresVivos {
			toychato3 = false
		}
		for toychato3 { //Esperar que todos envien sus numeros
		}
		fmt.Println(jugadoresVivos)
		/*for jugadoresRespuestas != jugadoresVivos {

		}

		ListoTodoNada++

		for ListoTodoNada != jugadoresVivos {

		}*/

		//fmt.Println(TodoNada)

		lado := index % 2

		pareja := 0

		if lado == 0 {
			pareja = TodoNada[index+1]
		} else {
			pareja = TodoNada[index-1]
		}

		valor_jugador := math.Abs(float64(eleccion_lider - int(mensaje)))
		valor_pareja := math.Abs(float64(eleccion_lider - pareja))

		fmt.Println(valor_jugador)
		if valor_jugador > valor_pareja {
			//muere
			fmt.Println(">> ***BANG!!**  El jugador ", jugador, " ha muerto")
			jugadoresVivos--
			estado = "Morir"
			ListoTodoNada2++

		} else if valor_jugador == valor_pareja {
			fmt.Println(">> El jugador ", jugador, " ha ganado el juego del Calamar")
			estado = "Ganar"
			terminar = true
			ListoTodoNada2++
		} else {
			///gana
			fmt.Println(">> El jugador ", jugador, " ha ganado el juego del Calamar")
			estado = "Ganar"
			terminar = true
			ListoTodoNada2++
		}

		for ListoTodoNada2 != len(jugadores) {

		}

		return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
	}

	if ganar {
		//jugadoresActivos--
		jugadoresEnviar++
		return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
	}

	jugadoresEnviar++
	return &pb.Estado{Estado: estado, Ronda: int32(ronda)}, nil
}

func solicitar_jugadas(jugadores []int) {
	fmt.Println(" ")
	for i := 0; i < len(jugadores); i++ {
		fmt.Print("[", jugadores[i], "] Jugador ", jugadores[i], "\n")
	}
	fmt.Println(" ")
	fmt.Print("Selecciones que hacer: ")
	var jugador_elegido int
	fmt.Scanln(&jugador_elegido)

	//Aqui habria que comunicarse con NameNode mediante gRPC para conseguir las jugadas del jugador seleccionado
}

/*
func juego1(jugadores,inputs){
	//Se reciben los inputs de los jugadores
}
*/
func Menu() {
	for cantidadJugadores != 16 {
	}
	time.Sleep(2 * time.Second)
	jugadores = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	sumaetapa1 = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	TodoNada = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	//Idea(ignorar)
	//parejas = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	var eleccion int

	menu := true
	for menu {
		fmt.Println("\nEste es el menu del Lider")
		fmt.Println(" ")
		fmt.Println("[1] Comenzar siguiente etapa de juego")
		fmt.Println("[2] Solicitar jugadas de algun jugador")
		fmt.Println(" ")
		fmt.Print("Selecciones que hacer: ")

		fmt.Scanln(&eleccion)

		if eleccion == 1 {
			etapa += 1
			esperaEtapa = false
			fmt.Println("Comenzando el siguiente juego...")
			if etapa == 1 {
				fmt.Println(">>> Los jugadores vivos son: ", jugadores)
				eleccion_lider = etapa1()
				fmt.Println(">>>>>>>> El lider ha escogido ", eleccion_lider)

			} else if etapa == 2 {
				eleccion_lider = etapa2()
				jugarEtapa = true
				jugadoresEnviar = 0
				fmt.Println(">>>>>>>> El lider ha escogido ", eleccion_lider)
				jugadores = findAndDelete(jugadoresEtapa1, -1)
				fmt.Println(">>> Los jugadores vivos son: ", jugadores)
				if jugadoresVivos%2 == 1 {
					matar := random_range(0, jugadoresVivos-1)
					jugadorMuerto = jugadores[matar]
					jugadores = remove(jugadores, matar)
				}
				jugadoresVivosEtapa2 = len(jugadores)

				primer_equipo = jugadores[:(jugadoresVivos / 2)]
				segundo_equipo = jugadores[(jugadoresVivos / 2):]

				fmt.Println("Primer equipo: ", primer_equipo)
				fmt.Println("Segundo equipo: ", segundo_equipo)

			} else if etapa == 3 {
				fmt.Println(">>> Los jugadores vivos son: ", jugadores)
				eleccion_lider = etapa3()
				jugarEtapa = true
				fmt.Println(">>>>>>>> Comienza la etapa 3: Todo o Nada ")
				fmt.Println(">>>>>>>> El lider ha escogido ", eleccion_lider)
				if jugadoresVivos%2 == 1 {
					matar := random_range(0, jugadoresVivos-1)
					jugadorMuerto = jugadores[matar]
					jugadores = remove(jugadores, matar)
				}

				TodoNada = TodoNada[:len(jugadores)]

			}
			time.Sleep(1000 * time.Millisecond)

			if terminar {
				os.Exit(0)
			}

			for jugarEtapa { //Esperar que todos terminen de jugar?
			}

		} else if eleccion == 2 {
			solicitar_jugadas(jugadores)
		}
	}
}

func main() {
	jugadores = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	jugadoresEtapa1 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterJugadoresServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	go Menu()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//Para obtener la lista de jugadores habria que implementar gRPC para saber que jugadores se quieren unir y agregarlos a la lista

	/*
		fmt.Println("Este es el menu del Lider")
		fmt.Println(" ")
		fmt.Println("[1] Comenzar siguiente etapa de juego")
		fmt.Println("[2] Solicitar jugadas de algun jugador")
		fmt.Println(" ")
		fmt.Print("Selecciones que hacer: ")
		var eleccion int
		fmt.Scanln(&eleccion)

		if eleccion == 1 {
			fmt.Println("Comenzando el siguiente juego...")
		} else if eleccion == 2 {
			solicitar_jugadas(jugadores)
		}

	*/

}
