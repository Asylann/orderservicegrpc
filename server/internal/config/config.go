package config

import (
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DataBaseSource string
	Port           string
}

func getConnection() string {
	// find if env variables are passed
	return fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("DATABASE"))
}

func LoadConfig() (Config, error) {
	source := getConnection()
	if source == "host= port= user= password= database= sslmode=disable" {
		return Config{}, errors.New("Env variables are not passed!")
	}
	return Config{DataBaseSource: source, Port: os.Getenv("PORT")}, nil
}
