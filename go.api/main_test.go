
package main

import (
    "os"
    "bytes"
    "testing"
    "time"
    "log"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"

)


var appUnderTest App

func TestMain(m *testing.M) {
    appUnderTest = App{}

    app_db_username := os.Getenv("APP_DB_USERNAME")
    if app_db_username == "" {
                                panic(errors.New("Env APP_DB_USERNAME must be set"))
                        }

    app_db_password := os.Getenv("APP_DB_PASSWORD")
    if app_db_password == "" {
                                panic(errors.New("Env APP_DB_PASSWORD must be set"))
                        }

    app_db_name := os.Getenv("APP_DB_NAME")
    if app_db_name == "" {
                                panic(errors.New("Env APP_DB_NAME must be set"))
                        }
    app_db_host := os.Getenv("APP_DB_HOST")
    if app_db_host == "" {
                                panic(errors.New("Env APP_DB_HOST must be set"))
                        }
    appUnderTest.Initialize(
        app_db_username,
        app_db_password,
        app_db_name,
        app_db_host)

    ensureTableExists()

    code := m.Run()

    clearTable()

    os.Exit(code)
}


func ensureTableExists() {
    if _, err := appUnderTest.DB.Exec(tableCreationQuery); err != nil {
        log.Fatal(err)
    }
}

func clearTable() {
    appUnderTest.DB.Exec("DELETE FROM people")
    appUnderTest.DB.Exec("ALTER SEQUENCE people_id_seq RESTART WITH 1")
}

func TestEmptyTable(t *testing.T) {
// outside the scope of the task
}

func TestGetNonExistentPeople(t *testing.T) {
    clearTable()

    req, _ := http.NewRequest("GET", "/hello/john", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusNotFound, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["error"] != "Person not found" {
        t.Errorf("Expected the 'error' key of the response to be set to 'Person not found'. Got '%s'", m["error"])
    }
}

func TestCreateFirstPerson(t *testing.T) {
    clearTable()

    payload := []byte(`{"date":"2000-01-01"}`)

    req, _ := http.NewRequest("POST", "/hello/john", bytes.NewBuffer(payload))
    response := executeRequest(req)

    checkResponseCode(t, http.StatusCreated, response.Code)

    var m map[string]interface{}
    json.Unmarshal(response.Body.Bytes(), &m)

    if m != nil {
        t.Errorf("Expected empty respons. Got '%v'", m)
    }

}

func TestGetPersonBirthdayInFiveDays(t *testing.T) {
    clearTable()
    addPersonBirthday("Alex",time.Now().AddDate(0, 0, +5))

    req, _ := http.NewRequest("GET", "/hello/Alex", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["message"] != "Hello, Alex! Your birthday is in 5 days" {
        t.Errorf("Expected the 'message' key of the response to be set to 'Hello, Alex! Your birthday is in 5 days'. Got '%s'", m["message"])
    }

}

func TestGetPersonBirthdayToday(t *testing.T) {
    clearTable()
    addPersonBirthday("Adam", time.Now())

    req, _ := http.NewRequest("GET", "/hello/Adam", nil)
    response := executeRequest(req)

    checkResponseCode(t, http.StatusOK, response.Code)

    var m map[string]string
    json.Unmarshal(response.Body.Bytes(), &m)
    if m["message"] != "Hello, Adam! Happy birthday!" {
        t.Errorf("Expected the 'message' key of the response to be set to 'Hello, Adam! Happy birthday!'. Got '%s'", m["message"])
    }

}

func TestUpdateBirthday(t *testing.T) {
// not requiere in the task, placeholder for future
}

func TestDeleteBirthday(t *testing.T) {
// not requiere in the task, placeholder for future
}

func TestNotUniqueNamesBirthday(t *testing.T) {
// not requiere in the task, placeholder for future
// part of problem is mittigated by adding index to the db
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()
    appUnderTest.Router.ServeHTTP(rr, req)

    return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
    if expected != actual {
        t.Errorf("Expected response code %d. Got %d\n", expected, actual)
    }
}

func addPersonBirthday(name string, birthday time.Time) {
    //db.Exec("INSERT INTO tablename VALUES (time_column) (($1));", time.Now())
    appUnderTest.DB.Exec(
        "INSERT INTO people(name, birthday) VALUES($1, $2)",
        name, birthday)

    return
}


const tableCreationQuery = `CREATE TABLE IF NOT EXISTS people
(
id SERIAL,
name TEXT NOT NULL,
birthday DATE  NOT NULL,
CONSTRAINT people_pkey PRIMARY KEY (id)
)`





