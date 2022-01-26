package main

import (
	"flag"
	"log"
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
	removeAll    = flag.Bool("remove-all", false, "remove all attachments before attaching the new file")
	drafts       = flag.Bool("drafts", false, "publish also to draft releases")
	preRelease   = flag.Bool("pre-release", false, "publish also to pre releases")
	releaseID    = flag.Int64("release-id", 0, "release ID to attach file")
	releaseTag   = flag.String("release-tag", "", "release tag to attach file")
)

func main() {
	flag.Parse()

	if *instance == "" {
		i, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_INSTANCE")
		if !ok {
			log.Fatal("incorrect arguments: no instance")
		}
		instance = &i
	}

	if *token == "" {
		t, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_TOKEN")
		if !ok {
			log.Fatal("incorrect arguments: no token")
		}
		token = &t
	}

	if *user == "" {
		u, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_USER")
		if !ok {
			log.Fatal("incorrect arguments: no user")
		}
		user = &u
	}

	if *repo == "" {
		r, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_REPO")
		if !ok {
			log.Fatal("incorrect arguments: no repo")
		}
		repo = &r
	}

	if *path == "" {
		p, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_PATH")
		if !ok {
			log.Fatal("incorrect arguments: no path")
		}
		path = &p
	}

	if *filename == "" {
		f, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_FILENAME")
		if !ok {
			p := strings.Split(*path, "/")
			f = p[len(p)-1]
		}
		filename = &f
	}

	if !*removeOthers {
		r, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_REMOVE_OTHERS")
		// only run this if it is set
		if ok {
			remove, err := strconv.ParseBool(r)
			if err != nil {
				log.Fatal(err)
			}
			removeOthers = &remove
		}
	}

	if !*removeAll {
		r, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_REMOVE_ALL")
		// only run this if it is set
		if ok {
			remove, err := strconv.ParseBool(r)
			if err != nil {
				log.Fatal(err)
			}
			removeAll = &remove
		}
	}

	if !*drafts {
		dEnv, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_DRAFTS")
		// only run this if it is set
		if ok {
			d, err := strconv.ParseBool(dEnv)
			if err != nil {
				log.Fatal(err)
			}
			drafts = &d
		}
	}

	preReleaseSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "pre-release" {
			preReleaseSet = true
		}
	})

	if !preReleaseSet {
		pEnv, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_PRE_RELEASE")
		// only run this if it is set
		if ok {
			p, err := strconv.ParseBool(pEnv)
			if err != nil {
				log.Fatal(err)
			}
			preRelease = &p
		} else {
			preRelease = nil
		}
	}

	releaseIDIsEnv := false
	if *releaseID == 0 {
		i, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_RELEASE_ID")
		if ok {
			releaseIDIsEnv = true
			i, err := strconv.ParseInt(i, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			releaseID = &i
		}
	}

	releaseTagIsEnv := false
	if *releaseTag == "" {
		t, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_RELEASE_TAG")
		releaseTagIsEnv = ok
		releaseTag = &t
	}

	useReleaseID := true
	if *releaseTag != "" && *releaseID != 0 {
		if releaseIDIsEnv == releaseTagIsEnv {
			log.Fatal("incorrect arguments: both release ID and tag set")
		}
		useReleaseID = !releaseIDIsEnv
	}

	c, err := gitea.NewClient(*instance, gitea.SetToken(*token))
	if err != nil {
		log.Fatal(err)
	}

	var release *gitea.Release
	if *releaseID != 0 && useReleaseID {
		r, _, err := c.GetRelease(*user, *repo, *releaseID)
		if err != nil {
			log.Fatal(err)
		}

		release = r
	} else if *releaseTag != "" {
		r, _, err := c.GetReleaseByTag(*user, *repo, *releaseTag)
		if err != nil {
			log.Fatal(err)
		}

		release = r
	} else {
		releases, _, err := c.ListReleases(*user, *repo, gitea.ListReleasesOptions{
			ListOptions: gitea.ListOptions{
				PageSize: 1,
			},
			IsDraft:      drafts,
			IsPreRelease: preRelease,
		})
		if err != nil {
			log.Fatal(err)
		}

		if len(releases) != 1 {
			log.Fatal("no releases")
		}
		release = releases[0]
	}

	if *removeOthers || *removeAll {
		for _, v := range release.Attachments {
			if *removeAll || v.Name == *filename {
				_, err := c.DeleteReleaseAttachment(*user, *repo, release.ID, v.ID)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	file, err := os.OpenFile(*path, os.O_RDONLY, 0o644)
	if err != nil {
		log.Fatal(err)
	}

	_, _, err = c.CreateReleaseAttachment(*user, *repo, release.ID, file, *filename)
	if err != nil {
		log.Fatal(err)
	}
}
