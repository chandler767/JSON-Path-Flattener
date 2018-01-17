package main

import (
	"os"
	"errors"
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"io/ioutil"
)

var (
	jsontoflatten string
	err error
)

func GetFromURL(path string) (string, error) { // Get JSON from URL and return body as string.
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

func OpenFromDisk(path string) (string, error) { // Open JSON file and return body as string.
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
		jsontoflatten = os.Args[1] // Expecting a path, URL, or string.
		if jsontoflatten != "" {
			if strings.HasPrefix(jsontoflatten, "https://") || strings.HasPrefix(jsontoflatten, "http://") { // Check for URL prefix.
				fmt.Println("Fetching JSON from URL.")
				jsontoflatten, err = GetFromURL(jsontoflatten)
				if err != nil {
					panic(err)
				}
			} else if strings.HasSuffix(jsontoflatten, ".json") { // If suffixed with '.json' then assumed to be a file on disk.
				fmt.Println("Fetching JSON from disk.")
				jsontoflatten, err = OpenFromDisk(jsontoflatten)
				if err != nil {
					panic(err)
				}
			}
			fmt.Println(jsontoflatten)
			return
	    }
	}
	fmt.Println(errors.New("A valid path, URL, or string is required."))
}
