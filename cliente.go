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

var i int
var p MemoriaProceso

func cliente() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	request := MemoriaProceso{
		ID:     -1,
		Actual: -1,
		State:  1,
	}
	err = gob.NewEncoder(c).Encode(request)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		c.Close()
		s, err := net.Listen("tcp", ":9998")
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			c, err := s.Accept()
			if err != nil {
				fmt.Println(err)
				continue
			}
			go handleServer(c)
			break
		}
		s.Close()

	}
	c.Close()
}

func handleServer(c net.Conn) {
	err := gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(p)
	go ProcesoCliente(p.ID, p.Actual)

}

func ProcesoCliente(id int, pActual int) {
	i = pActual
	for {
		fmt.Printf("id %d: %d \n", id, i)
		i = i + 1
		time.Sleep(time.Millisecond * 500)
	}

}

func main() {
	go cliente()

	defer enviarProceso()

	var input string
	fmt.Scanln(&input)
}

func enviarProceso() {
	fmt.Println("Enviar proceso:")
	p.Actual = i
	p.State = 0
	fmt.Println(p)
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}
