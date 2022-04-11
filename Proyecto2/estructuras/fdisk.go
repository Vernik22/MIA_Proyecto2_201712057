package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"unsafe"
)

func EjecutarComandoFdisk(datFormat PropFdisk) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando FDISK                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	mbrTemp := MBR{}
	var tamanoParticion int64 = 0

	if datFormat.setUnit == "m" {
		tamanoParticion = int64(datFormat.setSize) * 1024 * 1024
	} else if datFormat.setUnit == "k" {
		tamanoParticion = int64(datFormat.setSize) * 1024
	} else if datFormat.setUnit == "b" {
		tamanoParticion = int64(datFormat.setSize)
	}

	//verificar que no exista el nombre en las particiones principales

	if existeNombreParticion(datFormat.setPath, datFormat.setName) {

		if datFormat.setType == "p" {
			f, err := os.OpenFile(datFormat.setPath, os.O_RDWR, 0755)
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			defer f.Close()
			if err != nil {
				fmt.Println("No existe el archivo en la ruta")
			} else {
				if string(mbrTemp.Mbr_partition_1.Part_status[:]) == "0" || string(mbrTemp.Mbr_partition_2.Part_status[:]) == "0" || string(mbrTemp.Mbr_partition_3.Part_status[:]) == "0" || string(mbrTemp.Mbr_partition_4.Part_status[:]) == "0" {
					crearParticionPrimaria(datFormat, tamanoParticion)
				}
			}

		} else if datFormat.setType == "e" {

		} else if datFormat.setType == "l" {

		}

	} else {
		fmt.Println("ERROR: el nombre de la particion ya existe")
	}
	imprimirDatosDisco(datFormat.setPath)
}

func existeNombreParticion(path string, nombre string) bool {
	mbrComprobar := MBR{}
	var parts [4]Particion
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbrComprobar)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		parts[0] = mbrComprobar.Mbr_partition_1
		parts[1] = mbrComprobar.Mbr_partition_2
		parts[2] = mbrComprobar.Mbr_partition_3
		parts[3] = mbrComprobar.Mbr_partition_4
		for i := 0; i < 4; i++ {

			if string(parts[i].Part_type[:]) == "p" {

				var1 := nombreAString(parts[i].Part_name)
				if string(var1) == string(nombre) {
					fmt.Println("*Existe " + string(parts[i].Part_name[:]))
					return false
				}

			} else if string(parts[i].Part_type[:]) == "e" {
				var1 := nombreAString(parts[i].Part_name)
				if string(var1) == string(nombre) {
					fmt.Println("*Existe " + string(parts[i].Part_name[:]))
					return false
				} else {
					//revisar las particiones logicas
					ebrTemp := EBR{}
					f.Seek(parts[i].Part_start, 0)
					err = binary.Read(f, binary.BigEndian, &ebrTemp)
					for ebrTemp.Part_next != -1 {
						var2 := nombreAString(ebrTemp.Part_name)
						if string(var2) == nombre {
							return false
						} else {
							f.Seek(ebrTemp.Part_next, 0)
							err = binary.Read(f, binary.BigEndian, &ebrTemp)
						}
					}
				}
			}
		}
		return true
	}
	return false
}

func crearParticionPrimaria(datFormat PropFdisk, tamanoPart int64) {
	mbrTemp := MBR{}
	particionTemp := Particion{}
	var startPart int64 = int64(unsafe.Sizeof(mbrTemp))

	if datFormat.setFit == "bf" {
		copy(particionTemp.Part_fit[:], "b")
	} else if datFormat.setFit == "wf" {
		copy(particionTemp.Part_fit[:], "w")
	} else if datFormat.setFit == "ff" {
		copy(particionTemp.Part_fit[:], "f")
	}

	f, err := os.OpenFile(datFormat.setPath, os.O_RDWR, 0755)
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbrTemp)
	defer f.Close()
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		var parts [4]Particion

		//ver si todas las particiones estan vacias para crear la primer particion
		if string(mbrTemp.Mbr_partition_1.Part_status[:]) == "0" && string(mbrTemp.Mbr_partition_2.Part_status[:]) == "0" && string(mbrTemp.Mbr_partition_3.Part_status[:]) == "0" && string(mbrTemp.Mbr_partition_4.Part_status[:]) == "0" {
			//seteando los valores que llevara la nueva particion
			copy(particionTemp.Part_name[:], datFormat.setName)
			copy(particionTemp.Part_type[:], "p")
			particionTemp.Part_size = tamanoPart
			particionTemp.Part_start = startPart
			copy(particionTemp.Part_status[:], "1")
			//copiando particion creada
			mbrTemp.Mbr_partition_1 = particionTemp
			//ahora modificamos el mbr original en el archivo por el nuevo mbr con los datos de la particion creada
			f.Seek(0, 0)
			err = binary.Write(f, binary.BigEndian, mbrTemp)
			fmt.Println("Se creo la particion : " + datFormat.setName + " Correctamente")
		} else {
			parts[0] = mbrTemp.Mbr_partition_1
			parts[1] = mbrTemp.Mbr_partition_2
			parts[2] = mbrTemp.Mbr_partition_3
			parts[3] = mbrTemp.Mbr_partition_4
			if string(mbrTemp.Mbr_dsk_fit[:]) == "f" {

			} else if string(mbrTemp.Mbr_dsk_fit[:]) == "w" {

			} else if string(mbrTemp.Mbr_dsk_fit[:]) == "b" {

			}
		}

	}
}

