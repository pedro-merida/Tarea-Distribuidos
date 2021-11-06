package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	pb "example.com/go-jugador-grpc/jugador"
	"google.golang.org/grpc"
)

/*
func random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())

	numerito := rand.Intn(max-min) + min
	return numerito
}
*/
const (
	address = "dist37:50051"
)

var rondai int = 0
var rondaf int = 1
var wg sync.WaitGroup

func jugador(i int) {
	direccion := "dist37:50051"
	conn, err := grpc.Dial(direccion, grpc.WithInsecure(), grpc.WithBlock())
	vivo := 1

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewJugadoresClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	decision := "si"
	if len(os.Args) > 1 {
		decision = os.Args[1]
	}
	r, err := c.EnviarSolicitud(ctx, &pb.Peticion{Peticion: decision}) //Cambiar name
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetMensaje())

	rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa1"})

	if errr != nil {
		log.Fatalf("could not greet: %v", errr)
	}

	etapa := rr.GetJuego()

	/*for vivo == 1 {
		if etapa == "La Etapa 1 comenzara ahora" {
			vivo = etapa1Bot(c, i)
		}
	}*/
	puto := 1
	for vivo == 1 {
		time.Sleep(200 * time.Millisecond)
		if etapa == "La Etapa 1 comenzara ahora" && rondaf > rondai && puto == 1 {
			vivo, puto = etapa1Bot(c, i)
			rondai = rondaf

			if vivo == 1 && puto == 0 {
				//Siguiente etapa
				rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa2"})
				if errr != nil {
					log.Fatalf("could not greet: %v", errr)
				}

				etapa = rr.GetJuego()
				rondaf = 1
				rondai = 0
			}
		} else if etapa == "La Etapa 2 comenzara ahora" {
			vivo, puto = etapa2Bot(c, i)

			if vivo == 1 && puto == 0 {
				//Siguiente etapa
				rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa3"})
				if errr != nil {
					log.Fatalf("could not greet: %v", errr)
				}
				etapa = rr.GetJuego()
			}

			if vivo == 1 {
				//Siguiente etapa
				rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa3"})
				if errr != nil {
					log.Fatalf("could not greet: %v", errr)
				}
				etapa = rr.GetJuego()
			}
		} else {
			vivo, puto = etapa3Bot(c, i)

		}
	}
}
func jugadoresRestantes() {
	//jugadores := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	//var ce [15]pb.jugadoresClient
	for i := 0; i < 15; i++ {
		time.Sleep(200 * time.Millisecond)
		go jugador(i + 1)
	}

}

func etapa1Bot(c pb.JugadoresClient, jugador int) (int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rand.Seed(time.Now().UnixNano() + int64(jugador))
	numero := rand.Intn(9) + 1 //Arreglar el random
	//time.Sleep(250 * time.Millisecond)
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: int32(jugador)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	puto := 1
	estado := r.GetEstado()

	if int(r.GetRonda()) == -1 {
		puto = 0
	} else {
		rondaf = int(r.GetRonda())
	}

	if estado == "Vivir" {
		return 1, puto
	} else if estado == "Morir" {
		return 0, puto
	} else {
		return 0, puto
	}
}

func etapa1(c pb.JugadoresClient) (int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fmt.Println("Elije un número del 1 al 10")
	var numero int
	puto := 1
	fmt.Scanln(&numero)
	fmt.Println("Escogiste ", numero)

	//Enviar numero al lider
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: 16})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	estado := r.GetEstado()
	if int(r.GetRonda()) == -1 {
		//fmt.Println("Gana la etapa")
		puto = 0 //Gana la etapa
	} else {
		rondaf = int(r.GetRonda())
	}
	if estado == "Vivir" {
		//fmt.Println("Has sobrevivido... por ahora")
		return 1, puto
	} else if estado == "Morir" {
		//fmt.Println("RIP")
		return 0, puto
	} else {
		return 3, puto
	}
}

func etapa2Bot(c pb.JugadoresClient, jugador int) (int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rand.Seed(time.Now().UnixNano() + int64(jugador))
	numero := rand.Intn(3) + 1 //Arreglar el random
	//time.Sleep(250 * time.Millisecond)
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: int32(jugador)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	puto := 1
	estado := r.GetEstado()

	if int(r.GetRonda()) == -1 {
		puto = 0
	} else {
		rondaf = int(r.GetRonda())
	}

	if estado == "Vivir" {
		return 1, puto
	} else if estado == "Morir" {
		return 0, puto
	} else {
		return 0, puto
	}
}

func etapa2(c pb.JugadoresClient) int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fmt.Println("Elije un número del 1 al 4")
	var numero int
	fmt.Scanln(&numero)
	fmt.Println("Escogiste ", numero)

	//Enviar numero al lider
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: 16})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	estado := r.GetEstado()
	if estado == "Vivir" {
		//fmt.Println("Has sobrevivido... por ahora")
		return 1
	} else if estado == "Morir" {
		//fmt.Println("RIP")
		return 0
	} else {
		return 3
	}
}

