# go-base

Simple database engine created with golang just for fun

## Usage

Run command from command line:

```
go run main.go cli
```

Run command from file

```
go run main.go run file.sql
```

### Supported commands:

- Create database
  ```sql
  create db [blog]
  ```
- Drop database
  ```sql
  drop db [blog]
  ```
- Select database in sql file
  ```sql
  set db [blog]
  ```
- Create document:

  ```sql
  create doc [user]
  ```

- Insert data into specific document:

  ```sql
  insert [user] 'sojeb' 'sikder'
  ```
