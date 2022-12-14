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

	getenv := func(name string) (string, bool) {
		i, ok := syscall.Getenv("GITEA_RELEASE_ATTACHER_" + name)
		if ok {
			return i, true
		}

		return syscall.Getenv("PLUGIN_" + name)
	}

	errors := []string{}
	if *instance == "" {
		i, ok := getenv("INSTANCE")
		if !ok {
			errors = append(errors, "no instance")
		}
		instance = &i
	}

	if *token == "" {
		t, ok := getenv("TOKEN")
		if !ok {
			errors = append(errors, "no token")
		}
		token = &t
	}

	if *repo == "" {
		r, ok := getenv("REPO")
		if !ok {
			errors = append(errors, "no repository")
		}
		repo = &r
	}

	if *user == "" {
		u, ok := getenv("USER")
		if !ok {
			if strings.Contains(*repo, "/") {
				u = strings.Split(*repo, "/")[0]
				r := strings.Split(*repo, "/")[1]
				repo = &r
			} else {
				errors = append(errors, "no user")
			}
		}
		user = &u
	}

	if *path == "" {
		p, ok := getenv("PATH")
		if !ok {
			errors = append(errors, "no path")
		}
		path = &p
	}

	if *filename == "" {
		f, ok := getenv("FILENAME")
		if !ok {
			p := strings.Split(*path, "/")
			f = p[len(p)-1]
		}
		filename = &f
	}

	if !*removeOthers {
		r, ok := getenv("REMOVE_OTHERS")
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
		r, ok := getenv("REMOVE_ALL")
		// only run this if it is set
		if ok {
			remove, err := strconv.ParseBool(r)
			if err != nil {
				log.Fatal(err)
			}
			removeAll = &remove
		}
	}

	draftsSet := false
	preReleaseSet := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == "drafts" {
			draftsSet = true
		} else if f.Name == "pre-release" {
			preReleaseSet = true
		}
	})

	if !draftsSet {
		dEnv, ok := getenv("DRAFTS")
		// only run this if it is set
		if ok {
			d, err := strconv.ParseBool(dEnv)
			if err != nil {
				log.Fatal(err)
			}
			drafts = &d
		} else {
			drafts = nil
		}
	}

	if !preReleaseSet {
		pEnv, ok := getenv("PRE_RELEASE")
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
		i, ok := getenv("RELEASE_ID")
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
		t, ok := getenv("RELEASE_TAG")
		releaseTagIsEnv = ok
		releaseTag = &t
	}

	useReleaseID := true
	if *releaseTag != "" && *releaseID != 0 {
		if releaseIDIsEnv == releaseTagIsEnv {
			errors = append(errors, "both release ID and tag set")
		}
		useReleaseID = !releaseIDIsEnv
	}

	if len(errors) > 0 {
		log.Fatal("incorrect arguments: " + strings.Join(errors, ", "))
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
