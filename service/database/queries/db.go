package queries

import "database/sql"

// db è il riferimento al database condiviso usato da tutte le query
var db *sql.DB

// Init inizializza il package queries con la connessione al database.
// Deve essere chiamato una volta all'avvio dell'applicazione.
func Init(database *sql.DB) {
	db = database
}
