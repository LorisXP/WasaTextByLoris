package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversations gestisce GET /api/users/{userID}/conversations
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	// Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		// Ottieni la lista di conversazioni
		result, errGetConv := model.GetConversationByUserID(userId)

		// Se il gruppo è ottenuto correttamente
		if errGetConv == nil {
			ctx.Logger.Info("conversation obtained successfully")

			// if result è null o vuoto, allora 404
			if len(result) != 0 {

				// Costruisci la risposta JSON
				jsonBytes, marshalErr := json.Marshal(result)

				// Se il JSON viene creato correttamente
				if marshalErr == nil {

					// Tutto ok
					responseBody = jsonBytes
					outErr = nil

					statusCode = http.StatusOK // 200
					ctx.Logger.Info("conversation retrieved successfully")

				} else {
					statusCode = http.StatusInternalServerError
					outErr = marshalErr
					ctx.Logger.WithError(marshalErr).Error("error marshalling response")
				}
			} else {
				statusCode = http.StatusNotFound
				outErr = fmt.Errorf("no conversation found for this user")
				ctx.Logger.Info("no conversations found for this user")
			}

		} else {
			statusCode = http.StatusInternalServerError
			outErr = errGetConv
			ctx.Logger.WithError(errGetConv).Error("error during retrieving conversation")
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

// getConversation gestisce GET /api/users/{userID}/conversations/users/{conversationID}/messages
func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	// Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		conversationIdStr := ps.ByName("conversationID")
		conversationId, errConv := strconv.Atoi(conversationIdStr)

		// Se non ci sono errori nell'estrazione
		if errConv == nil || conversationId > 0 {
			ctx.Logger.Info("conversationID parsed successfully")

			// Ottieni prima la conversazione
			conversation, errConv := model.GetConversationByID(conversationId, true)
			ctx.Logger.Debug("Passed by GetConversationByID()")

			// Se conversation è ottenuto correttamente
			if errConv == nil {
				ctx.Logger.Info("conversation obtained successfully")

				// Ottieni la lista dei messaggi
				messages, errGetMess := model.GetListMessages(&conversation, userId)
				ctx.Logger.Debug("Passed by GetListMessages()")

				// Se hai ricevuto i messaggi
				if errGetMess == nil {
					// Costruisci la risposta JSON
					jsonBytes, marshalErr := json.Marshal(messages)

					// Se il JSON viene creato correttamente
					if marshalErr == nil {

						// Tutto ok
						responseBody = jsonBytes
						outErr = nil

						statusCode = http.StatusOK // 200
						ctx.Logger.Infof("messages list obtained successfully between users")

					} else {
						statusCode = http.StatusInternalServerError
						outErr = marshalErr
						ctx.Logger.WithError(marshalErr).Error("error marshalling response")
					}
				} else {
					statusCode = http.StatusInternalServerError
					outErr = errGetMess
					ctx.Logger.WithError(errGetMess).Error("error obtaining messages")
				}

			} else {
				statusCode = http.StatusNotFound
				outErr = errConv
				ctx.Logger.WithError(errConv).Error("conv not found")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errConv
			ctx.Logger.WithError(errConv).Error("invalid conversationId")
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

// getConversation gestisce GET /api/users/{userID}/conversations/users/{conversationID}/messages
func (rt *_router) getConversationGroups(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	// Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		conversationIdStr := ps.ByName("conversationID")
		conversationId, errConv := strconv.Atoi(conversationIdStr)

		// Se non ci sono errori nell'estrazione
		if errConv == nil || conversationId > 0 {
			ctx.Logger.Info("conversationID parsed successfully")

			// Ottieni prima la conversazione
			conversation, errConv := model.GetConversationByID(conversationId, false)
			ctx.Logger.Debug("Passed by GetConversationByID()")

			// Se la conversation è ottenuto correttamente
			if errConv == nil {
				ctx.Logger.Info("conversation obtained successfully")

				// Ottieni la lista dei messaggi
				messages, errGetMess := model.GetListMessages(&conversation, userId)
				ctx.Logger.Debug("Passed by GetListMessages()")

				// Se hai ricevuto i messaggi
				if errGetMess == nil {
					// Costruisci la risposta JSON
					jsonBytes, marshalErr := json.Marshal(messages)

					// Se il JSON viene creato correttamente
					if marshalErr == nil {

						// Tutto ok
						responseBody = jsonBytes
						outErr = nil

						statusCode = http.StatusOK // 200
						ctx.Logger.Infof("messages list obtained successfully between users and groups")

					} else {
						statusCode = http.StatusInternalServerError
						outErr = marshalErr
						ctx.Logger.WithError(marshalErr).Error("error marshalling response")
					}
				} else {
					statusCode = http.StatusInternalServerError
					outErr = errGetMess
					ctx.Logger.WithError(errGetMess).Error("error obtaining messages")
				}

			} else {
				statusCode = http.StatusNotFound
				outErr = errConv
				ctx.Logger.WithError(errConv).Error("conv not found")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errConv
			ctx.Logger.WithError(errConv).Error("invalid conversationId")
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

// createConversationUsers gestisce POST /api/users/{userID}/conversations/users
func (rt *_router) createConversationUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	// Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		// Decodifica il body JSON
		var reqBody struct {
			UserName string `json:"userName"`
		}
		decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

		// Se il JSON è valido, prosegui.
		if decodeErr == nil {
			ctx.Logger.Info("request decoded successfully")

			// Valida l'input di userName
			validationErr := model.ValidateInput(reqBody.UserName, 3, 15, `^[a-z]+[0-9]*$`, "string")

			// Se è corretto
			if validationErr == nil {
				ctx.Logger.Info("valid userName for new conversation")

				// Ottieni lo userID dello userName
				receiverId, errorReceiverId := model.GetUserIdByName(reqBody.UserName)
				ctx.Logger.Debug("passed by GetUserIdByName()")

				// Se la ricerca è stata effettuata senza errori
				if errorReceiverId == nil {
					ctx.Logger.Debug("research of userName ok")

					// Se hai trovato un riscontro
					if receiverId > 0 {
						ctx.Logger.Info("found userID by userName")

						// Crea la conversazione (o restituisci quella esistente)
						conversation, errCreate := model.CreateConversation(true, userId, receiverId)
						ctx.Logger.Debug("passed by CreateConversation()")

						// Se ci sei riuscito
						if errCreate == nil {
							ctx.Logger.Info("conversation created/retrieved successfully")

							// Costruisci la risposta JSON
							respData := struct {
								ConversationID int `json:"conversationID"`
							}{ConversationID: conversation.ConversationID}

							jsonBytes, marshalErr := json.Marshal(respData)

							// Se il JSON viene creato correttamente
							if marshalErr == nil {

								// Tutto ok
								responseBody = jsonBytes
								outErr = nil

								statusCode = http.StatusCreated // 201 - conversazione creata
								ctx.Logger.Info("new conversation created successfully")

							} else {
								statusCode = http.StatusInternalServerError
								outErr = marshalErr
								ctx.Logger.WithError(marshalErr).Error("error marshalling response")
							}

						} else {
							statusCode = http.StatusInternalServerError
							outErr = errCreate
							ctx.Logger.WithError(errCreate).Error("unable to create conversation")
						}

					} else {
						statusCode = http.StatusNotFound
						outErr = fmt.Errorf("userName %s doesn't exists on WasaText", reqBody.UserName)
						ctx.Logger.WithError(outErr)
					}

				} else {
					statusCode = http.StatusNotFound
					outErr = fmt.Errorf("userName %s not found on WasaText", reqBody.UserName)
					ctx.Logger.WithError(errorReceiverId).Infof("userName %s not found", reqBody.UserName)
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = validationErr
				ctx.Logger.WithError(validationErr).Error("invalid userName")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = decodeErr
			ctx.Logger.WithError(decodeErr).Error("error decoding request body")
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
