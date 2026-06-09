package controllers

import (
	"net/http"

	"pcc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Binding dari POST JSON
type StrukturPenggajian struct {
	Id          uint
	NamaPegawai string  `json:"namapegawai" binding:"required"`
	GajiPokok   float64 `json:"gajipokok" binding:"required"`
	JamLembur   int     `json:"jamlembur" binding:"required"`
}

type StrukturPenggajianHapus struct {
	Id uint `binding:"required"`
}

func TampilPenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelPenggajian []models.Penggajian
	hasil := db.Find(&modelPenggajian)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tampil Data",
			"kesalahan": kesalahan.Error(),
			"data":      nil,
		})
	}
}

func TambahPenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPenggajian StrukturPenggajian
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	// Hitung upah lembur, gaji kotor, pajak, gaji bersih
	upahLembur := float64(dataPenggajian.JamLembur) * 50000.0
	gajiKotor := dataPenggajian.GajiPokok + upahLembur

	var pajak float64
	if gajiKotor > 5000000 {
		pajak = gajiKotor * 0.05
	}

	gajiBersih := gajiKotor - pajak

	modelPenggajian := models.Penggajian{
		NamaPegawai: dataPenggajian.NamaPegawai,
		GajiPokok:   dataPenggajian.GajiPokok,
		JamLembur:   dataPenggajian.JamLembur,
		GajiKotor:   gajiKotor,
		Pajak:       pajak,
		GajiBersih:  gajiBersih,
	}

	hasil := db.Create(&modelPenggajian)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPenggajian,
		})
	}
}

func UbahPenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPenggajian StrukturPenggajian
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	// Hitung ulang
	upahLembur := float64(dataPenggajian.JamLembur) * 50000.0
	gajiKotor := dataPenggajian.GajiPokok + upahLembur

	var pajak float64
	if gajiKotor > 5000000 {
		pajak = gajiKotor * 0.05
	}

	gajiBersih := gajiKotor - pajak

	var modelPenggajian models.Penggajian
	db.First(&modelPenggajian, dataPenggajian.Id)
	modelPenggajian.NamaPegawai = dataPenggajian.NamaPegawai
	modelPenggajian.GajiPokok = dataPenggajian.GajiPokok
	modelPenggajian.JamLembur = dataPenggajian.JamLembur
	modelPenggajian.GajiKotor = gajiKotor
	modelPenggajian.Pajak = pajak
	modelPenggajian.GajiBersih = gajiBersih

	hasil := db.Save(&modelPenggajian)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data",
			"kesalahan": nil,
			"data":      modelPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelPenggajian,
		})
	}
}

func HapusPenggajian(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataPenggajian StrukturPenggajianHapus
	if err := c.ShouldBindJSON(&dataPenggajian); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelPenggajian models.Penggajian
	hasil := db.Delete(&modelPenggajian, dataPenggajian.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataPenggajian,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataPenggajian,
		})
	}
}
