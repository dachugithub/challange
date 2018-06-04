package main

import (
        "os"
        "errors"
        "fmt"
)

func main() {
    a := App{}
// initialize app using env vars
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
    app_service_port := os.Getenv("APP_SERVICE_PORT")
    if app_service_port == "" {
                                panic(errors.New("Env APP_SERVICE_PORT must be set"))
                        }
    a.Initialize(
        app_db_username,
        app_db_password,
        app_db_name,
        app_db_host)
    fmt.Println("Starting service")
    a.Run(app_service_port)
}
