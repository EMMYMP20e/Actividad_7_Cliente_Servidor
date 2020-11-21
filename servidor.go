package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type MemoriaProceso struct {
	ID     int
	Actual int
	State  int
}

var envioProceso MemoriaProceso
var indexes []int

//defer

func servidor() {

	cPrint := make(chan int)
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < 5; i++ {
		indexes = append(indexes, i)
		m := MemoriaProceso{
			ID:     i,
			Actual: 0,
			State:  0,
		}
		go Proceso(m, cPrint)
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleClient(c, cPrint)
	}
	s.Close()
}

func Proceso(m MemoriaProceso, cPrint chan int) {
	i := m.Actual
	seguir := true
	for {
		select {
		case msg := <-cPrint:
			if msg == m.ID {
				seguir = false
			}
		default:
		}
		if seguir {
			fmt.Printf("id %d: %d \n", m.ID, i)
			i = i + 1
			time.Sleep(time.Millisecond * 500)
		} else {
			envioProceso = MemoriaProceso{
				ID:     m.ID,
				Actual: i,
				State:  0,
			}
			return
		}
	}
}

func handleClient(c net.Conn, cPrint chan int) {
	var request MemoriaProceso
	err := gob.NewDecoder(c).Decode(&request)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(request)
	if request.State == 1 {
		fmt.Println("Peticion de Envío de proceso")
		if len(indexes) < 1 {
			fmt.Println("No hay procesos disponibles")
			return
		}
		i := indexes[0]
		indexes[0] = indexes[len(indexes)-1]
		indexes = indexes[:len(indexes)-1]
		for j := 0; j < len(indexes); j++ {
			cPrint <- i
		}
		time.Sleep(time.Millisecond * 750)
		c, err := net.Dial("tcp", ":9998")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = gob.NewEncoder(c).Encode(envioProceso)
		if err != nil {
			fmt.Println(err)
		}
		c.Close()
	} else {
		request.Actual++ //por delay
		fmt.Println("Proceso de vuelta ID:", request.ID)
		go Proceso(request, cPrint)
	}

}

func main() {
	go servidor()

	var input string
	fmt.Scanln(&input)
}
