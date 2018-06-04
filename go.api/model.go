
package main

import (
        "time"
        "database/sql"
        "fmt"
)

type person struct {
    ID    int     `json:"id"`
    Name  string  `json:"name"`
    Birthday time.Time `json:"dateOfBirth"`
}

type jsonPayload struct {
        Birthday string `json:"dateOfBirth"`
    }

func (p *person) createPersonBirthday(db *sql.DB) error {
  err := db.QueryRow(
        "INSERT INTO people(name, birthday) VALUES($1, $2) RETURNING id",
        p.Name, p.Birthday).Scan(&p.ID)

    if err != nil {
        return err
    }

    return nil

}

func (p *person) updatePersonBirthday(db *sql.DB) error {
  _, err :=
        db.Exec("UPDATE people SET name=$1, birthday=$2 WHERE name=$1",
            p.Name, p.Birthday)

    return err
}

func (p *person) checkPersonBirthday(db *sql.DB) error {
  fmt.Println("looking for", p.Name ,"'s birthday")
  return db.QueryRow("SELECT birthday FROM people WHERE name=$1",
        p.Name).Scan(&p.Birthday)
}


// helper for testing

func getPeople(db *sql.DB, start, count int) ([]person, error) {
  rows, err := db.Query(
        "SELECT id, name, birthday FROM people LIMIT $1 OFFSET $2",
        count, start)

    if err != nil {
        return nil, err
    }

    defer rows.Close()

    people := []person{}

    for rows.Next() {
        var p person
        if err := rows.Scan(&p.ID, &p.Name, &p.Birthday); err != nil {
            return nil, err
        }
        people = append(people, p)
    }

    return people, nil
}
