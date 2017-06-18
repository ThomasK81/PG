package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type linestats struct {
	Id         string
	Count      int
	Percentage float64
	Title      string
	Line       string
}

func strcount(s []string, key string) linestats {
	result := 0
	for i := range s {
		if s[i] == key {
			result = result + 1
		}
	}
	var floatresult float64
	switch {
	case result == 0:
		floatresult = float64(0)
	default:
		floatresult = float64(result) / float64(len(s)) * 100
	}
	return linestats{Count: result, Percentage: floatresult}
}

func boolcount(s []bool, key bool) linestats {
	result := 0
	for i := range s {
		if s[i] == key {
			result = result + 1
		}
	}
	var floatresult float64
	switch {
	case result == 0:
		floatresult = float64(0)
	default:
		floatresult = float64(result) / float64(len(s)) * 100
	}
	return linestats{Count: result, Percentage: floatresult}
}

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

type FolderContent struct {
	Trees []Tree `json:"tree"`
}

type Tree struct {
	Path string `json:"path"`
}

type RepoMeta struct {
	Path   string `json:"path"`
	GitUrl string `json:"git_url"`
	Size   int    `json:"size"`
}

func findrepometa(s string) []RepoMeta {
	data, err := getContent(s)
	if err != nil {
		fmt.Println(err)
		return []RepoMeta{}
	}
	var repo []RepoMeta
	var folders []RepoMeta
	err = json.Unmarshal(data, &repo)
	if err != nil {
		fmt.Println(err)
		return []RepoMeta{}
	}
	for i := range repo {
		if repo[i].Size == 0 {
			folders = append(folders, repo[i])
		}
	}
	return folders
}

func findtxturl(s string) []string {
	data, err := getContent(s)
	if err != nil {
		fmt.Println("connection went wrong")
		return []string{}
	}
	var repo FolderContent
	err = json.Unmarshal(data, &repo)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	result := []string{}
	for i := range repo.Trees {
		if strings.Contains(repo.Trees[i].Path, ".txt") {
			result = append(result, repo.Trees[i].Path)
		}
	}
	return result
}

func main() {
	volume := os.Args[1]
	repo := findrepometa("https://api.github.com/repos/OGL-PatrologiaGraecaDev/" + volume + "/contents/")
	for i := range repo {
		txts := findtxturl(repo[i].GitUrl)
		if len(txts) != 0 {
			version := repo[i].Path
			filename := "pg." + volume + "." + version + ".csv"
			fileOS, _ := os.Create(filename)
			defer fileOS.Close()
			writer := csv.NewWriter(fileOS)
			writer.Comma = '#'
			defer writer.Flush()
			for j := range txts {
				file := txts[j]
				pageurl := "https://raw.githubusercontent.com/OGL-PatrologiaGraecaDev/" + volume + "/master/" + version + "/" + file
				fmt.Println(pageurl)
				pagecontent, _ := getContent(pageurl)
				page := string(pagecontent)

				var b []byte
				var language_int int
				var language string
				var languages []string
				var charsize []int
				lines := strings.Split(page, "\n")
				for n := range lines {
					fields := strings.Fields(lines[n])
					languages = []string{}
					for i := range fields {
						charsize = []int{}
						language_int = 0
						b = []byte(fields[i])
						for len(b) > 0 {
							_, size := utf8.DecodeLastRune(b)
							charsize = append(charsize, size)
							b = b[:len(b)-size]
						}
						for j := range charsize {
							language_int = language_int + charsize[j]
						}
						language_int = language_int / len(charsize)
						switch {
						case language_int == 1:
							language = "Latin"
						case language_int == 2:
							language = "Greek"
						default:
							language = "Other"
						}
						languages = append(languages, language)
					}
					capinfo := []bool{}
					re := regexp.MustCompile("[^a-z|^A-Z]")
					teststring := re.ReplaceAllString(lines[n], "")
					b = []byte(teststring)
					for len(b) > 0 {
						char, size := utf8.DecodeLastRune(b)
						capinfo = append(capinfo, unicode.IsLower(char))
						b = b[:len(b)-size]
					}
					m := make(map[string]linestats)
					m["Latin"] = strcount(languages, "Latin")
					m["Greek"] = strcount(languages, "Greek")
					m["Other"] = strcount(languages, "Other")
					m["Caps"] = boolcount(capinfo, false)
					switch {
					case m["Caps"].Percentage > float64(80):
						m["Title"] = linestats{Title: lines[n]}
					default:
						m["Title"] = linestats{Title: ""}
					}
					m["Text"] = linestats{Line: lines[n]}
					IdStr := strings.Split(filename, ".csv")[0] + ":" + strings.Split(file, ".")[0] + "." + strconv.Itoa(n+1)
					m["Id"] = linestats{Id: IdStr}
					writer.Write([]string{m["Id"].Id, m["Text"].Line, strconv.Itoa(m["Latin"].Count), strconv.Itoa(m["Greek"].Count), strconv.Itoa(m["Other"].Count), m["Title"].Title})
				}
			}
		}
	}
}
