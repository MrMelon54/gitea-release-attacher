# gitea-release-attacher

Add an attachment to the latest Gitea release

## Usage

```
Usage of ./gitea-release-attacher:
  -drafts
    	publish also to draft releases
  -filename string
    	attachment filename
  -instance string
    	Gitea instance
  -path string
    	filepath to be attached
  -pre-release
    	publish also to pre releases
  -release-id int
    	release ID to attach file
  -release-tag string
    	release tag to attach file
  -remove-all
    	remove all attachments before attaching the new file
  -remove-others
    	remove other attachments with this name
  -repo string
    	repo name
  -token string
    	Gitea API token
  -user string
    	repo owner
```

### Environment variables

Setting any option using environment variables is supported. The environment variables have the scheme `GITEA_RELEASE_ATTACHER_*`. You can replace the the `*` with these values:

* `INSTANCE`
* `TOKEN`
* `USER`
* `REPO`
* `PATH`
* `FILENAME`
* `REMOVE_OTHERS`
* `REMOVE_ALL`
* `DRAFTS`
* `PRE_RELEASE`
* `RELEASE_ID`
* `RELEASE_TAG`

They have the same effects as the corresponding command line options, but the command line options are preferred.

### Docker

A Docker image is provided at [qwerty287/gitea-release-attacher](https://hub.docker.com/r/qwerty287/gitea-release-attacher). This has the application as entrypoint, it is recommended to use environment variables.

## Build

```sh
go build
```

If you would like to contribute, please format your code with `gofumpt` and make sure that it passes the `golangci-lint` linters.