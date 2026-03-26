package api

import (
	"encoding/json"
	"net/http"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLogin gestisce POST /api/auth
// Se l'utente esiste, restituisce 200 con il suo userID.
// Se non esiste, lo crea e restituisce 201 con il nuovo userID.
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int
	var responseBody []byte
	var outErr error

	// Decodifica il body JSON { "userName": "..." }
	var reqBody struct {
		UserName string `json:"userName"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

	// Se il JSON è valido, prosegui.
	if decodeErr == nil {

		// Valida l'input userName con il Validator
		validationErr := model.ValidateInput(reqBody.UserName, 3, 15, `^[a-z]+[0-9]*$`, "string")

		// Se il campo inserito è valido
		if validationErr == nil {

			// Chiama il model per autenticare o creare l'utente
			userID, photo, newUser, authErr := model.AuthUser(reqBody.UserName)

			// Se l'autenticazione va a buon fine
			if authErr == nil {

				// Costruisci la risposta JSON con lo userID e la foto
				respData := struct {
					UserID int    `json:"userID"`
					Photo  string `json:"photo,omitempty"`
				}{UserID: userID, Photo: photo}

				jsonBytes, marshalErr := json.Marshal(respData)

				// Se il JSON viene creato correttamente
				if marshalErr == nil {

					// Tutto ok
					responseBody = jsonBytes

					if newUser {
						statusCode = http.StatusCreated // 201
						ctx.Logger.Debugf("new user created: %s (ID: %d)", reqBody.UserName, userID)
						ctx.Logger.Infof("new user created: %s", reqBody.UserName)
					} else {
						statusCode = http.StatusOK // 200
						ctx.Logger.Debugf("user authenticated: %s (ID: %d)", reqBody.UserName, userID)
						ctx.Logger.Infof("user authenticated: %s", reqBody.UserName)
					}

				} else {
					statusCode = http.StatusInternalServerError
					outErr = marshalErr
					ctx.Logger.WithError(marshalErr).Error("error marshalling response")
				}

			} else {
				statusCode = http.StatusInternalServerError
				outErr = authErr
				ctx.Logger.WithError(authErr).Error("authentication error")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = validationErr
			ctx.Logger.WithError(validationErr).Error("invalid userName. Check limits on doc")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = decodeErr
		ctx.Logger.WithError(decodeErr).Error("invalid request body")
	}

	// Unico punto di uscita
	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(responseBody)
	}
}
