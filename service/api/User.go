package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// setMyUserName gestisce PATCH /api/users/{userID}/me/name
func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userId")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		// Decodifica il body JSON { "userName": "..." }
		var reqBody struct {
			UserName string `json:"userName"`
		}
		decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

		// Se il JSON è valido, prosegui.
		if decodeErr == nil {
			ctx.Logger.Info("request decoded successfully")

			// Valida l'input del nuovo userName con il Validator
			validationErr := model.ValidateInput(reqBody.UserName, 3, 15, `^[a-z]+[0-9]*$`, "string")
			
			//Se il campo inserito è valido
			if validationErr == nil {
				ctx.Logger.Info("request data is valid")

				//Ottieni prima l'utente
				user, errUsr := model.GetUser(userId)
				ctx.Logger.Debug("passed by GetUser()")
				
				//Se l'utente è ottenuto correttamente
				if errUsr == nil {
					ctx.Logger.Info("found user")

					//Aggiorna il nome
					errUpdateName := model.SetUserName(&user, reqBody.UserName)
					ctx.Logger.Debug("passed by SetUserName()")

					//Se il nome viene aggiornato correttamente
					if errUpdateName == nil {
						ctx.Logger.Info("userName udpated successfully")

						// Costruisci la risposta JSON
						respData := struct {}{}

						jsonBytes, marshalErr := json.Marshal(respData)
						
						//Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusOK // 200
							ctx.Logger.Debugf("new userName update successfully: %s (ID: %d)", reqBody.UserName, userId)
							ctx.Logger.Infof("new userName update successfully: %s", reqBody.UserName)

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else if errors.Is(errUpdateName, model.ErrUserNameAlreadyInUse) {
						statusCode = http.StatusConflict
						outErr = errUpdateName
						ctx.Logger.WithError(errUpdateName).Errorf("userName %s already in use", reqBody.UserName)
					} else {
						statusCode = http.StatusInternalServerError
						outErr = errUpdateName
						ctx.Logger.WithError(errUpdateName).Error("error updating userName")
					}


				} else {
					statusCode = http.StatusNotFound
					outErr = errUsr
					ctx.Logger.WithError(errUsr).Errorf("userID %d  not found", userId)
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

	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
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

// setMyPhoto gestisce PUT /api/users/{userID}/me/photo
func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userId")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		// Leggi il file "photo" dal form multipart
		// Limite: ~14MB (dimensione massima della foto codificata in base64 dalla YAML)
		parseErr := r.ParseMultipartForm(14 << 20)

		// Se il form è valido, prosegui.
		if parseErr == nil {
			ctx.Logger.Info("multipart form parsed successfully")

			// Estrai il file "photo" dal form
			file, _, fileErr := r.FormFile("photo")

			if fileErr == nil {
				defer file.Close()

				// Leggi il contenuto del file
				fileBytes, readErr := io.ReadAll(file)

				if readErr == nil {
					ctx.Logger.Info("photo file read successfully")

					// Codifica in base64
					photoBase64 := base64.StdEncoding.EncodeToString(fileBytes)

					//Ottieni prima l'utente
					user, errUsr := model.GetUser(userId)
					ctx.Logger.Debug("passed by GetUser()")

					//Se l'utente è ottenuto correttamente
					if errUsr == nil {
						ctx.Logger.Info("found user")

						//Aggiorna la foto
						errUpdatePhoto := model.SetPhoto(&user, photoBase64)
						ctx.Logger.Debug("passed by SetPhoto()")

						//Se la foto viene aggiornata correttamente
						if errUpdatePhoto == nil {
							ctx.Logger.Info("user profile picture udpated successfully")

							// Costruisci la risposta JSON
							respData := struct{}{}

							jsonBytes, marshalErr := json.Marshal(respData)

							//Se il JSON viene creato correttamente
							if marshalErr == nil {

								// Tutto ok
								responseBody = jsonBytes
								outErr = nil

								statusCode = http.StatusOK // 200
								ctx.Logger.Debugf("new user profile picture update successfully (ID: %d)", userId)
								ctx.Logger.Infof("new user profile picture update successfully")

							} else {
								statusCode = http.StatusInternalServerError
								outErr = marshalErr
								ctx.Logger.WithError(marshalErr).Error("error marshalling response")
							}

						} else {
							statusCode = http.StatusInternalServerError
							outErr = errUpdatePhoto
							ctx.Logger.WithError(errUpdatePhoto).Error("error updating user picture")
						}

					} else {
						statusCode = http.StatusNotFound
						outErr = errUsr
						ctx.Logger.WithError(errUsr).Errorf("userID %d not found", userId)
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = readErr
					ctx.Logger.WithError(readErr).Error("error reading photo file")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = fileErr
				ctx.Logger.WithError(fileErr).Error("missing or invalid 'photo' field in form")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = parseErr
			ctx.Logger.WithError(parseErr).Error("invalid multipart form")
		}

	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
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

// findUser gestisce GET /api/users/{userID}/other/{userName}
func (rt *_router) findUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	
	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userId")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		//Estrai lo userName e validalo
		userNameStr := ps.ByName("userName")
		userNameToSearch := strings.TrimSpace(userNameStr)

		// Valida l'input dello userName con il Validator
		validationErr := model.ValidateInput(userNameToSearch, 3, 15, `^[a-z]+[0-9]*$`, "string")

		//Se il campo inserito è valido
		if validationErr == nil {
			ctx.Logger.Info("request data is valid")

			//Ricevi i risultati
			users, errGetUsers := model.GetUsersByName(userNameToSearch)
			ctx.Logger.Debug("passed by GetUsersByName()")

			if errGetUsers == nil {
				ctx.Logger.Infof("found %d users matching '%s'", len(users), userNameToSearch)

				// Costruisci la risposta JSON come array di { userName, photo }
				type UserAndPhoto struct {
					UserName string `json:"userName"`
					Photo    string `json:"photo,omitempty"`
				}

				result := make([]UserAndPhoto, 0, len(users))
				for _, u := range users {
					result = append(result, UserAndPhoto{
						UserName: u.Name,
						Photo:    u.Photo,
					})
				}

				jsonBytes, marshalErr := json.Marshal(result)

				//Se il JSON viene creato correttamente
				if marshalErr == nil {

					// Tutto ok
					responseBody = jsonBytes
					outErr = nil
					statusCode = http.StatusOK // 200
					ctx.Logger.Infof("findUser completed: %d results", len(result))

				} else {
					statusCode = http.StatusInternalServerError
					outErr = marshalErr
					ctx.Logger.WithError(marshalErr).Error("error marshalling response")
				}

			} else if errors.Is(errGetUsers, model.ErrUserNotFound) {
				statusCode = http.StatusNotFound
				outErr = errGetUsers
				ctx.Logger.WithError(errGetUsers).Warnf("no users found matching '%s'", userNameToSearch)
			} else {
				statusCode = http.StatusInternalServerError
				outErr = errGetUsers
				ctx.Logger.WithError(errGetUsers).Error("error searching users")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = validationErr
			ctx.Logger.WithError(validationErr).Error("invalid userName. Check limits on doc")
		}

	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
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
