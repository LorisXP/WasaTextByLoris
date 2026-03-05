package dml

import "database/sql"

// db è il riferimento al database condiviso usato da tutte le operazioni DML.
var db *sql.DB

// Init inizializza il package dml con la connessione al database.
// Deve essere chiamato una sola volta all'avvio dell'applicazione.
func Init(database *sql.DB) {
	db = database
}
