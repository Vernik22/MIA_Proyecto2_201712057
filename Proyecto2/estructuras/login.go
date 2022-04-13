package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"unsafe"
)

func EjecutarComandoLogin(datLog PropLogin, listaDiscos *[100]Mount) bool {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando LOGIN                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if existeIdMount(datLog.setId, listaDiscos) {
		mbrTemp := MBR{}
		pathD := ""
		nombreParticion := ""
		for i := 0; i < 100; i++ {
			if datLog.setId == listaDiscos[i].IdMount {
				pathD = listaDiscos[i].Path
				nombreParticion = listaDiscos[i].NombreParticion
				break
			}
		}
		f, err := os.OpenFile(pathD, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe el archivo en la ruta")
		} else {
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			_, iniPart := returnDatosPart(mbrTemp, pathD, nombreParticion)
			superBlock := SupB{}
			f.Seek(iniPart, 0)
			err = binary.Read(f, binary.BigEndian, &superBlock)
			inodoUsers := Inodo{}
			f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(inodoUsers)), 0)
			err = binary.Read(f, binary.BigEndian, &inodoUsers)
			texto := BArchivo{}
			var userstxt string = ""
			for i := 0; i < 16; i++ {
				if inodoUsers.I_block[i] != -1 {
					f.Seek(superBlock.S_block_start+(inodoUsers.I_block[i]*int64(unsafe.Sizeof(BArchivo{}))), 0)
					err = binary.Read(f, binary.BigEndian, &texto)
					userstxt += string(texto.B_content[:])
				}
			}

			listaUsuarios := strings.Split(userstxt, "\n")
			for i := 0; i < len(listaUsuarios)-1; i++ {
				linea := strings.Split(listaUsuarios[i], ",")
				if linea[1] == "U" && linea[0] != "0" {
					if linea[3] == datLog.setUsuario && linea[4] == datLog.setPassword {
						fmt.Println("Bienvenido al sistema: " + datLog.setUsuario)
						return true
					}
				}
			}
		}
	} else {
		fmt.Println("ERROR: el ID de la particion no existe, o no esta montada")
	}
	return false
}
