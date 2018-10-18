# Go event sourcing hack project

## Running

You can start the project local by running.
```
make docker-start-dev
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

## Event sourcing

### Stream categories

* USER

### Commands

* Create new user

### Events

* User created (USER_CREATED)
