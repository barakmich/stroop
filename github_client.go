package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var upstream = flag.String("upstream", "origin", "Github user or org for upstream fork")

type tokenSource struct {
	token *oauth2.Token
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return t.token, nil
}

type stroopClient struct {
	github *github.Client
	remote *Remote
}

func startClient() *stroopClient {
	keyfile, err := os.Open(os.ExpandEnv("$HOME/.github.key"))
	if err != nil {
		log.Fatal(err)
	}
	r := bufio.NewReader(keyfile)
	key, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	ts := &tokenSource{
		&oauth2.Token{AccessToken: key},
	}
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return &stroopClient{
		github: github.NewClient(tc),
	}
}

func (c *stroopClient) MustDetectRemote() {
	remotes, _ := getRemotes()
	var remote *Remote
	for _, r := range remotes {
		if r.LocalName == *upstream {
			remote = r
			break
		}
	}
	if remote == nil {
		log.Fatal("Couldn't find upstream remote named", *upstream)
	}
	c.remote = remote
}

func (c *stroopClient) GetIssues() []github.Issue {
	var issues []github.Issue
	for i := 1; i < 4; i++ {
		opts := &github.IssueListByRepoOptions{
			ListOptions: github.ListOptions{
				PerPage: 50,
				Page:    i,
			},
		}
		i, _, err := c.github.Issues.ListByRepo(c.remote.User, c.remote.Repo, opts)
		if err != nil {
			log.Errorf("%s", err.Error())
			return []github.Issue{}
		}
		issues = append(issues, i...)
	}
	return issues
}

type Remote struct {
	LocalName string
	User      string
	Repo      string
}

func (r *Remote) GithubName() string {
	return fmt.Sprintf("%s/%s", r.User, r.Repo)
}

func getRemotes() ([]*Remote, error) {
	var out []*Remote
	cmd := exec.Command("git", "remote", "-v")
	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(bytes)), "\n")
	re, err := regexp.Compile("github.com[:/](.*?)/(.*?).git")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rehttp, err := regexp.Compile("github.com[:/](.*?)/(.*)")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for _, line := range lines {
		record := strings.Split(line, "\t")
		lname := record[0]
		record = strings.Split(record[1], " ")
		path := record[0]
		fetchpush := record[1]
		if fetchpush != "(fetch)" {
			continue
		}
		submatch := re.FindStringSubmatch(path)
		if submatch == nil {
			submatch = rehttp.FindStringSubmatch(path)
			if submatch == nil {
				continue
			}
		}
		out = append(out, &Remote{
			LocalName: lname,
			User:      submatch[1],
			Repo:      submatch[2],
		})

	}
	return out, nil
}
