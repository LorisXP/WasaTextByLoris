package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

// ErrUserNotInGroup viene restituito quando l'utente non è membro né admin del gruppo
var ErrUserNotInGroup = errors.New("user is not a member of the group")

// CreateGroup creates a group in WasaText with name, photo and the ID of the admin
func CreateGroup(name string, photo string, admin int) (entity.Group, error) {
	logrus.Debug("Entered in CreateGroup()")
	logrus.Infof("Creating a group named %s, by admin %d ", name, admin)

	//Imposta i default
	outErr := fmt.Errorf("Unable to create a new group named %s, by admin %d ", name, admin)
	var group entity.Group

	//Crea l'inserimento a DB
	name = strings.TrimSpace(name)
	groupID, err := dml.CreateGroup(name, photo, admin)
	logrus.Debug("Passed by CreateGroup()")

	//Se non ci sono errori, prosegui
	if err == nil && groupID != 0 {

		//Crea l'oggetto gruppo
		group.GroupID = groupID
		group.Name = name
		group.Photo = photo
		group.AdminID = admin

		outErr = nil

		logrus.Info("Group created successfully")

	} else {
		outErr = fmt.Errorf("error during creating group: %w", err)
		logrus.Error(outErr)
	}

	return group, outErr
}

// GetGroupBasic returns the raw entity.Group by groupID (without membership check or member list).
func GetGroupBasic(groupID int) (entity.Group, error) {
	logrus.Debug("Entered in GetGroupBasic()")
	logrus.Infof("Getting basic group info by id %d", groupID)

	group, err := queries.GetGroupByID(groupID)
	if err != nil {
		return entity.Group{}, fmt.Errorf("error getting group by id %d: %w", groupID, err)
	}

	return group, nil
}

// GetGroup returns group info by groupID, verificando che userID sia membro o admin.
func GetGroup(groupID int, userID int) (entity.GroupInfoResponse, error) {
	logrus.Debug("Entered in GetGroup()")
	logrus.Infof("Getting group by id %d, requested by userID %d", groupID, userID)

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve group by id %d ", groupID)
	var result entity.GroupInfoResponse

	//Verifica se l'utente fa parte del gruppo
	isMember, errCheck := queries.IsUserInGroup(userID, groupID)

	if errCheck != nil {
		outErr = fmt.Errorf("error checking membership for userID %d in groupID %d: %w", userID, groupID, errCheck)
		logrus.Error(outErr)
	} else if !isMember {
		outErr = fmt.Errorf("%w: userID %d in groupID %d", ErrUserNotInGroup, userID, groupID)
		logrus.WithField("userID", userID).Warnf("userID %d is not a member of groupID %d", userID, groupID)
	} else {
		//Ottieni le info del gruppo a DB
		group, err := queries.GetGroupByID(groupID)
		logrus.Debug("Passed by GetGroupByID()")

		//Se non ci sono errori, prosegui
		if err == nil && group.GroupID != 0 {

			//Ottieni il nome dell'admin
			admin, errAdmin := queries.GetUserByID(group.AdminID)
			logrus.Debug("Passed by GetUserByID() for admin")

			if errAdmin == nil {

				//Ottieni la lista dei membri
				members, errMembers := queries.GetGroupMembers(groupID)
				logrus.Debug("Passed by GetGroupMembers()")

				if errMembers == nil {

					//Costruisci la risposta
					result.GroupID = group.GroupID
					result.Name = group.Name
					result.Photo = group.Photo
					result.Administrator = admin.Name

					result.Members = make([]entity.MemberResponse, 0, len(members))
					for _, m := range members {
						result.Members = append(result.Members, entity.MemberResponse{
							UserName: m.Name,
							Photo:    m.Photo,
						})
					}

					outErr = nil
					logrus.Info("Group obtained successfully")

				} else {
					outErr = fmt.Errorf("error getting members for groupID %d: %w", groupID, errMembers)
					logrus.Error(outErr)
				}

			} else {
				outErr = fmt.Errorf("error getting admin name for adminID %d: %w", group.AdminID, errAdmin)
				logrus.Error(outErr)
			}

		} else {
			outErr = fmt.Errorf("error during obtaining group by id %d: %w", groupID, err)
			logrus.Error(outErr)
		}
	}

	return result, outErr
}

