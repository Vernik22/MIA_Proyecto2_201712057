package estructuras

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func EjecutarComandoMkdisk(datDiscos PropMkdisk) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando MKDISK                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	DirExist(datDiscos.setPath)
	if !archivoExiste(datDiscos.setPath) {
		dt := time.Now()
		mbrTemp := MBR{}
		copy(mbrTemp.mbr_fecha_creacion[:], []byte(dt.String()))
		mbrTemp.mbr_dsk_signature = int64(rand.Intn(200) + dt.Second())

		var file, _ = os.Create(datDiscos.setPath)

		defer file.Close()
		var buffer [1024]byte

		//simular un kb
		if datDiscos.setUnit == "k" {
			mbrTemp.mbr_tamano = int64(datDiscos.setSize) * 1024
			//se llena la variable buffer con ceros para que no este en null
			for i := 0; i < 1024; i++ {
				buffer[i] = '0'
			}
			file.Seek(0, 0)
			for i := 0; i < datDiscos.setSize; i++ {
				err := binary.Write(file, binary.BigEndian, buffer)
				if err != nil {
					log.Fatalln(err, datDiscos.setPath)
				}
			}
		} else if datDiscos.setUnit == "m" {
			mbrTemp.mbr_tamano = int64(datDiscos.setSize) * 1024 * 1024
			//se llena la variable buffer con ceros para que no este en null
			for i := 0; i < 1024; i++ {
				buffer[i] = '0'
			}
			file.Seek(0, 0)
			for i := 0; i < datDiscos.setSize*1024; i++ {
				err := binary.Write(file, binary.BigEndian, buffer)
				if err != nil {
					log.Fatalln(err, datDiscos.setPath)
				}
			}
		}

		if datDiscos.setFit == "ff" {
			copy(mbrTemp.mbr_dsk_fit[:], "f")
		} else if datDiscos.setFit == "wf" {
			copy(mbrTemp.mbr_dsk_fit[:], "w")
		} else if datDiscos.setFit == "bf" {
			copy(mbrTemp.mbr_dsk_fit[:], "b")
		}

		vacia := Particion{}
		copy(vacia.part_fit[:], "-")
		copy(vacia.part_name[:], "0")
		copy(vacia.part_type[:], "-")
		vacia.part_size = -1
		vacia.part_start = -1
		copy(vacia.part_status[:], "0")
		//asignar esta particion vacia al mbr creado
		mbrTemp.mbre_partition_1 = vacia
		mbrTemp.mbre_partition_2 = vacia
		mbrTemp.mbre_partition_3 = vacia
		mbrTemp.mbre_partition_4 = vacia

		//ahora se escribe el mbr creado en el archivo
		f, err := os.OpenFile(datDiscos.setPath, os.O_WRONLY, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln(err)
			}
		}()
		f.Seek(0, 0)
		err = binary.Write(f, binary.BigEndian, mbrTemp)
		if err != nil {
			log.Fatalln(err, datDiscos.setPath)
		}
		fmt.Println("Disco Creado Exitosamente")
		fmt.Println("")
	} else {
		fmt.Println("ERROR: el archivo que desea crear ya existe")
	}
}

func DirExist(dir string) {
	resultados := strings.Split(dir, "/")
	newpath := ""

	for i := 1; i < len(resultados)-1; i++ {
		//lenar el string con el path sin el disk.dk
		newpath += "/" + resultados[i]
	}
	//fmt.Println(newpath)
	executeComandMkdir(newpath)

}

func executeComandMkdir(path string) {
	cmd := exec.Command("mkdir", "-p", "{}", ":::", path)
	cmd.CombinedOutput()
}

func archivoExiste(ruta string) bool {
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		return false
	}
	return true
}
