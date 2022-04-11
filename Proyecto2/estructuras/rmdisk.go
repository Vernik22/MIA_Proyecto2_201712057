package estructuras

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func EjecutarComandoRmdisk(diskRem string) {
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("                            Comando RMDISK                         ")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("")
	if archivoExiste(diskRem) {
		var comando string = ""
		fmt.Println("El archivo existe en: " + diskRem)
		fmt.Println("**Se borrara el disco, ¿Esta de acuerdo? [s/n]")
		fmt.Print(">> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		comando = scanner.Text()

		if comando == "S" || comando == "s" {
			executeComandRmdisk(diskRem)
		} else {
			fmt.Println("** NO se borro el disco **")
		}

	} else {
		fmt.Println("ERROR: el archivo que desea eliminar no existe")
	}
}

func executeComandRmdisk(path string) {
	cmd := exec.Command("rm", path)
	cmd.CombinedOutput()
}
