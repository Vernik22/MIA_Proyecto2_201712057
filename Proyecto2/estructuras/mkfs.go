package estructuras

import "fmt"

func EjecutarComandoMkfs(datFS PropMkfs, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKFS                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if existeIdMount(datFS.setId, listaDiscos) {

	} else {
		fmt.Println("ERROR: el ID de la particion no existe, o no esta montada")
	}
}

func existeIdMount(idMount string, listaDiscos *[100]Mount) bool {
	for i := 0; i < 100; i++ {
		if idMount == listaDiscos[i].IdMount {
			return true
		}
	}
	return false
}
