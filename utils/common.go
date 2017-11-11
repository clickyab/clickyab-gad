package utils

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// TranslateError create an error wit another string as content if the first error is exist
func TranslateError(err error, str string) error {
	if err != nil {
		return errors.New(str)
	}

	return err
}

// IsTopicMatched is utility function to match for topic in amqp way.
// * means exactly one part, # means any part
func IsTopicMatched(pattern, topic string) bool {
	p := strings.Replace(pattern, "*", "[^.]+", -1)
	p = strings.Replace(p, "#", ".*", -1)
	p = "^" + p + "$"

	re, err := regexp.Compile(p)
	if err != nil {
		logrus.Warn("%s is translated to %s which is error : %s", pattern, p, err.Error())
		return false
	}

	return re.MatchString(topic)
}

// StringToDate try to convert string to date in several format
func StringToDate(s string) (time.Time, error) {
	return parseDateWith(s, []string{
		time.RFC3339,
		"2006-01-02T15:04:05", // iso8601 without timezone
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.ANSIC,
		time.UnixDate,
		time.RubyDate,
		"2006-01-02 15:04:05Z07:00",
		"02 Jan 06 15:04 MST",
		"2006-01-02",
		"02 Jan 2006",
	})
}

func parseDateWith(s string, dates []string) (d time.Time, e error) {
	for _, dateType := range dates {
		if d, e = time.Parse(dateType, s); e == nil {
			return
		}
	}
	return d, fmt.Errorf("Unable to parse date: %s", s)
}

// ReverseString reverse a string
func ReverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// AppendIfMissing append a key if the key is not already in slice
func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// ID Channel is for a new unique string,
// Used mainly at generating payment token
var ID = make(chan string)

func init() {
	// Make sure random generator is a bit fair random :)
	rand.Seed(int64(time.Now().Nanosecond()))

	go func() {
		h := sha1.New()
		c := []byte(time.Now().String() + fmt.Sprint(rand.Float64()))
		for {
			_, _ = h.Write(c)
			ID <- fmt.Sprintf("%x", h.Sum(nil))
		}
	}()
}
