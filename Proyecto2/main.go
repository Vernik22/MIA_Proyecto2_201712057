package main

import (
	L "MIA_P2/estructuras"
	"bufio"
	"fmt"
	"os"
)

func main() {
	var ListaDiscos [100]L.Mount
	LlenarListaDisco(&ListaDiscos)
	//LeerArchivo de Entrada
	var comando string = ""
	scanner := bufio.NewScanner(os.Stdin)
	for comando != "salir" {
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Println("                            Ingrese un comando                         ")
		fmt.Println("--------------------------------------------------------------------------------")
		fmt.Print(">> ")
		scanner.Scan()
		comando = scanner.Text()
		if comando != "" && comando != "salir" {
			L.LeerTexto(comando, &ListaDiscos)
		}

	}
}

func LlenarListaDisco(montadas *[100]L.Mount) {
	c := L.Mount{}
	c.NombreParticion = " "
	c.IdMount = " "
	copy(c.Estado[:], "0")
	copy(c.EstadoMks[:], "0")
	for i := 0; i < 100; i++ {
		montadas[i] = c
	}

}
