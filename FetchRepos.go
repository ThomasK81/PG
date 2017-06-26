package main

import (
	"encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os/exec"
  "os"
	"path/filepath"
	"regexp"
	"strings"
)

type Repo struct {
	Name   string `json:"name"`
}

func findRepos(s string) []Repo {
	data, err := getContent(s)
	if err != nil {
		fmt.Println(err)
		return []Repo{}
	}
	var repos []Repo
	err = json.Unmarshal(data, &repos)
	if err != nil {
		fmt.Println(err)
		return []Repo{}
	}
	return repos
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

func main() {
  var reponame []string
  repos := findRepos("https://api.github.com/users/OGL-PatrologiaGraecaDev/repos?per_page=100")
  for i:= range repos {
    reponame = append(reponame, repos[i].Name)
  }
  repos = findRepos("https://api.github.com/users/OGL-PatrologiaGraecaDev/repos?per_page=100&page=2")
  for i:= range repos {
    reponame = append(reponame, repos[i].Name)
  }
  cmd := "./test"
	storedfiles := checkExt(".csv")
	reponame = removefromRepo(reponame, storedfiles)

  for i := len(reponame)-1; i >= 0; i-- {
    args := []string{reponame[i]}
    if err := exec.Command(cmd, args...).Run(); err != nil {
      fmt.Fprintln(os.Stderr, err)
      os.Exit(1)
    }
    fmt.Println("Successfully called ./test", reponame[i])}
}

func checkExt(ext string) []string {
	pathS, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var files []string
	filepath.Walk(pathS, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(ext, f.Name())
			if err == nil && r {
        filename := strings.Join([]string{strings.Split(f.Name(), ".")[1], strings.Split(f.Name(), ".")[2]}, ".")
				files = append(files, filename)
			}
		}
		return nil
	})
	return files
}

func removeDuplicatesUnordered(elements []string) []string {
    encountered := map[string]bool{}

    // Create a map of all unique elements.
    for v:= range elements {
        encountered[elements[v]] = true
    }

    // Place all keys from the map into a slice.
    result := []string{}
    for key, _ := range encountered {
        result = append(result, key)
    }
    return result
}

func removefromRepo(input, control []string) []string {
	var result []string
	for i := range input {
		if strcontains(control, input[i]) == false {
			result = append(result, input[i])
		}
		}
		return result
}

func strcontains(s []string, e string) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
