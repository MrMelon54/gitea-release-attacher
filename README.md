# gitea-release-attacher

Add an attachment to the latest Gitea release

## Usage

```
Usage of ./gitea-release-attacher:
  -filename string
        attachment filename
  -instance string
        Gitea instance
  -path string
        filepath to be attached
  -remove-others
        remove other attachments with this name
  -repo string
        repo name
  -token string
        Gitea API token
  -user string
        repo owner
```

## Build

```sh
go build
```