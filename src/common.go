package main

import "os"

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
