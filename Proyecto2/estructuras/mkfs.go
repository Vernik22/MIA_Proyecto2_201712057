package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"time"
	"unsafe"
)

func EjecutarComandoMkfs(datFS PropMkfs, listaDiscos *[100]Mount) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKFS                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")

	if existeIdMount(datFS.setId, listaDiscos) {
		mbrTemp := MBR{}
		pathD := ""
		nombreParticion := ""
		for i := 0; i < 100; i++ {
			if datFS.setId == listaDiscos[i].IdMount {
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
			tamParticion, iniPart := returnDatosPart(mbrTemp, pathD, nombreParticion)
			superBloque := SupB{}
			inodoT := Inodo{}
			var n int64 = 0
			//Ext2
			n = (int64(tamParticion) - int64(unsafe.Sizeof(superBloque))) / (4 + int64(unsafe.Sizeof(inodoT)) + 3*64)
			//Numero de estructuras
			//cantBlockCarp := n
			//cantBlockArch := n
			cantidadInodos := n

			//Bitmaps
			inicioBitmapInodo := iniPart + int64(unsafe.Sizeof(superBloque))
			inicioBitmapBlockA := inicioBitmapInodo + cantidadInodos

			//Inicio bloques e inodos
			inicioInodo := inicioBitmapBlockA + 3*cantidadInodos
			inicioBloques := inicioInodo + cantidadInodos*int64(unsafe.Sizeof(inodoT))

			//iniciando valores de superbloque
			superBloque.S_filesystem_type = 2
			superBloque.S_inodes_count = cantidadInodos
			superBloque.S_blocks_count = cantidadInodos * 3
			superBloque.S_free_blocks_count = cantidadInodos * 3
			superBloque.S_free_inodes_count = cantidadInodos
			dt := time.Now()
			copy(superBloque.S_mtime[:], dt.String())
			superBloque.S_mnt_count = 1
			superBloque.S_magic = 0xEF53
			superBloque.S_inode_size = int64(unsafe.Sizeof(inodoT))
			superBloque.S_block_size = int64(unsafe.Sizeof(BArchivo{}))
			superBloque.S_first_ino = 0
			superBloque.S_first_blo = 0
			superBloque.S_bm_inode_start = inicioBitmapInodo
			superBloque.S_bm_block_start = inicioBitmapBlockA
			superBloque.S_inode_start = inicioInodo
			superBloque.S_block_start = inicioBloques

			f.Seek(iniPart, 0)
			err = binary.Write(f, binary.BigEndian, superBloque)
			inicializarBitmaps(pathD, iniPart, superBloque)
			crearRaiz(pathD, iniPart)
			fmt.Println("Se creo el sistema de archivos correctamente")
		}
	} else {
		fmt.Println("ERROR: el ID de la particion no existe, o no esta montada")
	}
}

func existeIdMount(idMount string, listaDiscos *[100]Mount) bool {
	for i := 0; i < 100; i++ {
		if idMount == listaDiscos[i].IdMount {
			return true
		}
	}
	return false
}

