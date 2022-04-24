package estructuras

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func EjecutarComandoMkfile(datFile PropMkfile, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKFILE                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	pathDisco := ""
	nombrePart := ""
	for i := 0; i < 100; i++ {
		if datFile.setIdMontada == listaDiscos[i].IdMount {
			pathDisco = listaDiscos[i].Path
			nombrePart = listaDiscos[i].NombreParticion
			break
		}
	}
	modificarArchivoFile(pathDisco, nombrePart, datFile)

}

func modificarArchivoFile(pathDisco string, nombrePart string, datFile PropMkfile) {
	rutaArchivo := strings.Split(datFile.setPath, "/")
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
		//inodoAnterior := Inodo{} +int64(unsafe.Sizeof(Inodo{}))
		f.Seek(superBlock.S_inode_start, 0)
		err = binary.Read(f, binary.BigEndian, &inodoTemp)

		existeCarpeta := false
		carp := len(rutaArchivo)

		for i := 1; i < carp; i++ {
			existeCarpeta = false
			if string(inodoTemp.I_type[:]) == "0" {
				for j := 0; j < 16; j++ {
					if inodoTemp.I_block[j] != -1 {
						carpetaComprobar := BCarpeta{}
						f.Seek(superBlock.S_block_start+(inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{}))), 0)
						err = binary.Read(f, binary.BigEndian, &carpetaComprobar)
						for k := 0; k < 4; k++ {
							var1 := nombreAStringFile(carpetaComprobar.B_content[k].B_name)
							if rutaArchivo[i] == var1 {
								//inodoAnterior = inodoTemp
								f.Seek(superBlock.S_inode_start+(int64(carpetaComprobar.B_content[k].B_inodo)*int64(unsafe.Sizeof(Inodo{}))), 0)
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
				if i == (carp - 1) {
					if datFile.setCont != "-" {
						inodoArchivoNuevo := Inodo{}
						for f := 0; f < 16; f++ {
							inodoArchivoNuevo.I_block[f] = -1
						}
						//inodoarchivo nuevo
						inodoArchivoNuevo.I_gid = 1
						inodoArchivoNuevo.I_uid = 1
						copy(inodoArchivoNuevo.I_type[:], "1")
						inodoArchivoNuevo.I_perm = 664
						dt := time.Now()
						copy(inodoArchivoNuevo.I_mtime[:], dt.String())
						copy(inodoArchivoNuevo.I_ctime[:], dt.String())
						copy(inodoArchivoNuevo.I_atime[:], dt.String())
						dat, err := ioutil.ReadFile(datFile.setCont)
						Check(err)

						inodoArchivoNuevo.I_size = int64(len(string(dat)))
						blocksAUsar := cantidadBloquesAUsar(string(dat))

						for h := 0; h < blocksAUsar; h++ {
							if h < 16 {
								inodoArchivoNuevo.I_block[h] = superBlock.S_first_blo
								superBlock.S_first_blo = superBlock.S_first_blo + 1
								superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1
								var otro byte = '1'
								var actual byte = '0'
								for j := 0; j < int(superBlock.S_blocks_count); j++ {

									f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Read(f, binary.BigEndian, &actual)

									if string(actual) == "0" {
										f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Write(f, binary.BigEndian, otro)
										break
									}
								}
							}
						}
						f.Seek(superBlock.S_inode_start+superBlock.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
						err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)
						var condi int = 0
						var posicion int = 0

						for inodoArchivoNuevo.I_block[condi] != -1 {
							contNumero := 0
							contenidoNuevo := ""

							for {
								if contNumero < 63 && posicion < len(string(dat)) {
									contenidoNuevo += string(dat[posicion])
									contNumero++
									posicion++
								} else {
									break
								}
							}
							archivoNuevo := BArchivo{}
							copy(archivoNuevo.B_content[:], contenidoNuevo)
							f.Seek(superBlock.S_block_start+inodoArchivoNuevo.I_block[condi]*int64(unsafe.Sizeof(BArchivo{})), 0)
							err = binary.Write(f, binary.BigEndian, archivoNuevo)
							condi++
						}
						for n := 0; n < 16; n++ {
							if inodoTemp.I_block[n] != -1 {
								carpetaComprobar := BCarpeta{}
								f.Seek(superBlock.S_block_start+inodoTemp.I_block[n]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

								for h := 0; h < 4; h++ {
									if carpetaComprobar.B_content[h].B_inodo == -1 {
										copy(carpetaComprobar.B_content[h].B_name[:], rutaArchivo[carp-1])
										carpetaComprobar.B_content[h].B_inodo = int32(superBlock.S_first_ino)
										f.Seek(superBlock.S_block_start+inodoTemp.I_block[n]*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Write(f, binary.BigEndian, carpetaComprobar)
										n = 20
										break
									}
								}
							} else {
								//crear nuevo bloque carpeta
								carpetaComprobar := BCarpeta{}
								for k := 0; k < 4; k++ {
									carpetaComprobar.B_content[k].B_inodo = -1
									copy(carpetaComprobar.B_content[k].B_name[:], " ")
								}
								copy(carpetaComprobar.B_content[0].B_name[:], rutaArchivo[carp-1])
								carpetaComprobar.B_content[0].B_inodo = int32(superBlock.S_first_ino)

								f.Seek(superBlock.S_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Write(f, binary.BigEndian, carpetaComprobar)

								inodoTemp.I_block[n] = superBlock.S_first_blo
								var otro byte = '1'
								var actual byte = '0'
								for j := 0; j < int(superBlock.S_blocks_count); j++ {

									f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Read(f, binary.BigEndian, &actual)

									if string(actual) == "0" {
										f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Write(f, binary.BigEndian, otro)
										break
									}
								}

								superBlock.S_first_blo = superBlock.S_first_blo + 1
								superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1
								f.Seek(superBlock.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

								dt := time.Now()
								copy(inodoTemp.I_mtime[:], dt.String())
								f.Seek(superBlock.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
								err = binary.Write(f, binary.BigEndian, inodoTemp)
								break
							}
						}
						superBlock.S_first_ino = superBlock.S_first_ino + 1
						superBlock.S_free_inodes_count = superBlock.S_free_inodes_count - 1
						f.Seek(iniPart, 0)
						err = binary.Write(f, binary.BigEndian, superBlock)

						var otro byte = '1'
						var actual byte = '0'
						for j := 0; j < int(superBlock.S_inodes_count); j++ {

							f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Read(f, binary.BigEndian, &actual)

							if string(actual) == "0" {
								f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Write(f, binary.BigEndian, otro)
								break
							}
						}
						fmt.Println("Se creo el Archivo con exito")
						break

					} else {
						inodoArchivoNuevo := Inodo{}
						for h := 0; h < 16; h++ {
							inodoArchivoNuevo.I_block[h] = -1
						}
						inodoArchivoNuevo.I_gid = 1
						inodoArchivoNuevo.I_uid = 1
						copy(inodoArchivoNuevo.I_type[:], "1")
						inodoArchivoNuevo.I_perm = 664
						dt := time.Now()
						copy(inodoArchivoNuevo.I_mtime[:], dt.String())
						copy(inodoArchivoNuevo.I_atime[:], dt.String())
						copy(inodoArchivoNuevo.I_ctime[:], dt.String())

						if datFile.setSize == 0 {
							inodoArchivoNuevo.I_size = 0

							for u := 0; u < 16; u++ {
								if inodoTemp.I_block[u] != -1 {
									carpetaComprobar := BCarpeta{}
									f.Seek(superBlock.S_block_start+inodoTemp.I_block[u]*int64(unsafe.Sizeof(BCarpeta{})), 0)
									err = binary.Read(f, binary.BigEndian, &carpetaComprobar)
									for p := 0; p < 4; p++ {
										if carpetaComprobar.B_content[p].B_inodo == -1 {
											copy(carpetaComprobar.B_content[p].B_name[:], rutaArchivo[carp-1])
											carpetaComprobar.B_content[p].B_inodo = int32(superBlock.S_first_ino)
											f.Seek(superBlock.S_block_start+inodoTemp.I_block[u]*int64(unsafe.Sizeof(BCarpeta{})), 0)
											err = binary.Write(f, binary.BigEndian, carpetaComprobar)
											u = 20
											break
										}
									}
								} else {
									//crear nuevo bloque carpeta
									carpetaComprobar := BCarpeta{}
									for k := 0; k < 4; k++ {
										carpetaComprobar.B_content[k].B_inodo = -1
										copy(carpetaComprobar.B_content[k].B_name[:], " ")
									}
									copy(carpetaComprobar.B_content[0].B_name[:], rutaArchivo[carp-1])
									carpetaComprobar.B_content[0].B_inodo = int32(superBlock.S_first_ino)
									f.Seek(superBlock.S_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
									err = binary.Write(f, binary.BigEndian, carpetaComprobar)

									inodoTemp.I_block[u] = superBlock.S_first_blo

									var otro byte = '1'
									var actual byte = '0'
									for j := 0; j < int(superBlock.S_blocks_count); j++ {

										f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Read(f, binary.BigEndian, &actual)

										if string(actual) == "0" {
											f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
											err = binary.Write(f, binary.BigEndian, otro)
											break
										}
									}
									superBlock.S_first_blo = superBlock.S_first_blo + 1
									superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1

									f.Seek(superBlock.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
									err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

									copy(inodoTemp.I_mtime[:], dt.String())
									f.Seek(superBlock.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
									err = binary.Write(f, binary.BigEndian, inodoTemp)
									break
								}
							}
							f.Seek(superBlock.S_inode_start+superBlock.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
							err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)

							superBlock.S_first_ino = superBlock.S_first_ino + 1
							superBlock.S_free_inodes_count = superBlock.S_free_inodes_count - 1
							f.Seek(iniPart, 0)
							err = binary.Write(f, binary.BigEndian, superBlock)

							var otro byte = '1'
							var actual byte = '0'
							for j := 0; j < int(superBlock.S_inodes_count); j++ {

								f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Read(f, binary.BigEndian, &actual)

								if string(actual) == "0" {
									f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Write(f, binary.BigEndian, otro)
									break
								}
							}

							fmt.Println("Se creo el Archivo con exito")
						} else {
							inodoArchivoNuevo.I_size = int64(datFile.setSize)
							blocksAUsar := cantidadBloquesAUsar1(datFile.setSize)
							if blocksAUsar < 16 {
								for u := 0; u < blocksAUsar; u++ {
									inodoArchivoNuevo.I_block[u] = superBlock.S_first_blo
									superBlock.S_first_blo = superBlock.S_first_blo + 1
									superBlock.S_blocks_count = superBlock.S_blocks_count - 1

									var otro byte = '1'
									var actual byte = '0'
									for j := 0; j < int(superBlock.S_blocks_count); j++ {

										f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Read(f, binary.BigEndian, &actual)

										if string(actual) == "0" {
											f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
											err = binary.Write(f, binary.BigEndian, otro)
											break
										}
									}

								}
								var condi int = 0
								posicion := 0
								for {
									if inodoArchivoNuevo.I_block[condi] != -1 {
										break
									}
									contNumero := 0
									contenidoNuevo := ""
									numero := 0
									for {
										if contNumero < 63 && posicion < datFile.setSize {
											contenidoNuevo += strconv.Itoa(numero)
											contNumero++
											numero++
											posicion++
										} else {
											break
										}

										if numero == 10 {
											numero = 0
										}
									}
									archivoNuevo := BArchivo{}
									copy(archivoNuevo.B_content[:], contenidoNuevo)
									f.Seek(superBlock.S_block_start+inodoArchivoNuevo.I_block[condi]*int64(unsafe.Sizeof(BArchivo{})), 0)
									err = binary.Write(f, binary.BigEndian, archivoNuevo)
									condi++

								}
								for h := 0; h < 16; h++ {
									if inodoTemp.I_block[h] != -1 {
										carpetaComprobar := BCarpeta{}
										f.Seek(superBlock.S_block_start+inodoTemp.I_block[h]*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Read(f, binary.BigEndian, &carpetaComprobar)
										for k := 0; k < 4; k++ {
											if carpetaComprobar.B_content[k].B_inodo == -1 {
												copy(carpetaComprobar.B_content[k].B_name[:], rutaArchivo[carp-1])
												carpetaComprobar.B_content[k].B_inodo = int32(superBlock.S_first_ino)
												f.Seek(superBlock.S_block_start+inodoTemp.I_block[h]*int64(unsafe.Sizeof(BCarpeta{})), 0)
												err = binary.Write(f, binary.BigEndian, carpetaComprobar)
												h = 20
												break
											}
										}

									} else {
										//crear nuevo bloque carpeta
										carpetaComprobar := BCarpeta{}
										for k := 0; k < 4; k++ {
											carpetaComprobar.B_content[k].B_inodo = -1
											copy(carpetaComprobar.B_content[k].B_name[:], " ")
										}
										carpetaComprobar.B_content[0].B_inodo = int32(superBlock.S_first_ino)
										copy(carpetaComprobar.B_content[0].B_name[:], rutaArchivo[carp-1])
										f.Seek(superBlock.S_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Write(f, binary.BigEndian, carpetaComprobar)

										inodoTemp.I_block[h] = superBlock.S_first_blo

										var otro byte = '1'
										var actual byte = '0'
										for j := 0; j < int(superBlock.S_blocks_count); j++ {

											f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
											err = binary.Read(f, binary.BigEndian, &actual)

											if string(actual) == "0" {
												f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
												err = binary.Write(f, binary.BigEndian, otro)
												break
											}
										}
										superBlock.S_first_blo = superBlock.S_first_blo + 1
										superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1

										f.Seek(superBlock.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

										dt := time.Now()
										copy(inodoTemp.I_mtime[:], dt.String())
										f.Seek(superBlock.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
										err = binary.Write(f, binary.BigEndian, carpetaComprobar)
										break

									}
								}

								f.Seek(superBlock.S_inode_start+superBlock.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
								err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)

								superBlock.S_first_ino = superBlock.S_first_ino + 1
								superBlock.S_free_inodes_count = superBlock.S_free_inodes_count - 1

								f.Seek(iniPart, 0)
								err = binary.Write(f, binary.BigEndian, superBlock)

								var otro byte = '1'
								var actual byte = '0'
								for j := 0; j < int(superBlock.S_inodes_count); j++ {

									f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Read(f, binary.BigEndian, &actual)

									if string(actual) == "0" {
										f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Write(f, binary.BigEndian, otro)
										break
									}
								}
								fmt.Println("Se creo el Archivo con exito")
							} else {
								fmt.Println("No se pudo crear el archivo")
							}
						}
						break
					}
				} else {
					if datFile.setR {
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
						copy(inodoArchivoNuevo.I_atime[:], dt.String())
						copy(inodoArchivoNuevo.I_ctime[:], dt.String())
						inodoArchivoNuevo.I_block[0] = superBlock.S_first_blo

						carpetaRaiz := BCarpeta{}

						copy(carpetaRaiz.B_content[0].B_name[:], ".")
						copy(carpetaRaiz.B_content[1].B_name[:], "..")
						copy(carpetaRaiz.B_content[2].B_name[:], " ")
						copy(carpetaRaiz.B_content[3].B_name[:], " ")

						carpetaRaiz.B_content[0].B_inodo = int32(superBlock.S_first_ino)

						carpetaComprobar := BCarpeta{}
						f.Seek(superBlock.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

						carpetaRaiz.B_content[1].B_inodo = carpetaComprobar.B_content[0].B_inodo
						carpetaRaiz.B_content[2].B_inodo = -1
						carpetaRaiz.B_content[3].B_inodo = -1

						f.Seek(superBlock.S_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Write(f, binary.BigEndian, carpetaRaiz)

						superBlock.S_first_blo = superBlock.S_first_blo + 1
						superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1

						for b := 0; b < 16; b++ {
							if inodoTemp.I_block[b] != -1 {

								carpetaComprobar1 := BCarpeta{}
								f.Seek(superBlock.S_block_start+inodoTemp.I_block[b]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &carpetaComprobar1)

								for k := 0; k < 4; k++ {
									if carpetaComprobar1.B_content[k].B_inodo == -1 {
										copy(carpetaComprobar1.B_content[k].B_name[:], rutaArchivo[i])
										carpetaComprobar1.B_content[k].B_inodo = int32(superBlock.S_first_ino)
										f.Seek(superBlock.S_block_start+inodoTemp.I_block[b]*int64(unsafe.Sizeof(BCarpeta{})), 0)
										err = binary.Write(f, binary.BigEndian, carpetaComprobar1)
										b = 20
										break
									}
								}

							} else {
								//crear nuevo bloque carpeta
								carpetaComprobar1 := BCarpeta{}
								for h := 0; h < 4; h++ {
									carpetaComprobar1.B_content[h].B_inodo = -1
									copy(carpetaComprobar1.B_content[h].B_name[:], " ")
								}
								copy(carpetaComprobar1.B_content[0].B_name[:], rutaArchivo[i])
								carpetaComprobar1.B_content[0].B_inodo = int32(superBlock.S_first_ino)

								f.Seek(superBlock.S_block_start+superBlock.S_first_blo*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Write(f, binary.BigEndian, carpetaComprobar1)

								inodoTemp.I_block[b] = superBlock.S_first_blo

								var otro byte = '1'
								var actual byte = '0'
								for j := 0; j < int(superBlock.S_blocks_count); j++ {

									f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
									err = binary.Read(f, binary.BigEndian, &actual)

									if string(actual) == "0" {
										f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
										err = binary.Write(f, binary.BigEndian, otro)
										break
									}
								}

								superBlock.S_first_blo = superBlock.S_first_blo + 1
								superBlock.S_free_blocks_count = superBlock.S_free_blocks_count - 1

								f.Seek(superBlock.S_block_start+inodoTemp.I_block[0]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &carpetaComprobar)

								dt := time.Now()
								copy(inodoTemp.I_mtime[:], dt.String())
								f.Seek(superBlock.S_inode_start+int64(carpetaComprobar.B_content[0].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
								err = binary.Write(f, binary.BigEndian, inodoTemp)
								break
							}
						}

						f.Seek(superBlock.S_inode_start+superBlock.S_first_ino*int64(unsafe.Sizeof(Inodo{})), 0)
						err = binary.Write(f, binary.BigEndian, inodoArchivoNuevo)

						superBlock.S_first_ino = superBlock.S_first_ino + 1
						superBlock.S_free_inodes_count = superBlock.S_free_inodes_count - 1

						f.Seek(iniPart, 0)
						err = binary.Write(f, binary.BigEndian, superBlock)

						var otro byte = '1'
						var actual byte = '0'
						for j := 0; j < int(superBlock.S_blocks_count); j++ {

							f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Read(f, binary.BigEndian, &actual)

							if string(actual) == "0" {
								f.Seek(superBlock.S_bm_block_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Write(f, binary.BigEndian, otro)
								break
							}
						}

						for j := 0; j < int(superBlock.S_inodes_count); j++ {

							f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
							err = binary.Read(f, binary.BigEndian, &actual)

							if string(actual) == "0" {
								f.Seek(superBlock.S_bm_inode_start+int64(j)*int64(unsafe.Sizeof(otro)), 0)
								err = binary.Write(f, binary.BigEndian, otro)
								break
							}
						}

						//inodoAnterior = inodoTemp
						inodoTemp = inodoArchivoNuevo

						fmt.Println("Se creo la carpeta: " + rutaArchivo[i])
					} else {
						fmt.Println("ERROR: no se encontro la ruta para crear el archivo")
					}
				}
			}
		}
	}
}

func nombreAStringFile(data [12]byte) string {
	var var1 string = ""
	for i := 0; i < 12; i++ {
		if data[i] != 0 {
			var1 += string(data[i])
		}
	}
	return var1
}

func cantidadBloquesAUsar(dat string) int {
	contador := 1
	noBloques := 0
	for i := 0; i < len(dat); i++ {
		if contador == 63 {
			noBloques += 1
			contador = 0
		}
		contador++
	}
	if len(dat)%63 != 0 {
		noBloques += 1
	}
	return noBloques
}
func cantidadBloquesAUsar1(dat int) int {
	contador := 1
	noBloques := 0
	for i := 0; i < dat; i++ {
		if contador == 63 {
			noBloques += 1
			contador = 0
		}
		contador++
	}
	if dat%63 != 0 {
		noBloques += 1
	}
	return noBloques
}
