# Go event sourcing hack project

## Running

Starting the database
```
docker-compose up
```

## Creating the events table
In another terminal paste the following.

```
mysql -P 3306 -h 127.0.0.1 -u root -ppassword events < create.sql
```

## Installing and starting the app
```
go get
go run *.go
```
