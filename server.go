package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"github.com/asafron/meetings-scheduler/db"
	"github.com/asafron/meetings-scheduler/controllers"
	"github.com/asafron/meetings-scheduler/config"
	"github.com/asafron/meetings-scheduler/auth"
	"errors"
	"github.com/asafron/meetings-scheduler/helpers"
)


func main() {
	// logging
	log.SetFormatter(&log.TextFormatter{})

	configWrapper := config.GetConfigWrapper()

	// db
	dal := initMongo(configWrapper.GetCurrent().MongoHost)
	log.Info("DB connection was established")
	defer dal.Close()

	authorizer := auth.NewAuthenticator(dal, config.GetConfigWrapper().GetCurrent().SessionKey)

	// controllers
	mc := controllers.NewMeetingsController(dal)
	ac := controllers.NewAdminController(dal)

	r := mux.NewRouter()
	r.Handle("/ws/version", requestQueueHandler(http.HandlerFunc(Version))).Methods("GET")

	// client
	r.Handle("/ws/meetings/addAvailableTime", requestQueueHandler(http.HandlerFunc(mc.AddAvailableTime))).Methods("POST")
	r.Handle("/ws/meetings/schedule", requestQueueHandler(http.HandlerFunc(mc.ScheduleMeeting))).Methods("POST")
	//r.Handle("/ws/meetings/getAllMeetings", requestQueueHandler(http.HandlerFunc(mc.GetAllMeetings))).Methods("GET")
	r.Handle("/ws/meetings/getAvailableMeetings", requestQueueHandler(http.HandlerFunc(mc.GetAvailableMeetings))).Methods("GET")
	//r.Handle("/ws/meetings/getScheduledMeetings", requestQueueHandler(http.HandlerFunc(mc.GetScheduledMeetings))).Methods("GET")

	// admin
	r.Handle("/applications", RecoverWrap(authorizer.AuthMiddleware(http.HandlerFunc(ac.ManagerGetAllMeetings)))).Methods("GET")
	r.Handle("/ws/admin/getAllMeetingStatus", requestQueueHandler(http.HandlerFunc(ac.ManagerGetAllMeetings))).Methods("POST")
	r.Handle("/ws/admin/cancelMeeting", requestQueueHandler(http.HandlerFunc(ac.ManagerCancelMeeting))).Methods("POST")

	// http setup
	http.Handle("/", &MyServer{r})
	log.Info("starting server, listening on port 4000...")
	err := http.ListenAndServe(":4000",nil)
	if err!=nil {
		panic(err)
	}
}

type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Credentials","true")
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, X-Sdk-Key, Authorization, Cache-control")
	}
	// Stop here if its Pre-flighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func requestQueueHandler(fn http.Handler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		fn.ServeHTTP(rw,req)
	}
}

type VersionResponse struct {
	Version string `json:"version"`
}

func Version(writer http.ResponseWriter, req *http.Request) {
	res := VersionResponse{ Version: "4"}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(js)
}

func initMongo(dbUrl string) *db.DAL{
	dal := db.NewDatabaseAccessor(dbUrl)
	err:= dal.Initialize()
	if err!=nil {
		panic(err)
	}
	return dal
}

func RecoverWrap(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = helpers.GeneralErrorInternal
				}
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(w, req)
	})
}