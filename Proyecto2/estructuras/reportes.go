package estructuras

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unsafe"
)

func EjecutarRepDisk(datRep PropRep, listaDiscos *[100]Mount) {
	var buffer bytes.Buffer

	buffer.WriteString("digraph G{\ntbl [\nshape=box\nlabel=<\n<table border='0' cellborder='2' width='100' height=\"30\" color='lightblue4'>\n<tr>")
	cuerpo := ""
	finG := "     </tr>\n</table>\n>];\n}"

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            REPORTE DISK                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if existeIdMount(datRep.setId, listaDiscos) {
		pathDisco := ""

		for i := 0; i < 100; i++ {
			if datRep.setId == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				break
			}
		}

		DirExist(datRep.setPath)

		mbrTemp := MBR{}
		f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe el archivo en la ruta")
		} else {
			porcentajeUtilizado := 0.0
			var EspacioUtilizado int64 = 0

			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			TamanoDisco := mbrTemp.Mbr_tamano

			var parts [4]Particion
			parts[0] = mbrTemp.Mbr_partition_1
			parts[1] = mbrTemp.Mbr_partition_2
			parts[2] = mbrTemp.Mbr_partition_3
			parts[3] = mbrTemp.Mbr_partition_4

			cuerpo += "<td height='30' width='75'> MBR </td>"

			for i := 0; i < 4; i++ {

				if string(parts[i].Part_type[:]) == "p" {

					porcentajeUtilizado = (float64(parts[i].Part_size) / float64(TamanoDisco)) * 100
					cuerpo += "<td height='30' width='75.0'>PRIMARIA <br/>" + nombreAString(parts[i].Part_name) + " <br/> Ocupado: " + strconv.Itoa(int(porcentajeUtilizado)) + "%</td>"
					EspacioUtilizado += parts[i].Part_size

				} else if string(parts[i].Part_type[:]) == "e" {
					EspacioUtilizado += parts[i].Part_size
					porcentajeUtilizado = (float64(parts[i].Part_size) / float64(TamanoDisco)) * 100
					cuerpo += "<td  height='30' width='15.0'>\n <table border='5'  height='30' WIDTH='15.0' cellborder='1'>\n  <tr>  <td height='60' colspan='100%'>EXTENDIDA <br/> " + nombreAString(parts[i].Part_name) + " <br/> Ocupado:" + strconv.Itoa(int(porcentajeUtilizado)) + "%</td>  </tr>\n<tr>"

					//revisar las particiones logicas
					ebrTemp := EBR{}
					f.Seek(parts[i].Part_start, 0)
					err = binary.Read(f, binary.BigEndian, &ebrTemp)

					var EspacioUtilizado1 int64 = 0
					cont := 0
					for ebrTemp.Part_next != -1 {
						EspacioUtilizado1 += ebrTemp.Part_size
						porcentajeUtilizado = (float64(ebrTemp.Part_size) / float64(parts[i].Part_size)) * 100
						cuerpo += "<td height='30'>EBR</td><td height='30'> Logica:  " + nombreAString(ebrTemp.Part_name) + " " + strconv.Itoa(int(porcentajeUtilizado)) + "%</td>"
						cont++
						f.Seek(ebrTemp.Part_next, 0)
						err = binary.Read(f, binary.BigEndian, &ebrTemp)

					}
					if parts[i].Part_size-EspacioUtilizado1 > 0 {
						porcentajeUtilizado = (float64(TamanoDisco-EspacioUtilizado1) / float64(TamanoDisco)) * 100
						cuerpo += "<td height='30' width='100%'>Libre: " + strconv.Itoa(int(porcentajeUtilizado)) + "%</td>"
					}

					cuerpo += "</tr>\n </table>\n</td>"

				} else if string(parts[i].Part_status[:]) == "0" {
					cuerpo += "<td height='30' width='75.0'>Libre</td>"
				}
			}
			if TamanoDisco-EspacioUtilizado > 0 {
				porcentajeUtilizado = (float64(TamanoDisco-EspacioUtilizado) / float64(TamanoDisco)) * 100
				cuerpo += "<td height='30' width='75.0'>Libre: " + strconv.Itoa(int(porcentajeUtilizado)) + "%</td>"
			}
			textoFinal := cuerpo + finG

			buffer.WriteString(textoFinal)
			datos := string(buffer.String())
			CreateArchivo(datRep.setPath, datos)

			fmt.Println("Reporte Disk creado correctamente")

		}
	} else {
		fmt.Println("El id no existe, o la particion no esta montada")
	}
}

