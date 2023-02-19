FROM golang:1.20 as builder

WORKDIR /app

RUN apt-get update -qq
RUN apt-get install -y -qq libtesseract-dev libleptonica-dev

COPY go.* ./
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o go-pgs-converter

FROM debian:11.6

RUN apt-get update -qq
RUN apt-get install -y -qq tesseract-ocr tesseract-ocr-eng tesseract-ocr-deu tesseract-ocr-fra tesseract-ocr-jpn ca-certificates
COPY --from=builder app/sample/ /sample/
COPY --from=builder /app/go-pgs-converter /go-pgs-converter
CMD ["/go-pgs-converter"]