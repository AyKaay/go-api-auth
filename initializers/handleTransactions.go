package initializers

import (
	"encoding/json"
	"io"
	"os"
	"user-auth-api/models"
)

func LoadTransactions() ([]models.Transaction, error) {
	file, err := os.Open("db.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var transactions []models.Transaction
	if err := json.Unmarshal(data, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func WriteTransactions(transactions []models.Transaction) error {
	data, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile("db.json", data, 0644) // 0644 to give permission to write to json
	if err != nil {
		return err
	}

	return nil
}
