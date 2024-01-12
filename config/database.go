package config

import (
	"fmt"
	"log"
	"my_vocab/models"

	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DB_USER     = "myvocab"
	DB_PASSWORD = "myvocab123"
	DB_NAME     = "myvocab_db"
	DB_HOST     = "vocab-db.cbuoaypgqh0v.us-east-1.rds.amazonaws.com"
	DB_PORT     = "3306"
)

// const (
// 	DB_USER     = "root"
// 	DB_PASSWORD = "flutter123"
// 	DB_NAME     = "myvocab_db"
// 	DB_HOST     = "127.0.0.1"
// 	DB_PORT     = "3306"
// )

// const (
// 	DB_USER     = "root"
// 	DB_PASSWORD = ""
// 	DB_NAME     = "myvocab_db"
// 	DB_HOST     = "localhost"
// 	DB_PORT     = "3306"
// )

var DB *gorm.DB

func InitDb() {
	dbInfo := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)
	var err error
	DB, err = gorm.Open(mysql.Open(dbInfo), &gorm.Config{})

	if err != nil {
		log.Print("Connection Failed : ", err)
	}

	if err == nil {
		MigrateDB()
		SeedDB()
		log.Println("Connected!")
	}

}

func MigrateDB() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.TypeVocab{})
	DB.AutoMigrate(&models.Vocab{})

}

func SeedDB() {
	timeNow := time.Now()

	var listTypeVocabs = []models.TypeVocab{
		{
			Type:        "Noun (Kata Benda)",
			Description: "Noun atau kata benda digunakan untuk nama orang, benda, hewan, tempat, dan ide atau konsep. Noun sendiri dapat dibagi lagi ke berbagai jenis, seperti countable, uncountable / mass, common, proper, concrete, abstract, dan collective noun.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Pronoun (Kata Ganti)",
			Description: "Pronoun adalah kata yang digunakan untuk menggantikan noun. Terdapat 8 jenis pronoun, yaitu personal, demonstrative, interrogative, indefinite, possessive, reciprocal, relative, reflexive, dan intensive pronoun.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Verb (Kata Kerja)",
			Description: "Verb adalah kata kerja yang digunakan untuk menunjukkan tindakan atau keadaan. Verb dapat dibagi ke dalam beberapa jenis, seperti action verb dan linking verb.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Adjective (Kata Sifat)",
			Description: "Adjective adalah kata sifat yang digunakan untuk memberi keterangan pada noun atau pronoun.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Adverb (Kata Keterangan)",
			Description: "Sebagai kata keterangan, fungsi adverb adalah untuk memberikan keterangan tambahan pada verb, adjective, atau adverb itu sendiri. Adverb juga bisa dikelompokkan menjadi beberapa jenis, seperti manner, degree, frequency, place, dan time.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Preposition (Kata Depan)",
			Description: "Fungsi preposition adalah untuk menunjukkan hubungan antara noun dan kata-kata lainnya dalam sebuah kalimat.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Conjunction (Kata Hubung)",
			Description: "Conjunction digunakan untuk menghubungkan dua kata, frasa, klausa hingga kalimat. Terdapat 3 jenis conjuction, yaitu coordinating, subordinating, dan correlative conjuction.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
		{
			Type:        "Interjection (Kata Seru)",
			Description: "Jenis kata yang satu ini biasanya digunakan untuk mengungkapkan emosi.",
			CreatedAt:   timeNow,
			UpdatedAt:   timeNow,
		},
	}

	for _, value := range listTypeVocabs {
		// check type vocab
		checkTypeVocab := models.TypeVocab{}
		DB.Model(models.TypeVocab{}).Where("type = ?", value.Type).First(&checkTypeVocab)
		if checkTypeVocab.Type == "" {
			err := DB.Save(&value).Error
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}
