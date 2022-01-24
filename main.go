package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"code.gitea.io/sdk/gitea"
)

var (
	instance     = flag.String("instance", "", "Gitea instance")
	token        = flag.String("token", "", "Gitea API token")
	user         = flag.String("user", "", "repo owner")
	repo         = flag.String("repo", "", "repo name")
	path         = flag.String("path", "", "filepath to be attached")
	filename     = flag.String("filename", "", "attachment filename")
	removeOthers = flag.Bool("remove-others", false, "remove other attachments with this name")
)

func main() {
	flag.Parse()

	if *instance == "" {
		i, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_INSTANCE")
		if !ok {
			fmt.Println("incorrect arguments: no instance")
			os.Exit(1)
		}
		instance = &i
	}

	if *token == "" {
		t, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_TOKEN")
		if !ok {
			fmt.Println("incorrect arguments: no token")
			os.Exit(1)
		}
		token = &t
	}

	if *user == "" {
		u, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_USER")
		if !ok {
			fmt.Println("incorrect arguments: no user")
			os.Exit(1)
		}
		user = &u
	}

	if *repo == "" {
		r, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_REPO")
		if !ok {
			fmt.Println("incorrect arguments: no repo")
			os.Exit(1)
		}
		repo = &r
	}

	if *path == "" {
		p, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_PATH")
		if !ok {
			fmt.Println("incorrect arguments: no path")
			os.Exit(1)
		}
		path = &p
	}

	if *filename == "" {
		f, _ := syscall.Getenv("GITEA_RELEASE_ATTACHER_FILENAME")
		filename = &f
	}

	if *removeOthers == false {
		r, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_REMOVE_OTHERS")
		// only run this if it is set
		if ok {
			remove, err := strconv.ParseBool(r)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			removeOthers = &remove
		}
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

	if *filename == "" {
		p := strings.Split(*path, "/")
		filename = &p[len(p)-1]
	}

	if *removeOthers {
		for _, v := range releases[0].Attachments {
			if v.Name == *filename {
				_, err := c.DeleteReleaseAttachment(*user, *repo, releases[0].ID, v.ID)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}

	file, err := os.OpenFile(*path, os.O_RDONLY, 0o644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_, _, err = c.CreateReleaseAttachment(*user, *repo, releases[0].ID, file, *filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
