package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func pushByKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	decoder := json.NewDecoder(r.Body)
	var dat map[string]interface{}
	if err := decoder.Decode(&dat); err != nil {
		http.Error(w, "post body malformed", http.StatusNotAcceptable)
		return
	}
	v, e := dat[key]
	if !e {
		http.Error(w, "post body does not have key", http.StatusNotAcceptable)
		return
	}
	myDB.Lock()
	_, exist := myDB.Q[key]
	if exist {
		myDB.Q[key] = append(myDB.Q[key].([]interface{}), v)
	} else {
		myDB.Q[key] = append([]interface{}{}, v)
	}
	myDB.need2sync = true
	myDB.Unlock()
	dbSync(false)
}

func popByKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	myDB.Lock()
	_, e := myDB.Q[key]
	if !e {
		myDB.Unlock()
		http.Error(w, "Queue is empty", http.StatusNotAcceptable)
		return
	}
	len := len(myDB.Q[key].([]interface{}))
	if len == 0 {
		myDB.Unlock()
		http.Error(w, "Queue is empty", http.StatusNotAcceptable)
		return
	}
	value := myDB.Q[key].([]interface{})[len-1]
	vType := fmt.Sprintf("%T", value)
	// pop
	myDB.Q[key] = myDB.Q[key].([]interface{})[:len-1]
	valueStruct := struct {
		Key   string
		Value interface{}
		Type  string
	}{
		key,
		value,
		vType,
	}
	myDB.need2sync = true
	myDB.Unlock()
	json.NewEncoder(w).Encode(valueStruct)
}

func lenQ(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	myDB.Lock()
	_, e := myDB.Q[key]
	if !e {
		myDB.Unlock()
		http.Error(w, "Queue is empty", http.StatusNotAcceptable)
		return
	}
	len := len(myDB.Q[key].([]interface{}))
	myDB.Unlock()
	json.NewEncoder(w).Encode(len)
}

func listQueues(w http.ResponseWriter, r *http.Request) {
	var allKeys []string
	myDB.Lock()
	for k := range myDB.Q {
		allKeys = append(allKeys, k)
	}
	myDB.Unlock()
	json.NewEncoder(w).Encode(allKeys)
}
