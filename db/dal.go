package db

import (
	"gopkg.in/mgo.v2"
	"errors"
	"github.com/asafron/meetings-scheduler/models"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/asafron/meetings-scheduler/helpers"
)

const dbName = "meeting-scheduler"
const dbCollectionMeetings = "meetings"

const dbFieldMeetingsDay = "day"
const dbFieldMeetingsMonth = "month"
const dbFieldMeetingsYear = "year"
const dbFieldMeetingsStartTime = "start_time"
const dbFieldMeetingsEndTime = "end_time"
const dbFieldMeetingsRepresentative = "representative"
const dbFieldMeetingsDisplayId = "display_id"

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
	// Applications indices
	meetingsCollection := dal.session.DB(dbName).C(dbCollectionMeetings)

	// Index to ensure no 2 meetings with the same time
	uniqueIndexes := [][]string {[]string{dbFieldMeetingsDay,
		dbFieldMeetingsMonth, dbFieldMeetingsYear,
		dbFieldMeetingsStartTime, dbFieldMeetingsEndTime, dbFieldMeetingsRepresentative}, []string{dbFieldMeetingsDisplayId}}
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

func (dal *DAL) InsertAvailableMeetingTime(day, month, year, startTime, endTime int, representative string) error {
	meeting := models.Meeting{}
	meeting.Id = bson.NewObjectId()
	meeting.DisplayId = helpers.RandStringBytesMaskImprSrc(8)
	meeting.Day = day
	meeting.Month = month
	meeting.Year = year
	meeting.StartTime = startTime
	meeting.EndTime = endTime
	meeting.Representative = representative
	meeting.UserName = ""
	meeting.UserEmail = ""
	meeting.UserPhone = ""
	meeting.UserSchool = ""
	meeting.UserIdNumber = ""
	meeting.UserPreferredSchoolDay = ""
	meeting.CreatedAt = time.Now().UTC()
	meeting.UpdatedAt = time.Now().UTC()

	err := dal.session.DB(dbName).C(dbCollectionMeetings).Insert(meeting)
	if (err != nil) {
		log.Error(err)
		return errors.New(fmt.Sprintf("meeting already exists on: %s", meeting.GetMeetingDateAsString()))
	}

	return nil
}

func (dal *DAL) GetMeetingByTime(day, month, year, startTime, endTime int) []models.Meeting {
	meetings := []models.Meeting{}
	query := bson.M{"day": day, "month": month, "year": year, "start_time": startTime, "end_time": endTime }
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Find(query).All(&meetings)
	if (err != nil) {
		return nil
	}
	return meetings
}

func (dal *DAL) UpdateMeetingDetails(day, month, year, startTime, endTime int, name, email, phone, school, idNumber, schoolDay string) error {
	allMeetings := dal.GetMeetingByTime(day, month, year, startTime, endTime)
	log.Info(fmt.Sprintf("found %d meetings", len(allMeetings)))
	if len(allMeetings) == 0 {
		return errors.New("no meetings at that time")
	}
	var meeting models.Meeting

	// check if meeting time available
	meetingAvailable := false
	for i := 0; i < len(allMeetings); i++ {
		meeting = allMeetings[i]
		log.Info(fmt.Sprintf("current check meeting has the UserName %s", meeting.UserName))
		if len(meeting.UserName) == 0 {
			meetingAvailable = true
			break
		}
	}
	if !meetingAvailable {
		return errors.New("no available meetings at that time")
	}

	// check if id number already exists
	allMeetings = dal.GetAllMeetings()
	idNumberExists := false
	for i := 0; i < len(allMeetings); i++ {
		mtg := allMeetings[i]
		if mtg.UserIdNumber == idNumber {
			log.Info(fmt.Sprintf("interation # %d", i))
			idNumberExists = true
			break
		}
	}
	if idNumberExists {
		return errors.New("a meeting with the following id number already exists")
	}

	log.Info(meeting.Id)

	colQuerier := bson.M{"_id" : meeting.Id}
	change := bson.M{"$set": bson.M{
		"user_name": name,
		"user_email" : email,
		"user_phone" : phone,
		"user_school" : school,
		"user_id_number" : idNumber,
		"user_preferred_school_day" : schoolDay,
		"updated_at": time.Now().UTC(),
	}}
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Update(colQuerier, change)
	if err != nil {
		return err
	}
	return nil
}

func (dal *DAL) GetAllMeetings() []models.Meeting {
	meetings := []models.Meeting{}
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Find(bson.M{}).All(&meetings)
	if err != nil {
		return nil
	}
	return meetings
}

func (dal *DAL) GetAvailableMeetings() []models.Meeting {
	meetings := dal.GetAllMeetings()
	availableMeetings := []models.Meeting{}
	if meetings != nil {
		for i := 0; i < len(meetings); i++ {
			if len(meetings[i].UserName) == 0 {
				availableMeetings = append(availableMeetings, meetings[i])
			}
		}
	}
	return availableMeetings
}

func (dal *DAL) GetScheduledMeetings() []models.Meeting {
	meetings := dal.GetAllMeetings()
	scheduledMeetings := []models.Meeting{}
	if meetings != nil {
		for i := 0; i < len(meetings); i++ {
			if len(meetings[i].UserName) > 0 {
				scheduledMeetings = append(scheduledMeetings, meetings[i])
			}
		}
	}
	return scheduledMeetings
}

func (dal *DAL) GetMeetingByDisplayId(displayId string) *models.Meeting {
	meeting := &models.Meeting{}
	query := bson.M{"display_id": displayId}
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Find(query).One(&meeting)
	if err == nil {
		return meeting
	}
	return nil
}

func (dal *DAL) CancelMeeting(displayId string) error {
	meeting := dal.GetMeetingByDisplayId(displayId)
	if meeting == nil {
		return errors.New("no such meeting")
	}

	colQuerier := bson.M{"_id" : meeting.Id}
	change := bson.M{"$set": bson.M{
		"user_name": "",
		"user_email" : "",
		"user_phone" : "",
		"user_school" : "",
		"user_id_number" : "",
		"user_preferred_school_day" : "",
		"updated_at": time.Now().UTC(),
	}}
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Update(colQuerier, change)
	return err
}