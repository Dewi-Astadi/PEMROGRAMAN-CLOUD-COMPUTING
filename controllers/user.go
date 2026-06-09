package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"pcc/models"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type StrukturUserTambah struct {
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StrukturUserUbah struct {
	Id       uint   `json:"id" binding:"required"`
	Nama     string `json:"nama" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type StrukturUserHapus struct {
	Id uint `json:"id" binding:"required"`
}

type StrukturLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func UserTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil tampil data",
		"data":   users,
	})
}

func UserTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserTambah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	modelUser := models.User{
		Nama:     dataUser.Nama,
		Username: dataUser.Username,
		Password: encryptedString,
	}

	hasil := db.Save(&modelUser)

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    true,
			"pesan":     "Berhasil tambah data",
			"kesalahan": nil,
			"data":      modelUser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal tambah Data",
			"kesalahan": hasil.Error.Error(),
			"data":      modelUser,
		})
	}
}

func UserUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserUbah
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelUser models.User
	cekUser := db.First(&modelUser, dataUser.Id)

	if cekUser.Error == nil {
		var sha = sha1.New()
		sha.Write([]byte(dataUser.Password))
		var encrypted = sha.Sum(nil)
		var encryptedString = fmt.Sprintf("%x", encrypted)

		modelUser.Nama = dataUser.Nama
		modelUser.Username = dataUser.Username
		modelUser.Password = encryptedString
		hasil := db.Save(&modelUser)

		if hasil.Error == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":    true,
				"pesan":     "Berhasil ubah data",
				"kesalahan": nil,
				"data":      modelUser,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":    false,
				"pesan":     "Gagal ubah Data",
				"kesalahan": hasil.Error.Error(),
				"data":      modelUser,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data Tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
			"data":      modelUser,
		})
	}
}

func UserHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserHapus
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":    false,
			"pesan":     "Gagal membaca Data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelUser models.User
	cekUser := db.First(&modelUser, dataUser.Id)

	if cekUser.Error == nil {
		hasil := db.Delete(&modelUser)

		if hasil.Error == nil {
			c.JSON(http.StatusOK, gin.H{
				"status":    true,
				"pesan":     "Berhasil hapus data",
				"kesalahan": nil,
				"data":      modelUser,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"status":    false,
				"pesan":     "Gagal hapus Data",
				"kesalahan": hasil.Error.Error(),
				"data":      modelUser,
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Data Tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
			"data":      modelUser,
		})
	}
}

func UserLogin(c *gin.Context) (any, error) {
	db := c.MustGet("db").(*gorm.DB)

	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" {
		username = c.Query("username")
		password = c.Query("password")
	}

	if username == "" || password == "" {
		return nil, jwtV3.ErrMissingLoginValues
	}

	var sha = sha1.New()
	sha.Write([]byte(password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	var modelUser models.User
	cekUser := db.Where("username = ?", username).
		Where("password = ?", encryptedString).
		First(&modelUser)

	if cekUser.Error == nil {
		return modelUser, nil
	}
	return nil, jwtV3.ErrFailedAuthentication
}
