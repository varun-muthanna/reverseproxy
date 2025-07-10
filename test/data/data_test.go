package data //test unexported functions

import (
	"os"
	"reflect"
	"testing"
)

// 2 cases of POST with valid and invalid JSON
func TestCreateFile(t *testing.T) {

	d1 := Data{ID: 1, Name: "Varun", Address: "Kotoor village"} //valid JSON

	d2 := Data{Address: "Kotoor village"} //invalid JSON

	err1 := d1.PostData()
	err2 := d2.PostData()

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
	var DataList, DataListT []*Data
	var err error

	_, err = GetData()

	if err != nil {
		t.Errorf("Error reading empty file, %v", err)
	}

	d1 := &Data{ID:1,Name: "Varun", Address: "Kotoor village", Phone_number: 122}
	d1.PostData()
	DataListT = append(DataListT, d1)

	d2 := &Data{ID:2,Name: "Mayank", Address: "Kotoor village", Phone_number: 122}
	d2.PostData()
	DataListT = append(DataListT, d2)
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

	d1 := &Data{
		ID:   1,
		Name: "Varun",
	}

	d1.PutData()
	defer os.Remove(FILE_PATH)

	dataList, err := GetData()

	if err != nil {
		t.Errorf("Error getting data")
	}

	if !reflect.DeepEqual(d1, dataList[0]) {
		t.Errorf("Data not matching in new PUT")
	}

	d1 = &Data{
		ID:   1,
		Name: "Mayank",
	}

	d1.PutData()
	dataList, err = GetData()

	if err != nil {
		t.Errorf("Error getting data")
	}

	if !reflect.DeepEqual(d1, dataList[0]) {
		t.Errorf("Data not matching in existing PUT")
	}

}

// 2 cases of DELETE with valid ID and invalid ID
func TestDeleteFile(t *testing.T) {

	d1 := &Data{ID: 1, Name: "Varun", Address: "Kotoor village"}

	d2 := &Data{ID: 2, Name: "Mayank", Address: "Kotoor village"}

	err1 := d1.PostData()
	err2 := d2.PostData()

	defer os.Remove(FILE_PATH)

	if err1 != nil || err2 != nil {
		t.Errorf("error posting data")
	}

	DeleteData(1)

	DataList, err := GetData()

	if err != nil {
		t.Errorf("error getting data")
	}

	if !reflect.DeepEqual(d2, DataList[0]) {
		t.Errorf("error in deletion")
	}

}
