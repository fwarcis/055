// Package logging...
package logging

import (
	"log"

)

func HandleForLogging(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

