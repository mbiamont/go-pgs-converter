# go-pgs-converter
Go PGS subtitle converter

## Warning

`go-pgs-converter` needs [Tesseract](https://github.com/tesseract-ocr/tessdoc) to be installed on the machine. 

An easy way to do it, is to run `go-pgs-converter` on Docker. See [Docker](#Docker)

## Usage

```go
import pgs "github.com/mbiamont/go-pgs-converter"

subs, err := pgs.ConvertToSubtitles("input.sup", &pgs.ConversionOptions{
    InputLanguage: "eng",
})
```

## Fix OCR detection

```go
subs, err := pgs.ConvertToSubtitles("input.sup", &pgs.ConversionOptions{
    InputLanguage: "eng",
    TextCorrection: func(s string) string {
        output := s
        strings.ReplaceAll(output, "|", "I")
        
        return output
    },
})
```

## Convert to another format

```go
subs, err := pgs.ConvertToSubtitles("input.sup", nil)

if err != nil {
    panic(err)
}

f, err := os.Create("output.srt")

if err != nil {
    panic(err)
}

err = subs.WriteToSRT(f)
```

## Docker

Example of a Dockerfile to run `go-pgs-converter`

```dockerfile
FROM golang:1.20 as builder

WORKDIR /app

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

COPY go.* ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o YOUR_PROJECT_NAME_HERE

FROM debian:11.6

RUN apt-get update -qq
RUN apt-get install -y -qq ca-certificates tesseract-ocr tesseract-ocr-eng tesseract-ocr-fra ## Add as many languages you need 

COPY --from=builder /app/YOUR_PROJECT_NAME_HERE /YOUR_PROJECT_NAME_HERE
CMD ["/YOUR_PROJECT_NAME_HERE"]
```