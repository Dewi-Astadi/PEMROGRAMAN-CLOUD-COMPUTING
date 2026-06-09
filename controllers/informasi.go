package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"pcc/models"
)

type StrukturInformasi struct {
	Id         uint
	Judul      string `binding:"required"`
	Konten     string `binding:"required"`
	UrlDokumen string `binding:"required"`
}

func InformasiTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var modelInformasi []models.Informasi
	hasil := db.Find(&modelInformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil Tampil data",
			"kesalahan": nil,
			"data":      modelInformasi,
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

func InformasiTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi
	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	modelInformasi := models.Informasi{
		Judul:      dataInformasi.Judul,
		Konten:     dataInformasi.Konten,
		UrlDokumen: dataInformasi.UrlDokumen,
	}
	hasil := db.Create(&modelInformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal Tambah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelInformasi,
		})
	}
}

func InformasiUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi
	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelInformasi models.Informasi
	db.First(&modelInformasi, dataInformasi.Id)
	modelInformasi.Judul = dataInformasi.Judul
	modelInformasi.Konten = dataInformasi.Konten
	modelInformasi.UrlDokumen = dataInformasi.UrlDokumen
	hasil := db.Save(&modelInformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil ubah data",
			"kesalahan": nil,
			"data":      modelInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal ubah Data",
			"kesalahan": kesalahan.Error(),
			"data":      modelInformasi,
		})
	}
}

func InformasiHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dataInformasi StrukturInformasi
	if err := c.ShouldBindJSON(&dataInformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelInformasi models.Informasi
	hasil := db.Delete(&modelInformasi, dataInformasi.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil hapus data",
			"kesalahan": nil,
			"data":      dataInformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal hapus Data",
			"kesalahan": kesalahan.Error(),
			"data":      dataInformasi,
		})
	}
}
