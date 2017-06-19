package main

import (
	"encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "os/exec"
  "os"
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
  args := []string{reponame[40]}
  if err := exec.Command(cmd, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println("Successfully called ./test")
}
