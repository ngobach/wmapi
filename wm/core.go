package wm

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var client http.Client

func init() {
	client = http.Client{}
}

func urlFor(c *Country) string {
	if c != nil {
		return WebRoot + "/country/" + string(*c)
	} else {
		return WebRoot
	}
}

func getSeries(doc string, name string) []int {
	rgxp := regexp.MustCompile("(?s)series.*?'" + name + "'.*?\\[(.*?)\\]")
	match := rgxp.FindStringSubmatch(doc)
	line := match[1]
	arr := strings.Split(line, ",")
	var rs []int
	for _, v := range arr {
		value, _ := strconv.Atoi(strings.TrimSpace(v))
		rs = append(rs, value)
	}
	return rs
}

func getFirstDay(doc string) (*time.Time, error) {
	rgxp := regexp.MustCompile("categories:\\s*\\[\"(.*?)\"")
	m := rgxp.FindStringSubmatch(doc)
	value, err := time.Parse("2006 Jan 02", "2020 "+m[1])
	if err != nil {
		panic(err)
		return nil, err
	}
	return &value, nil
}

func mustGetUpdatedAt(doc *goquery.Document) time.Time {
	s := doc.Find(".label-counter + div").Text()
	s = strings.TrimSpace(s)
	s = s[14:]
	t, e := time.Parse("January 02, 2006, 15:04 GMT", s)
	if e != nil {
		panic(e)
	}
	return t
}

func fastCall(url string) (*string, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch document: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected error code %d", resp.StatusCode)
	}
	doc, err := func(raw []byte, err error) (string, error) {
		return string(raw), err
	}(ioutil.ReadAll(resp.Body))
	if err != nil {
		err = resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("unable to read response: %w", err)
	}
	return &doc, nil
}

func castToInt(s string) int {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", "")
	v, e := strconv.Atoi(s)
	if e != nil {
		panic(e)
	}
	return v
}

func buildDailyCases(start time.Time, totals []int, news []int, actives []int) []DailyStatistic {
	result := make([]DailyStatistic, len(totals))
	for i := 0; i < len(totals); i++ {
		result[i].Date = start.Add(time.Hour * time.Duration(24*i))
		result[i].Total = totals[i]
		result[i].Active = actives[i]
		result[i].New = news[i]
	}
	if len(result) > 30 {
		result = result[len(result)-30:]
	}
	return result
}

func GetStatistics(country *Country) (*Report, error) {
	response, err := fastCall(urlFor(country))
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(*response))
	if err != nil {
		return nil, err
	}
	docTitle := doc.Find("title").Text()
	if strings.Contains(docTitle, "404") {
		return nil, fmt.Errorf("document not found: %s", docTitle)
	}
	rp := Report{
		Country:   country,
		UpdatedAt: mustGetUpdatedAt(doc),
	}
	selection := doc.Find(".maincounter-number")
	rp.Total = castToInt(selection.Eq(0).Text())
	rp.Deaths = castToInt(selection.Eq(1).Text())
	rp.Recovered = castToInt(selection.Eq(2).Text())
	if country != nil {
		totalCases := getSeries(*response, "Cases")
		dailyNewCases := getSeries(*response, "Daily Cases")
		activeCases := getSeries(*response, "Currently Infected")
		firstDay, err := getFirstDay(*response)
		if err != nil {
			return nil, fmt.Errorf("unable to obtain starting Date: %w", err)
		}
		rp.Days = buildDailyCases(*firstDay, totalCases, dailyNewCases, activeCases)
	}
	return &rp, nil
}
