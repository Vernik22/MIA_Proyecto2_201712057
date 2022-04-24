package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
	"time"
	"unsafe"
)

func EjecutarComandoMkdir(datDir PropMkdir, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKDIR                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	pathDisco := ""
	nombrePart := ""
	for i := 0; i < 100; i++ {
		if datDir.setIdMontada == listaDiscos[i].IdMount {
			pathDisco = listaDiscos[i].Path
			nombrePart = listaDiscos[i].NombreParticion
			break
		}
	}
	modificarArchivoDir(pathDisco, nombrePart, datDir)
}

func modificarArchivoDir(pathDisco string, nombrePart string, datDir PropMkdir) {

	rutaArchivo := strings.Split(datDir.setPath, "/")
	superBloque := SupB{}

	f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		mbrTemp := MBR{}
		err = binary.Read(f, binary.BigEndian, &mbrTemp)
		_, iniPart := returnDatosPart(mbrTemp, pathDisco, nombrePart)
		f.Seek(iniPart, 0)
		err = binary.Read(f, binary.BigEndian, &superBloque)

		inodoTemp := Inodo{}
		//inodoAnterior := Inodo{} +int64(unsafe.Sizeof(Inodo{}))

		f.Seek(superBloque.S_inode_start, 0)
		err = binary.Read(f, binary.BigEndian, &inodoTemp)

		//inodoAnterior = inodoTemp
		existeCarpeta := false
		carp := len(rutaArchivo)

		for i := 1; i < carp; i++ {
			existeCarpeta = false

			if string(inodoTemp.I_type[:]) == "0" {
				for j := 0; j < 16; j++ {
					if inodoTemp.I_block[j] != -1 {
						carpetaComprobar := BCarpeta{}
						f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

						for k := 0; k < 4; k++ {
							var1 := nombreAStringFile(carpetaComprobar.B_content[k].B_name)
							if rutaArchivo[i] == var1 {
								//inodoAnterior = inodoTemp
								f.Seek(superBloque.S_inode_start+int64(carpetaComprobar.B_content[k].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
								err = binary.Read(f, binary.BigEndian, &inodoTemp)
								existeCarpeta = true
								j = 20
								break
							}
						}
					}
				}
			}

			if !existeCarpeta {
				if i == carp-1 {
					f.Seek(iniPart, 0)
					err = binary.Read(f, binary.BigEndian, &superBloque)

					inodoArchivoNuevo := Inodo{}
					for n := 0; n < 16; n++ {
						inodoArchivoNuevo.I_block[n] = -1
					}

					inodoArchivoNuevo.I_gid = 1
					inodoArchivoNuevo.I_uid = 1
					copy(inodoArchivoNuevo.I_type[:], "0")
					inodoArchivoNuevo.I_perm = 664
					dt := time.Now()
					copy(inodoArchivoNuevo.I_mtime[:], dt.String())
					copy(inodoArchivoNuevo.I_ctime[:], dt.String())
					copy(inodoArchivoNuevo.I_atime[:], dt.String())
					inodoArchivoNuevo.I_block[0] = superBloque.S_first_blo

					carpetaRaiz := BCarpeta{}

					copy(carpetaRaiz.B_content[0].B_name[:], ".")
					copy(carpetaRaiz.B_content[1].B_name[:], "..")
					copy(carpetaRaiz.B_content[2].B_name[:], " ")
					copy(carpetaRaiz.B_content[3].B_name[:], " ")
					carpetaRaiz.B_content[0].B_inodo = int32(superBloque.S_first_ino)
					carpetaComprobar := BCarpeta{}
					f.Seek(superBloque.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
					err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

					carpetaRaiz.B_content[1].B_inodo = carpetaComprobar.B_content[0].B_inodo
					carpetaRaiz.B_content[2].B_inodo = -1
					carpetaRaiz.B_content[3].B_inodo = -1

					f.Seek(superBloque.S_block_start+superBloque.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
					err = binary.Write(f, binary.BigEndian, carpetaRaiz)

					var otro byte = '1'
					var actual byte = '0'
					for j := 0; j < int(superBloque.S_blocks_count); j++ {

						f.Seek(superBloque.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
						err = binary.Read(f, binary.BigEndian, &actual)

						if string(actual) == "0" {
							f.Seek(superBloque.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Write(f, binary.BigEndian, otro)
							break
						}
					}

					superBloque.S_first_blo = superBloque.S_first_blo + 1
					superBloque.S_free_blocks_count = superBloque.S_free_blocks_count - 1

					for j := 0; j < 16; j++ {
						if inodoTemp.I_block[j] != -1 {
							carpetaComprobar1 := BCarpeta{}
							f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
							err = binary.Read(f, binary.BigEndian, &carpetaComprobar1)

							for k := 0; k < 4; k++ {
								if carpetaComprobar1.B_content[k].B_inodo == -1 {
									copy(carpetaComprobar1.B_content[k].B_name[:], rutaArchivo[i])
									carpetaComprobar1.B_content[k].B_inodo = int32(superBloque.S_first_ino)
									f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
									err = binary.Write(f, binary.BigEndian, carpetaComprobar1)

									j = 20
									break
								}
							}
						} else {
							//crear nuevo bloque carpeta
							carpetaComprobar1 := BCarpeta{}
							for k := 0; k < 4; k++ {
								carpetaComprobar1.B_content[k].B_inodo = -1
								copy(carpetaComprobar1.B_content[k].B_name[:], " ")
							}
							copy(carpetaComprobar1.B_content[0].B_name[:], rutaArchivo[i])
							carpetaComprobar1.B_content[0].B_inodo = int32(superBloque.S_first_ino)
							f.Seek(superBloque.S_block_start+superBloque.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
							err = binary.Write(f, binary.BigEndian, carpetaComprobar1)

							inodoTemp.I_block[j] = superBloque.S_first_blo
							var otro byte = '1'
							var actual byte = '0'
							for h := 0; h < int(superBloque.S_blocks_count); h++ {

								f.Seek(superBloque.S_bm_block_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Read(f, binary.BigEndian, &actual)

								if string(actual) == "0" {
									f.Seek(superBloque.S_bm_block_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Write(f, binary.BigEndian, otro)
									break
								}
							}

							superBloque.S_first_blo = superBloque.S_first_blo + 1
							superBloque.S_free_blocks_count = superBloque.S_free_blocks_count - 1

							f.Seek(superBloque.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
							err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

							dt := time.Now()
							copy(inodoTemp.I_mtime[:], dt.String())

							f.Seek(superBloque.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
							err = binary.Write(f, binary.BigEndian, inodoTemp)
							break

						}
					}
					f.Seek(superBloque.S_inode_start+superBloque.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)

					otro = '1'
					for h := 0; h < int(superBloque.S_inodes_count); h++ {

						f.Seek(superBloque.S_bm_inode_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
						err = binary.Read(f, binary.BigEndian, &actual)

						if string(actual) == "0" {
							f.Seek(superBloque.S_bm_inode_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Write(f, binary.BigEndian, otro)
							break
						}
					}

					superBloque.S_first_ino = superBloque.S_first_ino + 1
					superBloque.S_free_inodes_count = superBloque.S_free_inodes_count - 1

					f.Seek(iniPart, 0)
					err = binary.Write(f, binary.BigEndian, superBloque)

					inodoTemp = inodoArchivoNuevo
					fmt.Println("Se creo la carpeta: " + rutaArchivo[i])
					break
				} else {
					if datDir.setP {
						inodoArchivoNuevo := Inodo{}
						for k := 0; k < 16; k++ {
							inodoArchivoNuevo.I_block[k] = -1
						}

						inodoArchivoNuevo.I_gid = 1
						inodoArchivoNuevo.I_uid = 1
						copy(inodoArchivoNuevo.I_type[:], "0")
						inodoArchivoNuevo.I_perm = 664
						dt := time.Now()
						copy(inodoArchivoNuevo.I_mtime[:], dt.String())
						copy(inodoArchivoNuevo.I_ctime[:], dt.String())
						copy(inodoArchivoNuevo.I_atime[:], dt.String())
						inodoArchivoNuevo.I_block[0] = superBloque.S_first_blo

						carpetaRaiz := BCarpeta{}

						copy(carpetaRaiz.B_content[0].B_name[:], ".")
						copy(carpetaRaiz.B_content[1].B_name[:], "..")
						copy(carpetaRaiz.B_content[2].B_name[:], " ")
						copy(carpetaRaiz.B_content[3].B_name[:], " ")
						carpetaRaiz.B_content[0].B_inodo = int32(superBloque.S_first_ino)
						carpetaComprobar := BCarpeta{}
						f.Seek(superBloque.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

						carpetaRaiz.B_content[1].B_inodo = carpetaComprobar.B_content[0].B_inodo
						carpetaRaiz.B_content[2].B_inodo = -1
						carpetaRaiz.B_content[3].B_inodo = -1

						f.Seek(superBloque.S_block_start+superBloque.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Write(f, binary.BigEndian, carpetaRaiz)

						superBloque.S_first_blo = superBloque.S_first_blo + 1
						superBloque.S_free_blocks_count = superBloque.S_free_blocks_count - 1

						for j := 0; j < 16; j++ {
							if inodoTemp.I_block[j] != -1 {
								carpetaComprobar1 := BCarpeta{}
								f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &carpetaComprobar1)

								for k := 0; k < 4; k++ {
									if carpetaComprobar1.B_content[k].B_inodo == -1 {
										copy(carpetaComprobar1.B_content[k].B_name[:], rutaArchivo[i])
										carpetaComprobar1.B_content[k].B_inodo = int32(superBloque.S_first_ino)
										f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Write(f, binary.BigEndian, carpetaComprobar1)
										j = 20
										break
									}
								}
							} else {
								//crear nuevo bloque carpeta
								carpetaComprobar1 := BCarpeta{}
								for k := 0; k < 4; k++ {
									copy(carpetaComprobar1.B_content[k].B_name[:], " ")
									carpetaComprobar1.B_content[k].B_inodo = -1
								}
								copy(carpetaComprobar1.B_content[0].B_name[:], rutaArchivo[i])
								carpetaComprobar1.B_content[0].B_inodo = int32(superBloque.S_first_ino)

								f.Seek(superBloque.S_block_start+superBloque.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Write(f, binary.BigEndian, carpetaComprobar1)

								inodoTemp.I_block[j] = superBloque.S_first_blo

								var otro byte = '1'
								var actual byte = '0'
								for h := 0; h < int(superBloque.S_blocks_count); h++ {

									f.Seek(superBloque.S_bm_block_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Read(f, binary.BigEndian, &actual)

									if string(actual) == "0" {
										f.Seek(superBloque.S_bm_block_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Write(f, binary.BigEndian, otro)
										break
									}
								}

								superBloque.S_first_blo = superBloque.S_first_blo + 1
								superBloque.S_free_blocks_count = superBloque.S_free_blocks_count - 1

								inodoTemp.I_mtime = inodoArchivoNuevo.I_mtime

								f.Seek(superBloque.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
								err = binary.Write(f, binary.BigEndian, inodoTemp)
								break

							}
						}

						f.Seek(superBloque.S_inode_start+superBloque.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
						err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)

						superBloque.S_first_ino = superBloque.S_first_ino + 1
						superBloque.S_free_inodes_count = superBloque.S_free_inodes_count - 1

						f.Seek(iniPart, 0)
						err = binary.Write(f, binary.BigEndian, superBloque)

						var otro byte = '1'
						var actual byte = '0'
						for h := 0; h < int(superBloque.S_inodes_count); h++ {

							f.Seek(superBloque.S_bm_inode_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Read(f, binary.BigEndian, &actual)

							if string(actual) == "0" {
								f.Seek(superBloque.S_bm_inode_start+int64(h)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Write(f, binary.BigEndian, otro)
								break
							}
						}
						inodoTemp = inodoArchivoNuevo

						fmt.Println("Se creo la carpeta: " + rutaArchivo[i])

					} else {
						fmt.Println("ERROR: no se encontro la ruta, parametro p no ingresado")
						break
					}
				}
			}
		}
	}
}
