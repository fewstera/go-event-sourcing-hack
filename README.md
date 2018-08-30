# Go event sourcing hack project

## Running

Starting the database
```
docker-compose up
```

### Creating the events table
In another terminal paste the following.

```
mysql -P 3306 -h 127.0.0.1 -u root -ppassword events < create.sql
```

### Installing and starting the app
```
go get
go run *.go
```

## Useful curl commands

Create a new user

```
curl http://localhost:8080/users -d '{"name": "Aidan Fewster", "age": 25 }'
```

Get a user

```
curl http://localhost:8080/users/{userId}
```

Increase a users age

```
curl http://localhost:8080/users/{userId}/increase-age -XPOST
```

Change a users name

```
curl http://localhost:8080/users/{userId} -XPATCH -d '{"name": "Aidan Wynne Fewster"}'
```

## Event sourcing

### Stream categories

* USER

### Commands

* Create new user
* Increase users age
* Change users name

### Events

* User created (USER_CREATED)
* User got older (USER_GOT_OLDER)
* User name changed (USER_NAME_CHANGED)
