package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Person struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

func StructToJSON(data any) (string, error) {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(dataJSON), nil
}

func JSONToStruct(jsonStr string, result any) error {
	return json.Unmarshal([]byte(jsonStr), result)
}

func WriteJSONFile(filename string, data any) error {
	dataJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, dataJSON, 0644)
}

func ReadJSONFile(filename string, result any) error {
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, result)
}

func main() {
	person := Person{Name: "Suryo", Age: 40, Email: "suryo@suryo.com"}

	jsonStr, err := StructToJSON(person)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	fmt.Println("Output:\n", jsonStr)

	var newPerson Person
	err = JSONToStruct(jsonStr, &newPerson)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	fmt.Println("Output :", newPerson)

	filename := "person.json"
	err = WriteJSONFile(filename, person)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	fmt.Println("File JSON :", filename)

	var filePerson Person
	err = ReadJSONFile(filename, &filePerson)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	fmt.Println("Output :", filePerson)
}
