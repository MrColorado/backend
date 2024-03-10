package logger

import (
	"errors"
)

var (
	config Configuration
)

func Init(c Configuration) error {
	if c.AppName == "" {
		return errors.New("logger: `AppName` configuration field cannot be empty")
	}
	config = c

	if err := initLogger(); err != nil {
		return err
	}

	return nil
}

func IsLoggerEnable() bool {
	return config.Logger != nil
}

func Stop() {
	stopLogger()

	config = Configuration{}
}