// LeaveGroup allows a user with their userID to leave a group by groupID
func LeaveGroup(userID int, group entity.Group) error {
	logrus.Debug("Entered in LeaveGroup()")
	logrus.Infof("userID %d leaving from groupID %d", userID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for userID %d to leave from group by id %d ", userID, group.GroupID)

	//Effettua l'operazione al DB
	err := dml.LeaveGroup(userID, group.GroupID)
	logrus.Debug("Passed by LeaveGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		outErr = nil
		logrus.Infof("userID %d left successfully from group by id %d", userID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error leaving groupID %d by userID %d: %w", group.GroupID, userID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// AddToGroup allows an admin of a group to add a list of users by their names
func AddToGroup(group entity.Group, userName []string) error {
	logrus.Debug("Entered in AddToGroup()")
	logrus.Infof("adminID %d adding users %v, in groupID %d", group.AdminID, userName, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d adding users %v, in groupID %d", group.AdminID, userName, group.AdminID)

	//Bisogna prima ottenere gli userID dai nomi e poi inserirli nel gruppo
	userID_list, err := queries.GetUsersIDByName(userName)
	logrus.Debug("Passed by GetUsersIDByName()")

	//Se non ci sono errori, prosegui
	if err == nil {
		logrus.Infof("Obtained correct userID by this userNames: %v ", userName)

		//Se le due liste sono uguali prosegui
		if len(userID_list) == len(userName) {
			logrus.Info("All userID found correctly")
			logrus.Debugf("The userID list is: %v", userID_list)

			//Effettua l'operazione al DB
			err = dml.AddToGroup(group.AdminID, userID_list, group.GroupID)
			logrus.Debug("Passed by AddToGroup()")

			//Se non ci sono errori, prosegui
			if err == nil {
				logrus.Infof("adminID %d added users %v to groupID %d successfully", group.AdminID, userName, group.GroupID)
				outErr = nil
			} else {
				outErr = fmt.Errorf("error during add userID list in the groupID %d: %w", group.AdminID, err)
				logrus.Error(outErr)
			}

		} else {
			outErr = fmt.Errorf("number of usernames differs from the number of userIDs: some usernames probably do not exist")
			logrus.Error(outErr)
		}

	} else {
		outErr = fmt.Errorf("error during getting userID by userName: %w", err)
		logrus.Error(outErr)
	}

	return outErr

}

// SetNameGroup allows an admin of a group to set the group name
func SetNameGroup(group *entity.Group, new_name string) error {
	logrus.Debug("Entered in SetNameGroup()")
	logrus.Infof("adminID %d setting the name of groupID %d", group.AdminID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d setting the name of groupID %d", group.AdminID, group.GroupID)

	//Effettua l'operazione al DB
	new_name = strings.TrimSpace(new_name)
	err := dml.UpdateNameGroup(group.GroupID, group.AdminID, new_name)
	logrus.Debug("Passed by UpdateNameGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		group.Name = new_name
		outErr = nil
		logrus.Infof("adminID %d set the name '%s' of groupID %d successfully", group.AdminID, new_name, group.GroupID)
	} else {
		outErr = fmt.Errorf("error setting new group name '%s' in groupID %d: %w", new_name, group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// SetGroupPhoto allows an admin of a group to set the group picture
func SetGroupPhoto(group *entity.Group, new_photo string) error {
	logrus.Debug("Entered in SetGroupPhoto()")
	logrus.Infof("adminID %d setting the photo of groupID %d", group.AdminID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d setting the photo of groupID %d", group.AdminID, group.GroupID)

	//Effettua l'operazione al DB
	err := dml.UpdatePhotoGroup(group.GroupID, group.AdminID, new_photo)
	logrus.Debug("Passed by UpdatePhotoGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		group.Photo = new_photo
		outErr = nil
		logrus.Infof("adminID %d set the photo of groupID %d successfully", group.AdminID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error setting new group photo in groupID %d: %w", group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// KickFromGroup allows an admin of a group to kick a user by their username
func KickFromGroup(group *entity.Group, userName string) error {
	//Bisogna prima ottenere lo userID corrispondente allo username e poi fare il kick

	logrus.Debug("Entered in KickFromGroup()")
	logrus.Infof("adminID %d is kicking out the user '%s' from groupID %d", group.AdminID, userName, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d to kick out the user '%s' from groupID %d", group.AdminID, userName, group.GroupID)

	//Ottieni lo userID
	user_id, err := queries.GetUserIDByName(userName)

	//Se non ci sono errori, prosegui
	if err == nil && user_id != 0 {
		logrus.Infof("found userID by his userName '%s'", userName)
		logrus.Debugf("The userID of userName '%s' is %d", userName, user_id)

		//Quindi effettua il kick
		err := dml.KickFromGroup(group.GroupID, group.AdminID, user_id)
		logrus.Debug("Passed by KickFromGroup()")

		//Se il remove ha funzionato
		if err == nil {
			outErr = nil
			logrus.Infof("adminID %d kicked userName '%s' from groupID %d successfully", group.AdminID, userName, group.GroupID)
		} else {
			outErr = fmt.Errorf("unable to kick userName '%s' from groupID %d: %w", userName, group.GroupID, err)
			logrus.Error(outErr)
		}
	} else {
		outErr = fmt.Errorf("error searching userID by '%s' to kick them from groupID %d: %w", userName, group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

// DeleteGroup allows an admin of a group to delete the entire group
func DeleteGroup(group *entity.Group) error {
	//La cancellazione del gruppo comporta la cancellazione di membri, eventi, messaggi, contenuti ecc.

	logrus.Debug("Entered in DeleteGroup()")
	logrus.Warningf("adminID %d is deleting groupID %d", group.AdminID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("unable for adminID %d to delete groupID %d", group.AdminID, group.GroupID)

	//Effettua l'operazione al DB
	err := dml.DeleteGroup(group.GroupID, group.AdminID)
	logrus.Debug("Passed by DeleteGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		//Azzeralo
		group = &entity.Group{}
		outErr = nil
		logrus.Infof("adminID %d deleted groupID %d successfully", group.AdminID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error deleting group %d: %w", group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}
