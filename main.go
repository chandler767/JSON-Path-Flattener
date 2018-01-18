package main

import (
	"os"
	"errors"
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"io/ioutil"
	"reflect"
	"github.com/json-iterator/go"
)

var (
	err error
	nestedjson map[string]interface{}
	flattenedjson map[string]string = make(map[string]string)
)

func flatten(jsontoflatten map[string]string, prefix string, val reflect.Value) {
	if val.Kind() == reflect.Interface {
		val = val.Elem()
	}
	switch val.Kind() {
	case reflect.Bool:
		if val.Bool() {
			jsontoflatten[prefix] = "true"
		} else {
			jsontoflatten[prefix] = "false"
		}
	case reflect.Int:
		jsontoflatten[prefix] = fmt.Sprintf("%d", val.Int())
	case reflect.Float64:
		jsontoflatten[prefix] = strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.2f", val), "0"), ".") // Remove trailing zeros.
	case reflect.Map:
		for _, key := range val.MapKeys() {
			if key.Kind() == reflect.Interface {
				key = key.Elem()
			}
			if key.Kind() != reflect.String {
				panic(fmt.Sprintf("Key is not string: %s", key))
			}
			flatten(jsontoflatten, fmt.Sprintf("%s.%s", prefix, key.String()), val.MapIndex(key))
		}
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			flatten(jsontoflatten, fmt.Sprintf("%s[%d]", prefix, i), val.Index(i))
		}
	case reflect.String:
		jsontoflatten[prefix] = val.String()
	default:
		if !val.IsValid() {
			jsontoflatten[prefix] = "nil" // Insert a nil as string.
		} else {
			panic(fmt.Sprintf("Unexpected JSON value: %s", val)) // Check JSON format.
		}
	}
}

func getfromurl(path string) (string, error) { // Get JSON from URL.
	response, err := http.Get(path)
	if err != nil {
		return "", err
	} else {
		if response.StatusCode == http.StatusOK {
			defer response.Body.Close()
			jsonbody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return "", err
			}
			response.Body.Close()
			return string(jsonbody), nil
		} else {
			return "", errors.New("Unexpected status code.")
		}
	}
}

func readfromdisk(path string) (string, error) { // Open JSON file.
	JSONfile, err := os.Open(path)
	if err != nil{
		return "", err
	}
	defer JSONfile.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(JSONfile)
	JSONfile.Close()
	return buf.String(), nil
}

func main() {
	if len(os.Args) > 1 {
		jsontoflatten := os.Args[1] // Expecting a path, URL, or string.
		if jsontoflatten != "" {
			if strings.HasPrefix(jsontoflatten, "https://") || strings.HasPrefix(jsontoflatten, "http://") { // Get JSON from URL.
				jsontoflatten, err = getfromurl(jsontoflatten)
			} else if strings.HasSuffix(jsontoflatten, ".json") { // Reading JSON from disk.
				jsontoflatten, err = readfromdisk(jsontoflatten)
			}
			if err != nil {
				panic(err)
			}
			var json = jsoniter.ConfigCompatibleWithStandardLibrary
			err = json.Unmarshal([]byte(jsontoflatten), &nestedjson)
			if err != nil {
				fmt.Println(errors.New("Could not unmarshal JSON. A valid path, URL, or string is required."))
				panic(err)
			}
			for key, val := range nestedjson { // Flatten.
				flatten(flattenedjson, key, reflect.ValueOf(val))
			}
			for key, val := range flattenedjson { // Print.
				fmt.Println(fmt.Sprintf(".%s=%s", key, val))
			}
			return
	    }
	}
	fmt.Println(errors.New("A valid path, URL, or string is required."))
}
