Steps to Run
------------

* Export Env variables
  * export DATABASE_DSN="file:employee.db?cache=shared&_fk=1"
  * export SERVER_ADDR=":8080"
* go run ./cmd/server/main.go
