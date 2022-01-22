package main

import (
	"fmt"
	"flag"
	"os"
	"strings"

	"code.gitea.io/sdk/gitea"
)

var (
	instance = flag.String("instance", "", "Gitea instance")
	token = flag.String("token", "", "Gitea API token")
	user = flag.String("user", "", "repo owner")
	repo = flag.String("repo", "", "repo name")
	path = flag.String("path", "", "filepath to be attached")
	filename = flag.String("filename", "", "attachment filename")
)

func main() {
	flag.Parse()
	
	if *instance == "" || *user == "" || *repo == "" || *path == "" {
		fmt.Println("incorrect arguments")
		os.Exit(1)
	}
	c, err := gitea.NewClient(*instance, gitea.SetToken(*token))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	releases, _, err := c.ListReleases(*user, *repo, gitea.ListReleasesOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 1,
		},
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(releases) != 1 {
		fmt.Println("no releases")
		os.Exit(1)
	}

	file, err := os.OpenFile(*path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *filename == "" {
		p := strings.Split(*path, "/")
		filename = &p[len(p) - 1]
	}
	_, _, err = c.CreateReleaseAttachment(*user, *repo, releases[0].ID, file, *filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}