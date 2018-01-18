# JSON Path Flattening

### How To Use
You can pass a filepath, URL, or string to flatten.
```
go run main.go ~/JSON-path-flattening/sample.json
go run main.go https://raw.githubusercontent.com/chandler767/JSON-path-flattening/master/sample.json
go run main.go "{\"server_ip\":\"192.168.0.1\",\"action\":\"stop\"}"
```

A JSON path is a string that uniquely identifies a subvalue inside a JSON 
value. For example, given the following JSON:

```
{
    "books": [
        {
            "title": "JSON and you",
            "pages": 234
        }
    ]
}
```

The path `.books` points at the array of books, `.books[0]` points at the first
book, and `.books[0].title` at the title of the first book.

Takes **any arbitrary JSON** as input and writes as 
output, line-by-line, the path to every scalar value in the JSON, an equals 
sign, and the value. The output for the example JSON given above would be:

```
.books[0].title=JSON and you
.books[0].pages=234
```

