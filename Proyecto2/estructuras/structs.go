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
	Estado          [1]byte
	EstadoMks       [1]byte
}

//Estructuras para discos y particiones
type Particion struct {
	part_status [1]byte
	part_type   [1]byte
	part_fit    [1]byte
	part_start  int64
	part_size   int64
	part_name   [16]byte
}

type MBR struct {
	mbr_tamano         int64
	mbr_fecha_creacion [19]byte
	mbr_dsk_signature  int64
	mbr_dsk_fit        [1]byte
	mbre_partition_1   Particion
	mbre_partition_2   Particion
	mbre_partition_3   Particion
	mbre_partition_4   Particion
}

type EBR struct {
	part_status [1]byte
	part_fit    [1]byte
	part_start  int64
	part_size   int64
	part_next   int64
	part_name   [16]byte
}

type SupB struct {
	s_filesystem_type   int64
	s_inodes_count      int64
	s_blocks_count      int64
	s_free_blocks_count int64
	s_free_inodes_count int64
	s_mtime             [19]byte
	s_mnt_count         int64
	s_magic             int64
	s_inode_size        int64
	s_block_size        int64
	s_first_ino         int64 //Primer inodo libre
	s_first_blo         int64 //Primer bloque libre
	s_bm_inode_start    int64 //guardara el inicio del bitmap de inodos
	s_bm_block_start    int64 //guarda el inicio del bitmap de bloques
	s_inode_start       int64 //guarda donde empiezan los inodos
	s_block_start       int64
}

type Inodo struct {
	i_uid   int64
	i_gid   int64
	i_size  int64
	i_atime [19]byte  //ultima vez que se leyo sin modificar
	i_ctime [19]byte  //fecha en la que se creo el inodo
	i_mtime [19]byte  //ultima vez que se modifico
	i_block [15]int64 //cantidad de bloques que hay; apunta hacia el bloque apuntador que tiene los 16 apuntadores
	i_type  [1]byte   //indica si es carpeta o archivo :: 1=Archivo  0=Carpeta
	i_perm  int64
}

type Content struct {
	b_name  [12]byte
	b_inodo int64
}

type BCarpeta struct {
	b_content [4]Content
}

type BArchivo struct {
	b_content [64]byte
}

type BApun struct {
	b_apuntadores [4]Content
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
