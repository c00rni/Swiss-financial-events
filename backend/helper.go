package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func respondWithoutContent(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	return
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, errorMessage string) error {
	type response struct {
		Error string `json:"error"`
	}

	return respondWithJSON(w, code, response{Error: errorMessage})
}

func cfasocietyDateFormater(rawDate string) (time.Time, time.Time, error) {
	rawDates := strings.Split(rawDate, " - ")
	rawDates = append(strings.Split(rawDates[0], ", "), rawDates[1])
	rawDates = strings.Split(strings.Join(rawDates, " "), " ")

	day, err := strconv.Atoi(rawDates[0])
	if err != nil {
		return time.Unix(0, 0), time.Unix(0, 0), err
	}

	months := map[string]time.Month{
		"Jan": time.January,
		"Feb": time.February,
		"Mar": time.March,
		"Apr": time.April,
		"May": time.May,
		"Jun": time.June,
		"Jul": time.July,
		"Aug": time.August,
		"Sep": time.September,
		"Oct": time.October,
		"Nov": time.November,
		"Dec": time.December,
	}
	month, ok := months[rawDates[1]]
	if !ok {
		return time.Unix(0, 0), time.Unix(0, 0), errors.New(fmt.Sprintf("'%v' not recognize as a valid month.", rawDates[1]))
	}

	year, err := strconv.Atoi(rawDates[2])
	if err != nil {
		return time.Unix(0, 0), time.Unix(0, 0), err
	}

	days := map[string]int{
		"Mon": 0,
		"Tue": 1,
		"Wed": 2,
		"Thu": 3,
		"Fri": 4,
		"Sat": 5,
		"Sun": 6,
	}
	startDayIndice, ok := days[rawDates[3]]
	var startHour int
	var startMin int
	var endHour int
	var endMin int
	var eventDays int

	if !ok {
		startTimeStr := strings.Split(rawDates[3], ":")
		startHour, err = strconv.Atoi(startTimeStr[0])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		startMin, err = strconv.Atoi(startTimeStr[1])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		endTimeStr := strings.Split(rawDates[4], ":")
		endHour, err = strconv.Atoi(endTimeStr[0])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		endMin, err = strconv.Atoi(endTimeStr[1])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}
	} else {
		endDayIndice := days[rawDates[5]]
		if startDayIndice < endDayIndice {
			eventDays = endDayIndice - startDayIndice
		} else {
			eventDays = 7 - startDayIndice + endDayIndice
		}
		startTimeStr := strings.Split(rawDates[4], ":")
		startHour, err = strconv.Atoi(startTimeStr[0])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		startMin, err = strconv.Atoi(startTimeStr[1])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		endTimeStr := strings.Split(rawDates[6], ":")
		endHour, err = strconv.Atoi(endTimeStr[0])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}

		endMin, err = strconv.Atoi(endTimeStr[1])
		if err != nil {
			return time.Unix(0, 0), time.Unix(0, 0), err
		}
	}

	startDate := time.Date(year, month, day, startHour, startMin, 0, 0, time.UTC)
	endDate := time.Date(year, month, day+eventDays, endHour, endMin, 0, 0, time.UTC)
	return startDate, endDate, nil
}

func decodeJSONBody(r *http.Request, ptr interface{}) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(ptr)
	if err != nil {
		return err
	}
	return nil
}
