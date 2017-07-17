package main

import (
    "fmt"
    "io/ioutil"
    "regexp"
    "strings"
)

func UniqStr(input []string) []string {
	u := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}

func main() {
    files, _ := ioutil.ReadDir("./data/")
    matchingnames := []string{}

    for _, f := range files {
      filename := f.Name()
      filenamesl := strings.Split(filename, ".")
      filename = strings.Join(filenamesl[:len(filenamesl)-3], ".") + "."
      matchingnames = append(matchingnames, filename)
      }
    unique := UniqStr(matchingnames)

    for i:= range unique {
    fmt.Println(unique[i])
    number := unique[i] + "[a-z]"
    fmt.Println(number)
    matchingnames = []string{}
    for _, f := range files {
      filename := f.Name()
    matched, _ := regexp.MatchString(number, filename)
    if  matched == true {
      filename = "data/" + filename
      matchingnames = append(matchingnames, filename)
    }
  }
  fmt.Println(matchingnames[0])
  fmt.Println(matchingnames[1:len(matchingnames)])
}
}
