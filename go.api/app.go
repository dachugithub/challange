
package main

import (
  "database/sql"
  "fmt"
  "time"
  "log"
  "encoding/json"
  "net/http"

  "github.com/gorilla/mux"
  _ "github.com/lib/pq"


)

// basic app struct

type App struct {
    Router *mux.Router
    DB     *sql.DB
}


// Initialize the application 

func (a *App) Initialize(user, password, dbname, host string) {
    connectionString :=
        fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable", user, password, dbname, host)

    var err error
    a.DB, err = sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }

    a.Router = mux.NewRouter()
    a.initializeRoutes()
 }

// Start the app

func (a *App) Run(addr string) {
    fmt.Printf("Starting starting to listen on the port %s", addr)
    log.Fatal(http.ListenAndServe(addr, a.Router))
 }


// get person birthday

func (a *App) getBirthday(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    name := vars["name"]
    fmt.Println("someone asked for", name, "'s birthday")
    // it would be good to make username validation here, regexp ^[a-zA-Z]+$ is very basic. Task did not sepcified scope of names though (starting with capital letters?).
    p := person{Name: name}
    if err := p.checkPersonBirthday(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Person not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    var m = map[string]string{}
    m["message"] =  p.Name + "'s birthday is on " + p.Birthday.Format("January-02")
    if(p.Birthday.Format("01-02") == time.Now().AddDate(0, 0, +5).Format("01-02")){
      m["message"] = "Hello, " + p.Name + "! Your birthday is in 5 days"
    }

    if(p.Birthday.Format("01-02") == time.Now().Format("01-02")){
      m["message"] = "Hello, " + p.Name + "! Happy birthday!"
    }

    respondWithJSON(w, http.StatusOK, m)
}

func (a *App) postBirthday(w http.ResponseWriter, r *http.Request) {
    var p person
    var m jsonPayload

    vars := mux.Vars(r)
    name := vars["name"]
    // it would be good to make username validation here, regexp ^[a-zA-Z]+$ is very basic. Task did not sepcified scope of names though (starting with capital letters?).
    p.Name = name
    decoder := json.NewDecoder(r.Body)

    if err := decoder.Decode(&m); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    fmt.Println("someone want's to update", name, "'s birthday to ", m.Birthday)
    // final values setup

    p.Name = name
    birthdayParsed, _ := time.Parse("2006-01-02", m.Birthday)
    p.Birthday  = birthdayParsed

    if err := p.checkPersonBirthday(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            fmt.Println( name, "'s birthday does not exist creating with date", m.Birthday)
            if err := p.createPersonBirthday(a.DB); err != nil {
                    respondWithError(w, http.StatusInternalServerError, "Error creating persons birthday")
            return
            }
            respondWithJSON(w, http.StatusCreated, "")
            return
        }
    }

    // reseting birthday to the parsed date
    p.Birthday  = birthdayParsed
    fmt.Println( p.Name, "'s birthday already exists updating to  ", p.Birthday)

    if err := p.updatePersonBirthday(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    // in description it was 201 No content, no content is 204.
    respondWithJSON(w, http.StatusCreated, "")
}
// function that will be used in the deployment
func (a *App) healthcheck(w http.ResponseWriter, r *http.Request) {
    respondWithJSON(w, http.StatusOK, "")
}


func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

// routes 

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/hello/{name:[a-zA-Z]+}", a.getBirthday).Methods("GET")
    a.Router.HandleFunc("/hello/{name:[a-zA-Z]+}", a.postBirthday).Methods("POST")
    a.Router.HandleFunc("/healthcheck", a.healthcheck).Methods("GET")
}


