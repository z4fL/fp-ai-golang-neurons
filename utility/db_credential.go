package utility

import (
	"os"
	"strconv"

	"github.com/z4fL/fp-ai-golang-neurons/model"
)

func GetDBCredential() (*model.DBCredential, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}

	return &model.DBCredential{
		Host:         os.Getenv("DB_HOST"),
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
		Port:         port,
		Schema:       os.Getenv("DB_SCHEMA"),
	}, nil
}
