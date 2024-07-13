package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	Port = "6969"
	fileName = "./static/quran-final.json"
)

type Quran struct {
	Surahs []Surah `json:"quran"`
}

type Surah struct {
	SurahNumber int `json:"surahNumber"`
	SurahName string `json:"surahName"`
	Ayahs []Ayah `json:"ayahs"`
}

type Ayah struct {
	AyahNumber int `json:"ayahNumber"`
	AyahText string `json:"ayahText"`
}

func (quran *Quran) getSurahByNumber(surahNumber int) Surah {
	if surahNumber < 1 || surahNumber > 114 {
		return Surah{}
	}
	return quran.Surahs[surahNumber-1]
}

func (quran *Quran) getAyah(surahNumber int, ayahNumber int) Ayah {
	surah := quran.getSurahByNumber(surahNumber)
	if surah.SurahNumber == 0 {
		return Ayah{}
	}
	if ayahNumber < 1 || ayahNumber > len(surah.Ayahs) {
		return Ayah{}
	}
	return surah.Ayahs[ayahNumber-1]
}

func (quran *Quran) getSurahLength(surahNumber int) int {
	// TODO: apply proper restrictions on surahNumber
	if surahNumber < 1 || surahNumber > 114 {
		return 0
	}
	surah := quran.getSurahByNumber(surahNumber)
	return len(surah.Ayahs)
}

func (quran *Quran) getSurahName(surahNumber int) string {
	surah := quran.getSurahByNumber(surahNumber)
	return surah.SurahName
}

func main() {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var quran Quran
	json.Unmarshal(byteValue, &quran)


	// serve
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/:surahNumber", func(c *gin.Context) {
		surahNumber, err := strconv.Atoi(c.Param("surahNumber"))
		if err != nil || surahNumber < 1 || surahNumber > 114 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid surah number",
			})
			return
		}

		surah := quran.getSurahByNumber(surahNumber)

		c.JSON(http.StatusOK, gin.H{
			"surah": surah,
		})
	})

	router.GET("/:surahNumber/:ayahNumber", func(c *gin.Context) {
		surahNumber, err := strconv.Atoi(c.Param("surahNumber"))
		if err != nil || surahNumber < 1 || surahNumber > 114 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid surah number",
			})
			return
		}

		ayahNumber, err := strconv.Atoi(c.Param("ayahNumber"))
		if err != nil || ayahNumber < 1 || ayahNumber > quran.getSurahLength(surahNumber) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid ayah number",
			})
			return
		}

		ayah := quran.getAyah(surahNumber, ayahNumber)
		surahName := quran.getSurahName(surahNumber)

		c.JSON(http.StatusOK, gin.H{
			"ayah": ayah,
			"surahName": surahName,
		})
	})

	router.GET("/random", func(c *gin.Context) {
		randomSurah := 1 + (rand.Int() % 114)
		randomAyah := 1 + (rand.Int() % quran.getSurahLength(randomSurah))
		
		ayah := quran.getAyah(randomSurah, randomAyah)
		surahName := quran.getSurahName(randomSurah)

		c.JSON(http.StatusOK, gin.H{
			"ayah": ayah,
			"surahName": surahName,
		})
	})

	log.Println("Server started at :"+Port)
	router.Run(":"+Port)
}
