# go-base

Simple database engine created with golang just for fun.

**Note: Currently query based db not working properly, any kind of contributions are welcome.**

### Todo:

- ~~Create cli (querydb)~~
- ~~Execute query file (querydb)~~
- implement writing data (querydb)
- implement b-tree for reading data (querydb)

## Usage

### In-memory key value based db

Run simple key value db from command line:

```
go run main.go mapdb
```

### query based db

Run query on command line:

```
go run main.go cli
```

Run query from file

```
go run main.go run file.sql
```

### Supported commands:

- Database oparations

  - Create database
    ```sql
    create db [blog]
    ```
  - Drop database
    ```sql
    drop db [blog]
    ```
  - Select database
    ```sql
    set db [blog]
    ```

- Document oparations

  - Create document:

    ```sql
    create doc [user]
    ```

  - Insert data into document:

    ```sql
    insert [user] {firstName} 'sojeb' {lastName} 'sikder'
    ```

## How go-base query based db works under the hood

In the first place go-base takes query. Then go-base engine splits query into statement using semi-clone (;) seperator. Each statements goes to tokenizer for generating tokens. All the tokens goes to perser. Perser perse the statement with brackets, quotation delimiter etc. And all the parsed data goes to main operation unit for processing.
