package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func EjecutarComandoMkgrp(datGrp PropMkgrp, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKGRP                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if datGrp.setUsuarioAct == "root" {
		pathDisco := ""
		nombrePart := ""
		for i := 0; i < 100; i++ {
			if datGrp.setIdMontada == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				nombrePart = listaDiscos[i].NombreParticion
				break
			}
		}
		modificarArchivo(pathDisco, nombrePart, datGrp.setName)

	} else {
		fmt.Println("Usuario incorrecto, no es el usuario root")
	}
}

func modificarArchivo(pathDisco string, nombrePart string, nombreG string) {
	cuenta := len(nombreG)
	//fmt.Println("Numero de caracteres del nombre: " + cuenta)
	if cuenta < 11 {
		mbrTemp := MBR{}
		f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe el archivo en la ruta")
		} else {
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			_, iniPart := returnDatosPart(mbrTemp, pathDisco, nombrePart)
			superBlock := SupB{}
			f.Seek(iniPart, 0)
			err = binary.Read(f, binary.BigEndian, &superBlock)

			inodoTemp := Inodo{}
			f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
			err = binary.Read(f, binary.BigEndian, &inodoTemp)

			texto := BArchivo{}
			userstxt := ""
			for i := 0; i < 16; i++ {
				if inodoTemp.I_block[i] != -1 {
					f.Seek(superBlock.S_block_start+(inodoTemp.I_block[i]*int64(unsafe.Sizeof(texto))), 0)
					err = binary.Read(f, binary.BigEndian, &texto)
					userstxt += archivoAString(texto.B_content)
				}
			}

			idGrupo := 1
			existeElGrupo := false
			listaUsuarios := strings.Split(userstxt, "\n")
			//fmt.Println(strings.Replace(userstxt, "\n", "?", -1))
			for i := 0; i < len(listaUsuarios)-1; i++ {
				linea := strings.Split(listaUsuarios[i], ",")
				if linea[1] == "G" && linea[0] != "0" && linea[0] != "" {
					//fmt.Println(linea[2])
					if linea[2] == nombreG {
						fmt.Println("ERROR: el nombre del grupo ya existe: " + nombreG)
						existeElGrupo = true
						break
					}
					idGrupo++
				} else if linea[1] == "G" && linea[0] == "0" && linea[0] != "" {
					idGrupo++
				}
			}
			if !existeElGrupo {

				inodoTemp1 := Inodo{}
				f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(inodoTemp1)), 0)
				err = binary.Read(f, binary.BigEndian, &inodoTemp1)
				contador := 1
				noBloques := 0

				cuenta += 5
				inodoTemp1.I_size += int64(cuenta)
				nuevoCont := strconv.Itoa(idGrupo) + ",G," + nombreG + "\n"
				userstxt = userstxt + nuevoCont
				for i := 0; i < int(inodoTemp1.I_size); i++ {
					if contador == 63 {
						noBloques += 1
						contador = 0
					}
					contador++
				}
				if inodoTemp1.I_size%63 != 0 {
					noBloques += 1
				}

				if noBloques == 1 {

					copy(texto.B_content[:], userstxt)
					f.Seek(superBlock.S_block_start+(inodoTemp.I_block[0]*int64(unsafe.Sizeof(texto))), 0)
					err = binary.Write(f, binary.BigEndian, texto)

					dt := time.Now()
					copy(inodoTemp1.I_mtime[:], dt.String())

					f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Write(f, binary.BigEndian, &inodoTemp1)
					fmt.Println("Se creo el grupo exitosamente: " + nombreG)
					fmt.Println(userstxt)
				} else {
					textoNuevo := BArchivo{}
					textEnBlock := ""
					if inodoTemp.I_block[noBloques-1] == -1 {
						//fmt.Println("SE creo otro bloque")
						copy(textoNuevo.B_content[:], nuevoCont)
						dt := time.Now()
						copy(inodoTemp1.I_mtime[:], dt.String())
						inodoTemp1.I_block[noBloques-1] = superBlock.S_first_blo

						f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
						err = binary.Write(f, binary.BigEndian, &inodoTemp1)

						f.Seek(superBlock.S_block_start+(inodoTemp1.I_block[noBloques-1]*int64(unsafe.Sizeof(texto))), 0)
						err = binary.Write(f, binary.BigEndian, &textoNuevo)

						//agregar un 1 al bitmap de bloques
						var otro byte = '1'
						f.Seek(superBlock.S_bm_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(otro)), 0)
						err = binary.Write(f, binary.BigEndian, otro)
						//actualizar superbloque
						superBlock.S_first_blo = superBlock.S_first_blo + 1
						superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1
						f.Seek(iniPart, 0)
						err = binary.Write(f, binary.BigEndian, superBlock)
						fmt.Println("Se creo el grupo exitosamente: " + nombreG)
						fmt.Println(userstxt)

					} else {
						//fmt.Println("SE modifico el otro bloque")
						f.Seek(superBlock.S_block_start+(inodoTemp.I_block[noBloques-1]*int64(unsafe.Sizeof(texto))), 0)
						err = binary.Read(f, binary.BigEndian, &textoNuevo)

						textEnBlock = archivoAString(textoNuevo.B_content)
						textEnBlock = textEnBlock + nuevoCont
						copy(textoNuevo.B_content[:], textEnBlock)
						f.Seek(superBlock.S_block_start+(inodoTemp.I_block[noBloques-1]*int64(unsafe.Sizeof(texto))), 0)
						err = binary.Write(f, binary.BigEndian, textoNuevo)

						dt := time.Now()
						copy(inodoTemp1.I_mtime[:], dt.String())

						f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
						err = binary.Write(f, binary.BigEndian, &inodoTemp1)
						fmt.Println("Se creo el grupo exitosamente: " + nombreG)
						fmt.Println(userstxt)

					}

				}

			}

		}
	} else {
		fmt.Println("ERROR: el nombre del grupo es mayor a 10 caracteres")
		fmt.Println(cuenta)
	}
}