func hayEspacio(tamPart int64, tamDisco int64) bool {
	if ((tamPart) > tamDisco) || (tamPart < 0) {
		fmt.Println("ERROR ---->EL Tamanio de la particion es mayor a el tamanio del disco o el tamanio es incorrecto")
		return false
	}
	return true
}

func nombreAString(data [16]byte) string {
	var var1 string = ""
	for i := 0; i < 16; i++ {
		if data[i] != 0 {
			var1 += string(data[i])
		}
	}
	return var1
}

func imprimirDatosDisco(path string) {
	mbrTemp := MBR{}
	f, err := os.OpenFile(path, os.O_RDWR, 0755)
	f.Seek(0, 0)
	err = binary.Read(f, binary.BigEndian, &mbrTemp)
	defer f.Close()
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		var parts [4]Particion
		parts[0] = mbrTemp.Mbr_partition_1
		parts[1] = mbrTemp.Mbr_partition_2
		parts[2] = mbrTemp.Mbr_partition_3
		parts[3] = mbrTemp.Mbr_partition_4
		fmt.Println("\n----------DATOS DEL DISCO-----")
		fmt.Println("")
		fmt.Println("Disk Name: " + strconv.Itoa(int(mbrTemp.Mbr_dsk_signature)))
		fmt.Println("Disk Size: " + strconv.Itoa(int(mbrTemp.Mbr_tamano)))
		fmt.Println("Disk Date: " + string(mbrTemp.Mbr_fecha_creacion[:]))
		for i := 0; i < 4; i++ {
			fmt.Println("")
			fmt.Println("PARTICION: " + strconv.Itoa(i+1))
			fmt.Println("Particion Status: " + string(parts[i].Part_status[:]))
			fmt.Println("Particion Type: " + string(parts[i].Part_type[:]))
			fmt.Println("Particion Fit: " + string(parts[i].Part_fit[:]))
			fmt.Println("Particion Start: " + strconv.Itoa(int(parts[i].Part_start)))
			fmt.Println("Particion Size: " + strconv.Itoa(int(parts[i].Part_size)))
			fmt.Println("Particion Name: " + string(parts[i].Part_name[:]))

			if string(parts[i].Part_type[:]) == "e" {
				fmt.Println("")
				fmt.Println("--------------Particiones Logicas----------------")
				ebrTemp := EBR{}
				f.Seek(parts[i].Part_start, 0)
				err = binary.Read(f, binary.BigEndian, &ebrTemp)
				for ebrTemp.Part_next != -1 {
					fmt.Println("Particion Status: " + string(ebrTemp.Part_status[:]))
					fmt.Println("Particion Next: " + strconv.Itoa(int(ebrTemp.Part_next)))
					fmt.Println("Particion Fit: " + string(ebrTemp.Part_fit[:]))
					fmt.Println("Particion Start: " + strconv.Itoa(int(ebrTemp.Part_start)))
					fmt.Println("Particion Size: " + strconv.Itoa(int(ebrTemp.Part_size)))
					fmt.Println("Particion Name: " + string(ebrTemp.Part_name[:]))

					f.Seek(ebrTemp.Part_next, 0)
					err = binary.Read(f, binary.BigEndian, &ebrTemp)

				}

				fmt.Println("--------------FIN Particiones Logicas----------------")
				fmt.Println("")

			}
		}
		fmt.Println("")
	}
}
