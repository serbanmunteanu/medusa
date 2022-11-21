package env

import (
    "github.com/joho/godotenv"
    log "github.com/sirupsen/logrus"
)

func Load(filenames ...string) {
    if err := godotenv.Load(filenames...); err != nil {
        log.Warning("Could not load config from .env.dev file: " + err.Error())
    }
}
