# CGI-Bin Zipcode

> A cgi-bin zipcode geocoding script. Written in Go.

Dataset from [Geo Names](http://www.geonames.org/) ((CC BY 4.0))

New datasets can be download from: http://download.geonames.org/export/zip

On first run, the script will turn the dataset into a .gob quick dataset loading.

After updating the dataset, remove the old .gob file and a new file will be created on next run.

Support for single country dataset (US). Multiple countries to be added in the future.

## Compiling

```bash
go build zipcode.go
```

Linux cross compile:

```bash
GOOS=linux GOARCH=amd64 go build -o bin/zipcode-amd64-linux zipcode.go
```

## Usage

Executed from cgi-bin folder.
Place 'US.txt' or generated gob dataset in same directory as binary.

### Request

```json
{
	"zipcode": "97211"
	"countryCode: "US"
}
```

### Repsonse

```json
{
	"zipCode": "97211",
	"placeName": "Portland",
	"adminName": "Oregon",
	"latitude": "45.5653",
	"longitude": "-122.6448"
}
```

## TODO

- [ ] Multi country support, specified by request
- [ ] Response error codes
- [ ] Clean up error handling

## Release History

- 1.0
  - Inital release

## Meta

Drew VanVlack – [@avanvlack](https://twitter.com/avanvlack) – A.Vanvlack@Gmail.com
Distributed under the MIT license. See `LICENSE` for more information.

[https://github.com/avanvlack](https://github.com/avanvlack/)
