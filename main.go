package main

import (
	"FinalGo/controller"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	fmt.Println(viper.Get("mysql.dsn"))
	dsn := viper.GetString("mysql.dsn")

	dialactor := mysql.Open(dsn)
	_, err = gorm.Open(dialactor)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connection successful")

	//Start Web Api
	controller.StartServer()
}
