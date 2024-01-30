package crow

import (
	"errors"
	"os"
)

var (
	config    Configuration
	useMetric bool
	useSentry bool
	useSlack  bool
)

func Init(c Configuration) error {
	if c.AppName == "" {
		return errors.New("crow: `AppName` configuration field cannot be empty")
	}
	config = c

	if os.Getenv("CROW_USE_METRIC") != "" {
		useMetric = true
	}

	if os.Getenv("CROW_USE_SENTRY") != "" {
		useSentry = true
	}

	if os.Getenv("CROW_USE_SLACK") != "" {
		useSlack = true
	}

	statsdHost := os.Getenv("CROW_METRIC_STATSD_HOST")
	statsdPort := os.Getenv("CROW_METRIC_STATSD_PORT")
	if statsdHost != "" && statsdPort != "" {
		config.Metric.StatsdURL = statsdHost + ":" + statsdPort
	}

	if err := initLogger(); err != nil {
		return err
	}
	if err := initMetric(); err != nil {
		return err
	}
	if err := initSlack(); err != nil {
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
	useMetric = false
	useSentry = false
	useSlack = false
}
