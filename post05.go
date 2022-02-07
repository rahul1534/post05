package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Decsription string
}

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan:", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

func AddUser(d Userdata) (int, error) {
	d.Username = strings.ToLower(d.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1, err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists", Username)
		return -1, err
	}

	insertStatement := `INSERT INTO "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID, err
	}

	insertStatement = `INSERT INTO "userdata" ("userid", "name", "surname", "description") values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Decsription)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1, err
	}
	return userID, nil
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT "username" from "users" WHERE id = %d`, id)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exists", id)
	}

	deleteStatement := `DELETE FROM "userdata" WHERE userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id", "username", "name", "surname", "description"
							FROM "users", "userdata"
							WHERE users.id = userdata.userid`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id int
		var username, name, surname, description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return Data, err
		}
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Decsription: description}
		Data = append(Data, temp)
	}
	defer rows.Close()
	return Data, nil
}

func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("user does not exist")
	}

	d.ID = userID
	updateStatement := `UPDATE "userdata" SET "name"=$1, "surname"=$2, "description"=$3 WHERE "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Decsription, d.ID)
	if err != nil {
		return err
	}

	return nil
}
