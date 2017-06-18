package main

import (
	"encoding/csv"
  "fmt"
  "strings"
  "io/ioutil"
  "io"
  "log"
)

type Version struct {
  Id []string
  Text []string
  Latin []string
  Greek []string
  Other []string
  Title []string
}

func main() {
  var output Version
  data, err := ioutil.ReadFile("data/pg.Vol.-1.coo.31924054872803_ocr.csv")
  if err != nil {
    fmt.Println(err)
  }
  str := string(data)
  reader := csv.NewReader(strings.NewReader(str))
  reader.Comma = '#'
  reader.LazyQuotes = true

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
    output.Id = append(output.Id, line[0])
    output.Text = append(output.Text, line[1])
    output.Latin = append(output.Latin, line[2])
    output.Greek = append(output.Greek, line[3])
    output.Other = append(output.Other, line[4])
    output.Title = append(output.Title, line[5])
	}
  fmt.Println(output.Title)
}
