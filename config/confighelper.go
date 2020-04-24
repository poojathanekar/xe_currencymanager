package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.AddConfigPath(".")      // optionally look for config in the working directory
	err := viper.ReadInConfig()
	// Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file %s ", err))
	}
}

//GetConfig method to get configs from config file
func GetConfig(keyName string) string {
	keyValue := viper.GetString(keyName)
	return keyValue
}
func GetStringSlice(keyName string) []string {
	keyValue := viper.GetStringSlice(keyName)
	return keyValue
}

// load environment variable
func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file", key)
	}
	return os.Getenv(key)
}
