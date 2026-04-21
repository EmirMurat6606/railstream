// Package gtfs is responsible for fetching realtime train data from the Belgian Mobility API
package gtfs

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/EmirMurat6606/railstream/internal/mqtt"

	scheduler "github.com/EmirMurat6606/railstream/internal/scheduler"

	"github.com/PuerkitoBio/goquery"
)

// StartRealtimeFetcher fetches the realtime NMBS api periodically
//
// The realtime data is fetched every 30 seconds from 6:30 to 12:30 every day
func StartRealtimeFetcher(ctx context.Context, publisher *mqtt.Publisher) error {

	loc, err := time.LoadLocation("Europe/Brussels")

	if err != nil {
		return err
	}

	window := scheduler.Window{
		StartHour:   6,
		StartMinute: 30,
		EndHour:     12,
		EndMinute:   30,
		Location:    loc,
	}

	fromName, toName := "Temse", "Sint-Niklaas"

	task := func(ctx context.Context) error {
		return fetchJourney(ctx, fromName, toName, publisher)
	}

	return scheduler.RunPeriodic(ctx, 30*time.Second, window, task)
}

func fetchJourney(ctx context.Context, fromFrName, toFrName string, publisher *mqtt.Publisher) error {

	from, ok := getStation(fromFrName)
	if !ok {
		return fmt.Errorf("vertrekstation niet gevonden: %s", fromFrName)
	}

	to, ok := getStation(toFrName)
	if !ok {
		return fmt.Errorf("aankomststation niet gevonden: %s", toFrName)
	}

	loc, _ := time.LoadLocation("Europe/Brussels")
	now := time.Now().In(loc)

	// Date and time are layout, not actual data and time
	dateStr := now.Format("02/01/2006")
	timeStr := now.Format("1504")

	jar, _ := cookiejar.New(nil)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	client := &http.Client{
		Jar:       jar,
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// -------- GET homepage (for token) --------

	req, err := http.NewRequest("GET", "https://www.belgiantrain.be/nl", nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "nl-BE,nl-NL;q=0.9,nl;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`)
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-full-version-list", `"Google Chrome";v="147.0.7727.102", "Not.A/Brand";v="8.0.0.0", "Chromium";v="147.0.7727.102"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-platform-version", `"19.0.0"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(string(bodyBytes)),
	)

	if err != nil {
		log.Fatal(err)
	}

	token, exists := doc.Find(`input[name="__RequestVerificationToken"]`).Attr("value")
	if !exists {
		log.Fatal("token not found")
	}

	// -------- POST route planner --------

	extraOptions := "%5B%7B%22name%22%3A%22Minimum+transfer+time%22%2C%22value%22%3A%220%22%7D%2C%7B%22name%22%3A%22First+mile%22%2C%22value%22%3A%22By+foot%22%7D%2C%7B%22name%22%3A%22Last+mile%22%2C%22value%22%3A%22By+foot%22%7D%2C%7B%22name%22%3A%22Transport+means%22%2C%22value%22%3A%22Train%22%7D%2C%7B%22name%22%3A%22Slower+trains%22%2C%22value%22%3A%22exclude%22%7D%5D"

	postBody :=
		"__RequestVerificationToken=" + token +
			"&Origin.Name=" + url.QueryEscape(from.Name) + // "Temse"
			"&Origin.Icon=nmbs-logo" +
			"&Origin.ExtId=" + from.ExtID + // "8894672"
			"&Origin.IsBelgian=true" +
			"&Origin.Reference=" + from.encodedReference() + // URL-encoded Hafas string
			"&Destination.Name=" + url.QueryEscape(to.Name) +
			"&Destination.Icon=nmbs-logo" +
			"&Destination.ExtId=" + to.ExtID +
			"&Destination.IsBelgian=true" +
			"&Destination.Reference=" + to.encodedReference() +
			"&DatePicker=" + url.QueryEscape(dateStr) +
			"&TimePicker=" + timeStr +
			"&BoardType=DepartureBoard" +
			"&Language=Dutch" +
			"&ExtraOptions=" + extraOptions +
			"&SaveExtraOptions=False" +
			"&IsInternationalTrip=False"

	body := strings.NewReader(postBody)

	req2, err := http.NewRequest(
		"POST",
		"https://www.belgiantrain.be/api/routeplanner/GetJourneySearchResult/",
		body,
	)

	if err != nil {
		log.Fatal(err)
	}

	req2.Header.Set("accept", "*/*")
	req2.Header.Set("accept-language", "nl-BE,nl-NL;q=0.9,nl;q=0.8,en-US;q=0.7,en;q=0.6")
	req2.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req2.Header.Set("origin", "https://www.belgiantrain.be")
	req2.Header.Set("priority", "u=1, i")
	req2.Header.Set("referer", "https://www.belgiantrain.be/nl/travel-info/prepare-for-your-journey/use-the-sncb-app")
	req2.Header.Set("sec-ch-ua", `"Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`)
	req2.Header.Set("sec-ch-ua-arch", `"x86"`)
	req2.Header.Set("sec-ch-ua-bitness", `"64"`)
	req2.Header.Set("sec-ch-ua-full-version-list", `"Google Chrome";v="147.0.7727.102", "Not.A/Brand";v="8.0.0.0", "Chromium";v="147.0.7727.102"`)
	req2.Header.Set("sec-ch-ua-mobile", "?0")
	req2.Header.Set("sec-ch-ua-model", `""`)
	req2.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req2.Header.Set("sec-ch-ua-platform-version", `"19.0.0"`)
	req2.Header.Set("sec-fetch-dest", "empty")
	req2.Header.Set("sec-fetch-mode", "cors")
	req2.Header.Set("sec-fetch-site", "same-origin")
	req2.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36")
	req2.Header.Set("x-requested-with", "XMLHttpRequest")

	postResp, err := client.Do(req2)
	if err != nil {
		return err
	}

	defer postResp.Body.Close()

	resultBody, err := io.ReadAll(postResp.Body)
	if err != nil {
		return err
	}

	parseDelays(resultBody)

	return nil

}

func parseDelays(body []byte) {

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Println("parse error:", err)
		return
	}

	delay := strings.TrimSpace(
		doc.Find(".planner__delay").First().Text(),
	)

	if delay == "" {
		log.Println("No delay detected")
		return
	}

	log.Println("Delay detected:", delay)
}
