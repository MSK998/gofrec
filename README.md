# gofrec

GOF(ile)REC(ord) is a multi record file parser that can read multiple different kinds of records based on an identifier supplied by the user.

## Example

### Types
First we need to set up our types that will be used by the parser. Each type must have **one** identfier tag with each field containing a length tag.

```go
type T1 struct {
	ID   string `Identifier:"001" Length:"3"`
	Name string `Length:"8"`
}

type T2 struct {
	ID   string `Identifier:"002" Length:"3"`
	Name string `Length:"8"`
}

type T3 struct {
	ID     string `Identifier:"003" Length:"3"`
	Number int    `Length:"8"`
}
```

### Instantiating the parser

To instantiate a parser for basic usage the folling information is required: 

1. An identification start and end in the same form as the "slice" operator eg `slice[0:3]` for index 0 up to (but excluding) 3
2. A slice of empty types that you want the parser to use to identify the records

```go
recordTypes := []interface{}{
    T1{},
    T2{},
    T3{},
}

parser := gofrec.Parser{
    RecordTypes: recordTypes,
    IdStart: 0,
    IdEnd: 3,
}
```

### Reading the file

The reading of the file is put up to the developer using this package, there are so many ways to do it so I didn't want to lock anyone into a single way.

Once the file has been read into the application there is an optional step of converting the bytes to lines with the `parser.BytesToLines(fileContents []bytes)` method. 

```go
bytes, err := os.ReadFile("example.txt")
if err != nil {
    panic(err)
}

numLines, err := parser.BytesToLines(bytes)
if err != nil {
    panic(err)
}
```

### Parsing the data

Finally, to parse all of these records into a type we call the `parser.Parse()` method as below.

```go
numRecords, err := parser.Parse()
if err != nil {
    panic(err)
}
```

At this point the parser will have all records in the `parser.Records` field and can be accessed like any other slice/array.