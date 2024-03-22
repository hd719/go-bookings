package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hd719/go-bookings/internal/models"
)

var theTests = []struct {
	name   string // name for the test
	url    string
	method string

	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"c", "/contact", "GET", http.StatusOK},
	// {"post-search-avail", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"post-search-avail-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-01-01"},
	// 	{key: "end", value: "2020-01-02"},
	// }, http.StatusOK},
	// {"mr-post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "john"},
	// 	{key: "last_name", value: "smith"},
	// 	{key: "email", value: "a@a.com"},
	// 	{key: "phone", value: "1111111111"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := GetRoutes()

	// Test server
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		// 2 type of tests GET and POST

		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("For %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Generals Quarters",
		},
	}

	req, err := http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		fmt.Println("Error in test TestRepository_Reservation")
		fmt.Println("Check to see if there is a reservation in the session")
		return
	}
	ctx := GetCtx(req)
	req = req.WithContext(ctx)

	// Simulating what we get back from a req/res life cycle when someone fires up a web browser
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything), reintialized req
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		fmt.Println("Error in test TestRepository_Reservation")
		fmt.Println("Check to see if there is a reservation in the session")
		return
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case where room is not in session (reset everything), reintialized req
	req, err = http.NewRequest("GET", "/make-reservation", nil)
	if err != nil {
		fmt.Println("Error in test TestRepository_Reservation")
		fmt.Println("Check to see if there is a reservation in the session")
		return
	}
	ctx = GetCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

// Create a fake session and adding it to the context
func GetCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}