func EjecutarRepTree(datRep PropRep, listaDiscos *[100]Mount) {
	var buffer bytes.Buffer

	buffer.WriteString("digraph G{ \n rankdir=LR;\n node [shape = record, style=filled, fillcolor=lightBlue]; \n graph [pad=0.5, nodesep=1.5, ranksep=2, splines=true,]; \n")
	cuerpo := ""
	finG := " \n}"

	superBloque := SupB{}

	var actual byte = '0'
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            REPORTE TREE                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if existeIdMount(datRep.setId, listaDiscos) {
		pathDisco := ""
		partName := ""
		for i := 0; i < 100; i++ {
			if datRep.setId == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				partName = listaDiscos[i].NombreParticion
				break
			}
		}

		DirExist(datRep.setPath)
		mbrTemp := MBR{}
		f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe el archivo en la ruta")
		} else {
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			_, iniPart := returnDatosPart(mbrTemp, pathDisco, partName)

			f.Seek(iniPart, 0)
			err = binary.Read(f, binary.BigEndian, &superBloque)

			contadorIno := 0
			contador := 0
			actualBlo := 0

			//Graficando los inodos
			for i := 0; i < int(superBloque.S_inodes_count); i++ {
				f.Seek(superBloque.S_bm_inode_start+int64(i)*int64(unsafe.Sizeof(actual)), 0)
				err = binary.Read(f, binary.BigEndian, &actual)

				if string(actual) == "1" {
					inodoTemp := Inodo{}
					f.Seek(superBloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Read(f, binary.BigEndian, &inodoTemp)
					cuerpo += "inodo" + strconv.Itoa(contadorIno) + " [shape=record,label=\"<f0> Inodo " + strconv.Itoa(contadorIno) + "|i_uid : " + strconv.Itoa(int(inodoTemp.I_uid)) + "|i_gid : " + strconv.Itoa(int(inodoTemp.I_gid)) + "|i_size : " + strconv.Itoa(int(inodoTemp.I_size)) + "|i_atime : " + string(inodoTemp.I_atime[:]) + "|i_ctime : " + string(inodoTemp.I_ctime[:]) + "|i_mtime : " + string(inodoTemp.I_mtime[:]) + "|i_type : " + string(inodoTemp.I_type[:])
					contadorIno++

					for j := 0; j < 16; j++ {
						cuerpo += "|<f" + strconv.Itoa(j+1) + "> i_block" + strconv.Itoa(j+1) + " : " + strconv.Itoa(int(inodoTemp.I_block[j]))
					}
					cuerpo += "|i_perm : " + strconv.Itoa(int(inodoTemp.I_perm)) + "\"];"

				}
			}

			//Graficando Bloques

			for i := 0; i < int(superBloque.S_inodes_count); i++ {
				f.Seek(superBloque.S_bm_inode_start+int64(i)*int64(unsafe.Sizeof(actual)), 0)
				err = binary.Read(f, binary.BigEndian, &actual)

				if string(actual) == "1" {
					inodoTemp := Inodo{}
					f.Seek(superBloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Read(f, binary.BigEndian, &inodoTemp)

					for j := 0; j < 16; j++ {
						if inodoTemp.I_block[j] != -1 {
							if string(inodoTemp.I_type[:]) == "0" {
								bloTemp := BCarpeta{}
								f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &bloTemp)

								cuerpo += "\n node [shape = record, style=filled, fillcolor=orange]; \n"
								cuerpo += "struct" + strconv.Itoa(int(inodoTemp.I_block[j])) + " [shape=record, width = 3,label=\"<f0> Bloque Carpeta" + strconv.Itoa(int(inodoTemp.I_block[j])) + "|{B_name | B_inodo}"

								for k := 0; k < 4; k++ {
									cuerpo += "|{" + nombreAStringFile(bloTemp.B_content[k].B_name) + "|<f" + strconv.Itoa(int(k+1)) + "> " + strconv.Itoa(int(bloTemp.B_content[k].B_inodo)) + "}"
								}
								contador++
								cuerpo += "\"];"
							} else {
								bloTemp := BArchivo{}
								f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &bloTemp)

								cuerpo += "\n node [shape = record, style=filled, fillcolor=gray]; \n"
								cuerpo += "struct" + strconv.Itoa(int(inodoTemp.I_block[j])) + " [shape=record,label=\"<f0> Bloque Archivo" + strconv.Itoa(int(inodoTemp.I_block[j])) + "|"
								//posibleCambio

								if strings.Contains(archivoAString(bloTemp.B_content), "\n") {
									saltos := strings.Split(archivoAString(bloTemp.B_content), "\n")
									for z := 0; z < len(saltos); z++ {
										cuerpo += saltos[z] + "\\n"
									}
								} else {
									cuerpo += archivoAString(bloTemp.B_content)
								}
								contador++
								cuerpo += "\"];"

							}
						}
					}
				}
			}

			cuerpo += "\n \n"

			//creando las lineas de los inodos a los bloques
			for i := 0; i < int(superBloque.S_inodes_count); i++ {
				f.Seek(superBloque.S_bm_inode_start+int64(i)*int64(unsafe.Sizeof(actual)), 0)
				err = binary.Read(f, binary.BigEndian, &actual)

				if string(actual) == "1" {
					inodoTemp := Inodo{}
					f.Seek(superBloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Read(f, binary.BigEndian, &inodoTemp)

					for j := 0; j < 16; j++ {
						if inodoTemp.I_block[j] != -1 {
							cuerpo += "inodo" + strconv.Itoa(i) + ":f" + strconv.Itoa(j+1) + "-> struct" + strconv.Itoa(int(inodoTemp.I_block[j])) + ":f0   \n"
						}
					}
				}
			}
			//creando las lineas para los bloques a los inodos
			for i := 0; i < int(superBloque.S_inodes_count); i++ {
				f.Seek(superBloque.S_bm_inode_start+int64(i)*int64(unsafe.Sizeof(actual)), 0)
				err = binary.Read(f, binary.BigEndian, &actual)

				if string(actual) == "1" {
					inodoTemp := Inodo{}
					f.Seek(superBloque.S_inode_start+int64(i)*int64(unsafe.Sizeof(Inodo{})), 0)
					err = binary.Read(f, binary.BigEndian, &inodoTemp)
					for j := 0; j < 16; j++ {
						if inodoTemp.I_block[j] != -1 {
							if string(inodoTemp.I_type[:]) == "0" {
								bloTemp := BCarpeta{}
								f.Seek(superBloque.S_block_start+inodoTemp.I_block[j]*int64(unsafe.Sizeof(BCarpeta{})), 0)
								err = binary.Read(f, binary.BigEndian, &bloTemp)

								for k := 0; k < 4; k++ {
									var1 := nombreAStringFile(bloTemp.B_content[k].B_name)
									if bloTemp.B_content[k].B_inodo != -1 && var1 != "." && var1 != ".." {
										cuerpo += "struct" + strconv.Itoa(int(inodoTemp.I_block[j])) + ":f" + strconv.Itoa(k+1) + "-> inodo" + strconv.Itoa(int(bloTemp.B_content[k].B_inodo)) + ":f0  \n"
										actualBlo++
									}
								}

							} else {
								actualBlo++
							}
						}
					}
				}
			}

			textoFinal := cuerpo + finG

			buffer.WriteString(textoFinal)
			datos := string(buffer.String())
			CreateArchivo(datRep.setPath, datos)

			fmt.Println("Reporte Tree creado correctamente")
		}

	} else {
		fmt.Println("El id no existe, o la particion no esta montada")
	}
}

