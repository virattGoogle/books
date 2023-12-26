package storage

import(
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)


type Config struct{

	Host    string
	Port    string
	Password  string
	User      string
	DBName    string
	SSLMode   string
}

func NewConnection(config *Config)(*gorm.DB,error){
	dsn := fmt.Sprintf("host=%s port=%s password=%s dbname=%s user=%s sslmode=%s", config.Host,config.Port,config.Password,config.DBName,config.User,config.SSLMode)

	db , err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err !=nil {
		return db,err
	}

	return db,nil
}