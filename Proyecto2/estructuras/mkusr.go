package estructuras

import "fmt"

func EjecutarComandoMkusr(datUsr PropMkusr, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKUSR                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	if datUsr.setUsuarioAct == "root" {
		pathDisco := ""
		nombrePart := ""
		for i := 0; i < 100; i++ {
			if datUsr.setIdMontada == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				nombrePart = listaDiscos[i].NombreParticion
				break
			}
		}
		modificarArchivoUsr(pathDisco, nombrePart, datUsr.setUsuario, datUsr.setPassword, datUsr.setGrp)

	} else {
		fmt.Println("Usuario incorrecto, no es el usuario root")
	}
}
func modificarArchivoUsr(pathDisco string, nombrePart string, nombreUsr string, passW string, nombreG string) {

}

func EjecutarComandoRmusr(datUsr PropMkgrp, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando RMUSR                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	if datUsr.setUsuarioAct == "root" {
		pathDisco := ""
		nombrePart := ""
		for i := 0; i < 100; i++ {
			if datUsr.setIdMontada == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				nombrePart = listaDiscos[i].NombreParticion
				break
			}
		}
		modificarArchivoRmUsr(pathDisco, nombrePart, datUsr.setName)

	} else {
		fmt.Println("Usuario incorrecto, no es el usuario root")
	}
}

func modificarArchivoRmUsr(pathDisco string, nombrePart string, nombreUsr string) {

}
