package dml

import (
	"fmt"
)

func CreateGroup(name string, photo string, admin int) (int, error) {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert the new group %s in the database", name)
	groupID := 0

	//dml := "UPDATE ..."

	//err = UpdateDataToDatabase(query)

	return groupID, outErr
}

func LeaveGroup(userID int, groupID int) error {
	//Imposta i default
	outErr := fmt.Errorf("Unable to remove row in Members table for groupID %d performing by userID %d", groupID, userID)

	//dml := "DELETE ..."

	//err = DeleteDataToDatabase(query)

	return outErr
}

func AddToGroup(adminID int, userID []int, groupID int) error {
	//Imposta i default
	outErr := fmt.Errorf("Unable to insert row in Members table for groupID %d performing by adminID %d", groupID, adminID)

	//dml := "INSERT ..."

	//err = InsertDataToDatabase(query)

	return outErr
}

func UpdateNameGroup(groupID int, adminID int, new_name string) error {
	//Imposta i default
	outErr := fmt.Errorf("For UpdateNameGroup(), unable to update the row of groupID %d in the database", groupID)

	//dml := "UPDATE ..."

	//err = UpdateDataToDatabase(query)

	return outErr
}

func UpdatePhotoGroup(groupID int, adminID int, new_photo string) error {
	//Imposta i default
	outErr := fmt.Errorf("For UpdatePhotoGroup(),unable to update the row of groupID %d in the database", groupID)

	//dml := "UPDATE ..."

	//err = UpdateDataToDatabase(query)

	return outErr
}

func DeleteGroup(groupID int, userID int) error {
	//La cancellazione del gruppo comporta la cancellazione di membri, eventi, messaggi, contenuti ecc.
	//Imposta i default
	outErr := fmt.Errorf("For DeleteGroup(), unable to remove row in Groups table for groupID %d performing by userID %d", groupID, userID)

	//dml := "DELETE ..."

	//err = DeleteDataToDatabase(query)

	return outErr
}
