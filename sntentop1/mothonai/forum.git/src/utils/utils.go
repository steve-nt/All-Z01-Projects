package utils

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LogDebug(v any) {
	log.Printf("Debug: %#v", v)
}

func LogInfo(s string) {
	log.Printf("Info: %s", s)
}

func ConvertStringToTime(timeString string) (time.Time, error) {
	timestamp, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return ConvertInt64ToTime(timestamp), nil
}

func ConvertTimeToString(t time.Time) string {
	return t.String()
}

func ConvertInt64ToTime(i int64) time.Time {
	return time.Unix(i, 0)
}

func GetCurrentTimestamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func GetFunctionName() error {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return errors.New("")
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return errors.New("")
	}
	return errors.New(f.Name())
}

// Returns the hash (string) from password or error
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func StringToInt64(str string) (int64, error) {
	var err error
	var ok bool
	var n int64
	ok, err = regexp.MatchString(`^\d+$`, str)
	if !ok {
		return n, fmt.Errorf("Regex mismatch")
	}
	if err != nil {
		return n, err
	}
	n, err = strconv.ParseInt(str, 10, 64)
	if err != nil {
		return n, err
	}
	return n, err
}