func EjecutarRepFile(datRep PropRep, listaDiscos *[100]Mount) {
	var buffer bytes.Buffer

	buffer.WriteString("digraph G{ \n rankdir=LR; \n")
	cuerpo := ""
	finG := " \n}"

	superBloque := SupB{}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            REPORTE FILE                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	rutaArchivo := strings.Split(datRep.setRuta, "/")
	if existeIdMount(datRep.setId, listaDiscos) {
		pathDisco := ""
		partName := ""
		for i := 0; i < 100; i++ {
			if datRep.setId == listaDiscos[i].IdMount {
				pathDisco = listaDiscos[i].Path
				partName = listaDiscos[i].NombreParticion
				break
			}
		}

		DirExist(datRep.setPath)
		mbrTemp := MBR{}
		f, err := os.OpenFile(pathDisco, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe el archivo en la ruta")
		} else {
			f.Seek(0, 0)
			err = binary.Read(f, binary.BigEndian, &mbrTemp)
			_, iniPart := returnDatosPart(mbrTemp, pathDisco, partName)

			f.Seek(iniPart, 0)
			err = binary.Read(f, binary.BigEndian, &superBloque)

			inodoTemp := Inodo{}

			f.Seek(superBloque.S_inode_start, 0)
			err = binary.Read(f, binary.BigEndian, &inodoTemp)

			existeArchivo := false
			existeCarpeta := false
			carp := len(rutaArchivo)

			cuerpo += "struct0 [shape=record,label=\"<f0> Archivo | nombre : " + rutaArchivo[carp-1] + " | Contenido | "

			//posible cambio en carp +1
			for i := 1; i < carp+1; i++ {
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
									f.Seek(superBloque.S_inode_start+int64(carpetaComprobar.B_content[k].B_inodo)*int64(unsafe.Sizeof(Inodo{})), 0)
									err = binary.Read(f, binary.BigEndian, &inodoTemp)

									existeCarpeta = true
									j = 20
									break
								}

							}
						}
					}
				} else {
					fmt.Println("Se encotro el archivo :)")
					existeArchivo = true
					break
				}
				if !existeCarpeta {
					fmt.Println("ERROR: no se encontro la ruta del archivo")
					break
				}
			}
			if existeArchivo {
				for i := 0; i < 16; i++ {
					if inodoTemp.I_block[i] != -1 {
						bloTemp := BArchivo{}
						f.Seek(superBloque.S_block_start+inodoTemp.I_block[i]*int64(unsafe.Sizeof(BCarpeta{})), 0)
						err = binary.Read(f, binary.BigEndian, &bloTemp)

						//posible cambio
						saltos := strings.Split(archivoAString(bloTemp.B_content), "\n")
						for z := 0; z < len(saltos); z++ {
							cuerpo += saltos[z] + "\\n"
						}

					}
				}
			}
			cuerpo += " \"];"

			textoFinal := cuerpo + finG
			buffer.WriteString(textoFinal)
			datos := string(buffer.String())
			CreateArchivo(datRep.setPath, datos)

			fmt.Println("Reporte File creado correctamente")
		}
	} else {
		fmt.Println("El id no existe, o la particion no esta montada")
	}

}

func CreateArchivo(path string, data string) {
	fmt.Println(int64(unsafe.Sizeof(BCarpeta{})))
	fmt.Println(int64(unsafe.Sizeof(BArchivo{})))
	propiedades := strings.Split(path, "/")
	nombreArchivo := propiedades[len(propiedades)-1]
	f, err := os.Create(path[0:len(path)-len(nombreArchivo)] + nombreArchivo[0:len(nombreArchivo)-4] + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
	cmd := exec.Command("dot", "-Tpdf", path[0:len(path)-len(nombreArchivo)]+nombreArchivo[0:len(nombreArchivo)-4]+".txt", "-o", path)
	cmd.CombinedOutput()

}
