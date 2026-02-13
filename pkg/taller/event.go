package taller

import (
	"github.com/celsiainternet/elvis/console"
	"github.com/celsiainternet/elvis/event"
	"github.com/celsiainternet/elvis/et"
)

func initEvents() {	
	err := event.Stack("<channel>", eventAction)
	if err != nil {
		console.Error(err)
	}

}

func eventAction(m event.EvenMessage) {
	data, err := et.ToJson(m.Data)
	if err != nil {
		console.Error(err)
	}

	console.Log("eventAction", data)
}