func archivoAString(data [64]byte) string {
	var var1 string = ""
	for i := 0; i < 64; i++ {
		if data[i] != 0 {
			var1 += string(data[i])
		}
	}
	return var1
}

//-----------------------------------RMGRP----------------
func EjecutarComandoRmgrp(datGrp PropMkgrp, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando RMGRP                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if datGrp.setUsuarioAct == "root" {
		pathDisco := ""
		nombrePart := ""
		for i := 0; i < 100; i++ {
			if datGrp.setIdMontada == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				nombrePart = listaDiscos[i].NombreParticion
				break
			}
		}
		modificarArchivoRM(pathDisco, nombrePart, datGrp.setName)

	} else {
		fmt.Println("Usuario incorrecto, no es el usuario root")
	}
}

func modificarArchivoRM(pathDisco string, nombrePart string, nombreG string) {
	mbrTemp := MBR{}
	f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		f.Seek(0, 0)
		err = binary.Read(f, binary.BigEndian, &mbrTemp)
		_, iniPart := returnDatosPart(mbrTemp, pathDisco, nombrePart)
		superBlock := SupB{}
		f.Seek(iniPart, 0)
		err = binary.Read(f, binary.BigEndian, &superBlock)

		inodoTemp := Inodo{}
		f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
		err = binary.Read(f, binary.BigEndian, &inodoTemp)

		texto := BArchivo{}
		userstxt := ""
		usuarios := ""
		otraCadena := ""
		existeElGrupo := false
		for i := 0; i < 16; i++ {
			if inodoTemp.I_block[i] != -1 {
				f.Seek(superBlock.S_block_start+(inodoTemp.I_block[i]*int64(unsafe.Sizeof(texto))), 0)
				err = binary.Read(f, binary.BigEndian, &texto)
				userstxt = archivoAString(texto.B_content)
				otraCadena = ""
				listaUsuarios := strings.Split(userstxt, "\n")
				//fmt.Println(strings.Replace(userstxt, "\n", "?", -1))
				for f := 0; f < len(listaUsuarios)-1; f++ {
					if !existeElGrupo {
						linea := strings.Split(listaUsuarios[f], ",")
						if linea[1] == "G" && linea[0] != "0" && linea[0] != "" {
							//fmt.Println(linea[2])
							if linea[2] == nombreG {
								nuevaLinea := "0,G," + linea[2]
								listaUsuarios[f] = nuevaLinea
								//otraCadena += listaUsuarios[f] + "\n"
								fmt.Println("Se elimino el grupo ya existe: " + nombreG)
								existeElGrupo = true
								//break
							}

						} else if linea[1] == "G" && linea[0] == "0" && linea[0] != "" {
							if linea[2] == nombreG {
								fmt.Println("ERROR: el grupo esta actualmente removido")
							}
						}
					}
					otraCadena += listaUsuarios[f] + "\n"
					usuarios += listaUsuarios[f] + "\n"
				}
				if existeElGrupo {
					copy(texto.B_content[:], otraCadena)
					f.Seek(superBlock.S_block_start+(inodoTemp.I_block[i]*int64(unsafe.Sizeof(texto))), 0)
					err = binary.Write(f, binary.BigEndian, &texto)
					fmt.Println(usuarios)
					break
				}

			}

		}
	}
}
