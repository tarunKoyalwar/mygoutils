package io

import (
	"fmt"
	"io/ioutil"
)

//Read File data and handle error
func ReadFileData(filename string) []byte {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Problem while Reading file : ", filename)
		panic(err)
	}
	return dat
}

// A simple function to handle error with msg
func HandleError(er error, msg string) {
	if er != nil {
		fmt.Println(msg)
		panic(er)
	}
}
