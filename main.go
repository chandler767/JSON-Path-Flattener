package main

import (
	"os"
	"errors"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
)

var (
	jsontoflatten string
	err error
)

func GetFromURL(path string) (string, error) { // Download JSON file and return body as string.
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

func main() {
	if len(os.Args) > 1 {
		jsontoflatten = os.Args[1] // Expecting a path, URL, or string.
		if jsontoflatten != "" {
			if strings.HasPrefix(jsontoflatten, "https://") || strings.HasPrefix(jsontoflatten, "http://") { // Check for URL prefix
				fmt.Println("Fetching JSON file from URL.")
				jsontoflatten, err = GetFromURL(jsontoflatten)
				if err != nil {
					panic(err)
				}
			} else if strings.HasSuffix(jsontoflatten, ".json") { // If suffixed with '.json' then assumed to be a file on disk.
				fmt.Println("Fetching JSON file from disk.")



			}
			fmt.Println(jsontoflatten)
			return
	    }
	}
	fmt.Println(errors.New("A valid path, URL, or string is required."))
}
