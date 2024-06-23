package events

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/imedcl/manager-api/pkg/config"
	"github.com/imedcl/manager-api/pkg/data"
	"github.com/sirupsen/logrus"
)

type EventLog struct {
	UserNickname string `json:"user_nickname"`
	Resource     string `json:"resource"`
	Event        string `json:"event"`
	Sensor       string `json:"sensor"`
	PersonDni    string `json:"person_dni"`
	Params       string `json:"params"`
}

var db *data.DB

func Start(newDb *data.DB) {
	db = newDb
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

func PrettyParams(data interface{}) (string, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return strings.ReplaceAll(string(val), "\"", ""), nil
}

func (event *EventLog) Write() {
	config.LogfileInit()
	res, err := PrettyStruct(event)
	if err != nil {
		logrus.Fatal(err)
	}
	log.Println(res)

	newLog := &data.Log{
		Module:       event.Resource,
		UserNickname: event.UserNickname,
		Sensor:       event.Sensor,
		PersonDni:    event.PersonDni,
		Detail:       event.Event,
		Params:       event.Params,
	}

	db.CreateLog(newLog)

}
