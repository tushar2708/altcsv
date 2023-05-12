# Altcsv
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Ftushar2708%2Faltcsv.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Ftushar2708%2Faltcsv?ref=badge_shield)

[![GitHub license](https://img.shields.io/github/license/mashape/apistatus.svg)]()


[![Maintainability](https://api.codeclimate.com/v1/badges/926ce49973984e9aac06/maintainability)](https://codeclimate.com/github/tushar2708/altcsv/maintainability)


An alternative to standard CSV implementation of Go (<https://github.com/golang/go/tree/master/src/encoding/csv>),
and modified by Tushar Dwivedi (<https://github.com/tushar2708/altcsv>),
to support some more delimiters, and a few non-standard(?) CSV formats.

In addition to standard encoding/csv functionality,

* You can use a custom quoting characters for reading and writing CSV, by giving setting `Quote` to a desired rune.
* You can force quote all fields, by setting `AllQuotes` flag as `true`

## Install

`go get github.com/tushar2708/altcsv`

Or, using glide:

`glide get github.com/tushar2708/altcsv`

## Usage

### Reading CSV file

```go
fileRdr, _ = os.Open("/tmp/custom_csv_file.txt")
csvRdr := altcsv.NewReader(fileRdr)
csvRdr.Comma = ';'    // use ; as "comma"
csvRdr.Quote = '|'    // use | as "quote"
content := csvReader.ReadAll()
```

### Writing CSV file

```go
headers = []string{"hero_name", "alter_ego", "identity"}
fileWtr, _ := os.Create("/tmp/all_quotes_csv_file.txt")
csvWtr := altcsv.NewWriter(csvH)
csvWtr.Quote = '|'        // use | as "quote"
csvWtr.AllQuotes = true   // surround each field with '|'
csvWtr.Write(headers)
csvWtr.Write([]string{"Spider-Man", "Peter Parker", "Secret Identity"})
csvWtr.Write([]string{"Captain America", "Steven Rogers", "Public Identity"})
csvWtr.Write([]string{"Thor", "Thor Odinson", "No dual Identity"})
csvWtr.Flush()
fileWtr.Close()
```

## To Do

* [x] Supporting custom `Quote` character.
* [x] Supporting `AllQuotes` flag.
* [ ] Supporting reading or writing headers.
* [ ] Reading & writing Row as a map[string]string, with header fields as keys.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftushar2708%2Faltcsv.svg?type=small)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftushar2708%2Faltcsv?ref=badge_small)
