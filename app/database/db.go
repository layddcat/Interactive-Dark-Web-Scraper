package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=db user=cti_user password=cti_pass dbname=cti_db port=5432 sslmode=disable"
	var db *gorm.DB
	var err error

	for i := 0; i < 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		fmt.Printf("Veritabanı henüz hazır değil, bekleniyor... (Deneme %d/5)\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Veritabanına 5 deneme sonunda bağlanılamadı: ", err)
	}

	db.AutoMigrate(&IntelligenceData{})
	DB = db
	fmt.Println("Veritabanı bağlantısı başarılı ve tablolar senkronize edildi.")
}