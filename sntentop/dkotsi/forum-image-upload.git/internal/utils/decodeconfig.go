package utils

import (
	"encoding/json"
	"forum-image-upload/internal/backend/models"
	"log"
	"os"
)

func DecodeConf() (conf models.Config) {

	configfile, err := os.Open("../../forum-image-upload_config.json")
	if err != nil {
		log.Println(err)
	}
	defer configfile.Close()

	decoder := json.NewDecoder(configfile)
	if err := decoder.Decode(&conf); err != nil {
		log.Println(err)
	}
	return conf
}
