# _Dream Theater_

![Project Image](./project.gif)

---

### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [Author Info](#author-info)

---

## Description

Dream Theater lets it's user to check out which movies avaiable right now and gives brief information about them. You can register and buy a ticket right now!

## Technologies

### Main Technologies

- [Go](https://go.dev/)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/)

### Libraries

- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)
- [golang/mock](https://github.com/golang/mock)
- [google/uuid](https://github.com/google/uuid)
- [lib/pq](https://github.com/lib/pq)
- [spf13/viper](https://github.com/spf13/viper)
- [stretchr/testify](https://github.com/stretchr/testify)
- [crypto](https://golang.org/x/crypto)

[Back To The Top](#dream-theater)

---

## How To Use

### Tools

- [Go](https://go.dev/dl/)
- [DataGrip](https://www.jetbrains.com/datagrip/download/#section=mac)
- [golang-migrate/migrate](https://github.com/golang-migrate/migrate)

### Setup Database

- Create Theater Database in psql console

```
CREATE DATABASE $name$;
```

- Migrations runs inside of the program

- For migrations down create a bash script down.sh

```
#!/bin/bash

migrate -path internal/db/migration -database "postgresql://$username$:$password$@localhost:5432/Theater?sslmode=disable" -verbose down
```

- Grant executable

```
chmod +x down.sh
```

- And run

```
make down
```

### Generate Database functions

- Generate SQL CRUD functions

```
make sqlc
```

- Generate mockdb

```
make mock
```

### Run tests

- Create a bash script test.sh

```
#!/bin/bash

DB_SOURCE="postgresql://$username$:$password$@localhost:5432/Theater?sslmode=disable" go test -v -cover ./...
```

- Grant executable

```
chmod +x test.sh
```

```
make test
```

### Start App

- To Start App create a bash script start.sh

```
#!/bin/bash

DB_SOURCE="postgresql://$username$:$password$@localhost:5432/Theater?sslmode=disable" TOKEN_SYMMETRIC_KEY=$32 character string$ go run cmd/web/main.go
```

- Grant executable

```
chmod +x start.sh
```

```
make server
```

### Give it a try

#### Routes

- Check the documentation from [here](https://documenter.getpostman.com/view/21428220/VUqpsxJX)

Don't forget to copy your access token for authentication required routes after logging in!

[Back To The Top](#dream-theater)

---

## Author Info

- Twitter - [@dev_bck](https://twitter.com/dev_bck)

[Back To The Top](#dream-theater)
