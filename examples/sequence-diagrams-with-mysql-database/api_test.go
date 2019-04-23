package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/steinfletcher/apitest"
)

// This test requires a mysql database to run

func TestGetUser_With_Default_Report_Formatter(t *testing.T) {
	username := uuid.NewV4().String()[0:7]

	DBSetup(func(db *sqlx.DB) {
		q := "INSERT INTO users (username, is_contactable) VALUES (?, ?)"
		db.MustExec(q, username, true)
	})

	apiTest("gets the user").
		Debug().
		Mocks(getUserMock(username)).
		Get("/user").
		Query("name", username).
		Expect(t).
		Status(http.StatusOK).
		Header("Content-Type", "application/json").
		Body(fmt.Sprintf(`{"name": "%s", "is_contactable": true}`, username)).
		End()
}

func getUserMock(username string) *apitest.Mock {
	return apitest.NewMock().
		Get("http://users/api/user").
		Query("id", username).
		RespondWith().
		Body(fmt.Sprintf(`{"name": "%s", "id": "1234"}`, username)).
		Status(http.StatusOK).
		End()
}

func apiTest(name string) *apitest.APITest {
	recorder := apitest.NewTestRecorder()
	recordingDriver := apitest.WrapWithRecorder("mysql", recorder)
	sql.Register("mysqlWithRecorder", recordingDriver)

	testDB, err := sqlx.Connect("mysqlWithRecorder", dbAddr)
	if err != nil {
		panic(err)
	}

	app := newApp(testDB)

	return apitest.New(name).
		Recorder(recorder).
		Report(apitest.SequenceDiagram()).
		Handler(app.Router)
}
