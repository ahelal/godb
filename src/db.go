package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sync"
	"time"
)

var ticker *time.Ticker
var myDB DB

type DB struct {
	sync.Mutex
	need2sync bool
	M         map[string]interface{}
	Q         map[string]interface{}
}

func dbRead(key string) (string, interface{}, bool) {
	myDB.Lock()
	value, exist := myDB.M[key]
	myDB.Unlock()
	if !exist {
		return "", "", false
	}
	vType := fmt.Sprintf("%T", value)
	return vType, value, true
}

func dbWrite(key string, value interface{}) {
	myDB.Lock()
	myDB.M[key] = value
	myDB.need2sync = true
	myDB.Unlock()
	dbSync(false)
}

func dbSync(CalledByTicker bool) {
	if !CalledByTicker && Config.DB.Sync != 0 {
		return
	}
	myDB.Lock()
	if !myDB.need2sync {
		myDB.Unlock()
		return
	}
	jsonDB, err := json.Marshal(myDB)
	panicOnError(err)
	myDB.Unlock()
	jsonEnc, err := doEncrypt(Config.DB.EncKey, jsonDB)
	panicOnError(err)
	err = ioutil.WriteFile(Config.DB.File, jsonEnc, 0755)
	panicOnError(err)
	fmt.Println("Synced")
	myDB.need2sync = false
}

func dbTicker() {
	fmt.Println("Setup ticker every", Config.DB.Sync)
	mal := time.Duration(1000 * Config.DB.Sync)
	ticker = time.NewTicker(time.Millisecond * mal)
	go func() {
		for range ticker.C {
			dbSync(true)
		}
		fmt.Println("Ticker stopped")
	}()
}

func initDB() {
	myDB.M = make(map[string]interface{})
	myDB.Q = make(map[string]interface{})
	if Config.DB.Sync == -1 {
		fmt.Println("DB init in memory use only")
		return
	}
	fExist, err := exists(Config.DB.File)
	panicOnError(err)
	if fExist {
		content, err := ioutil.ReadFile(Config.DB.File)
		panicOnError(err)
		decryptedContent, err := doDecrypt(Config.DB.EncKey, content)
		panicOnError(err)
		err = json.Unmarshal(decryptedContent, &myDB)
		panicOnError(err)
		fmt.Println("Loading DB from file")
	} else {
		fmt.Println("Init a new DB from file")
	}
	if Config.DB.Sync > 0 {
		dbTicker()
	}
}
