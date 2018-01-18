package main

import (
	"os"
	"errors"
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"reflect"
	"strconv"
)

var (
	err error
	nestedjson map[string]interface{}
	flattenedjson map[string]interface{} = make(map[string]interface{})
)

func flatten(jsontoflatten map[string]interface{}, key string, flattenedjson *map[string]interface{}) {
	for rkey, val := range jsontoflatten {
		fkey := key+rkey
		valref := reflect.ValueOf(val)
		if valref.Kind() == reflect.Interface {
			valref = valref.Elem() // Get the value
		}
		switch valref.Kind() {
			case reflect.Bool: 
				if valref.Bool() {
					(*flattenedjson)[fkey] = "true"
				} else {
					(*flattenedjson)[fkey] = "false"
				}
			case reflect.Int:
				(*flattenedjson)[fkey] = fmt.Sprintf("%d", val)
			case reflect.Float64:
				(*flattenedjson)[fkey] = fmt.Sprintf("%f", val)
			case reflect.String:
				(*flattenedjson)[fkey] = val.(string)
			case reflect.Slice:
				for i := 0; i<len(val.([]interface{})); i++ {
					if _, ok := val.([]string); ok {
						(*flattenedjson)[string(i)] = val.(string)
					} else if _, ok := val.([]int); ok {
						(*flattenedjson)[string(i)] = val.(int)
					} else {
						flatten(val.([]interface{})[i].(map[string]interface{}), rkey+"["+strconv.Itoa(i)+"].", flattenedjson)
					}
				}
			default:
				if !valref.IsValid() {
					(*flattenedjson)[fkey] = "nil"
				} else {
					panic(fmt.Sprintf("Unexpected JSON value: %s", val))
				}
		}
	}
}

func getfromurl(path string) (string, error) { // Get JSON from URL and return body as string.
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

func readfromdisk(path string) (string, error) { // Open JSON file and return body as string.
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
			err = json.Unmarshal([]byte(jsontoflatten), &nestedjson) 
			if err != nil {
				fmt.Println(errors.New("Could not unmarshal JSON. A valid path, URL, or string is required."))
				panic(err)
			}
			flatten(nestedjson, ".", &flattenedjson) // Flatten.
			for k, v := range flattenedjson { // Print.
				fmt.Println(fmt.Sprintf("%s=%s", k, v))
			}
			return
	    }
	}
	fmt.Println(errors.New("A valid path, URL, or string is required."))
}
