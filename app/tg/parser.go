package tg

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

type Parser interface {
	ParseFile(filename string, messagesChan chan<- ArchivedMessage) error
}

type TgArchivedHTMLParser struct {
	log *zap.Logger
}

func NewTgArchivedHTMLParser(log *zap.Logger) *TgArchivedHTMLParser {
	return &TgArchivedHTMLParser{
		log: log,
	}
}

func (p *TgArchivedHTMLParser) ParseFile(filename string, messagesChan chan<- ArchivedMessage) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}

	doc.Find("div.message.default").Each(func(i int, s *goquery.Selection) {
		id, exists := s.Attr("id")
		if !exists {
			return
		}

		dateStr, exists := s.Find("div.pull_right.date.details").Attr("title")
		if !exists {
			return
		}
		// "21.11.2024 19:20:37 UTC+03:00"
		parts := strings.Split(dateStr, " ")
		if len(parts) < 3 {
			return
		}

		dateTimeStr := parts[0] + " " + parts[1] // "21.11.2024 19:20:37"
		parsedTime, err := time.Parse("02.01.2006 15:04:05", dateTimeStr)
		if err != nil {
			return
		}

		tzStr := parts[2] // "UTC+03:00"
		offset, err := parseTZOffset(tzStr)
		if err != nil {
			return
		}
		loc := time.FixedZone(tzStr, offset)

		datetime := time.Date(
			parsedTime.Year(), parsedTime.Month(), parsedTime.Day(),
			parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(),
			0, loc,
		)

		var tags []string
		s.Find("div.text a").Each(func(i int, a *goquery.Selection) {
			text := strings.TrimPrefix(a.Text(), "#")
			tags = append(tags, text)
		})

		messagesChan <- ArchivedMessage{
			MessageID: id,
			Datetime:  datetime,
			Tags:      tags,
		}
	})

	return nil
}

// parseTZOffset парсит строку часового пояса формата "UTC±HH:MM" и возвращает смещение в секундах.
func parseTZOffset(offsetStr string) (int, error) {
	re := regexp.MustCompile(`^UTC([+-])(\d{2}):(\d{2})$`)
	matches := re.FindStringSubmatch(offsetStr)
	if matches == nil {
		return 0, fmt.Errorf("invalid timezone format: %s", offsetStr)
	}

	sign := matches[1]
	hours, _ := strconv.Atoi(matches[2])
	minutes, _ := strconv.Atoi(matches[3])

	// Проверяем диапазон значений
	if hours > 14 || (hours == 14 && minutes > 0) || minutes >= 60 {
		return 0, fmt.Errorf("invalid timezone values: %s", offsetStr)
	}

	totalSeconds := hours*3600 + minutes*60
	if sign == "-" {
		totalSeconds = -totalSeconds
	}

	return totalSeconds, nil
}
