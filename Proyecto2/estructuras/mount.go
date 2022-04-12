package estructuras

import (
	"fmt"
	"strconv"
)

func EjecutarComandoMount(datMount PropMount, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MOUNT                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	vError := true
	for i := 0; i < 100; i++ {
		if listaDiscos[i].NombreParticion == datMount.setName && listaDiscos[i].Path == datMount.setPath {
			fmt.Println("ERROR: Ya esta montada la particion solicidada")
			vError = false
			break
		}
	}
	if vError {
		if archivoExiste(datMount.setPath) {

			if !existeNombreParticion(datMount.setPath, datMount.setName) {
				nueva := Mount{}
				copy(nueva.EstadoMks[:], "0")
				nueva.IdMount = generarId(datMount, listaDiscos)
				nueva.Path = datMount.setPath
				nueva.NombreParticion = datMount.setName
				for i := 0; i < 100; i++ {
					if listaDiscos[i].IdMount == " " {
						listaDiscos[i] = nueva
						break
					}
				}

				fmt.Println("Particion montada con exito!")
			} else {
				fmt.Println("ERROR: el nombre de la particion no existe")
			}
		} else {
			fmt.Println("ERROR: el archivo que desea montar no existe")
		}
	}
}

func generarId(datMount PropMount, listaDiscos *[100]Mount) string {
	var l byte = 65
	var n int = 1
	var idTemp string = "57"

	for i := 0; i < 100; i++ {
		if listaDiscos[i].NombreParticion != " " {
			if listaDiscos[i].Path == datMount.setPath {
				l++
			} else {
				n++
			}
		}
	}
	idTemp = idTemp + strconv.Itoa(n) + string(l)
	fmt.Println("ID de la particion montada: " + idTemp)
	return idTemp
}
