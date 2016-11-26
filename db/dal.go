package db

import (
	"gopkg.in/mgo.v2"
	"errors"
	"github.com/asafron/meetings-scheduler/models"
	"gopkg.in/mgo.v2/bson"
	"time"
	"fmt"
)

const dbName = "meeting-scheduler"
const dbCollectionMeetings = "meetings"

const dbFieldMeetingsDay = "day"
const dbFieldMeetingsMonth = "month"
const dbFieldMeetingsYear = "year"
const dbFieldMeetingsStartTime = "start_time"
const dbFieldMeetingsEndTime = "end_time"

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
		dbFieldMeetingsStartTime, dbFieldMeetingsEndTime}}
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

func (dal *DAL) InsertAvailableMeetingTime(day, month, year, startTime, endTime int) error {
	meeting := models.Meeting{}
	meeting.Id = bson.NewObjectId()
	meeting.Day = day
	meeting.Month = month
	meeting.Year = year
	meeting.StartTime = startTime
	meeting.EndTime = endTime
	meeting.UserName = ""
	meeting.UserEmail = ""
	meeting.UserPhone = ""
	meeting.CreatedAt = time.Now().UTC()
	meeting.UpdatedAt = time.Now().UTC()

	err := dal.session.DB(dbName).C(dbCollectionMeetings).Insert(meeting)
	if (err != nil) {
		return errors.New(fmt.Sprintf("meeting already exists on: %s", meeting.GetMeetingDateAsString()))
	}

	return nil
}

func (dal *DAL) GetMeetingByTime(day, month, year, startTime, endTime int) *models.Meeting {
	meeting := models.Meeting{}
	query := bson.M{"day": day, "month": month, "year": year, "start_time": startTime, "end_time": endTime }
	err := dal.session.DB(dbName).C(dbCollectionMeetings).Find(query).Limit(1).One(&meeting)
	if (err != nil) {
		return nil
	}
	return &meeting
}

func (dal *DAL) UpdateMeetingDetails(day, month, year, startTime, endTime int, name, email, phone string) error {
	meeting := dal.GetMeetingByTime(day, month, year, startTime, endTime)
	if meeting == nil {
		return errors.New("no available meeting at that time")
	}
	if len(meeting.UserName) > 0 {
		return errors.New("meeting time already taken")
	}
	colQuerier := bson.M{"_id" : meeting.Id}
	change := bson.M{"$set": bson.M{
		"user_name": name,
		"user_email" : email,
		"user_phone" : phone,
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