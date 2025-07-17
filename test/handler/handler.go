package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/varun-muthanna/loadbalancer/test/data"
)

const Key string = "JSON"

func GetRouter(w http.ResponseWriter, r *http.Request) {

	DataList, err := data.GetData()

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in getting data %v", err)))
		return
	}

	bytes, err := json.Marshal(DataList)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in marshalling data %v", err)))
		return
	}

	w.Write(bytes)
}

func PutRouter(w http.ResponseWriter, r *http.Request) {

	bytes := r.Context().Value(Key).([]byte)

	err := data.PutData(bytes)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in PUT of data, %v", err)))
	}

}

func PostRouter(w http.ResponseWriter, r *http.Request) {
	bytes := r.Context().Value(Key).([]byte)

	err := data.PostData(bytes)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in POST of data, %v", err)))
	}

}

func DeleteRouter(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r) //map[string]string
	id := vars["id"]

	userId, err := strconv.Atoi(id)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in finding of ID, %v", err)))
	}

	err = data.DeleteData(userId)

	if err != nil {
		w.Write([]byte(fmt.Sprintf("error in DELETE of data, %v", err)))
	}

}

func MiddleWareHandler(next http.Handler) http.Handler { //incase of no mux we call MiddleWare(Handler) so next is passed and it returns to satisfy the Handler Interface as be passed in the Server field
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//LOGIC FOR DATA HANDLING IN CASE OF PUT/POST

		bytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			w.Write([]byte(fmt.Sprintf("Unable to read json, %v", err)))
			return
		}

		ctx := context.WithValue(context.Background(), Key, bytes)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
