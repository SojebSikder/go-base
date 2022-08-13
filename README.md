# go-base

Simple database engine created with golang just for fun

## Usage

run command from command line:

```
go run main.go cli
```

run command from file

```
go run main.go run file.sql
```

### Supported command:

- Create database
  ```sql
  create db [blog]
  ```
- Create document:

  ```sql
  create [user]
  ```

- Insert data into specific document:

  ```sql
  insert [user] 'sojeb' 'sikder'
  ```
