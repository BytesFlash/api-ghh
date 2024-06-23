package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
)

type ErrorType struct {
	Error string `json:"error"`
}

type PasswordValidation struct {
	MinLength    bool `json:"min_length"`
	Uppercase    bool `json:"uppercase"`
	Lowercase    bool `json:"lowercase"`
	SymbolNumber bool `json:"symbol_number"`
}

var logFile *os.File

func SetError(errorString string) *ErrorType {
	returnError := &ErrorType{
		Error: errorString,
	}
	return returnError
}

func Password(pass string) (bool, PasswordValidation) {
	var (
		minLengthCondition                 = true
		uppercase, lowercase, symbolNumber bool
		minLength                          uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			uppercase = true
			minLength++
		case unicode.IsLower(char):
			lowercase = true
			minLength++
		case unicode.IsNumber(char) || unicode.IsPunct(char) || unicode.IsSymbol(char):
			symbolNumber = true
			minLength++
		}
	}
	if minLength < 6 {
		minLengthCondition = false
	}
	data := PasswordValidation{
		MinLength:    minLengthCondition,
		Uppercase:    uppercase,
		Lowercase:    lowercase,
		SymbolNumber: symbolNumber,
	}
	if !uppercase || !lowercase || !symbolNumber || !minLengthCondition {
		return false, data
	}
	return true, data
}

func LogfileInit() {

	t := time.Now()
	fecha := fmt.Sprintf("%d-%02d-%02d.log",
		t.Year(), t.Month(), t.Day())
	if logFile != nil {
		if logFile.Name() != fecha {
			logFile.Close()
		}
	}

	if _, err := os.Stat("./logs/"); os.IsNotExist(err) {
		os.MkdirAll("./logs/", 0700)
	}

	logFile, err := os.OpenFile("./logs/"+fecha, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)

	if err != nil {
		logrus.Println("Failed to create the log file")
		logrus.Println(err)
		os.Exit(1)
	}

	log.SetOutput(logFile)
}

// Unique File Name
func GenerateUniqueFilename(filename string) string {
	extension := filepath.Ext(filename)
	filenameWithoutExt := filename[:len(filename)-len(extension)]
	uniqueFilename := fmt.Sprintf("%s-%s%s", filenameWithoutExt, generateRandomString(8), extension)
	return uniqueFilename
}

// Random string
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func PostRequest(url string, contentType string, data string) (bdyString string) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("apikey", "3BEGTrD88UluPRjSsBY4LaAU1qXI7AAB")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	bdyString = bytes.NewBuffer(body).String()
	return bdyString

}

func GetRequest(url string) (bdyString string) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("apikey", "3BEGTrD88UluPRjSsBY4LaAU1qXI7AAB")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	bdyString = bytes.NewBuffer(body).String()
	return bdyString

}