func crearRaiz(pathD string, iniPart int64) {
	dt := time.Now()
	f, err := os.OpenFile(pathD, os.O_RDWR, os.ModePerm)
	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		f.Seek(iniPart, 0)
		superB := SupB{}
		err = binary.Read(f, binary.BigEndian, &superB)
		var otro byte = '1'
		//llenar el primer bit del inodo de la carpeta raiz
		f.Seek(superB.S_bm_inode_start, 0)
		err = binary.Write(f, binary.BigEndian, otro)
		inodoTemp := Inodo{}
		for i := 0; i < 16; i++ {
			inodoTemp.I_block[i] = -1
		}
		copy(inodoTemp.I_atime[:], dt.String())
		copy(inodoTemp.I_ctime[:], dt.String())
		copy(inodoTemp.I_mtime[:], dt.String())
		inodoTemp.I_block[0] = 0
		inodoTemp.I_gid = 1
		inodoTemp.I_uid = 1
		inodoTemp.I_size = 27
		inodoTemp.I_perm = 777
		copy(inodoTemp.I_type[:], "0")
		f.Seek(superB.S_inode_start, 0)
		err = binary.Write(f, binary.BigEndian, inodoTemp)
		fmt.Println("Se creo la carpeta raiza (/)")

		f.Seek(superB.S_bm_inode_start+int64(unsafe.Sizeof(otro)), 0)
		err = binary.Write(f, binary.BigEndian, otro)
		inodoTemp.I_block[0] = 1
		inodoTemp.I_gid = 1
		inodoTemp.I_uid = 1
		inodoTemp.I_size = 27
		inodoTemp.I_perm = 777
		copy(inodoTemp.I_type[:], "1")
		f.Seek(superB.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
		err = binary.Write(f, binary.BigEndian, inodoTemp)
		//otro en bitmap de inodo
		f.Seek(superB.S_bm_inode_start+int64(unsafe.Sizeof(otro)), 0)
		err = binary.Write(f, binary.BigEndian, otro)
		//escribir 1 en el bitmap de carpeta y escribir carpeta
		f.Seek(superB.S_bm_block_start, 0)
		err = binary.Write(f, binary.BigEndian, otro)
		carpetRaiz := BCarpeta{}
		copy(carpetRaiz.B_content[0].B_name[:], ".")
		copy(carpetRaiz.B_content[1].B_name[:], "..")
		copy(carpetRaiz.B_content[2].B_name[:], "users.txt")
		carpetRaiz.B_content[0].B_inodo = 0
		carpetRaiz.B_content[1].B_inodo = 0
		carpetRaiz.B_content[2].B_inodo = 1
		carpetRaiz.B_content[3].B_inodo = -1
		f.Seek(superB.S_block_start, 0)
		err = binary.Write(f, binary.BigEndian, carpetRaiz)
		//escribir 1 en el bitmap de archivo y escribir el root en el users.txt
		f.Seek(superB.S_bm_block_start+int64(unsafe.Sizeof(otro)), 0)
		err = binary.Write(f, binary.BigEndian, otro)
		contenidoUsers := BArchivo{}
		copy(contenidoUsers.B_content[:], "1,G,root\n1,U,root,root,123\n")
		f.Seek(superB.S_block_start+int64(unsafe.Sizeof(BArchivo{})), 0)
		err = binary.Write(f, binary.BigEndian, contenidoUsers)
		fmt.Println("Se creo el archivo /users.txt")
		superB.S_free_blocks_count = superB.S_free_blocks_count - 2
		superB.S_free_inodes_count = superB.S_free_inodes_count - 2
		superB.S_first_blo = 2
		superB.S_first_ino = 2
		//actualizar superbloque
		f.Seek(iniPart, 0)
		err = binary.Write(f, binary.BigEndian, superB)
	}
}

func inicializarBitmaps(pathD string, iniPart int64, superB SupB) {
	f, err := os.OpenFile(pathD, os.O_RDWR, os.ModePerm)

	if err != nil {
		fmt.Println("No existe el archivo en la ruta")
	} else {
		var otro byte = '0'
		f.Seek(superB.S_bm_inode_start, 0)
		for i := 0; i < int(superB.S_inodes_count); i++ {
			err = binary.Write(f, binary.BigEndian, otro)
		}
		f.Seek(superB.S_bm_block_start, 0)
		for i := 0; i < int(superB.S_blocks_count); i++ {
			err = binary.Write(f, binary.BigEndian, otro)
		}

	}
}

func returnDatosPart(mbrTemp MBR, pathD string, nombrePart string) (int64, int64) {
	var parts [4]Particion
	parts[0] = mbrTemp.Mbr_partition_1
	parts[1] = mbrTemp.Mbr_partition_2
	parts[2] = mbrTemp.Mbr_partition_3
	parts[3] = mbrTemp.Mbr_partition_4
	noEncontrada := true

	var1 := nombreAString(mbrTemp.Mbr_partition_1.Part_name)
	var2 := nombreAString(mbrTemp.Mbr_partition_2.Part_name)
	var3 := nombreAString(mbrTemp.Mbr_partition_3.Part_name)
	var4 := nombreAString(mbrTemp.Mbr_partition_4.Part_name)
	if var1 == nombrePart || var2 == nombrePart || var3 == nombrePart || var4 == nombrePart {
		for i := 0; i < 4; i++ {
			var5 := nombreAString(parts[i].Part_name)
			if var5 == nombrePart {
				return parts[i].Part_size, parts[i].Part_start
				noEncontrada = false
				break
			}
		}
	}

	if noEncontrada {
		ebrExTemp := EBR{}
		f, err := os.OpenFile(pathD, os.O_RDWR, os.ModePerm)
		if err != nil {
			fmt.Println("No existe la ruta")
			return 0, 0
		}
		for i := 0; i < 4; i++ {
			if string(parts[i].Part_type[:]) == "e" {
				f.Seek(parts[i].Part_start, 0)
				err = binary.Read(f, binary.BigEndian, &ebrExTemp)
				break
			}
		}
		for ebrExTemp.Part_next != -1 {
			var2 := nombreAString(ebrExTemp.Part_name)
			if string(var2) == nombrePart {
				return ebrExTemp.Part_size, ebrExTemp.Part_start
			} else {
				f.Seek(ebrExTemp.Part_next, 0)
				err = binary.Read(f, binary.BigEndian, &ebrExTemp)
			}
		}
	}
	return 0, 0
}
