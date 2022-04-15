package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
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
		inodoAnterior := Inodo{}
		f.Seek(superBlock.S_inode_start+int64(unsafe.Sizeof(Inodo{})), 0)
		err = binary.Read(f, binary.BigEndian, &inodoTemp)

		existeCarpeta := false
		carp := len(rutaArchivo)
	}
}
