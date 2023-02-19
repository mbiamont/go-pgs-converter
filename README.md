# go-pgs-converter
Go PGS subtitle converter

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