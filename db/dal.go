package db

import (
	"gopkg.in/mgo.v2"
	"github.com/asafron/meetings-scheduler/models"
	"gopkg.in/mgo.v2/bson"
	"time"
	log "github.com/Sirupsen/logrus"
	"github.com/asafron/meetings-scheduler/helpers"
)

// Database
const dbName = "meeting-scheduler"

// Collections
const dbCollectionUsers = "users"
const dbCollectionEvents = "events"

// Fields
const dbFieldUsersEmail = "email"
const dbFieldUsersDisplayId = "display_id"

const dbFieldEventsDisplayId = "display_id"

type DAL struct {
	session *mgo.Session
}

func NewDatabaseAccessor(url string) *DAL {
	session, err := mgo.Dial(url)
	if err!=nil {
		panic(err)
	}
	return &DAL{session : session}
}

func (dal *DAL) Initialize() (error) {
	usersCollection := dal.session.DB(dbName).C(dbCollectionUsers)
	uniqueIndexes := [][]string {[]string{dbFieldUsersEmail}, []string{dbFieldUsersDisplayId}}
	for _, element := range uniqueIndexes {
		index := mgo.Index {
			Key: element,
			Unique: true,
			DropDups: false,
			Background: false,
			Sparse: false,
		}
		err := usersCollection.EnsureIndex(index)
		if err != nil {
			return err
		}
	}

	meetingsCollection := dal.session.DB(dbName).C(dbCollectionEvents)
	uniqueIndexes = [][]string {[]string{dbFieldEventsDisplayId}}
	for _, element := range uniqueIndexes {
		index := mgo.Index {
			Key: element,
			Unique: true,
			DropDups: false,
			Background: false,
			Sparse: false,
		}
		err := meetingsCollection.EnsureIndex(index)
		if err != nil {
			return err
		}
	}

	return nil
}

func(dal *DAL) Close() {
	dal.session.Close()
}

/* Users */

func (dal *DAL) FindActiveUserByEmail(email string)  (*models.User, error) {
	user := models.User{}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Find(bson.M{"email": email, "status" : models.USER_CONFIRMED}).One(&user)
	if (err != nil) {
		return &user, helpers.AuthenticationErrorLoginUserNotExists
	}
	return &user, nil
}

func (dal *DAL) FindAnyUserByEmail(email string)  (*models.User, error) {
	user := models.User{}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Find(bson.M{"email": email}).One(&user)
	if (err != nil) {
		return &user, helpers.AuthenticationErrorLoginUserNotExists
	}
	return &user, nil
}

func (dal *DAL) InsertUser(email string,hash []byte, firstName string , lastName string, confirmationToken string) error {
	user := models.User{
		Id: bson.NewObjectId(),
		DisplayId: helpers.RandStringBytesMaskImprSrc(8),
		Email: email,
		FirstName: firstName,
		LastName: lastName,
		Hash: hash,
		ConfirmationToken: confirmationToken,
		ConfirmationTokenStatus: models.CONFIRMATION_TOKEN_VALID,
		Confirmed: false,
		Status: models.USER_NOT_CONFIRMED,
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC()}

	err := dal.session.DB(dbName).C(dbCollectionUsers).Insert(user)
	if (err != nil) {
		log.Warn(err)
		return err
	}
	return  nil
}

func (dal *DAL) FindUserByConfirmationToken(confirmationToken string, email string)  (*models.User, error) {
	user := models.User{}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Find(bson.M{"confirmation_token": confirmationToken, "confirmation_token_status" : models.CONFIRMATION_TOKEN_VALID, "email" : email }).One(&user)
	if (err != nil) {
		return &user, helpers.AuthenticationErrorLoginUserNotExists
	}
	return &user, nil
}

func (dal *DAL) FindUserByRecoveryToken(recoveryToken string, email string)  (*models.User, error) {
	user := models.User{}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Find(bson.M{"email" : email, "recovery_token": recoveryToken, "recovery_token_status" : models.RECOVER_TOKEN_VALID, "recovery_token_expiry" : bson.M{ "$gt" : time.Now().UTC()} }).One(&user)
	if (err != nil) {
		return &user, helpers.AuthenticationErrorLoginUserNotExists
	}
	return &user, nil
}

func (dal *DAL) UpdateUserConfirmation(userId bson.ObjectId, userStatus models.UserStatusType, confirmationTokenStatus models.ConfirmationTokenStatusType, confirmed bool) (error) {
	colQueried := bson.M{"_id" : userId}
	change := bson.M{"$set": bson.M{
		"confirmation_token_status": confirmationTokenStatus,
		"status" : userStatus,
		"updated_at": time.Now().UTC(),
		"confirmed" : confirmed}}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Update(colQueried, change)
	if err != nil {
		return err
	}
	return nil
}

