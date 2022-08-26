package logger

import "log"

func PANIC(message string, err error) {
	if err != nil {
		log.Panicln(message, err)
	}
}

func INFO(message string, data interface{}) {
	log.Panicln(message, data)
}
