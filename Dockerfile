FROM golang

RUN go install codeberg.org/qwerty287/gitea-release-attacher@v1.0.0

CMD /go/bin/gitea-release-attacher