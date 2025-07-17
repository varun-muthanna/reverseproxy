package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	"github.com/go-playground/validator"
)

var mu sync.RWMutex //cannot be const - initialized at runtime

const (
	FILE_PATH string = "data.txt"
) //cannot reassign set at compile time

type Data struct {
	ID           int    `json:"ID" validate:"required"`
	Name         string `json:"name" validate:"required"`
	Address      string `json:"add"` //unexported lowercase by json/encoding
	Phone_number int    `json:"Ph"`
}

func ValidateJSON(d *Data) error {
	val := validator.New()
	err := val.Struct(d)
	return err
}

func GetData() ([]*Data, error) {

	mu.RLock()
	defer mu.RUnlock()

	var d []*Data
	f, err := os.OpenFile(FILE_PATH, os.O_RDONLY, 0644)

	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No file found")
			return []*Data{}, nil
		}
		return nil, err
	}

	defer f.Close()

	byte, err := ioutil.ReadAll(f)

	if len(byte) == 0 {
		return []*Data{}, nil
	}

	if err != nil {
		fmt.Println("Error reading file")
		return nil, err
	}

	json.Unmarshal(byte, &d) //json to struct
	return d, nil
}

func PostData(b []byte )error {

	var d *Data = &Data{}
	json.Unmarshal(b,&d)

	mu.Lock()
	defer mu.Unlock()

	var dataList []*Data = []*Data{}
	var err error

	if err := ValidateJSON(d); err != nil {
		fmt.Println("Error validating json")
		return err
	}

	f, err := os.OpenFile(FILE_PATH, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		fmt.Println("Error opening file")
		return err
	}
	bytesE, err := ioutil.ReadAll(f)

	defer f.Close() //read only so can be deffered also  as os.Write is write

	if err != nil {
		fmt.Println("Error reading file")
		return err
	}

	if len(bytesE) > 0 {
		json.Unmarshal(bytesE, &dataList)
	}

	dataList = append(dataList, d)

	bytes, err := json.Marshal(dataList) //struct to JSON

	if err != nil {
		fmt.Println("Error marshalling json")
		return err
	}

	os.WriteFile(FILE_PATH, bytes, 0644)
	return nil
}

func PutData(b []byte) error {

	var d *Data = &Data{}

	json.Unmarshal(b,&d)

	dataList, err := GetData()
	var isExisting bool = false

	if err != nil {
		fmt.Println("Error getting data")
		return err
	}

	for i, _ := range dataList {
		if dataList[i].ID == d.ID {
			mu.Lock()
			isExisting = true
			dataList[i] = d
			mu.Unlock()
		}
	}

	if !isExisting {
		err := PostData(b)
		if err != nil {
			fmt.Println("Error posting new put data")
			return err
		}
	} else {
		json, err := json.Marshal(dataList)
		if err != nil {
			fmt.Println("Error in marshalling JSON of PUT request")
			return err
		}
		os.WriteFile(FILE_PATH, json, 0644)
	}

	return nil
}

func DeleteData(ID int) error {

	DataList, err := GetData()
	var DataListNew []*Data
	var isInValid = true

	if err != nil {
		fmt.Println("Error getting data")
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	for i, _ := range DataList {
		if DataList[i].ID == ID {
			isInValid = false
		} else {
			DataListNew = append(DataListNew, DataList[i])
		}
	}

	if isInValid {
		return errors.New("invalid ID, not present in file")
	}

	bytes, err := json.Marshal(DataListNew)

	if err != nil {
		fmt.Println("Error marshalling JSON")
		return err
	}

	os.WriteFile(FILE_PATH, bytes, 0644)

	return nil
}
