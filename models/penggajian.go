package models

import "gorm.io/gorm"

type Penggajian struct {
	gorm.Model
	NamaPegawai string  `json:"namapegawai"`
	GajiPokok   float64 `json:"gajipokok"`
	JamLembur   int     `json:"jamlembur"`
	GajiKotor   float64 `json:"gajikotor"`
	Pajak       float64 `json:"pajak"`
	GajiBersih  float64 `json:"gajibersih"`
}
