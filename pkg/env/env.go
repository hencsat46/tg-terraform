package env

import (
	"bufio"
	"log"
	"os"
	"strings"

	errir "github.com/VanLavr/tg-bot/pkg/error"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

type EnvConfigurator struct {
	path string
}

func New(path string) *EnvConfigurator {
	return &EnvConfigurator{
		path: path,
	}
}

func loadEnvFromFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func GetTGURL(path string) (string, error) {
	loadEnvFromFile(path)
	value := os.Getenv("TGREQUEST")
	if value == "" {
		return "", errir.ErrBadEnvLoading
	}

	log.Println("env has been loaded")
	return value, nil
}

func GetPort(path string) (string, error) {
	loadEnvFromFile(path)
	value := os.Getenv("PORT")
	if value == "" {
		return "", errir.ErrBadEnvLoading
	}

	log.Println("env has been loaded")
	return value, nil
}
