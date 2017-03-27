package log_test

import (
	"ggs/log"
)

func Example() {
	name := "ggs"

	log.Debug("My name is %v", name)
	log.Info("My name is %v", name)
	log.Error("My name is %v", name)
}
