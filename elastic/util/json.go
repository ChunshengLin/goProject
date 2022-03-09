package util

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"unicode"
)


type JsonSnakeCase struct {
	Value interface{}
}

// MarshalJSON
func (j JsonSnakeCase) MarshalJSON() ([]byte, error) {
 	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	var wordBarrierRegex = regexp.MustCompile(`(\w)([A-Z])`)

	marshalled, err := json.Marshal(j.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled, 
		func(match []byte) []byte {
			return bytes.ToLower(wordBarrierRegex.ReplaceAll(
				match, 
				[]byte(`${1}_${2}`),
			))
		},
	)

	return converted, err
}

type JsonCamelCase struct {
	Value interface{}
}

func (c JsonCamelCase) MarshalJSON() ([]byte, error) {
	var keyMatchRegex = regexp.MustCompile(`\"(\w+)\":`)
	marshalled, err := json.Marshal(c.Value)
	converted := keyMatchRegex.ReplaceAllFunc(
		marshalled,
		func(match []byte) []byte {
			matchStr := string(match)
			key := matchStr[1 : len(matchStr)-2]
			resKey := Lcfirst(Case2Camel(key))
			return []byte(`"` + resKey + `":`)
		},
	)

 	return converted, err
}

func SnakeData(data interface{}) map[string]interface{} {
	snakedData, _ := json.Marshal(JsonSnakeCase{Value: data})

	var res map[string]interface{}
	_ = json.Unmarshal(snakedData, &res)

	return res
}

func CamelData(data interface{}) []byte {
	res, _ := json.Marshal(JsonCamelCase{Value: data})

	return res
}


func Case2Camel(name string) string {
	name = strings.Replace(name, "_", " ", -1)
	name = strings.Title(name)

	return strings.Replace(name, " ", "", -1)
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}

	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}