func (dal *DAL) UpdateUserPassword(userId bson.ObjectId, hash []byte, recoveryTokenStatus models.RecoverTokenStatusType, recoveryTokenExpiry time.Time) (error) {
	colQueried := bson.M{"_id" : userId, "recovery_token_expiry" : bson.M{ "$gt" : time.Now().UTC()}, "recovery_token_status" : models.RECOVER_TOKEN_VALID}
	change := bson.M{"$set": bson.M{
		"recovery_token_status": recoveryTokenStatus,
		"recovery_token_expiry" : recoveryTokenExpiry,
		"updated_at": time.Now().UTC(),
		"hash" : hash}}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Update(colQueried, change)
	if err != nil {
		return err
	}
	return nil
}

func (dal *DAL) UpdateUserRecovery(userId bson.ObjectId, recoveryToken string, recoveryTokenStatus models.RecoverTokenStatusType, recoveryTokenExpiry time.Time) error {
	colQueried := bson.M{"_id" : userId}
	change := bson.M{"$set": bson.M{
		"updated_at": time.Now().UTC(),
		"recovery_token_status": recoveryTokenStatus,
		"recovery_token_expiry" : recoveryTokenExpiry,
		"recovery_token": recoveryToken }}
	err := dal.session.DB(dbName).C(dbCollectionUsers).Update(colQueried, change)
	if err != nil {
		return err
	}
	return nil
}

/* Events */

func (dal *DAL) GetEventsForUser(displayId string) *[]models.Event {
	events := []models.Event{}
	err := dal.session.DB(dbName).C(dbCollectionEvents).Find(bson.M{"admin_user": displayId}).All(&events)
	if err != nil {
		log.Info(err)
	}
	return &events
}

func (dal *DAL) InsertEvent(name string, adminUser string, slots []models.Slot, meetings []models.Meeting) error {
	event := models.Event{
		Id: bson.NewObjectId(),
		DisplayId: helpers.RandStringBytesMaskImprSrc(8),
		Name: name,
		AdminUser: adminUser,
		Slots: slots,
		Meetings: meetings,
		CreatedAt:time.Now().UTC(),
		UpdatedAt:time.Now().UTC()}

	err := dal.session.DB(dbName).C(dbCollectionEvents).Insert(event)
	if (err != nil) {
		log.Fatal(err)
		return err
	}
	return  nil
}

func (dal *DAL) UpdateEvent(displayId string, name string, adminUser string, slots []models.Slot, meetings []models.Meeting) error {
	slotsToDb := []models.Slot{}
	for _, element := range slots {
		if len(element.Id) == 0 {
			element.Id = bson.NewObjectId()
			element.DisplayId = helpers.RandStringBytesMaskImprSrc(8)
		}
		slotsToDb = append(slotsToDb, element)
	}

	colQueried := bson.M{"display_id" : displayId}
	change := bson.M{"$set": bson.M{
		"name": name,
		"admin_user" : adminUser,
		"slots" : slotsToDb,
		"meetings" : meetings,
		"updated_at": time.Now().UTC()}}
	err := dal.session.DB(dbName).C(dbCollectionEvents).Update(colQueried, change)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (dal *DAL) RemoveEvent(displayId string) error {
	colQueried := bson.M{"display_id" : displayId}
	err := dal.session.DB(dbName).C(dbCollectionEvents).Remove(colQueried)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (dal *DAL) GetEventByDisplayId(displayId string) (*models.Event, error) {
	event := models.Event{}
	err := dal.session.DB(dbName).C(dbCollectionEvents).Find(bson.M{"display_id": displayId}).One(&event)
	if err != nil {
		log.Info(err)
		return nil, err
	}
	return &event, nil
}

func (dal *DAL) RemoveSlotFromEvent(eventDisplayId string, displayId string) error {
	event, err := dal.GetEventByDisplayId(eventDisplayId)
	if err != nil {
		log.Fatal(err)
		return err
	}

	indexToDelete := -1
	for index, element := range event.Slots {
		if element.DisplayId == displayId {
			indexToDelete = index
			break
		}
	}

	if indexToDelete == -1 {
		log.Fatal(helpers.SlotsErrorNotFound)
		return helpers.SlotsErrorNotFound
	}

	event.Slots = append(event.Slots[:indexToDelete], event.Slots[indexToDelete+1:]...)

	err = dal.UpdateEvent(event.DisplayId, event.Name, event.AdminUser, event.Slots, event.Meetings)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}