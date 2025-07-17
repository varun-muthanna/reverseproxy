package data //test unexported functions

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

// 2 cases of POST with valid and invalid JSON
func TestCreateFile(t *testing.T) {

	d1 := []byte(`{"ID": 1, "name": "Varun", "add": "Kotoor village"}`) //valid JSON

	d2 := []byte(`{"add": "Kotoor village"}`) //invalid JSON

	err1 := PostData(d1)
	err2 := PostData(d2)

	defer os.Remove(FILE_PATH)

	if err1 != nil {
		t.Errorf("Error posting valid data - valid data not accepted, %v", err1)
	}

	if err2 == nil {
		t.Errorf("Error posting invalid data - invalid data accepted, %v", err2)
	}
}

// 2 cases of GET with no file and one with file
func TestGetFile(t *testing.T) {
	var Data1, Data2 *Data
	var DataList, DataListT []*Data
	var err error

	_, err = GetData()

	if err != nil {
		t.Errorf("Error reading empty file, %v", err)
	}

	d1 := []byte(`{"ID":1, "name": "Varun", "add": "Kotoor village", "Ph": 122}`)
	PostData(d1)

	json.Unmarshal(d1,&Data1)
	DataListT = append(DataListT, Data1)

	d2 := []byte(`{"ID":2,"name": "Mayank", "add": "Kotoor village", "Ph": 122}`)
	PostData(d2)

	json.Unmarshal(d2,&Data2)
	DataListT = append(DataListT, Data2)
	defer os.Remove(FILE_PATH)

	DataList, err = GetData()

	if err != nil {
		t.Errorf("Error reading file, %v", err)
	}

	if !reflect.DeepEqual(DataListT, DataList) {
		t.Errorf("Error  in file content, %v", err)
	}
}

// 2 cases of PUT with old ID and new ID
func TestEditFile(t *testing.T) {
	var Data1, Data2 *Data
	d1 := []byte(`{"ID":1,"name": "Varun"}`)

	PutData(d1)
	defer os.Remove(FILE_PATH)

	dataList, err := GetData()

	json.Unmarshal(d1,&Data1)

	if err != nil {
		t.Errorf("Error getting data")
	}

	if !reflect.DeepEqual(Data1, dataList[0]) {
		t.Errorf("Data not matching in new PUT")
	}

	d2 := []byte(`{"ID":1,"name": "Mayank"}`)

	PutData(d2)
	dataList, err = GetData()

	json.Unmarshal(d2,&Data2)

	if err != nil {
		t.Errorf("Error getting data")
	}

	if !reflect.DeepEqual(Data2, dataList[0]) {
		t.Errorf("Data not matching in existing PUT")
	}

}

// 2 cases of DELETE with valid ID and invalid ID
func TestDeleteFile(t *testing.T) {
	var  Data2 *Data
	d1 := []byte(`{"ID": 1, "name": "Varun", "add": "Kotoor village"}`)

	d2 := []byte(`{"ID": 2, "name": "Mayank","add": "Kotoor village"}`)
	json.Unmarshal(d2,&Data2)

	err1 := PostData(d1)
	err2 := PostData(d2)

	defer os.Remove(FILE_PATH)

	if err1 != nil || err2 != nil {
		t.Errorf("error posting data")
	}

	DeleteData(1)

	DataList, err := GetData()

	if err != nil {
		t.Errorf("error getting data")
	}

	if !reflect.DeepEqual(Data2, DataList[0]) {
		t.Errorf("error in deletion")
	}

}
