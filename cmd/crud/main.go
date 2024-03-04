package main

import (
	"flag"
	"log"
	"os"
	"test-crud-goods/internal/server"
	"test-crud-goods/internal/utils/consts"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

var (
	envFile string
)

func init() {

	args := os.Args
	if len(args) == 1 {
		log.Fatal(color.RedString("Didn`t pass app mode: dev or prod"))
		return
	}
	mode := args[1]

	var env_mode_file string
	if mode == consts.DEV_MODE {
		log.Println(color.BlueString("run in development mode"))
		os.Setenv("APP_MODE", consts.DEV_MODE)
		env_mode_file = consts.DEV_MODE + ".env"
	}
	if mode == consts.PROD_MODE {
		log.Println(color.BlueString("run in production mode"))
		os.Setenv("APP_MODE", consts.DEV_MODE)
		env_mode_file = consts.PROD_MODE + ".env"
	}

	flag.StringVar(&envFile, "env", env_mode_file, "path to env file")
}

func main() {
	flag.Parse()

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal(color.RedString("Error loading .env file"))
	}
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
