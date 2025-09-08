FROM alpinelinux/golang
WORKDIR /usr/src/app

#go.sum не требуется в данном проекте
COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -v -o app ./

CMD ["./app"]

