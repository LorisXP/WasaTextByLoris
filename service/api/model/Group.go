package model

import (
	"fmt"

	"github.com/LorisXP/WasaTextByLoris/service/api/entity"
	"github.com/LorisXP/WasaTextByLoris/service/database/dml"
	"github.com/LorisXP/WasaTextByLoris/service/database/queries"
	"github.com/sirupsen/logrus"
)

/*
Create a group in WasaText with name, photo and the ID of the admin
*/
func CreateGroup(name string, photo string, admin int) (entity.Group, error) {
	logrus.Debug("Entered in CreateGroup()")
	logrus.Infof("Creating a group named %s, by admin %d ", name, admin)

	//Imposta i default
	outErr := fmt.Errorf("Unable to create a new group named %s, by admin %d ", name, admin)
	var group entity.Group

	//Crea l'inserimento a DB
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

		logrus.Info("Group created succesfully")

	} else {
		outErr = fmt.Errorf("error during creating group: %w", err)
		logrus.Error(outErr)
	}

	return group, outErr
}

/*
Return a group struct in WasaText by his groupID
*/
func GetGroup(groupID int) (entity.Group, error) {
	logrus.Debug("Entered in GetGroup()")
	logrus.Infof("Getting group by id %d", groupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable to retrieve group by id %d ", groupID)

	//Ottieni le info a DB
	group, err := queries.GetGroupByID(groupID)
	logrus.Debug("Passed by GetGroupByID()")

	//Se non ci sono errori, prosegui
	if err == nil && group.GroupID != 0 {
		outErr = nil
		logrus.Info("Group obtained succesfully")
	} else {
		outErr = fmt.Errorf("error during obtaining group by id %d: %w", groupID, err)
		logrus.Error(outErr)
	}

	return group, outErr
}

/*
Allows a user with him userID to leave from group by groupID
*/
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
		logrus.Infof("userID %d leave succesfully from group by id %d ", userID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error leaving groupID %d by userID %d: %w", group.GroupID, userID, err)
		logrus.Error(outErr)
	}

	return outErr
}

/*
Allows an admin of group to add a list of user by his name
*/
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
				logrus.Infof("adminID %d adding SUCCESFULLY users %v, in groupID %d", group.AdminID, userName, group.AdminID)
				outErr = nil
			} else {
				outErr = fmt.Errorf("error during add userID list in the groupID %d: %w", group.AdminID, err)
				logrus.Error(outErr)
			}

		} else {
			outErr = fmt.Errorf("number of usernames different from the number of userIDs, that username probably doesn't exist")
			logrus.Error(outErr)
		}

	} else {
		outErr = fmt.Errorf("error during getting userID by userName: %w", err)
		logrus.Error(outErr)
	}

	return outErr

}

/*
Allows an admin of group to set the group name
*/
func SetNameGroup(group *entity.Group, new_name string) error {
	logrus.Debug("Entered in SetNameGroup()")
	logrus.Infof("adminID %d setting the name of groupID %d", group.AdminID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d setting the name of groupID %d", group.AdminID, group.GroupID)

	//Effettua l'operazione al DB
	err := dml.UpdateNameGroup(group.GroupID, group.AdminID, new_name)
	logrus.Debug("Passed by UpdateNameGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		group.Name = new_name
		outErr = nil
		logrus.Infof("adminID %d setting SUCCESFULLY the name '%s' of groupID %d", group.AdminID, new_name, group.GroupID)
	} else {
		outErr = fmt.Errorf("error setting new group name '%s' in groupID %d: %w", new_name, group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

/*
Allows an admin of group to set the group picture
*/
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
		logrus.Infof("adminID %d setting SUCCESFULLY the photo of groupID %d", group.AdminID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error setting new group photo in groupID %d: %w", group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

/*
Allows an admin of group to kick out a user by his username
*/
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
			logrus.Infof("adminID %d kick out SUCCESFULLY the userName '%s' from groupID %d", group.AdminID, userName, group.GroupID)
		} else {
			outErr = fmt.Errorf("Unable to kick out the userName '%s' from groupID %d",userName, group.GroupID)
			logrus.Error(outErr)
		}
	} else {
		outErr = fmt.Errorf("error during search userID by '%s' for kick out him in groupID %d: %w", userName, group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}

/*
Allows an admin of group to move to trash all the group
*/
func DeleteGroup(group *entity.Group) error {
	//La cancellazione del gruppo comporta la cancellazione di membri, eventi, messaggi, contenuti ecc.

	logrus.Debug("Entered in DeleteGroup()")
	logrus.Warningf("adminID %d move to trash the groupID %d", group.AdminID, group.GroupID)

	//Imposta i default
	outErr := fmt.Errorf("Unable for adminID %d move to trash the groupID %d", group.AdminID, group.GroupID)

	//Effettua l'operazione al DB
	err := dml.DeleteGroup(group.GroupID, group.AdminID)
	logrus.Debug("Passed by DeleteGroup()")

	//Se non ci sono errori, prosegui
	if err == nil {
		//Azzeralo
		group = &entity.Group{}
		outErr = nil
		logrus.Infof("adminID %d moved SUCCESFULLY to trash the groupID %d", group.AdminID, group.GroupID)
	} else {
		outErr = fmt.Errorf("error move to trash the group %d: %w", group.GroupID, err)
		logrus.Error(outErr)
	}

	return outErr
}
