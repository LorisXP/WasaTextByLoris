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

// sendMessage gestisce POST /api/users/{userID}/conversations/users/{conversationID}/messages
func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			// Decodifica il body JSON
			var reqBody struct {
				Content     string `json:"content"`
				ContentType string `json:"content_type"`
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

			// Se il JSON è valido, prosegui.
			if decodeErr == nil {
				ctx.Logger.Info("request decoded successfully")

				// Valida l'input di content_type
				validationContentTypeErr := model.ValidateInput(reqBody.ContentType, 3, 5, `^(gif|text|photo)$`, "string")

				// Se è corretto
				if validationContentTypeErr == nil {
					ctx.Logger.Info("valid content_type for new message")

					validationContentErr := fmt.Errorf("enum not recognized: %s", reqBody.ContentType)

					// Valida anche content
					if reqBody.ContentType == "text" {
						ctx.Logger.Info("validating text...")

						// Valida l'input di content
						validationContentErr = model.ValidateInput(reqBody.Content, 0, 4095, `^.*?$`, "string")

					} else if reqBody.ContentType == "gif" || reqBody.ContentType == "photo" {
						ctx.Logger.Info("validating media...")

						// Valida l'input di content
						validationContentErr = model.ValidateInput(reqBody.Content, 0, 13981013, `^[A-Za-z0-9+/]+={0,2}$`, "string")

					} else {
						outErr = validationContentErr
						ctx.Logger.WithError(validationContentErr).Error("invalid enum in content_type")
					}

					// Se è corretto
					if validationContentErr == nil {
						ctx.Logger.Info("request validated successfully")

						createdMessage, errCreateMessage := model.CreateMessage(userId, conversationId, true, reqBody.Content, reqBody.ContentType)

						// Se ci sei riuscito
						if errCreateMessage == nil {
							ctx.Logger.Info("new message created successfully")

							// Costruisci la risposta JSON
							respData := struct {
								MessageID int `json:"messageID"`
							}{MessageID: createdMessage.MessageID}

							jsonBytes, marshalErr := json.Marshal(respData)

							// Se il JSON viene creato correttamente
							if marshalErr == nil {

								// Tutto ok
								responseBody = jsonBytes
								outErr = nil

								statusCode = http.StatusCreated // 201
								ctx.Logger.Debugf("new message created successfully %s (ID: %d)", reqBody.Content, userId)
								ctx.Logger.Info("new message created successfully")

							} else {
								statusCode = http.StatusInternalServerError
								outErr = marshalErr
								ctx.Logger.WithError(marshalErr).Error("error marshalling response")
							}

						} else {
							statusCode = http.StatusInternalServerError
							outErr = errCreateMessage
							ctx.Logger.WithError(errCreateMessage).Error("error during creating message")
						}
					} else {
						statusCode = http.StatusBadRequest
						outErr = validationContentTypeErr
						ctx.Logger.WithError(validationContentTypeErr).Error("invalid message content")
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = validationContentTypeErr
					ctx.Logger.WithError(validationContentTypeErr).Error("invalid enum in content_type")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = decodeErr
				ctx.Logger.WithError(decodeErr).Error("invalid request body")
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

// sendMessageGroup gestisce POST /api/users/{userID}/conversations/groups/{conversationID}/messages
func (rt *_router) sendMessageGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			// Decodifica il body JSON
			var reqBody struct {
				Content     string `json:"content"`
				ContentType string `json:"content_type"`
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

			// Se il JSON è valido, prosegui.
			if decodeErr == nil {
				ctx.Logger.Info("request decoded successfully")

				// Valida l'input di content_type
				validationContentTypeErr := model.ValidateInput(reqBody.ContentType, 3, 5, `^(gif|text|photo)$`, "string")

				// Se è corretto
				if validationContentTypeErr == nil {
					ctx.Logger.Info("valid content_type for new message")

					validationContentErr := fmt.Errorf("enum not recognized: %s", reqBody.ContentType)

					// Valida anche content
					if reqBody.ContentType == "text" {
						ctx.Logger.Info("validating text...")

						// Valida l'input di content
						validationContentErr = model.ValidateInput(reqBody.Content, 0, 4095, `^.*?$`, "string")

					} else if reqBody.ContentType == "gif" || reqBody.ContentType == "photo" {
						ctx.Logger.Info("validating media...")

						// Valida l'input di content
						validationContentErr = model.ValidateInput(reqBody.Content, 0, 13981013, `^[A-Za-z0-9+/]+={0,2}$`, "string")

					} else {
						outErr = validationContentErr
						ctx.Logger.WithError(validationContentErr).Error("invalid enum in content_type")
					}

					// Se è corretto
					if validationContentErr == nil {
						ctx.Logger.Info("request validated successfully")

						createdMessage, errCreateMessage := model.CreateMessage(userId, conversationId, false, reqBody.Content, reqBody.ContentType)

						// Se ci sei riuscito
						if errCreateMessage == nil {
							ctx.Logger.Info("new message to group created successfully")

							// Costruisci la risposta JSON
							respData := struct {
								MessageID int `json:"messageID"`
							}{MessageID: createdMessage.MessageID}

							jsonBytes, marshalErr := json.Marshal(respData)

							// Se il JSON viene creato correttamente
							if marshalErr == nil {

								// Tutto ok
								responseBody = jsonBytes
								outErr = nil

								statusCode = http.StatusCreated // 201
								ctx.Logger.Debugf("new message to group created successfully %s (ID: %d)", reqBody.Content, userId)
								ctx.Logger.Info("new message to group created successfully")

							} else {
								statusCode = http.StatusInternalServerError
								outErr = marshalErr
								ctx.Logger.WithError(marshalErr).Error("error marshalling response")
							}

						} else {
							statusCode = http.StatusInternalServerError
							outErr = errCreateMessage
							ctx.Logger.WithError(errCreateMessage).Error("error during creating message to group")
						}
					} else {
						statusCode = http.StatusBadRequest
						outErr = validationContentTypeErr
						ctx.Logger.WithError(validationContentTypeErr).Error("invalid message content")
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = validationContentTypeErr
					ctx.Logger.WithError(validationContentTypeErr).Error("invalid enum in content_type")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = decodeErr
				ctx.Logger.WithError(decodeErr).Error("invalid request body")
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

// forwardMessage gestisce POST /api/users/{userID}/conversations/users/{conversationID}/messages/{messageID}
func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")

				// Decodifica il body JSON per source_type
				var reqBody struct {
					SourceType string `json:"source_type"`
				}
				decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)
				if decodeErr != nil || (reqBody.SourceType != "users" && reqBody.SourceType != "groups") {
					statusCode = http.StatusBadRequest
					outErr = fmt.Errorf("invalid or missing source_type: must be 'users' or 'groups'")
					ctx.Logger.WithError(outErr).Error("invalid source_type")
				} else {
					sourceBetweenUsers := reqBody.SourceType == "users"

					// Inoltra il messaggio
					forwardedMessage, errForwardedMessage := model.ForwardMessage(userId, conversationId, true, sourceBetweenUsers, messageId)

					// Se ci sei riuscito
					if errForwardedMessage == nil {
						ctx.Logger.Info("message forwarded successfully")

						// Costruisci la risposta JSON
						respData := struct {
							MessageID        int `json:"messageID"`
							ReplyToMessageID int `json:"replyToMessageID"`
						}{MessageID: forwardedMessage.MessageID, ReplyToMessageID: messageId}

						jsonBytes, marshalErr := json.Marshal(respData)

						// Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusCreated // 201
							ctx.Logger.Infof("new message (%d) forwarded successfully (%d)", forwardedMessage.MessageID, messageId)

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errForwardedMessage
						ctx.Logger.WithError(errForwardedMessage).Error("error during forwarding message")
					}
				} // chiude else source_type

			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
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

// forwardMessageGroup gestisce POST /api/users/{userID}/conversations/groups/{conversationID}/messages/{messageID}
func (rt *_router) forwardMessageGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")

				// Decodifica il body JSON per source_type
				var reqBody struct {
					SourceType string `json:"source_type"`
				}
				decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)
				if decodeErr != nil || (reqBody.SourceType != "users" && reqBody.SourceType != "groups") {
					statusCode = http.StatusBadRequest
					outErr = fmt.Errorf("invalid or missing source_type: must be 'users' or 'groups'")
					ctx.Logger.WithError(outErr).Error("invalid source_type")
				} else {
					sourceBetweenUsers := reqBody.SourceType == "users"

					// Inoltra il messaggio
					forwardedMessage, errForwardedMessage := model.ForwardMessage(userId, conversationId, false, sourceBetweenUsers, messageId)

					// Se ci sei riuscito
					if errForwardedMessage == nil {
						ctx.Logger.Info("message forwarded successfully to group")

						// Costruisci la risposta JSON
						respData := struct {
							MessageID        int `json:"messageID"`
							ReplyToMessageID int `json:"replyToMessageID"`
						}{MessageID: forwardedMessage.MessageID, ReplyToMessageID: messageId}

						jsonBytes, marshalErr := json.Marshal(respData)

						// Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusCreated // 201
							ctx.Logger.Infof("new message (%d) forwarded successfully (%d) to group", forwardedMessage.MessageID, messageId)

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errForwardedMessage
						ctx.Logger.WithError(errForwardedMessage).Error("error during forwarding message to group")
					}
				} // chiude else source_type

			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
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

// deleteMessage gestisce DELETE /api/users/{userID}/conversations/users/{conversationID}/messages/{messageID}
func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")

				// Cancella il messaggio
				errDeleteMessage := model.DeleteMessage(userId, conversationId, true, messageId)

				// Se ci sei riuscito
				if errDeleteMessage == nil {
					ctx.Logger.Info("message deleted successfully")

					// Costruisci la risposta JSON
					respData := struct {
						MessageID int `json:"messageID"`
					}{MessageID: messageId}

					jsonBytes, marshalErr := json.Marshal(respData)

					// Se il JSON viene creato correttamente
					if marshalErr == nil {

						// Tutto ok
						responseBody = jsonBytes
						outErr = nil

						statusCode = http.StatusOK // 200
						ctx.Logger.Infof(" message (%d) deleted successfully", messageId)

					} else {
						statusCode = http.StatusInternalServerError
						outErr = marshalErr
						ctx.Logger.WithError(marshalErr).Error("error marshalling response")
					}

				} else {
					statusCode = http.StatusInternalServerError
					outErr = errDeleteMessage
					ctx.Logger.WithError(errDeleteMessage).Error("error during deleting message")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
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

// deleteMessage gestisce DELETE /api/users/{userID}/conversations/groups/{conversationID}/messages/{messageID}
func (rt *_router) deleteMessageGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")

				// Cancella il messaggio
				errDeleteMessage := model.DeleteMessage(userId, conversationId, false, messageId)

				// Se ci sei riuscito
				if errDeleteMessage == nil {
					ctx.Logger.Info("message of group deleted successfully")

					// Costruisci la risposta JSON
					respData := struct {
						MessageID int `json:"messageID"`
					}{MessageID: messageId}

					jsonBytes, marshalErr := json.Marshal(respData)

					// Se il JSON viene creato correttamente
					if marshalErr == nil {

						// Tutto ok
						responseBody = jsonBytes
						outErr = nil

						statusCode = http.StatusOK // 200
						ctx.Logger.Infof(" message of group (%d) deleted successfully", messageId)

					} else {
						statusCode = http.StatusInternalServerError
						outErr = marshalErr
						ctx.Logger.WithError(marshalErr).Error("error marshalling response")
					}

				} else {
					statusCode = http.StatusInternalServerError
					outErr = errDeleteMessage
					ctx.Logger.WithError(errDeleteMessage).Error("error during deleting message of group")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
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
