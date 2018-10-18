FROM golang:1.11

WORKDIR /app

RUN apt-get update && apt-get install -y inotify-tools

# Install golint package
RUN go get -u golang.org/x/lint/golint

ADD cmd/ /app/cmd
ADD pkg/ /app/pkg
ADD Makefile /app/Makefile
ADD create.sql /app/create.sql
ADD go.mod /app/go.mod
ADD go.sum /app/go.sum

RUN make

CMD /app/eventsourcing-hack

EXPOSE 8000
