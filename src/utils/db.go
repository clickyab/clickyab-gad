package utils

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// unquoted array values must not contain: (" , \ { } whitespace NULL)
	// and must be at least one char
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	// quoted array values are surrounded by double quotes, can be any
	// character except " or \, which must be backslash escaped:
	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	// an array value may be either quoted or unquoted:
	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	// Array values are separated with a comma IF there is more than one value:
	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

//StringSlice is simple slice to handle array in postgresql
type StringSlice []string

// PARSING ARRAYS
// SEE http://www.postgresql.org/docs/9.4/static/arrays.html#ARRAYS-IO
// Arrays are output within {} and a delimiter, which is a comma for most
// postgres types (; for box)
//
// Individual values are surrounded by quotes:
// The array output routine will put double quotes around element values if
// they are empty strings, contain curly braces, delimiter characters,
// double quotes, backslashes, or white space, or match the word NULL.
// Double quotes and backslashes embedded in element values will be
// backslash-escaped. For numeric data types it is safe to assume that double
// quotes will never appear, but for textual data types one should be prepared
// to cope with either the presence or absence of quotes.
// Parse the output string from the array type.
// Regex used: (((?P<value>(([^",\\{}\s(NULL)])+|"([^"\\]|\\"|\\\\)*")))(,)?)
func parseArray(array string) []string {
	var results []string
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		results = append(results, s)
	}
	return results
}

// Scan Implements sql.Scanner for the String slice type
// Scanners take the database value (in this case as a byte slice)
// and sets the value of the type.  Here we cast to a string and
// do a regexp based parse
func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return error(errors.New("Scan source was not []bytes"))
	}

	asString := string(asBytes)
	parsed := parseArray(asString)
	(*s) = StringSlice(parsed)

	return nil
}

// Value A very experimental string slice to postgres array. BUG use it with care, its not tested so much
func (s StringSlice) Value() (driver.Value, error) {
	b, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	// the first character is [ and the last is ] so change them to {}
	res := fmt.Sprintf("{%s}", string(b[1:len(b)-1]))
	return []byte(res), nil
}

// BuildPgPlaceHolder try to build a place holder and arguments for pq library (postgres)
func BuildPgPlaceHolder(start int, params ...interface{}) ([]string, []interface{}) {
	var res []string
	var p []interface{}

	for i := range params {
		res = append(res, fmt.Sprintf("$%d", start+i))
		p = append(p, params[i])
	}

	return res, p
}

// EscapeText is my attempt to implement the escape-string for postgres query for some
// trick on insert/duplicate
//func EscapeText(text string) string {
//	var buf []byte
//	escapeNeeded := false
//	startPos := 0
//	var c byte
//
//	// check if we need to escape
//	for i := 0; i < len(text); i++ {
//		c = text[i]
//		if c == '\\' || c == '\n' || c == '\r' || c == '\t' || c == '\'' {
//			escapeNeeded = true
//			startPos = i
//			break
//		}
//	}
//	if !escapeNeeded {
//		return string(append(buf, text...))
//	}
//
//	// copy till first char to escape, iterate the rest
//	result := append(buf, text[:startPos]...)
//	for i := startPos; i < len(text); i++ {
//		c = text[i]
//		switch c {
//		case '\'':
//			result = append(result, '\\', '\'')
//		case '\\':
//			result = append(result, '\\', '\\')
//		case '\n':
//			result = append(result, '\\', 'n')
//		case '\r':
//			result = append(result, '\\', 'r')
//		case '\t':
//			result = append(result, '\\', 't')
//		default:
//			result = append(result, c)
//		}
//	}
//	return string(result)
//}
