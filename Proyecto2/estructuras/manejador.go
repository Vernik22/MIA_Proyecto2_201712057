package estructuras

import (
	"container/list"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var hayInicioSesion bool = false
var usuarioIniciado = User{nombreUsuario: "-", idPartMontada: "-"}

func LeerTexto(data string, listaDiscos *[100]Mount) {
	//fmt.Println("desde LeerTexto: " + data)
	//fmt.Println(string(listaDiscos[0].Estado[:]))
	//para leer la cadena enviada
	ListaComandos := list.New()
	lineaComando := strings.Split(data, "\n")
	var com Comando
	for i := 0; i < len(lineaComando); i++ {
		if len(lineaComando[i]) > 0 {
			EsComentario := lineaComando[i][0:1]
			if EsComentario != "#" {
				comando := lineaComando[i]
				//ahora lo separo por espacios ejemplo: mkdisk -path -size
				propiedades := strings.Split(string(comando), " ")
				nombreComando := propiedades[0]
				com.Nombre = strings.ToLower(nombreComando)
				propiedadesTemp := make([]Propiedad, len(propiedades)-1)
				for f := 0; f < len(propiedadesTemp); f++ {
					propiedadesTemp[f].Nombre = "|"
				}
				for j := 1; j < len(propiedades); j++ {
					if propiedades[j] == "" || propiedades[j] == " " || propiedades[j] == "#" {
						continue
					} else {
						if strings.Contains(propiedades[j], "=") {
							if strings.Contains(propiedades[j], "#") {
								quitComen := strings.Split(propiedades[j], "#")
								valor_propiedad_Comando := strings.Split(quitComen[0], "=")
								propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
								propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
							} else if strings.Contains(propiedades[j], "\"") {
								valor_propiedad_Comando := strings.Split(propiedades[j], "=")
								propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
								propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
								for f := j + 1; f < len(propiedades); f++ {
									if strings.Contains(propiedades[f], "\"") {
										propiedadesTemp[j-1].Valor += " " + propiedades[f]
										break
									} else {
										propiedadesTemp[j-1].Valor += " " + propiedades[f]
									}
								}
							} else {
								valor_propiedad_Comando := strings.Split(propiedades[j], "=")
								propiedadesTemp[j-1].Nombre = valor_propiedad_Comando[0]
								propiedadesTemp[j-1].Valor = valor_propiedad_Comando[1]
							}

						} else if propiedades[j] == "-r" || propiedades[j] == "-R" {
							propiedadesTemp[j-1].Nombre = propiedades[j]
						} else if propiedades[j] == "-p" || propiedades[j] == "-P" {
							propiedadesTemp[j-1].Nombre = propiedades[j]
						}
					}
				}
				com.Propiedades = propiedadesTemp
				//agregando el comando a la lista de comandos
				ListaComandos.PushBack(com)
			} else {
				fmt.Println("Es un comentario")
			}
		}

	}

	listaComandosValidos(ListaComandos, listaDiscos)
}

func listaComandosValidos(ListaComandos *list.List, listaDiscos *[100]Mount) {
	/*
		c := Mount{}
		copy(c.Estado[:], "2")
		listaDiscos[0] = c

		fmt.Println("ejecutando comandos validos")
		fmt.Println(string(listaDiscos[0].Estado[:]))
		fmt.Println(string(listaDiscos[1].Estado[:]))
	*/
	for element := ListaComandos.Front(); element != nil; element = element.Next() {
		comandoTemp := element.Value.(Comando)
		//lista de propiedades de comando
		nombreComando := comandoTemp.Nombre
		if nombreComando == "mkdisk" {
			parametrosValidos := true
			flagFit := true  //opcional
			flagUnit := true //opcional
			flagSize := true //obligatorio
			flagPath := true //obligatorio

			mkdisk := PropMkdisk{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-size" {
					s, _ := strconv.Atoi(comandoTemp.Propiedades[f].Valor)
					flagSize = false
					if s > 0 {
						mkdisk.setSize = s
					} else {
						fmt.Println("ERROR: tamano erroneo, es negativo o cero, intente de nuevo")
						parametrosValidos = false
						flagSize = true
						break
					}
				} else if nombreProp == "-fit" {
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					flagFit = false
					if var1 == "bf" || var1 == "ff" || var1 == "wf" {
						mkdisk.setFit = var1
					} else {
						fmt.Println("ERROR: fit erroneo, no soportado")
						parametrosValidos = false
						flagFit = true
						break
					}

				} else if nombreProp == "-unit" {
					flagUnit = false
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "m" || var1 == "k" {
						mkdisk.setUnit = var1
					} else {
						fmt.Println("ERROR: unit erroneo, no soportado")
						parametrosValidos = false
						flagUnit = true
						break
					}
				} else if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						mkdisk.setPath = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						mkdisk.setPath = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
				}
			}
			if flagFit {
				mkdisk.setFit = "ff"
			}
			if flagUnit {
				mkdisk.setUnit = "m"
			}
			if flagSize == false && flagPath == false {
				parametrosValidos = false
			}

			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				EjecutarComandoMkdisk(mkdisk)
			}

		} else if nombreComando == "rmdisk" {
			parametrosValidos := true
			flagPath := true //obligatorio
			discoRem := ""
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						discoRem = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						discoRem = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
				}
			}
			if flagPath == false {
				parametrosValidos = false
			}

			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				EjecutarComandoRmdisk(discoRem)
			}

		} else if nombreComando == "fdisk" {
			parametrosValidos := true
			flagFit := true
			flagUnit := true
			flagType := true

			flagName := true
			flagSize := true
			flagPath := true

			fdisk := PropFdisk{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-size" {
					s, _ := strconv.Atoi(comandoTemp.Propiedades[f].Valor)
					flagSize = false
					if s > 0 {
						fdisk.setSize = s
					} else {
						fmt.Println("ERROR: tamano erroneo, es negativo o cero, intente de nuevo")
						parametrosValidos = false
						flagSize = true
						break
					}

				} else if nombreProp == "-unit" {
					flagUnit = false
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "m" || var1 == "k" || var1 == "b" {
						fdisk.setUnit = var1
					} else {
						fmt.Println("ERROR: unit erroneo, no soportado")
						parametrosValidos = false
						flagUnit = true
						break
					}
				} else if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						fdisk.setPath = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						fdisk.setPath = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
				} else if nombreProp == "-type" {
					flagType = false
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "p" || var1 == "e" || var1 == "l" {
						fdisk.setType = var1
					} else {
						fmt.Println("ERROR: type erroneo, no soportado")
						parametrosValidos = false
						flagType = true
						break
					}
				} else if nombreProp == "-fit" {
					flagFit = false
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "bf" || var1 == "ff" || var1 == "wf" {
						fdisk.setFit = var1
					} else {
						fmt.Println("ERROR: fit erroneo, no soportado")
						parametrosValidos = false
						flagFit = true
						break
					}
				} else if nombreProp == "-name" {
					flagName = false
					fdisk.setName = comandoTemp.Propiedades[f].Valor
				}
			}

			if flagFit {
				fdisk.setFit = "wf"
			}
			if flagUnit {
				fdisk.setUnit = "k"
			}
			if flagType {
				fdisk.setType = "p"
			}

			if flagPath == false && flagName == false && flagSize == false {
				parametrosValidos = false
			}
			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				EjecutarComandoFdisk(fdisk)
			}

		} else if nombreComando == "mount" {
			parametrosValidos := true
			flagPath := true
			flagName := true

			mount := PropMount{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						mount.setPath = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						mount.setPath = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
				} else if nombreProp == "-name" {
					flagName = false
					mount.setName = comandoTemp.Propiedades[f].Valor
				}

			}
			if flagPath == false && flagName == false {
				parametrosValidos = false
			}

			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				EjecutarComandoMount(mount, listaDiscos)
			}

		} else if nombreComando == "mkfs" {
			parametrosValidos := true
			flagId := true
			flagType := true

			mkfs := PropMkfs{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-id" {
					var1 := strings.ToUpper(comandoTemp.Propiedades[f].Valor)
					flagId = false
					mkfs.setId = var1
				} else if nombreProp == "-type" {
					flagType = false
					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "fast" || var1 == "full" {
						mkfs.setType = var1
					} else {
						fmt.Println("ERROR: type erroneo, no soportado")
						parametrosValidos = false
						flagType = true
						break
					}
				}

			}
			if flagType {
				mkfs.setType = "full"
			}
			if flagId == false {
				parametrosValidos = false
			}
			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				EjecutarComandoMkfs(mkfs, listaDiscos)
			}
		} else if nombreComando == "login" {
			parametrosValidos := true
			flagUser := true
			flagPass := true
			flagId := true

			login := PropLogin{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-id" {
					var1 := strings.ToUpper(comandoTemp.Propiedades[f].Valor)
					flagId = false
					login.setId = var1
				} else if nombreProp == "-usuario" {
					flagUser = false
					login.setUsuario = comandoTemp.Propiedades[f].Valor
				} else if nombreProp == "-password" {
					flagPass = false
					login.setPassword = comandoTemp.Propiedades[f].Valor
				}

			}
			if flagId == false && flagPass == false && flagUser == false {
				parametrosValidos = false
			}
			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				if !hayInicioSesion {
					sesionOk := EjecutarComandoLogin(login, listaDiscos)
					if sesionOk {
						hayInicioSesion = true
						usuarioIniciado.idPartMontada = login.setId
						usuarioIniciado.nombreUsuario = login.setUsuario
					} else {
						fmt.Println("ERROR: Usuario no encontrado ")
						fmt.Println(" ")
					}
				} else {
					fmt.Println("ERROR: YA hay una sesion iniciada ")
					fmt.Println(" ")
				}
			}
		} else if nombreComando == "logout" {
			if hayInicioSesion {
				fmt.Println("--------------------------------------------------------------------------------")
				fmt.Println("                            Ejecutar LOGOUT                         ")
				fmt.Println("--------------------------------------------------------------------------------")
				hayInicioSesion = false
				fmt.Println("Se cerro la sesion de: " + usuarioIniciado.nombreUsuario)
				fmt.Println(" ")
				usuarioIniciado.idPartMontada = "-"
				usuarioIniciado.nombreUsuario = "-"

			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}
		} else if nombreComando == "mkgrp" {
			if hayInicioSesion {
				parametrosValidos := true
				flagName := true

				mkgrp := PropMkgrp{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-name" {
						flagName = false
						mkgrp.setName = comandoTemp.Propiedades[f].Valor
						break
					}
				}
				if flagName == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					mkgrp.setUsuarioAct = usuarioIniciado.nombreUsuario
					mkgrp.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoMkgrp(mkgrp, listaDiscos)
				}
			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "rmgrp" {
			if hayInicioSesion {
				parametrosValidos := true
				flagName := true

				rmgrp := PropMkgrp{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-name" {
						flagName = false
						rmgrp.setName = comandoTemp.Propiedades[f].Valor
						break
					}
				}
				if flagName == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					rmgrp.setUsuarioAct = usuarioIniciado.nombreUsuario
					rmgrp.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoRmgrp(rmgrp, listaDiscos)
				}

			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "mkusr" {
			if hayInicioSesion {
				parametrosValidos := true
				flagUser := true
				flagPass := true
				flagGrp := true

				mkusr := PropMkusr{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-grp" {
						flagGrp = false
						mkusr.setGrp = comandoTemp.Propiedades[f].Valor
					} else if nombreProp == "-usuario" {
						flagUser = false
						mkusr.setUsuario = comandoTemp.Propiedades[f].Valor
					} else if nombreProp == "-pwd" {
						flagPass = false
						mkusr.setPassword = comandoTemp.Propiedades[f].Valor
					}

				}
				if flagGrp == false && flagPass == false && flagUser == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					mkusr.setUsuarioAct = usuarioIniciado.nombreUsuario
					mkusr.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoMkusr(mkusr, listaDiscos)
				}

			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "rmusr" {
			if hayInicioSesion {
				parametrosValidos := true
				flagName := true

				rmusr := PropMkgrp{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-usuario" {
						flagName = false
						rmusr.setName = comandoTemp.Propiedades[f].Valor
						break
					}
				}
				if flagName == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					rmusr.setUsuarioAct = usuarioIniciado.nombreUsuario
					rmusr.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoRmusr(rmusr, listaDiscos)
				}
			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "mkfile" {
			if hayInicioSesion {
				parametrosValidos := true
				flagPath := true
				flagR := true
				flagSize := true
				flagCont := true

				mkfile := PropMkfile{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-path" {
						if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
							conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
							mkfile.setPath = conc[1]
							fmt.Println("Path: " + conc[1])
						} else {
							mkfile.setPath = comandoTemp.Propiedades[f].Valor
						}
						flagPath = false
					} else if nombreProp == "-r" {
						flagR = false
						mkfile.setR = true

					} else if nombreProp == "-size" {
						s, _ := strconv.Atoi(comandoTemp.Propiedades[f].Valor)
						flagSize = false
						if s > 0 {
							mkfile.setSize = s
						} else {
							fmt.Println("ERROR: tamano erroneo, es negativo o cero, intente de nuevo")
							parametrosValidos = false
							flagSize = true
							break
						}
					} else if nombreProp == "-cont" {
						if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
							conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
							mkfile.setCont = conc[1]
							fmt.Println("Path: " + conc[1])
						} else {
							mkfile.setCont = comandoTemp.Propiedades[f].Valor
						}
						flagCont = false
					}
				}
				if flagSize {
					mkfile.setSize = 0
				}
				if flagCont {
					mkfile.setCont = "-"
				}
				if flagR {
					mkfile.setR = false
				}

				if flagPath == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					mkfile.setUsuarioAct = usuarioIniciado.nombreUsuario
					mkfile.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoMkfile(mkfile, listaDiscos)
				}

			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "mkdir" {
			if hayInicioSesion {
				parametrosValidos := true
				flagPath := true
				flagP := true

				mkdir := PropMkdir{}
				for f := 0; f < len(comandoTemp.Propiedades); f++ {
					nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
					if nombreProp == "-path" {
						if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
							conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
							mkdir.setPath = conc[1]
							fmt.Println("Path: " + conc[1])
						} else {
							mkdir.setPath = comandoTemp.Propiedades[f].Valor
						}
						flagPath = false
					} else if nombreProp == "-p" {
						flagP = false
						mkdir.setP = true

					}
				}

				if flagP {
					mkdir.setP = false
				}

				if flagPath == false {
					parametrosValidos = false
				}
				if parametrosValidos {
					fmt.Println("--->Parametros Invalidos ")
				} else {
					mkdir.setUsuarioAct = usuarioIniciado.nombreUsuario
					mkdir.setIdMontada = usuarioIniciado.idPartMontada
					EjecutarComandoMkdir(mkdir, listaDiscos)
				}
			} else {
				fmt.Println("ERROR: NO hay una sesion iniciada ")
				fmt.Println(" ")
			}

		} else if nombreComando == "pause" {
			fmt.Println("Pause. Presione una tecla para continuar...")
			fmt.Scanln()
		} else if nombreComando == "exec" {
			parametrosValidos := true
			flagPath := true
			pathExec := ""
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						pathExec = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						pathExec = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
					break
				}
			}
			if flagPath == false {
				parametrosValidos = false
			}
			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				ComandoExec(pathExec, listaDiscos)
			}
		} else if nombreComando == "rep" {
			parametrosValidos := true
			flagPath := true
			flagName := true
			flagId := true
			flagRuta := true

			reporte := PropRep{}
			for f := 0; f < len(comandoTemp.Propiedades); f++ {
				nombreProp := strings.ToLower(comandoTemp.Propiedades[f].Nombre)
				if nombreProp == "-path" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						reporte.setPath = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						reporte.setPath = comandoTemp.Propiedades[f].Valor
					}
					flagPath = false
				} else if nombreProp == "-name" {

					var1 := strings.ToLower(comandoTemp.Propiedades[f].Valor)
					if var1 == "disk" || var1 == "tree" || var1 == "file" {
						reporte.setName = comandoTemp.Propiedades[f].Valor
						flagName = false
					}

				} else if nombreProp == "-id" {
					flagId = false
					var1 := strings.ToUpper(comandoTemp.Propiedades[f].Valor)
					reporte.setId = var1

				} else if nombreProp == "-ruta" {
					if strings.Contains(comandoTemp.Propiedades[f].Valor, "\"") {
						conc := strings.Split(comandoTemp.Propiedades[f].Valor, "\"")
						reporte.setRuta = conc[1]
						fmt.Println("Path: " + conc[1])
					} else {
						reporte.setRuta = comandoTemp.Propiedades[f].Valor
					}
					flagRuta = false
				}
			}
			if flagPath == false && flagName == false && flagId == false {
				parametrosValidos = false
			}
			if flagRuta {
				reporte.setRuta = "-"
			}
			if parametrosValidos {
				fmt.Println("--->Parametros Invalidos ")
			} else {
				if reporte.setName == "disk" {
					EjecutarRepDisk(reporte, listaDiscos)
				} else if reporte.setName == "tree" {
					EjecutarRepTree(reporte, listaDiscos)
				} else if reporte.setName == "file" {
					if !flagRuta {
						EjecutarRepFile(reporte, listaDiscos)
					} else {
						fmt.Println("ERROR: se necesita una Ruta para este reporte ")
					}

				}
			}

		} else {
			fmt.Println("No se reconoce el comando ")

		}

	}
}

func ComandoExec(data string, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando Exec                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("Path: " + data + "\n")
	dat, err := ioutil.ReadFile(data)
	Check(err)
	LeerTexto(string(dat), listaDiscos)
}

func Check(e error) {
	if e != nil {
		fmt.Println("Error")
	}
}
