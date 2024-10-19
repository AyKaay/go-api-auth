package initializers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"user-auth-api/models"
)

type DBData struct {
	Users        []models.User        `json:"users"`
	Transactions []models.Transaction `json:"transactions"`
}

// LoadUsers loads users from db.json
func LoadUsers() ([]models.User, error) {
	file, err := os.Open("db.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var dbData DBData
	if err := json.Unmarshal(b, &dbData); err != nil {
		return nil, err
	}

	return dbData.Users, nil
}

// WriteDBData writes the entire DBData structure back to db.json
func WriteDBData(dbData DBData) error {
	data, err := json.MarshalIndent(dbData, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile("db.json", data, 0644) // 0644 for file permissions
}

// AddUser adds a new user to the db.json file
func AddUser(newUser models.User) error {
	// Load existing users
	users, err := LoadUsers()
	if err != nil {
		return err
	}

	// Append the new user
	users = append(users, newUser)

	// Create the updated DBData structure
	dbData := DBData{
		Users:        users,
		Transactions: []models.Transaction{}, // If you need transactions, you can load them here too
	}

	// Write the updated DBData back to db.json
	return WriteDBData(dbData)
}
