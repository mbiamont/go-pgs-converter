package go_pgs_converter

import (
	"bytes"
	sub "github.com/asticode/go-astisub"
	"github.com/disintegration/imaging"
	"github.com/mbiamont/go-pgs-parser/displaySet"
	"github.com/mbiamont/go-pgs-parser/pgs"
	"github.com/otiai10/gosseract/v2"
	"image"
	"image/color"
	"image/png"
	"strings"
	"time"
)

const defaultMaxSubtitleDuration = (time.Second * 2) + (time.Millisecond * 500)

type ConversionOptions struct {
	//Language of the PGS input (ISO 639-2)
	InputLanguage string

	//Maximum duration for a subtitle to appear on screen (default is 2.5 seconds)
	MaxSubtitleDuration time.Duration

	//Text correction function to reprocess the text after OCR
	TextCorrection func(string) string
}

func ConvertToSubtitles(inputFilePath string, options *ConversionOptions) (*sub.Subtitles, error) {
	parser := pgs.NewPgsParser()
	ocr := gosseract.NewClient()
	defer ocr.Close()

	if options != nil && len(options.InputLanguage) != 0 {
		err := ocr.SetLanguage(options.InputLanguage)

		if err != nil {
			return nil, err
		}
	}

	subtitles := sub.NewSubtitles()

	err := parser.ParsePgsFile(inputFilePath, func(index int, startTime time.Duration, data displaySet.ImageData) error {
		processedImage := processImageForOcr(data.Image)
		buf := new(bytes.Buffer)

		err := png.Encode(buf, processedImage)

		if err != nil {
			return err
		}

		err = ocr.SetImageFromBytes(buf.Bytes())

		if err != nil {
			return err
		}

		text, err := ocr.Text()

		if err != nil {
			return err
		}

		if options != nil && options.TextCorrection != nil {
			text = options.TextCorrection(text)
		}

		maxSubtitleDuration := defaultMaxSubtitleDuration

		if options != nil && options.MaxSubtitleDuration != 0 {
			maxSubtitleDuration = options.MaxSubtitleDuration
		}

		appendSubtitle(subtitles, text, startTime, index, maxSubtitleDuration)
		return nil
	})

	return subtitles, err
}

func appendSubtitle(subtitle *sub.Subtitles, text string, startTime time.Duration, index int, maxSubtitleDuration time.Duration) {
	text = strings.TrimSpace(text)
	if len(text) == 0 {
		return
	}

	item := &sub.Item{}
	item.Index = index
	item.StartAt = startTime
	item.EndAt = startTime + maxSubtitleDuration

	item.Lines = append(item.Lines, sub.Line{Items: []sub.LineItem{{Text: text}}})

	if len(subtitle.Items) > 0 {
		previousEndAt := subtitle.Items[len(subtitle.Items)-1].EndAt

		if previousEndAt > item.EndAt {
			subtitle.Items[len(subtitle.Items)-1].EndAt = item.StartAt
		}
	}

	subtitle.Items = append(subtitle.Items, item)
}

func processImageForOcr(img image.Image) image.Image {
	img = imaging.AdjustBrightness(img, -100)
	img = imaging.AdjustContrast(img, -100)

	dst := imaging.New(img.Bounds().Dx()+100, img.Bounds().Dy()+100, color.Transparent)
	dst = imaging.Paste(dst, img, image.Pt(50, 50))

	return dst
}
