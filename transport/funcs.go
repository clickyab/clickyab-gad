package transport

import "time"

//KeyGenDaily function  generate a redis key
func KeyGenDaily(prefix, value string) string {
	date := time.Now().Format("060102")
	return prefix + DELIMITER + value + DELIMITER + date
}
