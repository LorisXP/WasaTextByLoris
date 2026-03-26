package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// commentMessage gestisce:
//
//	POST /api/users/:userID/conversations/users/:conversationID/messages/:messageID/comments
//	POST /api/users/:userID/conversations/groups/:conversationID/messages/:messageID/comments
func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var statusCode int
	var responseBody []byte
	var outErr error
	ctx.Logger.Debug("default init ok")

	// Determina il tipo dalla URL
	betweenUsers := strings.Contains(r.URL.Path, "/conversations/users/")

	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		messageIdStr := ps.ByName("messageID")
		messageId, errMess := strconv.Atoi(messageIdStr)
		if errMess == nil || messageId > 0 {
			ctx.Logger.Info("messageId parsed successfully")

			var reqBody struct {
				Reaction string `json:"reaction"`
			}
			decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)
			if decodeErr == nil {
				ctx.Logger.Info("request decoded successfully")

				validationErr := model.ValidateInput(reqBody.Reaction, 1, 4, `^.*?$`, "string")
				if validationErr == nil {
					ctx.Logger.Info("request validated successfully")

					comment, errComment := model.CreateComment(userId, messageId, betweenUsers, reqBody.Reaction)
					if errComment == nil {
						ctx.Logger.Info("new reaction added successfully")

						respData := struct {
							CommentID int `json:"commentID"`
						}{CommentID: comment.CommentID}
						jsonBytes, marshalErr := json.Marshal(respData)
						if marshalErr == nil {
							responseBody = jsonBytes
							statusCode = http.StatusOK
							ctx.Logger.Infof("reaction '%s' added to messageID %d (betweenUsers: %t)", reqBody.Reaction, messageId, betweenUsers)
						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}
					} else {
						statusCode = http.StatusInternalServerError
						outErr = errComment
						ctx.Logger.WithError(errComment).Error("error during adding reaction")
					}
				} else {
					statusCode = http.StatusBadRequest
					outErr = validationErr
					ctx.Logger.WithError(validationErr).Error("invalid reaction")
				}
			} else {
				statusCode = http.StatusBadRequest
				outErr = decodeErr
				ctx.Logger.WithError(decodeErr).Error("invalid request body")
			}
		} else {
			statusCode = http.StatusBadRequest
			outErr = errMess
			ctx.Logger.WithError(errMess).Error("invalid messageId")
		}
	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
	}

	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_, _ = w.Write(responseBody)
	}
}

// deleteComment gestisce:
//
//	DELETE /api/users/:userID/conversations/users/:conversationID/messages/:messageID/comments/:commentID
//	DELETE /api/users/:userID/conversations/groups/:conversationID/messages/:messageID/comments/:commentID
func (rt *_router) deleteComment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var statusCode int
	var outErr error
	ctx.Logger.Debug("default init ok")

	// Determina il tipo dalla URL
	betweenUsers := strings.Contains(r.URL.Path, "/conversations/users/")

	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		commentIdStr := ps.ByName("commentID")
		commentId, errComm := strconv.Atoi(commentIdStr)
		if errComm == nil || commentId > 0 {
			ctx.Logger.Info("commentID parsed successfully")

			messageIdStr := ps.ByName("messageID")
			messageId, errMess := strconv.Atoi(messageIdStr)
			if errMess == nil || messageId > 0 {
				ctx.Logger.Info("messageId parsed successfully")

				errDelete := model.DeleteComment(userId, commentId, messageId, betweenUsers)
				if errDelete == nil {
					statusCode = http.StatusNoContent
					ctx.Logger.Infof("commentID %d deleted from messageID %d (betweenUsers: %t)", commentId, messageId, betweenUsers)
				} else {
					statusCode = http.StatusInternalServerError
					outErr = errDelete
					ctx.Logger.WithError(errDelete).Error("error during deleting comment")
				}
			} else {
				statusCode = http.StatusBadRequest
				outErr = errMess
				ctx.Logger.WithError(errMess).Error("invalid messageId")
			}
		} else {
			statusCode = http.StatusBadRequest
			outErr = errComm
			ctx.Logger.WithError(errComm).Error("invalid commentID")
		}
	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
	}

	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
	}
}
