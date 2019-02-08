package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func keyValue(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getByKey(w, r)
	case "POST":
		setByKey(w, r)
	case "DELETE":
		deleteByKey(w, r)
	}
}
func getByKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	vType, value, exist := dbRead(key)
	if exist {
		// fmt.Fprintf(w, "Key: "+value.(string))
		value := struct {
			Key   string
			Value interface{}
			Type  string
		}{
			key,
			value,
			vType,
		}
		json.NewEncoder(w).Encode(value)
	} else {
		http.Error(w, "Unknown key", http.StatusNotFound)
	}
}

func setByKey(w http.ResponseWriter, r *http.Request) {
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
	dbWrite(key, v)
	msg := fmt.Sprintf("key %s set!\n", key)
	w.Write([]byte(msg))
}

func deleteByKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var msg string
	_, e := myDB.M[key]
	if e {
		myDB.Lock()
		delete(myDB.M, key)
		myDB.need2sync = true
		myDB.Unlock()
		msg = fmt.Sprintf("key %s deleted!\n", key)
	} else {
		msg = fmt.Sprintf("No key %s to delete!\n", key)
	}
	w.Write([]byte(msg))
}

func incByKey(w http.ResponseWriter, r *http.Request) {
	var value int
	vars := mux.Vars(r)
	key := vars["key"]
	myDB.Lock()
	v, exist := myDB.M[key]
	if exist {
		switch fmt.Sprintf("%T", v) {
		case "int":
			value = v.(int)
		case "float32":
			value = int(v.(float32))
		case "float64":
			value = int(v.(float64))
		default:
			myDB.Unlock()
			msg := fmt.Sprintf("key '%s' is not an int or float.", key)
			http.Error(w, msg, http.StatusNotAcceptable)
			return
		}
	}
	if r.Method == "GET" || r.Method == "POST" {
		value++
	}
	if r.Method == "DELETE" {
		value--
	}
	myDB.need2sync = true
	myDB.M[key] = value
	myDB.Unlock()
	dbSync(false)
	w.Write([]byte("key inc!\n"))
}

func listKey(w http.ResponseWriter, r *http.Request) {
	var allKeys []string
	myDB.Lock()
	for k := range myDB.M {
		allKeys = append(allKeys, k)
	}
	myDB.Unlock()
	json.NewEncoder(w).Encode(allKeys)
}
