package estructuras

import (
	"container/list"
	"fmt"
	"strings"
)

func LeerTexto(data string, listaDiscos *[100]Mount) {
	fmt.Println("desde LeerTexto: " + data)
	fmt.Println(string(listaDiscos[0].Estado[:]))
	//para leer la cadena enviada
	ListaComandos := list.New()
	lineaComando := strings.Split(data, "\n")
	var com Comando
	for i := 0; i < len(lineaComando); i++ {
		EsComentario := lineaComando[i][0:1]
		if EsComentario != "#" {
			comando := lineaComando[i]
			//ahora lo separo por espacios ejemplo: mkdisk -path -size
			propiedades := strings.Split(string(comando), " ")
			nombreComando := propiedades[0]
			com.Nombre = strings.ToLower(nombreComando)
			propiedadesTemp := make([]Propiedad, len(propiedades)-1)
			for f := 0; f < len(propiedadesTemp); f++ {
				propiedadesTemp[f].Nombre = "|"
			}
			for j := 1; j < len(propiedades); j++ {
				if propiedades[j] == "" || propiedades[j] == " " || propiedades[j] == "#" {
					continue
				} else {
					if strings.Contains(propiedades[j], "=") {
						if strings.Contains(propiedades[j], "#") {
							quitComen := strings.Split(propiedades[j], "#")
							valor_propiedad_Comando := strings.Split(quitComen[0], "=")
							propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
							propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
						} else if strings.Contains(propiedades[j], "\"") {
							valor_propiedad_Comando := strings.Split(propiedades[0], "=")
							propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
							propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
							for f := j + 1; f < len(propiedades); f++ {
								if strings.Contains(propiedades[f], "\"") {
									propiedadesTemp[j-1].Valor += " " + propiedades[f]
									break
								} else {
									propiedadesTemp[j-1].Valor += " " + propiedades[f]
								}
							}
						} else {
							valor_propiedad_Comando := strings.Split(propiedades[j], "=")
							propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
							propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
						}

					} else if propiedades[j] == "-r" || propiedades[j] == "-R" {
						propiedadesTemp[j-1].Nombre = propiedades[j]
					} else if propiedades[j] == "-p" || propiedades[j] == "-P" {
						propiedadesTemp[j-1].Nombre = propiedades[j]
					}
				}
			}
			com.Propiedades = propiedadesTemp
			//agregando el comando a la lista de comandos
			ListaComandos.PushBack(com)
		} else {
			fmt.Println("Es un comentario")
		}

	}

	listaComandosValidos(ListaComandos, listaDiscos)
}

func listaComandosValidos(ListaComandos *list.List, listaDiscos *[100]Mount) {

	c := Mount{}
	copy(c.Estado[:], "2")
	listaDiscos[0] = c

	fmt.Println("ejecutando comandos validos")
	fmt.Println(string(listaDiscos[0].Estado[:]))
	fmt.Println(string(listaDiscos[1].Estado[:]))
}
