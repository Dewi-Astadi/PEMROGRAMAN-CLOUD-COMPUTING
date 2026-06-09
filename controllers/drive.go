package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"pcc/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const googleScriptURL = "https://script.google.com/macros/s/AKfycbwsXsLdiHEZk3_m7nvPrpJJ_te0G528qpGqFOtmOt_n-Spfey-my367v0Euj-mU-8z9/exec"

func DriveUpload(c *gin.Context) {
	fileName := c.PostForm("fileName")
	if fileName == "" {
		fileName = c.PostForm("filename")
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "Gagal mengambil file dari form-data: " + err.Error(),
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")

	fileOpen, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": err.Error()})
		return
	}
	defer fileOpen.Close()

	fileData, _ := ioutil.ReadAll(fileOpen)
	data := base64.StdEncoding.EncodeToString(fileData)

	postBody, _ := json.Marshal(map[string]string{
		"filename": fileName,
		"fileName": fileName,
		"mimetype": mimeType,
		"mimeType": mimeType,
		"data":     data,
	})

	requestBody := bytes.NewBuffer(postBody)

	client := &http.Client{}
	req, err := http.NewRequest("POST", googleScriptURL, requestBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": false, "pesan": err.Error()})
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"kode_error": "ERR-DRIVE",
			"pesan":      "Gagal Terhubung ke Google Script",
			"error":      err.Error(),
		})
		return
	}
	defer res.Body.Close()

	hasilBody, _ := ioutil.ReadAll(res.Body)
	hasilString := string(hasilBody)

	fmt.Println("=== RESPONSE DARI GOOGLE SCRIPT ===")
	fmt.Println("Status Code:", res.StatusCode)
	fmt.Println("Response Body:", hasilString)
	fmt.Println("===================================")

	var hasilJson map[string]interface{}
	_ = json.Unmarshal([]byte(hasilString), &hasilJson)

	var rowsAffected int64
	if hasilJson != nil && hasilJson["fileId"] != nil {
		db := c.MustGet("db").(*gorm.DB)
		dokumenBaru := models.Dokumen{
			NamaDokumen: fileName,
			FileId:      hasilJson["fileId"].(string),
			FileUrl:     hasilJson["downloadUrl"].(string),
		}
		hasilDokumen := db.Create(&dokumenBaru)
		rowsAffected = hasilDokumen.RowsAffected
	}

	c.JSON(http.StatusOK, gin.H{
		"status":       true,
		"pesan":        "Proses kirim selesai",
		"data":         hasilJson,
		"tersimpan":    rowsAffected,
		"raw_response": hasilString,
	})
}

func DriveTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var dokumen []models.Dokumen
	db.Find(&dokumen)
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Tampil",
		"data":   dokumen,
	})
}

func DriveUnduh(c *gin.Context) {
	id := c.Param("id")

	res, err := http.Get(googleScriptURL + "?id=" + id)
	if err != nil {
		c.JSON(500, gin.H{
			"status": false,
			"pesan":  "Gagal Unduh",
		})
		return
	}

	hasilBody, _ := ioutil.ReadAll(res.Body)
	hasilString := string(hasilBody)

	var hasilJson map[string]interface{}
	json.Unmarshal([]byte(hasilString), &hasilJson)

	fileBase64 := hasilJson["file"].(string)
	mimeType := hasilJson["mimeType"].(string)

	file, _ := base64.StdEncoding.DecodeString(fileBase64)

	c.Writer.Header().Set("Content-Type", mimeType)
	c.Writer.Write(file)
}
