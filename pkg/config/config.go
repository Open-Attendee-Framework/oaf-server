package config

import (
	"encoding/json"
	"errors"
	"os"
)

//Config is the struct to save and load the config file
type Config struct {
	Host               string             `json:"host"`
	Port               string             `json:"port"`
	DatabaseConnection DatabaseConnection `json:"databaseConnection"`
}

//DatabaseConnection holds the database driver aswell as the connection string
type DatabaseConnection struct {
	Driver     string `json:"driver"`
	Connection string `json:"connection"`
}

//LoadConfig accepts a filepath and tries to load a config file from there
func LoadConfig(filepath string) (Config, error) {
	var res Config
	file, err := os.Open(filepath)
	if err != nil {
		return res, errors.New("Error opening file: " + err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&res)
	if err != nil {
		return res, errors.New("Error decoding file: " + err.Error())
	}
	return res, nil
}
