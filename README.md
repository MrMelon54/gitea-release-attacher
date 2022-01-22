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

### CI/CD

You can find a sample configuration that publishes the binary for this repository to the latest release at https://codeberg.org/qwerty287/gitea-release-attacher/src/branch/main/.woodpecker/build.yml#L9-L14. However, this does not use the `latest` file from the releases it published, instead it compiles it from source.

## Build

```sh
go build
```

If you would like to contribute, please format your code with `gofumpt` and make sure that it passes the `golangci-lint` linters.