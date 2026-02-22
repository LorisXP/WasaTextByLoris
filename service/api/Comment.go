package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// commentMessage gestisce POST /api/comments/users/:userID/messages/:messageID
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		messageIdStr := ps.ByName("messageID")
		messageId, errConv := strconv.Atoi(messageIdStr)

		//Se non ci sono errori nell'estrazione
		if errConv == nil || messageId > 0 {
			ctx.Logger.Info("messageId parsed successfully")
			
			// Decodifica il body JSON
			var reqBody struct {
				Reaction string `json:"reaction"`
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

			// Se il JSON è valido, prosegui.
			if decodeErr == nil {
				ctx.Logger.Info("request decoded successfully")

				//Valida l'input di content_type
				validationEmojiErr := model.ValidateInput(reqBody.Reaction, 1, 2, `^.*?$`, "string")
				
				//Se è corretto
				if validationEmojiErr == nil {
					ctx.Logger.Info("request validated successfully")

					comment, errComment := model.CreateComment(userId, messageId, true, reqBody.Reaction)
					
					//Se ci sei riuscito
					if errComment == nil {
						ctx.Logger.Info("new reaction added successfully")

						// Costruisci la risposta JSON
						respData := struct {
							CommentID int `json:"commentID"`
						}{CommentID: comment.CommentID}

						jsonBytes, marshalErr := json.Marshal(respData)
						
						//Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusOK // 200
							ctx.Logger.Debugf("new reaction added successfully %s (messageID: %d)", reqBody.Reaction, messageId)
							ctx.Logger.Info("new reaction added successfully")

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errComment
						ctx.Logger.WithError(errComment).Error("error during adding reaction to message")
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = validationEmojiErr
					ctx.Logger.WithError(validationEmojiErr).Error("invalid reaction")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = decodeErr
				ctx.Logger.WithError(decodeErr).Error("invalid request body")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errConv
			ctx.Logger.WithError(errConv).Error("invalid messageId")
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

// commentMessageGroup gestisce POST /api/comments/groups/:groupID/messages/:messageID
func (rt *_router) commentMessageGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		messageIdStr := ps.ByName("messageID")
		messageId, errConv := strconv.Atoi(messageIdStr)

		//Se non ci sono errori nell'estrazione
		if errConv == nil || messageId > 0 {
			ctx.Logger.Info("messageId parsed successfully")
			
			// Decodifica il body JSON
			var reqBody struct {
				Reaction string `json:"reaction"`
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

			// Se il JSON è valido, prosegui.
			if decodeErr == nil {
				ctx.Logger.Info("request decoded successfully")

				//Valida l'input di content_type
				validationEmojiErr := model.ValidateInput(reqBody.Reaction, 1, 2, `^.*?$`, "string")
				
				//Se è corretto
				if validationEmojiErr == nil {
					ctx.Logger.Info("request validated successfully")

					comment, errComment := model.CreateComment(userId, messageId, false, reqBody.Reaction)
					
					//Se ci sei riuscito
					if errComment == nil {
						ctx.Logger.Info("new reaction added successfully to message in group")

						// Costruisci la risposta JSON
						respData := struct {
							CommentID int `json:"commentID"`
						}{CommentID: comment.CommentID}

						jsonBytes, marshalErr := json.Marshal(respData)
						
						//Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusOK // 200
							ctx.Logger.Debugf("new reaction added successfully to message in group %s (messageID: %d)", reqBody.Reaction, messageId)
							ctx.Logger.Info("new reaction added successfully to message in group")

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errComment
						ctx.Logger.WithError(errComment).Error("error during adding reaction to message of group")
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = validationEmojiErr
					ctx.Logger.WithError(validationEmojiErr).Error("invalid reaction")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = decodeErr
				ctx.Logger.WithError(decodeErr).Error("invalid request body")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errConv
			ctx.Logger.WithError(errConv).Error("invalid messageId")
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

// deleteComment gestisce DELETE /api/comments/users/:userID/messages/:messageID/:commentID
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		commentIdStr := ps.ByName("commentID")
		commentId, errComm := strconv.Atoi(commentIdStr)

		//Se non ci sono errori nell'estrazione
		if errComm == nil || commentId > 0 {
			ctx.Logger.Info("commentID parsed successfully")

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")
				
				//Cancella il messaggio
				errDeleteComment := model.DeleteComment(userId, commentId, messageId, true)
				
				//Se ci sei riuscito
				if errDeleteComment == nil {
					ctx.Logger.Info("comment deleted successfully")

					// Tutto ok
					outErr = nil

					statusCode = http.StatusNoContent // 200
					ctx.Logger.Infof(" message (%d) deleted successfully", messageId)


				} else {
					statusCode = http.StatusInternalServerError
					outErr = errDeleteComment
					ctx.Logger.WithError(errDeleteComment).Error("error during deleting comment")
				}
				
			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errComm
			ctx.Logger.WithError(errComm).Error("invalid commentID")
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

// deleteCommentGroup gestisce DELETE /api/comments/groups/:groupID/messages/:messageID/:commentID
func (rt *_router) deleteCommentGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		commentIdStr := ps.ByName("commentID")
		commentId, errComm := strconv.Atoi(commentIdStr)

		//Se non ci sono errori nell'estrazione
		if errComm == nil || commentId > 0 {
			ctx.Logger.Info("commentID parsed successfully")

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)

			// Se il JSON è valido, prosegui.
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("params validated successfully")
				
				//Cancella il commento
				errDeleteComment := model.DeleteComment(userId, commentId, messageId, false)
				
				//Se ci sei riuscito
				if errDeleteComment == nil {
					ctx.Logger.Info("comment deleted successfully in group")

					// Tutto ok
					outErr = nil

					statusCode = http.StatusNoContent // 200
					ctx.Logger.Infof(" message (%d) deleted successfully", messageId)


				} else {
					statusCode = http.StatusInternalServerError
					outErr = errDeleteComment
					ctx.Logger.WithError(errDeleteComment).Error("error during deleting comment in group")
				}
				
			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errComm
			ctx.Logger.WithError(errComm).Error("invalid commentID")
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

