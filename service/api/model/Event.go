package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

/*
Restituisce un oggetto del tipo:
[

	{
		"actor" : `userName coinvolto nell'evento. JOIN con userID per recuperare il nome :string`,
		"action" : `colonna type del db :string `,
		"timestamp": `datetime dell'evento :string`
	},

]
*/
func GetEventsByGroupID(groupID int) ([]map[string]interface{}, error) {
	logrus.Debug("Entered in GetEventsByGroupID()")
	logrus.Infof("Getting events for groupID %d", groupID)

	// Imposta i default
	var outErr error
	var events []map[string]interface{}

	// Ottieni gli eventi dal DB
	events, err := queries.GetEventsByGroupID(groupID)

	// Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("Events for groupID %d obtained successfully", groupID)
	} else {
		outErr = fmt.Errorf("error during obtaining events for groupID %d: %w", groupID, err)
		logrus.Error(outErr)
	}

	return events, outErr
}
