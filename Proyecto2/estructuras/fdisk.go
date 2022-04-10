package estructuras

import (
	"fmt"
	"unsafe"
)

func EjecutarComandoFdisk(datFormat PropFdisk) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando FDISK                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	mbrTemp := MBR{}
	particionTemp := Particion{}
	var startPart int64 = int64(unsafe.Sizeof(mbrTemp))
	var tamanoParticion int64 = 0

	if datFormat.setUnit == "m" {
		tamanoParticion = int64(datFormat.setSize) * 1024 * 1024
	} else if datFormat.setUnit == "k" {
		tamanoParticion = int64(datFormat.setSize) * 1024
	} else if datFormat.setUnit == "b" {
		tamanoParticion = int64(datFormat.setSize)
	}

	if datFormat.setFit == "bf" {
		copy(particionTemp.part_fit[:], "b")
	} else if datFormat.setFit == "wf" {
		copy(particionTemp.part_fit[:], "w")
	} else if datFormat.setFit == "ff" {
		copy(particionTemp.part_fit[:], "f")
	}
	//verificar que no exista el nombre en las particiones principales
	var parts [4]Particion
	if existeNombreParticion(datFormat.setPath) {
		if datFormat.setType == "p" {

		} else if datFormat.setType == "e" {

		} else if datFormat.setType == "l" {

		}

	} else {
		fmt.Println("ERROR: el nombre de la particion ya existe")
	}

}

func existeNombreParticion(path string) bool {
	return true
}