func etapa3Bot(c pb.JugadoresClient, jugador int) (int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	rand.Seed(time.Now().UnixNano() + int64(jugador))
	numero := rand.Intn(9) + 1 //Arreglar el random
	//time.Sleep(250 * time.Millisecond)
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: int32(jugador)})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	puto := 1
	estado := r.GetEstado()

	if estado == "Vivir" {
		return 1, puto
	} else if estado == "Morir" {
		return 0, puto
	} else {
		return 0, puto
	}
}

func etapa3(c pb.JugadoresClient) int {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	fmt.Println("Elije un número del 1 al 10")
	var numero int
	fmt.Scanln(&numero)
	fmt.Println("Escogiste ", numero)

	//Enviar numero al lider
	r, err := c.EnviarJugada(ctx, &pb.Numero{Numero: int32(numero), Jugador: 16})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	estado := r.GetEstado()
	if estado == "Vivir" {
		//fmt.Println("Has sobrevivido... por ahora")
		return 1
	} else if estado == "Morir" {
		//fmt.Println("RIP")
		return 0
	} else {
		return 3
	}
}

func main() {
	vivo := 1
	scanner := bufio.NewScanner(os.Stdin)

	/*
		fmt.Println("Bienvenido Jugador 1, como te llamas?")
		fmt.Print(">>>")
		scanner.Scan()
		fmt.Println("Te sere sincero Jugador 1, no me importa en lo más minimo cual es tu nombre")
		fmt.Print(">>>...")
		scanner.Scan()
		fmt.Println("Vayamos al punto, estas aqui porque necesitas dinero, verdad?")
		fmt.Println("No tienes que responder, todos llegan aqui por lo mismo")
		fmt.Print(">>>...")
		scanner.Scan()
		fmt.Println("Como parte de las formalidades, antes de empezar debo preguntarte....")
	*/
	fmt.Println("Deseas participar en 'El Juego del Calamar'?")
	var decision string
	fmt.Println(">>> [si]  [no]")
	scanner.Scan()
	decision = scanner.Text()
	fmt.Println("Has dicho que", decision)
	if decision == "si" {
		fmt.Println(">>>>>>>>>> El Jugador 1 va a participar en los juegos <<<<<<<<<<")
		fmt.Println("Muy bien. En unos instantes estaras participando....")
		//Enviar peticion a lider
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()

		c := pb.NewJugadoresClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		defer cancel()

		if len(os.Args) > 1 {
			decision = os.Args[1]
		}

		go jugadoresRestantes()

		r, err := c.EnviarSolicitud(ctx, &pb.Peticion{Peticion: decision}) //Cambiar name
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		log.Printf("%s", r.GetMensaje())

		time.Sleep(2 * time.Second)

		rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa1"})

		if errr != nil {
			log.Fatalf("could not greet: %v", errr)
		}

		etapa := rr.GetJuego()
		log.Printf("%s", rr.GetJuego())
		puto := 1
		for vivo == 1 {
			if etapa == "La Etapa 1 comenzara ahora" && rondaf > rondai && puto == 1 {
				fmt.Println("Comienza la ronda ", rondaf)
				//rondai = rondaf
				vivo, puto = etapa1(c)

				if vivo == 1 {
					if puto == 1 {
						fmt.Println("Has sobrevivido... por ahora")
					} else {
						fmt.Println("Has ganado la etapa")
						//Consultar por siguiente etapa

						rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa2"})
						if errr != nil {
							log.Fatalf("could not greet: %v", errr)
						}

						etapa = rr.GetJuego()
						fmt.Println(etapa)
						rondai = 0
						rondaf = 1
					}
				} else if vivo == 0 {
					fmt.Println("RIP")
				} else {
					fmt.Println("¡Has ganado el calamar!")
				}
			} else if etapa == "La Etapa 2 comenzara ahora" {
				fmt.Println(etapa)
				//rondai = rondaf
				vivo = etapa2(c)

				if vivo == 1 {
					fmt.Println("Has ganado la etapa")

					rr, errr := c.EnviarInicio(ctx, &pb.Inicio{Inicio: "Etapa3"})
					if errr != nil {
						log.Fatalf("could not greet: %v", errr)
					}

					etapa = rr.GetJuego()
					fmt.Println(etapa)
				} else if vivo == 0 {
					fmt.Println("RIP")
				} else {
					fmt.Println("¡Has ganado el calamar!")
				}
			} else {
				vivo = etapa3(c)

				if vivo == 3 {
					fmt.Println("¡Has ganado el calamar!")
					time.Sleep(time.Second)
				} else {
					fmt.Println("RIP")
				}
			}
		}
	} else {
		fmt.Println("Chale..")
		fmt.Println("Vaya vaya... despues de llegar hasta aqui te arrepientes...")
	}

}
