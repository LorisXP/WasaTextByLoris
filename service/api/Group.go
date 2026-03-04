package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/LorisXP/WasaTextByLoris/service/api/model"
	"github.com/LorisXP/WasaTextByLoris/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// createGroup gestisce POST /api/users/{userID}/groups
func (rt *_router) createGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

		// Decodifica il body JSON { "userName": "..." }
		var reqBody struct {
			Name  string `json:"name"`
			Photo string `json:"photo"`
		}
		decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

		// Se il JSON è valido, prosegui.
		if decodeErr == nil {
			ctx.Logger.Info("request decoded successfully")

			// Valida l'input
			validationNameErr := model.ValidateInput(reqBody.Name, 3, 15, `^.*?$`, "string")
			validationPhotoErr := model.ValidateInput(reqBody.Photo, 0, 13981013, `^[A-Za-z0-9+/]+={0,2}$`, "string")

			//Se il campo name inserito è valido
			if validationNameErr == nil {
				ctx.Logger.Info("'name' in request is valid")

				//Se il campo photo inserito è valido
				if validationPhotoErr == nil {
					ctx.Logger.Info("'photo' in request is valid")

					//Crea il grupppo
					group, errGrp := model.CreateGroup(reqBody.Name, reqBody.Photo, userId)
					ctx.Logger.Debug("passed by CreateGroup()")

					//Se il gruppo è creato correttamente
					if errGrp == nil {
						ctx.Logger.Info("created group")

						// Costruisci la risposta JSON
						respData := struct {
							GroupID int `json:"groupID"`
						}{GroupID: group.GroupID}

						jsonBytes, marshalErr := json.Marshal(respData)

						//Se il JSON viene creato correttamente
						if marshalErr == nil {

							// Tutto ok
							responseBody = jsonBytes
							outErr = nil

							statusCode = http.StatusCreated // 201
							ctx.Logger.Infof("group created successfully: %s (ID: %d) by userID %d", reqBody.Name, group.GroupID, userId)

						} else {
							statusCode = http.StatusInternalServerError
							outErr = marshalErr
							ctx.Logger.WithError(marshalErr).Error("error marshalling response")
						}

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errGrp
						ctx.Logger.WithError(errGrp).Errorf("impossible for userID %d to create group", userId)
					}
				} else {
					statusCode = http.StatusBadRequest
					outErr = validationPhotoErr
					ctx.Logger.WithError(validationPhotoErr).Error("invalid photo. Check limits on doc")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = validationNameErr
				ctx.Logger.WithError(validationNameErr).Error("invalid name. Check limits on doc")
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

// getGroupInfo gestisce GET /api/users/{userID}/groups/{groupID}
func (rt *_router) getGroupInfo(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

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

		groupIdStr := ps.ByName("groupID")
		groupId, errGrp := strconv.Atoi(groupIdStr)

		//Se non ci sono errori nell'estrazione
		if errGrp == nil || groupId > 0 {
			ctx.Logger.Info("groupId parsed successfully")

			//Ottieni il gruppo (verifica che lo userID sia membro o admin)
			group, errGetGrp := model.GetGroup(groupId, userId)
			ctx.Logger.Debug("passed by GetGroup()")

			//Se il gruppo è ottenuto correttamente
			if errGetGrp == nil {
				ctx.Logger.Info("group obtained successfully")

				// Costruisci la risposta JSON
				jsonBytes, marshalErr := json.Marshal(group)

				//Se il JSON viene creato correttamente
				if marshalErr == nil {

					// Tutto ok
					responseBody = jsonBytes
					outErr = nil

					statusCode = http.StatusOK // 200
					ctx.Logger.Infof("group info retrieved successfully: groupID %d by userID %d", groupId, userId)

				} else {
					statusCode = http.StatusInternalServerError
					outErr = marshalErr
					ctx.Logger.WithError(marshalErr).Error("error marshalling response")
				}

			} else if errors.Is(errGetGrp, model.ErrUserNotInGroup) {
				statusCode = http.StatusForbidden // 403
				outErr = errGetGrp
				ctx.Logger.WithError(errGetGrp).Warnf("userID %d is not a member of groupID %d", userId, groupId)
			} else {
				statusCode = http.StatusNotFound
				outErr = errGetGrp
				ctx.Logger.WithError(errGetGrp).Error("group not found")
			}

		} else {
			statusCode = http.StatusUnauthorized
			outErr = errGrp
			ctx.Logger.WithError(errGrp).Error("invalid groupId")
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

// leaveGroup gestisce DELETE /api/users/{userID}/groups/{groupID}
func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		groupIdStr := ps.ByName("groupID")
		groupId, errGrp := strconv.Atoi(groupIdStr)

		//Se non ci sono errori nell'estrazione
		if errGrp == nil || groupId > 0 {
			ctx.Logger.Info("groupId parsed successfully")

			//Ottieni il gruppo
			group, errGetGrp := model.GetGroupBasic(groupId)
			ctx.Logger.Debug("passed by GetGroupBasic()")

			//Se il gruppo è ottenuto correttamente
			if errGetGrp == nil {
				ctx.Logger.Info("group obtained successfully")

				//Verifica che lo userID sia membro o admin del gruppo
				_, errMember := model.GetGroup(groupId, userId)
				ctx.Logger.Debug("passed by GetGroup() for membership check")

				if errMember == nil {

					//Fai uscire l'utente
					errLeaveGroup := model.LeaveGroup(userId, group)
					ctx.Logger.Debug("passed by LeaveGroup()")

					//Se l'operazione va a buon fine
					if errLeaveGroup == nil {

						//Tutto ok
						outErr = nil

						statusCode = http.StatusNoContent
						ctx.Logger.Infof("userID %d successfully left groupID %d", userId, groupId)

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errLeaveGroup
						ctx.Logger.WithError(errLeaveGroup).Errorf("impossible for userID %d to leave groupID %d", userId, groupId)
					}

				} else if errors.Is(errMember, model.ErrUserNotInGroup) {
					statusCode = http.StatusForbidden
					outErr = errMember
					ctx.Logger.WithError(errMember).Warnf("userID %d is not a member of groupID %d", userId, groupId)
				} else {
					statusCode = http.StatusInternalServerError
					outErr = errMember
					ctx.Logger.WithError(errMember).Error("error checking membership")
				}

			} else {
				statusCode = http.StatusNotFound
				outErr = errGetGrp
				ctx.Logger.WithError(errGetGrp).Error("group not found")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = errGrp
			ctx.Logger.WithError(errGrp).Error("invalid groupId")
		}

	} else {
		statusCode = http.StatusUnauthorized
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
	}

	// Unico punto di uscita (204 No Content: nessun body)
	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
	}
}

// addToGroup gestisce POST /api/groups/{groupID}/users
func (rt *_router) addToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	groupIdStr := ps.ByName("groupID")
	groupId, errGrp := strconv.Atoi(groupIdStr)

	//Se non ci sono errori nell'estrazione
	if errGrp == nil || groupId > 0 {
		ctx.Logger.Info("groupId parsed successfully")

		// Decodifica il body JSON { "userName": "..." }
		var reqBody struct {
			UserId    int      `json:"userID"`
			UserNames []string `json:"userNames"`
		}
		decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

		//Se il decode va a buon fie
		if decodeErr == nil {
			ctx.Logger.Info("request decoded successfully")

			//Preleva lo userId
			userId := reqBody.UserId

			// Valida l'input
			if userId > 0 {
				ctx.Logger.Info("userId parsed successfully")

				// Verifica che il Bearer token corrisponda al userID nel body
				if userId != ctx.BearerUserID {
					statusCode = http.StatusForbidden
					outErr = fmt.Errorf("Bearer token does not match userID in the request body")
					ctx.Logger.WithError(outErr).Warn("authorization mismatch")
				} else {

					// Valida tutti gli userNames
					var validationErr error
					for _, uName := range reqBody.UserNames {
						validationErr = model.ValidateInput(uName, 3, 15, `^[a-z]+[0-9]*$`, "string")
						if validationErr != nil {
							ctx.Logger.WithError(validationErr).Errorf("invalid userName: %s", uName)
							break
						}
					}

					if validationErr == nil {
						ctx.Logger.Info("all userNames are valid")

						// Ottieni il gruppo
						group, errGetGrp := model.GetGroupBasic(groupId)
						ctx.Logger.Debug("passed by GetGroupBasic()")

						// Se il gruppo è ottenuto correttamente
						if errGetGrp == nil {
							ctx.Logger.Info("group obtained successfully")

							// Verifica che lo userID sia admin del gruppo
							if group.AdminID == userId {
								ctx.Logger.Infof("userID is admin of groupID %d", groupId)

								errAddToGroup := model.AddToGroup(group, reqBody.UserNames)
								ctx.Logger.Debug("passed by AddToGroup()")

								// Se l'operazione va a buon fine
								if errAddToGroup == nil {

									// Tutto ok
									outErr = nil
									statusCode = http.StatusNoContent
									ctx.Logger.Infof("users successfully added to groupID %d", groupId)

								} else {
									statusCode = http.StatusInternalServerError
									outErr = errAddToGroup
									ctx.Logger.WithError(errAddToGroup).Errorf("impossible to add users in groupID %d", groupId)
								}

							} else {
								statusCode = http.StatusUnauthorized
								outErr = fmt.Errorf("userID %d is not admin of groupID %d", userId, groupId)
								ctx.Logger.WithError(outErr).Warn("user is not admin")
							}

						} else {
							statusCode = http.StatusNotFound
							outErr = errGetGrp
							ctx.Logger.WithError(errGetGrp).Error("group not found")
						}

					} else {
						statusCode = http.StatusBadRequest
						outErr = validationErr
						ctx.Logger.WithError(validationErr).Error("invalid userName in list")
					}
				}
			} else {
				statusCode = http.StatusBadRequest
				outErr = fmt.Errorf("invalid userId. must be > 0")
				ctx.Logger.WithError(outErr).Error("invalid userId")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = decodeErr
			ctx.Logger.WithError(decodeErr).Error("invalid request body")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = errGrp
		ctx.Logger.WithError(errGrp).Error("invalid groupId")
	}

	// Unico punto di uscita (204 No Content: nessun body)
	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
	}
}

// setGroupName gestisce PATCH /api/groups/{groupID}/name
func (rt *_router) setGroupName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai groupId dai parameters
	groupIdStr := ps.ByName("groupID")
	groupId, errGrp := strconv.Atoi(groupIdStr)

	//Se non ci sono errori nell'estrazione
	if errGrp == nil || groupId > 0 {
		ctx.Logger.Info("groupId parsed successfully")

		// Decodifica il body JSON
		var reqBody struct {
			UserId int    `json:"userID"`
			Name   string `json:"name"`
		}
		decodeErr := json.NewDecoder(r.Body).Decode(&reqBody)

		// Se il JSON è valido, prosegui.
		if decodeErr == nil {
			ctx.Logger.Info("request decoded successfully")

			// Verifica che il Bearer token corrisponda al userID nel body
			if reqBody.UserId != ctx.BearerUserID {
				statusCode = http.StatusForbidden
				outErr = fmt.Errorf("Bearer token does not match userID in the request body")
				ctx.Logger.WithError(outErr).Warn("authorization mismatch")
			} else {

				// Valida l'input del nuovo groupName con il Validator
				validationErr := model.ValidateInput(reqBody.Name, 3, 15, `^.*?$`, "string")

				//Se il campo inserito è valido
				if validationErr == nil {
					ctx.Logger.Info("request data is valid")

					// Ottieni il gruppo
					group, errGetGrp := model.GetGroupBasic(groupId)
					ctx.Logger.Debug("passed by GetGroupBasic()")

					// Se il gruppo è ottenuto correttamente
					if errGetGrp == nil {
						ctx.Logger.Info("group obtained successfully")

						// Verifica che lo userID sia admin del gruppo
						if group.AdminID == reqBody.UserId {
							ctx.Logger.Infof("userID is admin of groupID %d", groupId)

							//Prova ad impostare il nome
							errSetNameGroup := model.SetNameGroup(&group, reqBody.Name)
							ctx.Logger.Debug("passed by SetNameGroup()")

							//Se va a buon fine
							if errSetNameGroup == nil {
								ctx.Logger.Info("groupName updated successfully")

								// Costruisci la risposta JSON
								respData := struct{}{}
								jsonBytes, marshalErr := json.Marshal(respData)

								//Se il JSON viene creato correttamente
								if marshalErr == nil {

									// Tutto ok
									responseBody = jsonBytes
									outErr = nil

									statusCode = http.StatusOK // 200
									ctx.Logger.Debugf("groupName update successfully: %s (ID: %d)", reqBody.Name, groupId)
									ctx.Logger.Infof("groupName update successfully: %s", reqBody.Name)

								} else {
									statusCode = http.StatusInternalServerError
									outErr = marshalErr
									ctx.Logger.WithError(marshalErr).Error("error marshalling response")
								}

							} else {
								statusCode = http.StatusInternalServerError
								outErr = errSetNameGroup
								ctx.Logger.WithError(outErr).Errorf("error during set a new name '%s' to groupID %d", reqBody.Name, groupId)
							}

						} else {
							statusCode = http.StatusUnauthorized
							outErr = fmt.Errorf("userID %d is not admin of groupID %d", reqBody.UserId, groupId)
							ctx.Logger.WithError(outErr).Warn("user is not admin")
						}

					} else {
						statusCode = http.StatusNotFound
						outErr = errGetGrp
						ctx.Logger.WithError(errGetGrp).Error("group not found")
					}

				} else {
					statusCode = http.StatusBadRequest
					outErr = validationErr
					ctx.Logger.WithError(validationErr).Error("invalid groupName. Check limits on doc")
				}

			} // chiude else del controllo Bearer

		} else {
			statusCode = http.StatusBadRequest
			outErr = decodeErr
			ctx.Logger.WithError(decodeErr).Error("invalid request body")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = errGrp
		ctx.Logger.WithError(errGrp).Error("invalid groupId")
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

// setGroupPhoto gestisce PUT /api/groups/{groupID}/photo
func (rt *_router) setGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var responseBody []byte = nil
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai groupId dai parameters
	groupIdStr := ps.ByName("groupID")
	groupId, errGrp := strconv.Atoi(groupIdStr)
	//Se non ci sono errori nell'estrazione
	if errGrp == nil || groupId > 0 {
		ctx.Logger.Info("groupId parsed successfully")

		// Leggi il form multipart
		// Limite: ~14MB (dimensione massima della foto codificata in base64 dalla YAML)
		parseErr := r.ParseMultipartForm(14 << 20)

		// Se il form è valido, prosegui.
		if parseErr == nil {
			ctx.Logger.Info("multipart form parsed successfully")

			// Estrai lo userID dal form
			userIdStr := r.FormValue("userID")
			userId, errUsr := strconv.Atoi(userIdStr)

			if errUsr == nil && userId > 0 {
				ctx.Logger.Info("userId parsed successfully from form")

				// Verifica che il Bearer token corrisponda al userID nel form
				if userId != ctx.BearerUserID {
					statusCode = http.StatusForbidden
					outErr = fmt.Errorf("Bearer token does not match userID in the form")
					ctx.Logger.WithError(outErr).Warn("authorization mismatch")
				} else {

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

							// Ottieni il gruppo
							group, errGetGrp := model.GetGroupBasic(groupId)
							ctx.Logger.Debug("passed by GetGroupBasic()")

							// Se il gruppo è ottenuto correttamente
							if errGetGrp == nil {
								ctx.Logger.Info("group obtained successfully")

								// Verifica che lo userID sia admin del gruppo
								if group.AdminID == userId {
									ctx.Logger.Infof("userID %d is admin of groupID %d", userId, groupId)

									// Aggiorna la foto del gruppo
									errUpdatePhoto := model.SetGroupPhoto(&group, photoBase64)
									ctx.Logger.Debug("passed by SetGroupPhoto()")

									// Se la foto viene aggiornata correttamente
									if errUpdatePhoto == nil {
										ctx.Logger.Info("group photo updated successfully")

										// Costruisci la risposta JSON
										respData := struct{}{}
										jsonBytes, marshalErr := json.Marshal(respData)

										// Se il JSON viene creato correttamente
										if marshalErr == nil {

											// Tutto ok
											responseBody = jsonBytes
											outErr = nil

											statusCode = http.StatusOK // 200
											ctx.Logger.Debugf("group photo updated successfully (groupID: %d) by userID %d", groupId, userId)
											ctx.Logger.Infof("group photo updated successfully")

										} else {
											statusCode = http.StatusInternalServerError
											outErr = marshalErr
											ctx.Logger.WithError(marshalErr).Error("error marshalling response")
										}

									} else {
										statusCode = http.StatusInternalServerError
										outErr = errUpdatePhoto
										ctx.Logger.WithError(errUpdatePhoto).Error("error updating group photo")
									}

								} else {
									statusCode = http.StatusUnauthorized
									outErr = fmt.Errorf("userID %d is not admin of groupID %d", userId, groupId)
									ctx.Logger.WithError(outErr).Warn("user is not admin")
								}

							} else {
								statusCode = http.StatusNotFound
								outErr = errGetGrp
								ctx.Logger.WithError(errGetGrp).Error("group not found")
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

				} // chiude else del controllo Bearer

			} else {
				statusCode = http.StatusBadRequest
				outErr = fmt.Errorf("invalid or missing userID in form")
				ctx.Logger.WithError(outErr).Error("invalid userId in form")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = parseErr
			ctx.Logger.WithError(parseErr).Error("invalid multipart form")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = errGrp
		ctx.Logger.WithError(errGrp).Error("invalid groupId")
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

// kickUserFromGroup gestisce DELETE /api/users/{userID}/groups/{groupID}/member/{userName}
func (rt *_router) kickUserFromGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai userId dai parameters
	userIdStr := ps.ByName("userID")
	userId, errUsr := strconv.Atoi(userIdStr)
	//Se non ci sono errori nell'estrazione
	if errUsr == nil || userId > 0 {
		ctx.Logger.Info("userId parsed successfully")

		groupIdStr := ps.ByName("groupID")
		groupId, errGrp := strconv.Atoi(groupIdStr)

		//Se non ci sono errori nell'estrazione
		if errGrp == nil || groupId > 0 {
			ctx.Logger.Info("groupId parsed successfully")

			//Verifica se lo userName è valido
			userName := ps.ByName("userName")
			validationErr := model.ValidateInput(userName, 3, 15, `^[a-z]+[0-9]*$`, "string")

			//Se la validazione dello userName è andata a buon fine
			if validationErr == nil {
				ctx.Logger.Info("userName parsed successfully")

				//Ottieni il gruppo
				group, errGetGrp := model.GetGroupBasic(groupId)
				ctx.Logger.Debug("passed by GetGroupBasic()")

				//Se il gruppo è ottenuto correttamente
				if errGetGrp == nil {
					ctx.Logger.Info("group obtained successfully")

					// Verifica che lo userID sia admin del gruppo
					if group.AdminID == userId {
						ctx.Logger.Infof("userID %d is admin of groupID %d", userId, groupId)

						// Ottieni il nome dell'admin per verificare che non stia kickando se stesso
						admin, errAdmin := model.GetUser(userId)
						ctx.Logger.Debug("passed by GetUser() for admin name")

						if errAdmin == nil {

							// Verifica che lo userName da kickare non sia l'admin stesso
							if admin.Name != userName {
								ctx.Logger.Infof("userName '%s' is not the admin, proceeding with kick", userName)

								//Fai uscire l'utente
								errKickGroup := model.KickFromGroup(&group, userName)
								ctx.Logger.Debug("passed by KickFromGroup()")

								//Se l'operazione va a buon fine
								if errKickGroup == nil {

									//Tutto ok
									outErr = nil

									statusCode = http.StatusNoContent
									ctx.Logger.Infof("user %s successfully kicked from groupID %d", userName, groupId)

								} else {
									statusCode = http.StatusInternalServerError
									outErr = errKickGroup
									ctx.Logger.WithError(errKickGroup).Errorf("impossible to kick '%s' from groupID %d", userName, groupId)
								}

							} else {
								statusCode = http.StatusBadRequest
								outErr = fmt.Errorf("admin cannot kick themselves from groupID %d", groupId)
								ctx.Logger.WithError(outErr).Warn("admin tried to kick themselves")
							}

						} else {
							statusCode = http.StatusInternalServerError
							outErr = errAdmin
							ctx.Logger.WithError(errAdmin).Error("error retrieving admin user info")
						}

					} else {
						statusCode = http.StatusUnauthorized
						outErr = fmt.Errorf("userID %d is not admin of groupID %d", userId, groupId)
						ctx.Logger.WithError(outErr).Warn("user is not admin")
					}

				} else {
					statusCode = http.StatusNotFound
					outErr = errGetGrp
					ctx.Logger.WithError(errGetGrp).Error("group not found")
				}

			} else {
				statusCode = http.StatusBadRequest
				outErr = validationErr
				ctx.Logger.WithError(validationErr).Error("invalid userName")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = errGrp
			ctx.Logger.WithError(errGrp).Error("invalid groupId")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = errUsr
		ctx.Logger.WithError(errUsr).Error("invalid userId")
	}

	// Unico punto di uscita (204 No Content: nessun body)
	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
	}
}

// removeGroup gestisce DELETE /api/groups/{groupID}/users/{userID}
func (rt *_router) removeGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta i default
	var statusCode int = http.StatusInternalServerError
	var outErr error = nil
	ctx.Logger.Debug("default init ok")

	// Estrai groupId dai parameters
	groupIdStr := ps.ByName("groupID")
	groupId, errGrp := strconv.Atoi(groupIdStr)

	//Se non ci sono errori nell'estrazione
	if errGrp == nil || groupId > 0 {
		ctx.Logger.Info("groupId parsed successfully")

		// Estrai userId dai parameters
		userIdStr := ps.ByName("userID")
		userId, errUsr := strconv.Atoi(userIdStr)

		//Se non ci sono errori nell'estrazione
		if errUsr == nil || userId > 0 {
			ctx.Logger.Info("userId parsed successfully")

			// Ottieni il gruppo
			group, errGetGrp := model.GetGroupBasic(groupId)
			ctx.Logger.Debug("passed by GetGroupBasic()")

			// Se il gruppo è ottenuto correttamente
			if errGetGrp == nil {
				ctx.Logger.Info("group obtained successfully")

				// Verifica che lo userID sia admin del gruppo
				if group.AdminID == userId {
					ctx.Logger.Infof("userID %d is admin of groupID %d", userId, groupId)

					// Elimina il gruppo
					errDelete := model.DeleteGroup(&group)
					ctx.Logger.Debug("passed by DeleteGroup()")

					// Se l'operazione va a buon fine
					if errDelete == nil {

						// Tutto ok
						outErr = nil
						statusCode = http.StatusNoContent // 204
						ctx.Logger.Infof("groupID %d deleted successfully by adminID %d", groupId, userId)

					} else {
						statusCode = http.StatusInternalServerError
						outErr = errDelete
						ctx.Logger.WithError(errDelete).Errorf("impossible to delete groupID %d", groupId)
					}

				} else {
					statusCode = http.StatusUnauthorized
					outErr = fmt.Errorf("userID %d is not admin of groupID %d", userId, groupId)
					ctx.Logger.WithError(outErr).Warn("user is not admin")
				}

			} else {
				statusCode = http.StatusNotFound
				outErr = errGetGrp
				ctx.Logger.WithError(errGetGrp).Error("group not found")
			}

		} else {
			statusCode = http.StatusBadRequest
			outErr = errUsr
			ctx.Logger.WithError(errUsr).Error("invalid userId")
		}

	} else {
		statusCode = http.StatusBadRequest
		outErr = errGrp
		ctx.Logger.WithError(errGrp).Error("invalid groupId")
	}

	// Unico punto di uscita (204 No Content: nessun body)
	if outErr != nil {
		http.Error(w, outErr.Error(), statusCode)
	} else {
		w.WriteHeader(statusCode)
	}
}
