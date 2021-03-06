package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/eador/bookings/internal/models"
)

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
	{"search availability", "/search-availability", "get", http.StatusOK},
	{"non-existant", "/bad-url", "get", http.StatusNotFound},
	{"login", "/user/login", "get", http.StatusOK},
	{"logout", "/user/logout", "get", http.StatusOK},
	{"dashboard", "/admin/dashboard", "get", http.StatusOK},
	{"new reservations", "/admin/reservations-new", "get", http.StatusOK},
	{"all reservations", "/admin/reservations-all", "get", http.StatusOK},
	{"show reservations", "/admin/reservations/new/1/show", "get", http.StatusOK},
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
	if rr.Code != http.StatusSeeOther {
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

	if rr.Code != http.StatusSeeOther {
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

	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "Smith")
	postedData.Add("phone", "12345")
	postedData.Add("email", "john@smith.com")

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
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
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for missing reservation in session
	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for error inserting reservation
	reservation.RoomID = 2
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for error inserting room restriction
	reservation.RoomID = 1000
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Test for invalid data
	postedData.Del("first_name")
	postedData.Add("first_name", "J")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))
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

func TestRepository_ReservationSummary(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reservation-summary", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	handler := http.HandlerFunc(Repo.ReservationSummary)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("error got status %d", rr.Code)
	}

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	req, _ = http.NewRequest("GET", "/reservation-summary", nil)
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handlers returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

var chooseRoomsTests = []struct {
	name string
	url  string
	code int
}{
	{"Available", "/choose-room/1", http.StatusSeeOther},
	{"Bad Params", "/choose-room/fish", http.StatusSeeOther},
}

func TestRepository_ChooseRoom(t *testing.T) {

	for _, i := range chooseRoomsTests {
		reservation := models.Reservation{
			RoomID: 1,
			Room: models.Room{
				ID:       1,
				RoomName: "General's Quarters",
			},
		}
		req, _ := http.NewRequest("GET", i.url, nil)
		ctx := GetCtx(req)
		session.Put(ctx, "reservation", reservation)
		req = req.WithContext(ctx)
		req.RequestURI = i.url

		handler := http.HandlerFunc(Repo.ChooseRoom)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != i.code {
			t.Errorf("Error in test %s expected status %d but got %d", i.name, i.code, rr.Code)
		}
	}

	// Test no reservation
	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := GetCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"

	handler := http.HandlerFunc(Repo.ChooseRoom)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Error in test No Reservation expected status %d but got %d", http.StatusSeeOther, rr.Code)
	}
}

var availabilityTests = []struct {
	name    string
	reqBody io.Reader
	code    int
	message string
}{
	{"Available", strings.NewReader("start=2050-01-01&end=2050-01-02"), http.StatusOK, " "},
	{"No Form", nil, http.StatusSeeOther, "can't parse form"},
	{"Bad Start", strings.NewReader("start=Bad&end=2050-01-02"), http.StatusSeeOther, "can't parse start date"},
	{"Bad End", strings.NewReader("start=2050-01-01&end=BAD"), http.StatusSeeOther, "can't parse end date"},
	{"DB Error", strings.NewReader("start=1050-01-01&end=1050-01-02"), http.StatusSeeOther, "can't access database"},
	{"No Rooms", strings.NewReader("start=1050-01-01&end=2050-01-02"), http.StatusSeeOther, "No Availability"},
}

func TestRepository_PostAvailability(t *testing.T) {
	for _, i := range availabilityTests {
		req, _ := http.NewRequest("POST", "/search-availability", i.reqBody)
		ctx := GetCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		handler := http.HandlerFunc(Repo.PostAvailability)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if rr.Code != i.code {
			t.Errorf("Error in test %s expected code %d got %d", i.name, i.code, rr.Code)
		}
		if session.Get(ctx, "error") != nil {
			if fmt.Sprintf("%s", session.Get(ctx, "error")) != i.message {
				t.Errorf("Error in test %s expected message %s got %s", i.name, i.message, session.Get(ctx, "error"))
			}
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

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectecLocation   string
}{
	{
		"valid credentials",
		"me@here.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid credentials",
		"jack@nible.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid data",
		"j",
		http.StatusOK,
		`action="/user/login`,
		"",
	},
}

func TestLogin(t *testing.T) {
	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := GetCtx(req)
		req = req.WithContext(ctx)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectecLocation != "" {
			actualLoc, _ := rr.Result().Location()
			if actualLoc.String() != e.expectecLocation {
				t.Errorf("failed %s: expected location %s but got %s", e.name, e.expectecLocation, actualLoc.String())
			}
		}

		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected html %s but got %s", e.name, e.expectedHTML, html)
			}
		}
	}
}
