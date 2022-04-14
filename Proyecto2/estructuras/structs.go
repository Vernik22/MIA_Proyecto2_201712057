package estructuras

import "fmt"

func Formas() {
	fmt.Println("Hola Estructuras")
}

type Propiedad struct {
	Nombre string
	Valor  string
}

type Comando struct {
	Nombre      string
	Propiedades []Propiedad
}

type Mount struct {
	NombreParticion string
	IdMount         string
	Path            string
	EstadoMks       [1]byte
}

//Estructuras para discos y particiones
type Particion struct {
	Part_status [1]byte
	Part_type   [1]byte
	Part_fit    [1]byte
	Part_start  int64
	Part_size   int64
	Part_name   [16]byte
}

type MBR struct {
	Mbr_tamano         int64
	Mbr_fecha_creacion [19]byte
	Mbr_dsk_signature  int64
	Mbr_dsk_fit        [1]byte
	Mbr_partition_1    Particion
	Mbr_partition_2    Particion
	Mbr_partition_3    Particion
	Mbr_partition_4    Particion
}

type EBR struct {
	Part_status [1]byte
	Part_fit    [1]byte
	Part_start  int64
	Part_size   int64
	Part_next   int64
	Part_name   [16]byte
}

type SupB struct {
	S_filesystem_type   int64
	S_inodes_count      int64
	S_blocks_count      int64
	S_free_blocks_count int64
	S_free_inodes_count int64
	S_mtime             [19]byte
	S_mnt_count         int64
	S_magic             int64
	S_inode_size        int64
	S_block_size        int64
	S_first_ino         int64 //Primer inodo libre
	S_first_blo         int64 //Primer bloque libre
	S_bm_inode_start    int64 //guardara el inicio del bitmap de inodos
	S_bm_block_start    int64 //guarda el inicio del bitmap de bloques
	S_inode_start       int64 //guarda donde empiezan los inodos
	S_block_start       int64
}

type Inodo struct {
	I_uid   int64
	I_gid   int64
	I_size  int64
	I_atime [19]byte  //ultima vez que se leyo sin modificar
	I_ctime [19]byte  //fecha en la que se creo el inodo
	I_mtime [19]byte  //ultima vez que se modifico
	I_block [16]int64 //cantidad de bloques que hay; apunta hacia el bloque apuntador que tiene los 16 apuntadores
	I_type  [1]byte   //indica si es carpeta o archivo :: 1=Archivo  0=Carpeta
	I_perm  int64
}

type Content struct {
	B_name  [12]byte
	B_inodo int64
}

type BCarpeta struct {
	B_content [4]Content
}

type BArchivo struct {
	B_content [64]byte
}

type User struct {
	nombreUsuario string
	idPartMontada string
}

//
func BytesNombreParticion(data [16]byte) string {
	return string(data[:])
}
func ConvertData(data [64]byte) string {
	return string(data[:])
}

// encapsular los comandos
type PropMkdisk struct {
	setSize int
	setFit  string
	setUnit string
	setPath string
}

type PropFdisk struct {
	setSize int
	setUnit string
	setPath string
	setType string
	setFit  string
	setName string
}

type PropMount struct {
	setPath string
	setName string
}

type PropMkfs struct {
	setId   string
	setType string
}

type PropLogin struct {
	setUsuario  string
	setPassword string
	setId       string
}

type PropMkgrp struct {
	setUsuarioAct string
	setIdMontada  string
	setName       string
}

type PropMkusr struct {
	setUsuario    string
	setPassword   string
	setGrp        string
	setUsuarioAct string
}

type PropMkfile struct {
	setPath       string
	setSize       int
	setR          bool
	setCont       string
	setUsuarioAct string
}

type PropMkdir struct {
	setP          bool
	setPath       string
	setUsuarioAct string
}

type PropRep struct {
	setName string
	setPath string
	setId   string
	setRuta string
}
