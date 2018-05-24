package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Owner represents the owner payload in Github
type Owner struct {
	Login string `json:"login"`
}

// Repo represents the repository payload in Github
type Repo struct {
	Name      string `json:"name,omitempty"`
	CloneURL  string `json:"clone_url,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Owner     Owner  `json:"owner,omitempty"`
	Message   string `json:"message,omitempty"`
}

func main() {

	outputFile := flag.String("o", "repos.csv", "The output csv destination")
	flag.Parse()

	var repoNames = askRepos()
	fmt.Printf("Fetching %d repositories...\n", len(repoNames))

	file, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalf("unable to create repos.csv: %v", err)
	}
	defer file.Close()

	var repos = fetchRepos(repoNames)
	writeCSV(file, repos)
	writeCSV(os.Stdout, repos)

	fmt.Println("Done. Press ctrl + c to cancel.")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func writeCSV(out io.Writer, repos []*Repo) {
	var header = []string{"updated_at", "name", "login", "clone_url"}
	w := csv.NewWriter(out)
	defer w.Flush()

	// Write header to CSV
	w.Write(header)

	for _, repo := range repos {
		// Write to CSV
		w.Write([]string{repo.UpdatedAt, repo.Name, repo.Owner.Login, repo.CloneURL})
	}
}

func askRepos() []string {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(`Enter the repository name, e.g. kubernetes/charts:`)

	var repos []string
	for scanner.Scan() {
		if strings.ToLower(scanner.Text()) == "y" {
			break
		}
		repo := strings.TrimSpace(scanner.Text())
		if repo != "" {
			repos = append(repos, repo)
		}
		fmt.Println("Enter repository name, or [Y] to proceed:")
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	if len(repos) == 0 {
		log.Fatal("Error: You must provide at least one Github repository")
	}

	return repos
}

func fetchRepo(name string) (*Repo, error) {
	url := fmt.Sprintf("http://api.github.com/repos/%s", name)
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repo Repo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		return nil, err
	}

	if repo.Message == "Not Found" {
		return nil, nil
	}
	return &repo, nil
}

func fetchRepos(repoNames []string) []*Repo {
	var repos []*Repo
	for _, name := range repoNames {
		repo, err := fetchRepo(name)
		if err != nil {
			fmt.Printf("unable to fetch %s: %v\n", name, err)
			continue
		}

		if repo == nil {
			fmt.Printf("cannot find repository \"%s\"\n", name)
			continue
		}
		repos = append(repos, repo)
	}
	return repos
}
