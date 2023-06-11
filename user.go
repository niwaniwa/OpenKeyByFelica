package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	rootPath string = "./data/"
)

type User struct {
	ID          string `json:"id"`
	IDM         []byte `json:"idm"`
	Name        string `json:"name"`
	LastLogging string `json:"last_logging"`
	StNum       string `json:"st_num"`
}

var (
	debugUserData = []User{
		{ID: "1", IDM: []byte{1, 16, 3, 16, 197, 20, 106, 38}, Name: "hogepiyo", LastLogging: time.Now().String(), StNum: "2211101"},
	}
	userData []User
)

func ReadUserData() bool {
	if Debug {
		userData = debugUserData
		return false
	}

	var raw []User

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		file, err := os.ReadFile(path)

		fmt.Printf("path: %#v\n", path)

		if err != nil {
			return nil
		}

		var usr User

		if err := json.Unmarshal(file, &usr); err != nil {
			fmt.Println(err)
			return nil
		}

		raw = append(raw, usr)

		return nil
	})

	if err != nil {
		panic(err)
	}

	userData = raw

	return true

}

func RegisterUser(name string, stNum string) {

}

func SaveUserDatum(array []User) {
	for _, datum := range array {
		file, _ := os.Create(rootPath + datum.ID + ".json")
		defer file.Close()
		_ = json.NewEncoder(file).Encode(datum)
	}
}

func SaveUserData(data User) {
	file, _ := os.Create(rootPath + data.ID + ".json")
	defer file.Close()
	_ = json.NewEncoder(file).Encode(data)
}

func Contains(s []User, idm []byte) bool {
	for _, a := range s {
		if bytes.Equal(a.IDM, idm) {
			return true
		}
	}
	return false
}