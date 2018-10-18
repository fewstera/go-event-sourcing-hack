# Go event sourcing hack project

## Running

You can start the project local by running.
```
make docker-start-dev
```

## Useful curl commands

Get all users

```
curl http://localhost:8080/users/
```

Get a user

```
curl http://localhost:8080/users/:id
```

Create a new user

```
curl http://localhost:8080/users -d '{"name": "Aidan Fewster", "age": 25 }'
```

Deposit cash into a users account

```
curl http://localhost:8080/users/:id/deposit -d '{"version":1, "amount": 500.00}'
```

## Event sourcing

### Stream categories

* USER

### Commands

* Create new user
* Deposit cash

### Events

* User created (USER_CREATED)
* Deposited cash (DEPOSITED)
