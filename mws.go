/*
	All rights reserved. (c) 2021
*/
package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

const (
	hashCost = 10

	urlAdd          = "/add"
	urlPostAdd      = "/post/add"
	urlHello        = "/hello"
	urlWeeklyUse    = "/week"
	urlRegister     = "/signup"
	urlLogin        = "/login"
	urlLogout       = "/logout"
	urlPostLogin    = "/post/login"
	urlPostRegister = "/post/signup"
	urlPrivacy      = "/privacy"
	urlTerms        = "/terms"
	urlSupport      = "/support"
	tmplBase        = "templates/"
	tmplIndex       = tmplBase + "index.html"
	tmplAdd         = tmplBase + "add.html"
	tmplHello       = tmplBase + "hello.html"
	tmplWeeklyUse   = tmplBase + "weekly.html"
	tmplRegister    = tmplBase + "signup.html"
	tmplLogin       = tmplBase + "signin.html"
	tmplParts       = tmplBase + "parts.html"
	tmplPrivacy     = tmplBase + "privacy.html"
	tmplTerms       = tmplBase + "terms.html"
	tmplSupport     = tmplBase + "support.html"
)

// MedicineData holds all medicine database columns
type MedicineData struct {
	ID       int
	Name     string
	Desc     string
	Producer string
	Size     string
}

// MedicineEntryData holds instance data
type MedicineEntryData struct {
	ID          int
	MedicineID  int
	EntryDate   string
	FinalDate   string
	Name        string
	Producer    string
	Description string
}

// MedicineAlarmedEntryData holds instance data
type MedicineAlarmedEntryData struct {
	ID          int
	MedicineID  int
	EntryDate   string
	FinalDate   string
	Name        string
	Producer    string
	Description string
	Alarm       string
}

// MedicineUseAlarmEntryData holds instance data
type MedicineUseAlarmEntryData struct {
	ID         int
	MedicineID int
	Name       string
	Size       string
	Count      string
	Hour       string
}

// MedicineListingData holds all listing data
type MedicineListingData struct {
	Alarmed    []MedicineAlarmedEntryData
	Expired    []MedicineEntryData
	NotExpired []MedicineEntryData
}

// MedicineWeekListingData holds all listing data
type MedicineWeekListingData struct {
	Mon []MedicineUseAlarmEntryData
	Tue []MedicineUseAlarmEntryData
	Wed []MedicineUseAlarmEntryData
	Thu []MedicineUseAlarmEntryData
	Fri []MedicineUseAlarmEntryData
	Sat []MedicineUseAlarmEntryData
	Sun []MedicineUseAlarmEntryData
}

// AlarmData holds alarm entry data
type AlarmData struct {
	EntryID     int
	Time        int
	TimeType    string
	BeforeAfter string
}

// UseAlarmData holds alarm entry data
type UseAlarmData struct {
	EntryID int
	Hour    string
	Mon     string
	Tue     string
	Wed     string
	Thu     string
	Fri     string
	Sat     string
	Sun     string
}

type ByHour []MedicineUseAlarmEntryData

func (a ByHour) Len() int           { return len(a) }
func (a ByHour) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByHour) Less(i, j int) bool { return a[i].Hour < a[j].Hour }

var router = mux.NewRouter()
var db *sql.DB
var tmpl = make(map[string]*template.Template)

func getDate() string {
	current := time.Now().UTC()
	return current.Format("2006-01-02 15:04:05 -0700")
}

func getDateDaysAfter(days int) string {
	return time.Now().UTC().AddDate(0, 0, days).Format("2006-01-02 15:04:05 -0700")
}

func getDateWeeksAfter(weeks int) string {
	return time.Now().UTC().AddDate(0, 0, 7*weeks).Format("2006-01-02 15:04:05 -0700")
}

func getDateMonthsAfter(months int) string {
	return time.Now().UTC().AddDate(0, months, 0).Format("2006-01-02 15:04:05 -0700")
}

func getDateYearsAfter(years int) string {
	return time.Now().UTC().AddDate(years, 0, 0).Format("2006-01-02 15:04:05 -0700")
}

func getMedicineNameFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT name FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineNameFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineSizeFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT size FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineSizeFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineSizeTypeFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT size_type FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineSizeTypeFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineCountFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT med_count FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineCountFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineTypeFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT type FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineTypeFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineProducerFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT producer FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineTypeFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func getMedicineDescFromID(medID int) (medName string) {
	result := db.QueryRow("SELECT description FROM medicine WHERE medicine_id=$1", medID)
	err := result.Scan(&medName)

	if err != nil {
		fmt.Printf("ERROR getMedicineTypeFromID(%d): %s\n", medID, err)
		return
	}

	return medName
}

