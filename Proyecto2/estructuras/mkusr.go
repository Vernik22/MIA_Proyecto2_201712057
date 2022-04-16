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

	cuenta := len(nombreUsr)
	//fmt.Println("Numero de caracteres del nombre: " + cuenta)
	if cuenta < 11 {
		cuentaPWd := len(passW)
		if cuentaPWd < 11 {
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

				existeElGrupo := false
				listaUsuarios := strings.Split(userstxt, "\n")
				for i := 0; i < len(listaUsuarios)-1; i++ {
					linea := strings.Split(listaUsuarios[i], ",")
					if linea[1] == "G" && linea[0] != "0" && linea[0] != "" {
						if linea[2] == nombreG {
							existeElGrupo = true
							break
						}

					}
				}
				if existeElGrupo {
					idUsr := 1
					existeElUsr := false
					listaUsuarios1 := strings.Split(userstxt, "\n")
					for i := 0; i < len(listaUsuarios1)-1; i++ {
						linea := strings.Split(listaUsuarios1[i], ",")
						if linea[1] == "U" && linea[0] != "0" && linea[0] != "" {
							if linea[3] == nombreUsr {
								fmt.Println("ERROR: el nombre del usuario ya existe: " + nombreUsr)
								existeElUsr = true
								break
							}
							idUsr++
						} else if linea[1] == "U" && linea[0] == "0" && linea[0] != "" {
							idUsr++
						}
					}

					if !existeElUsr {
						inodoTemp1 := Inodo{}
						f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(inodoTemp1)), 0)
						err = binary.Read(f, binary.BigEndian, &inodoTemp1)
						contador := 1
						noBloques := 0

						nuevoCont := strconv.Itoa(idUsr) + ",U," + nombreG + "," + nombreUsr + "," + passW + "\n"
						cuenta = len(nuevoCont)
						inodoTemp1.I_size += int64(cuenta)
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
							fmt.Println("Se creo el usuario exitosamente: " + nombreUsr)
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
								fmt.Println("Se creo el usuario exitosamente: " + nombreUsr)
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
								fmt.Println("Se creo el usuario exitosamente: " + nombreUsr)
								fmt.Println(userstxt)

							}

						}

					}
				} else {
					fmt.Println("ERROR: el nombre del grupo no existe: " + nombreG)
				}

			}
		} else {
			fmt.Println("ERROR: el password del usuario es mayor a 10 caracteres")
			fmt.Println(cuentaPWd)
		}

	} else {
		fmt.Println("ERROR: el nombre del usuario es mayor a 10 caracteres")
		fmt.Println(cuenta)
	}
}

//-----------------------------------------RMUSR------------------------------
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
						if linea[1] == "U" && linea[0] != "0" && linea[0] != "" {
							//fmt.Println(linea[2])
							if linea[3] == nombreUsr {
								nuevaLinea := "0,U," + linea[2] + "," + linea[3] + "," + linea[4]
								listaUsuarios[f] = nuevaLinea
								//otraCadena += listaUsuarios[f] + "\n"
								fmt.Println("Se elimino el usuario ya existe: " + nombreUsr)
								existeElGrupo = true
								//break
							}

						} else if linea[1] == "U" && linea[0] == "0" && linea[0] != "" {
							if linea[3] == nombreUsr {
								fmt.Println("ERROR: el usuario esta actualmente removido")
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
