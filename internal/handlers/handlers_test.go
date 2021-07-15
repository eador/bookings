package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/eador/bookings/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "get", http.StatusOK},
	{"about", "/about", "get", http.StatusOK},
	{"gq", "/generals-quarters", "get", http.StatusOK},
	{"ms", "/majors-suite", "get", http.StatusOK},
	{"contact", "/contact", "get", http.StatusOK},
	{"search-availability", "/search-availability", "get", http.StatusOK},
	/*
		{"make-reservation", "/make-reservation", "get", []postData{}, http.StatusOK},


		{"post-search-avail-json", "/search-availability-json", "post", []postData{
			{key: "start", value: "2020-01-02"},
			{key: "end", value: "2020-02-02"},
		}, http.StatusOK},
		{"post-search-avail", "/search-availability", "post", []postData{
			{key: "start", value: "2020-01-02"},
			{key: "end", value: "2020-02-02"},
		}, http.StatusOK},
		{"post-make-reservation", "/make-reservation", "post", []postData{
			{key: "first_name", value: "John"},
			{key: "last_name", value: "Smith"},
			{key: "email", value: "j.smith@smith.com"},
			{key: "phone", value: "(555) 867-5309"},
		}, http.StatusOK},*/
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Reservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session ( reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test with non-existant room
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	reqBody := "first_name=John"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12345")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing post body
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for missing reservation in session
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for error inserting reservation
	reservation.RoomID = 2
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for error inserting room restriction
	reservation.RoomID = 1000
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// Test for invalid data
	reqBody = "first_name=J"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Smith")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=12345")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@smith.com")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostReservation handler returned wrong response code for invalid data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
}

var availabilityJSONTests = []struct {
	name    string
	reqBody io.Reader
	json    jsonResponse
}{
	{"Empty Body", nil, jsonResponse{
		OK:      false,
		Message: "Internal Server Error",
	}},
	{"Bad Start Date", strings.NewReader("start=bad&end=2050-01-02&room_id=1"), jsonResponse{
		OK:      false,
		Message: "Can not parse start date",
	}},
	{"Bad End Date", strings.NewReader("start=2050-01-02&end=bad&room_id=1"), jsonResponse{
		OK:      false,
		Message: "Can not parse end date",
	}},
	{"Available", strings.NewReader("start=2050-01-01&end=2050-01-02&room_id=1"), jsonResponse{
		OK:        true,
		Message:   "",
		StartDate: "2050-01-01",
		EndDate:   "2050-01-02",
		RoomID:    "1",
	}},
	{"Unavailable", strings.NewReader("start=2050-01-01&end=2050-01-02&room_id=2"), jsonResponse{
		OK:        false,
		Message:   "",
		StartDate: "2050-01-01",
		EndDate:   "2050-01-02",
		RoomID:    "2",
	}},
	{"Database Error", strings.NewReader("start=2050-01-01&end=2050-01-02&room_id=200"), jsonResponse{
		OK:      false,
		Message: "Error connecting to database",
	}},
}

func TestRepository_AvailabilityJSON(t *testing.T) {
	for _, e := range availabilityJSONTests {
		req, _ := http.NewRequest("POST", "/search-availability-json", e.reqBody)
		ctx := GetCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler := http.HandlerFunc(Repo.AvailabilityJSON)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		var j jsonResponse
		err := json.Unmarshal(rr.Body.Bytes(), &j)
		if err != nil {
			t.Error("failed to parse json")
		}
		if j.StartDate != e.json.StartDate {
			t.Errorf("Result JSON for test %s does not match: expected %s, actual %s", e.name, e.json.StartDate, j.StartDate)
		}
		if j.EndDate != e.json.EndDate {
			t.Errorf("Result JSON for test %s does not match: expected %s, actual %s", e.name, e.json.EndDate, j.EndDate)
		}
		if j.Message != e.json.Message {
			t.Errorf("Result JSON for test %s does not match: expected %s, actual %s", e.name, e.json.Message, j.Message)
		}
		if j.OK != e.json.OK {
			t.Errorf("Result JSON for test %s does not match: expected %t, actual %t", e.name, e.json.OK, j.OK)
		}
		if j.RoomID != e.json.RoomID {
			t.Errorf("Result JSON for test %s does not match: expected %s, actual %s", e.name, e.json.RoomID, j.RoomID)
		}
	}

}

func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