func main() {
	var err error

	// Database
	if db, err = sql.Open("sqlite3", "./mws.db"); err != nil {
		panic(err)
	}

	// Prepare templates
	tmpl[tmplIndex] = template.Must(template.ParseFiles(tmplIndex, tmplParts))
	tmpl[tmplAdd] = template.Must(template.ParseFiles(tmplAdd, tmplParts))
	tmpl[tmplHello] = template.Must(template.ParseFiles(tmplHello, tmplParts))
	tmpl[tmplWeeklyUse] = template.Must(template.ParseFiles(tmplWeeklyUse, tmplParts))
	tmpl[tmplRegister] = template.Must(template.ParseFiles(tmplRegister, tmplParts))
	tmpl[tmplLogin] = template.Must(template.ParseFiles(tmplLogin, tmplParts))
	tmpl[tmplPrivacy] = template.Must(template.ParseFiles(tmplPrivacy, tmplParts))
	tmpl[tmplTerms] = template.Must(template.ParseFiles(tmplTerms, tmplParts))
	tmpl[tmplSupport] = template.Must(template.ParseFiles(tmplSupport, tmplParts))

	// Function pages
	router.HandleFunc(urlLogin, loginHandler)
	router.HandleFunc(urlLogout, LogoutHandler)
	router.HandleFunc(urlRegister, registerHandler)
	router.HandleFunc(urlPostLogin, postLoginHandler).Methods("POST")
	router.HandleFunc(urlPostRegister, postRegisterHandler).Methods("POST")

	// Pages
	router.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		// Check login status
		if getUserName(request) == "" {
			http.Redirect(response, request, urlHello, 302)
			return
		}

		// Get alarms
		var alarms []AlarmData

		row, err := db.Query("SELECT entry_id, timer, timer_type, before_after FROM expire_alarms WHERE user_id=$1", getUserID(getUserName(request)))
		if err != nil {
			panic(err)
		}
		defer row.Close()

		for row.Next() {
			var entryID int
			var timer int
			var timerType sql.NullString
			var beforeAfter sql.NullString

			err = row.Scan(&entryID, &timer, &timerType, &beforeAfter)
			if err != nil {
				panic(err)
			}

			alarms = append(alarms, AlarmData{EntryID: entryID, Time: timer, TimeType: timerType.String, BeforeAfter: beforeAfter.String})
		}

		// Get entries
		var listingData MedicineListingData

		row, err = db.Query("SELECT entry_id, medicine_id, entry_date, expire_date FROM entries WHERE user_id=$1 ORDER BY expire_date ASC", getUserID(getUserName(request)))
		if err != nil {
			panic(err)
		}
		defer row.Close()

		for row.Next() {
			var id int
			var medicineID int
			var entryDate sql.NullString
			var finalDate sql.NullString

			err = row.Scan(&id, &medicineID, &entryDate, &finalDate)
			if err != nil {
				panic(err)
			}

			// Find alarm
			var alarmStr string
			var myAlarmDate string

			for i := range alarms {
				if alarms[i].EntryID == id {
					alarmStr = fmt.Sprintf("%d %s %s", alarms[i].Time, alarms[i].TimeType, alarms[i].BeforeAfter)

					if alarms[i].TimeType == "Day" {
						myAlarmDate = getDateDaysAfter(alarms[i].Time)
					} else if alarms[i].TimeType == "Week" {
						myAlarmDate = getDateWeeksAfter(alarms[i].Time)
					} else if alarms[i].TimeType == "Month" {
						myAlarmDate = getDateMonthsAfter(alarms[i].Time)
					} else if alarms[i].TimeType == "Year" {
						myAlarmDate = getDateYearsAfter(alarms[i].Time)
					}

					break
				}
			}

			// Separate them
			if finalDate.Valid {
				myEntryDate, err := time.Parse("2006-01-02 15:04:05 -0700", entryDate.String)
				if err != nil {
					return
				}

				outEntryDate := myEntryDate.Format("02/01/2006 15:04")

				myFinalDate, err := time.Parse("2006-01-02 15:04:05 -0700", finalDate.String)
				if err != nil {
					return
				}

				outFinalDate := myFinalDate.Format("02/01/2006 15:04")

				if finalDate.String < getDate() {
					listingData.Expired = append(listingData.Expired, MedicineEntryData{ID: id, MedicineID: medicineID, EntryDate: outEntryDate, FinalDate: outFinalDate, Name: getMedicineNameFromID(medicineID), Producer: getMedicineProducerFromID(medicineID), Description: getMedicineDescFromID(medicineID)})
				} else if finalDate.String < myAlarmDate {
					listingData.Alarmed = append(listingData.Alarmed, MedicineAlarmedEntryData{ID: id, MedicineID: medicineID, EntryDate: outEntryDate, FinalDate: outFinalDate, Name: getMedicineNameFromID(medicineID), Producer: getMedicineProducerFromID(medicineID), Description: getMedicineDescFromID(medicineID), Alarm: alarmStr})
				} else {
					listingData.NotExpired = append(listingData.NotExpired, MedicineEntryData{ID: id, MedicineID: medicineID, EntryDate: outEntryDate, FinalDate: outFinalDate, Name: getMedicineNameFromID(medicineID), Producer: getMedicineProducerFromID(medicineID), Description: getMedicineDescFromID(medicineID)})
				}
			}
		}

		// Execute template with prepared data
		err = tmpl[tmplIndex].Execute(response, listingData)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlWeeklyUse, func(response http.ResponseWriter, request *http.Request) {
		// Check login status
		if getUserName(request) == "" {
			http.Redirect(response, request, urlHello, 302)
			return
		}

		// Get use alarms
		var useAlarms []UseAlarmData

		row, err := db.Query("SELECT entry_id, mon, tue, wed, thu, fri, sat, sun, hour FROM use_alarms WHERE user_id=$1 ORDER BY hour ASC", getUserID(getUserName(request)))
		if err != nil {
			panic(err)
		}
		defer row.Close()

		for row.Next() {
			var entryID int
			var hour sql.NullString
			var mon sql.NullString
			var tue sql.NullString
			var wed sql.NullString
			var thu sql.NullString
			var fri sql.NullString
			var sat sql.NullString
			var sun sql.NullString

			err = row.Scan(&entryID, &mon, &tue, &wed, &thu, &fri, &sat, &sun, &hour)
			if err != nil {
				panic(err)
			}

			useAlarms = append(useAlarms, UseAlarmData{EntryID: entryID, Mon: mon.String, Tue: tue.String, Wed: wed.String, Thu: thu.String, Fri: fri.String, Sat: sat.String, Sun: sun.String, Hour: hour.String})
		}

		// Get entries
		var weekListData MedicineWeekListingData

		row, err = db.Query("SELECT entry_id, medicine_id, entry_date, expire_date FROM entries WHERE user_id=$1", getUserID(getUserName(request)))
		if err != nil {
			panic(err)
		}
		defer row.Close()

		for row.Next() {
			var id int
			var medicineID int
			var entryDate sql.NullString
			var finalDate sql.NullString

			err = row.Scan(&id, &medicineID, &entryDate, &finalDate)
			if err != nil {
				panic(err)
			}

			// Find alarm
			var myAlarm UseAlarmData

			for i := range useAlarms {
				if useAlarms[i].EntryID == id {
					myAlarm = useAlarms[i]

					break
				}
			}

			// Separate them
			if myAlarm.Mon == "on" {
				weekListData.Mon = append(weekListData.Mon, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Tue == "on" {
				weekListData.Tue = append(weekListData.Tue, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Wed == "on" {
				weekListData.Wed = append(weekListData.Wed, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Thu == "on" {
				weekListData.Thu = append(weekListData.Thu, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Fri == "on" {
				weekListData.Fri = append(weekListData.Fri, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Sat == "on" {
				weekListData.Sat = append(weekListData.Sat, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}

			if myAlarm.Sun == "on" {
				weekListData.Sun = append(weekListData.Sun, MedicineUseAlarmEntryData{ID: id, MedicineID: medicineID, Name: getMedicineNameFromID(medicineID), Size: fmt.Sprintf("%s %s", getMedicineSizeFromID(medicineID), getMedicineSizeTypeFromID(medicineID)), Count: fmt.Sprintf("%s %s", getMedicineCountFromID(medicineID), getMedicineTypeFromID(medicineID)), Hour: myAlarm.Hour})
			}
		}

		// Sort
		sort.Sort(ByHour(weekListData.Mon))
		sort.Sort(ByHour(weekListData.Tue))
		sort.Sort(ByHour(weekListData.Wed))
		sort.Sort(ByHour(weekListData.Thu))
		sort.Sort(ByHour(weekListData.Fri))
		sort.Sort(ByHour(weekListData.Sat))
		sort.Sort(ByHour(weekListData.Sun))

		// Execute template with prepared data
		err = tmpl[tmplWeeklyUse].Execute(response, weekListData)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlAdd, func(response http.ResponseWriter, request *http.Request) {
		// Check login status
		if getUserName(request) == "" {
			http.Redirect(response, request, urlHello, 302)
			return
		}

		err := tmpl[tmplAdd].Execute(response, nil)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlPostAdd, func(response http.ResponseWriter, request *http.Request) {
		// Check if user logged in
		if getUserName(request) == "" {
			http.Redirect(response, request, "/", 302)
			return
		}

		// Get the form data
		name := request.FormValue("medicineName")
		firm := request.FormValue("medicineFirm")
		expDate := request.FormValue("medicineExpDate")
		count := request.FormValue("entryCount")
		desc := request.FormValue("medicineDescription")
		size := request.FormValue("medicineSizePerBox")
		sizeType := request.FormValue("medicineSizeType")
		medCount := request.FormValue("medicineCountPerBox")
		medType := request.FormValue("medicineType")

		expName := request.FormValue("expireAlarmName")
		expTime := request.FormValue("expireAlarmTime")
		expType := request.FormValue("expireAlarmTimeType")
		expBeforeAfter := request.FormValue("expireAlarmBeforeAfter")
		expAction := request.FormValue("expireAlarmAction")

		mon := request.FormValue("useAlarmMonday")
		tue := request.FormValue("useAlarmTuesday")
		wed := request.FormValue("useAlarmWednesday")
		thu := request.FormValue("useAlarmThursday")
		fri := request.FormValue("useAlarmFriday")
		sat := request.FormValue("useAlarmSaturday")
		sun := request.FormValue("useAlarmSunday")
		useTime := request.FormValue("useAlarmTime")

		// Check if any of necessary fields are empty
		if name == "" || firm == "" || expDate == "" ||
			count == "" || size == "" || sizeType == "" ||
			medCount == "" || medType == "" || expName == "" ||
			expTime == "" || expType == "" || expBeforeAfter == "" ||
			expAction == "" || useTime == "" {
			http.Redirect(response, request, urlAdd, 302)
			return
		}

		// Turn expiration date into proper time data
		realExpDate, err := time.Parse("02/01/2006 15:04", expDate)

		if err != nil {
			fmt.Println("error 1")
			http.Redirect(response, request, urlAdd, 302)
			return
		}

		// Turn alarm clock into proper time data (only for checking)
		_, err = time.Parse("15:04", useTime)

		if err != nil {
			fmt.Println("error 2")
			http.Redirect(response, request, urlAdd, 302)
			return
		}

		// Insert the data into DB
		redirectTarget := "/"
		userID := getUserID(getUserName(request))

		// Prepare medicine data
		medicineSQLStatement := `INSERT INTO medicine(user_id,name,producer,description,size,size_type,med_count,type) VALUES(?,?,?,?,?,?,?,?)`
		medicineStatement, err := db.Prepare(medicineSQLStatement)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		medResult, err := medicineStatement.Exec(userID, name, firm, desc, size, sizeType, medCount, medType)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		medID, err := medResult.LastInsertId()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Prepare entry data
		entrySQLStatement := `INSERT INTO entries(medicine_id,user_id,entry_date,expire_date) VALUES(?,?,?,?)`
		entryStatement, err := db.Prepare(entrySQLStatement)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		entryResult, err := entryStatement.Exec(medID, userID, getDate(), realExpDate.Format("2006-01-02 15:04:05 -0700"))
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		entryID, err := entryResult.LastInsertId()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Prepare expire alarm
		expireSQLStatement := `INSERT INTO expire_alarms(entry_id,user_id,timer,timer_type,before_after,action) VALUES(?,?,?,?,?,?)`
		expireStatement, err := db.Prepare(expireSQLStatement)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = expireStatement.Exec(entryID, userID, expTime, expType, expBeforeAfter, expAction)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Prepare use alarm
		useSQLStatement := `INSERT INTO use_alarms(entry_id,user_id,mon,tue,wed,thu,fri,sat,sun,hour) VALUES(?,?,?,?,?,?,?,?,?,?)`
		useStatement, err := db.Prepare(useSQLStatement)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		_, err = useStatement.Exec(entryID, userID, mon, tue, wed, thu, fri, sat, sun, useTime)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		// Redirect user to login page
		fmt.Println("DEBUG: Successful entry record")

		http.Redirect(response, request, redirectTarget, 302)
	})

	router.HandleFunc(urlHello, func(response http.ResponseWriter, request *http.Request) {
		// Check login status
		if getUserName(request) != "" {
			http.Redirect(response, request, "/", 302)
			return
		}

		err := tmpl[tmplHello].Execute(response, nil)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlPrivacy, func(response http.ResponseWriter, request *http.Request) {
		err := tmpl[tmplPrivacy].Execute(response, nil)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlTerms, func(response http.ResponseWriter, request *http.Request) {
		err := tmpl[tmplTerms].Execute(response, nil)

		if err != nil {
			return
		}
	})

	router.HandleFunc(urlSupport, func(response http.ResponseWriter, request *http.Request) {
		err := tmpl[tmplSupport].Execute(response, nil)

		if err != nil {
			return
		}
	})

	// File server
	router.PathPrefix("/res/").Handler(http.StripPrefix("/res/", http.FileServer(http.Dir("static"))))

	// Server
	http.ListenAndServe(":8090", router)

	// Shutting down
	db.Close()
}
