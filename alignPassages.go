package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"io/ioutil"
	"regexp"
)

type Textversion struct {
	ID              []string
	Version, Volume []string
	Page, Line      []int
	Text            []string
	Latin           []int
	Greek           []int
	Other           []int
	Header          []bool
}

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

func Parse(text string, n int) map[string]int {
	chars := []rune(strings.Repeat(" ", n))
	table := make(map[string]int)

	for _, letter := range strings.Join(strings.Fields(text), " ") + " " {
		chars = append(chars[1:], letter)

		ngram := string(chars)
		if _, ok := table[ngram]; ok {
			table[ngram]++
		} else {
			table[ngram] = 1
		}
	}

	return table
}

func readPG(file string) Textversion {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println("could not open file")
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = '#'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = 6

	var textversion Textversion
	version, volume := "", ""
	page, vline := 0, 0

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		textversion.ID = append(textversion.ID, line[0])
		version = strings.Split(line[0], ":")[0]
		version = strings.Split(version, "_ocr")[0]
		version_string := strings.Split(version, ".")
		version = strings.Join([]string{version_string[len(version_string)-2], version_string[len(version_string)-1]}, ".")
		textversion.Version = append(textversion.Version, version)
		volume = strings.Split(line[0], ":")[0]
		volume = strings.Split(volume, "Vol.-")[1]
		volume_string := strings.Split(volume, ".")
		volume_string = volume_string[:len(volume_string)-2]
		volume = strings.Join(volume_string, ".")
		textversion.Volume = append(textversion.Volume, volume)
		page_string := strings.Split(line[0], ":")[1]
		page_string = strings.Split(page_string, ".")[0]
		page, _ = strconv.Atoi(page_string)
		textversion.Page = append(textversion.Page, page)
		line_string := strings.Split(line[0], ":")[1]
		line_string = strings.Split(line_string, ".")[1]
		vline, _ = strconv.Atoi(line_string)
		textversion.Line = append(textversion.Line, vline)
		textversion.Text = append(textversion.Text, line[1])
		latinwords, _ := strconv.Atoi(line[2])
		textversion.Latin = append(textversion.Latin, latinwords)
		greekwords, _ := strconv.Atoi(line[3])
		textversion.Greek = append(textversion.Greek, greekwords)
		other, _ := strconv.Atoi(line[4])
		textversion.Other = append(textversion.Other, other)
		header := false
		if line[5] != "" {
			header = true
		}
		textversion.Header = append(textversion.Header, header)
	}
	return textversion
}

func testsimilarity(textversion2 Textversion, textversion Textversion, window int) {
	filename := "Vol" + textversion.Volume[0] + "_" + textversion.Version[0] + "_matchedWith_" + textversion2.Version[0] + ".csv"
	fileOS, _ := os.Create(filename)
	defer fileOS.Close()
	writer := csv.NewWriter(fileOS)
	writer.Comma = '#'
	defer writer.Flush()

	for tester := range textversion.Text {
		var starttester int
		var biggest int
		switch {
		case tester < window:
			starttester = tester + 0
		default:
			starttester = tester * len(textversion2.ID) / len(textversion.ID)
		}
		text := textversion.Text[tester]
		scores := []float64{}
		comparetext := ""
		testngrams := []string{}
		comparengrams := []string{}
		for ngram, _ := range Parse(text, 3) {
			testngrams = append(testngrams, ngram)
		}
		var start, end int
		switch {
		case starttester-window <= 0:
			start = 0
			end = starttester + window
		case starttester >= len(textversion2.Text):
			start = len(textversion2.Text) - window
			end = len(textversion2.Text) - 1
		case starttester+window >= len(textversion2.Text)-1:
			start = starttester - window
			end = len(textversion2.Text) - 1
		default:
			end = starttester + window
			start = starttester - window
		}
		for i := range textversion2.Text[start:end] {
			count := 0
			comparetext = textversion2.Text[i+start]
			comparengrams = []string{}
			for ngram, _ := range Parse(comparetext, 3) {
				comparengrams = append(comparengrams, ngram)
			}
			for j := range testngrams {
				for k := range comparengrams {
					if testngrams[j] == comparengrams[k] {
						count = count + 1
					}
				}
			}
			score := float64(count) / float64(len(testngrams))
			scores = append(scores, score)
		}
		var n float64
		for i, v := range scores {
			if v > n {
				n = scores[i]
				biggest = i + start
			}
		}
		switch {
		case n > 0.75:
			tested := float64(tester) / float64(len(textversion.ID)) * float64(100)
			pertested := strconv.FormatFloat(tested, 'f', 2, 32)
			writer.Write([]string{textversion.ID[tester], textversion2.ID[biggest], strconv.FormatFloat(n, 'f', 3, 32)})
			fmt.Println("------------------------------------")
			fmt.Println("Tested:", pertested, "%")
		default:
			writer.Write([]string{textversion.ID[tester], "", strconv.FormatFloat(0.0, 'f', 3, 32)})
		}
	}
}

func main() {
	window := 300
	files, _ := ioutil.ReadDir("./data/")
	matchingnames := []string{}

	for _, f := range files {
		filename := f.Name()
		filenamesl := strings.Split(filename, ".")
		filename = strings.Join(filenamesl[:len(filenamesl)-3], ".") + "."
		matchingnames = append(matchingnames, filename)
	}
	unique := UniqStr(matchingnames)

	for i := range unique {
		fmt.Println(unique[i])
		number := unique[i] + "[a-z]"
		fmt.Println(number)
		matchingnames = []string{}
		for _, f := range files {
			filename := f.Name()
			matched, _ := regexp.MatchString(number, filename)
			if matched == true {
				filename = "data/" + filename
				matchingnames = append(matchingnames, filename)
			}
		}
		fmt.Println(matchingnames[0])
		textversion := readPG(matchingnames[0])
		fmt.Println(matchingnames[1:len(matchingnames)])
		textcollection := matchingnames[1:len(matchingnames)]
		var pgCollection []Textversion
		for i := range textcollection {
			textversion2 := readPG(textcollection[i])
			pgCollection = append(pgCollection, textversion2)
		}
		fmt.Println("Data is read. Begin testing. This will take a while...")
		for i := range pgCollection {
			textversion2 := pgCollection[i]
			go testsimilarity(textversion2, textversion, window)
		}
	}
	var input string
	fmt.Scanln(&input)
	fmt.Println("Done.")
}
