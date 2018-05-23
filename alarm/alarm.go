package alarm

import (
	"log"

	"fmt"

	"github.com/spf13/viper"
)

func Alarm(server string) {
	var alarms []Config
	err := viper.UnmarshalKey("alarms", &alarms)
	if err != nil {
		log.Println("Can't parse alarm config!")
		return
	}
	for _, a := range alarms {
		switch a.Type {
		case "smtp":
			var s SMTPConfig
			err := viper.UnmarshalKey("smtp", &s)
			if err != nil {
				log.Println("Can't parse smtp config!")
				return
			}
			s.Send(a.Targets, fmt.Sprintf("Server %s check failed!", server))

		default:
			log.Println("Alarm type not support!")

		}
	}
}
