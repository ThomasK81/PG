package main

import (
	"encoding/csv"
  "fmt"
  "strings"
  "io/ioutil"
  "io"
  "log"
  "github.com/aebruno/nwalgo"
)

type Version struct {
  Id []string
  Text []string
  Latin []string
  Greek []string
  Other []string
  Title []string
}

func maxfloat(floatslice []float64) int {
  max := floatslice[0]
  maxindex := 0
  for i, value := range floatslice {
    if value > max {
      max = value
      maxindex = i
    }
  }
  return maxindex
}

func loadCSV(s string) Version{
  var output Version
  data, err := ioutil.ReadFile(s)
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
  return output
}

func main() {
  output := loadCSV("data/pg.Vol.-1.hvd.32044054121090_ocr.csv")
  output2 := loadCSV("data/pg.Vol.-1.hvd.32044015466733_ocr.csv")
  var score_range []float64
	var text_range []string
	var id_range []string
  var index int
	for j:= range output.Text {
		score_range = []float64{}
		id_range = []string{}
		text_range = []string{}
		start := 0
		end := 0
		switch {
		case j - 200 < 0:
			start = 0
		default:
			start = j - 200
		}
		switch {
		case j + 200 >= len(output2.Text) - 1:
			end = len(output2.Text) - 1
		default:
			end = j + 200
		}
  for i:= start; i < end; i++ {
    aln1, _, score := nwalgo.Align(output.Text[j], output2.Text[i], 1, -1, -1)
    var f float64 = float64(score) / float64(len(aln1))
    score_range = append(score_range, f)
		id_range = append(id_range, output2.Id[i])
		text_range = append(text_range, output2.Text[i])
  }
  index = maxfloat(score_range)
	switch{
	case score_range[index] > 0.5:
		fmt.Println("----------------------------------")
		fmt.Println(output.Id[j])
	  fmt.Println(output.Text[j])
		fmt.Println(id_range[index])
	  fmt.Println(text_range[index])
		fmt.Println("Score:", score_range[index])
	default:
		fmt.Print(".")
	}
}
}
