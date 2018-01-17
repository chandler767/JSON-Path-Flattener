# JSON-path-flattening

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

Write a program that takes **any arbitrary JSON** as input and writes as 
output, line-by-line, the path to every scalar value in the JSON, an equals 
sign, and the value. The output for the example JSON given above would be:

```
.books[0].title=JSON and you
.books[0].pages=234
```