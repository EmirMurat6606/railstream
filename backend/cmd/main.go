package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func check(err error){
	if err != nil{
		panic(err)
	}
}

func main() {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
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
		log.Fatal("token niet gevonden")
	}

	fmt.Println(token)

	tokenString := "__RequestVerificationToken=" + token + "&Origin.Name=Temse&Origin.Icon=nmbs-logo&Origin.ExtId=8894672&Origin.IsBelgian=true&Origin.Reference=A%3D1%40O%3DTemse%40X%3D4221352%40Y%3D51126086%40U%3D80%40L%3D8894672%40B%3D1%40p%3D1776635490%40&Destination.Name=Sint-Niklaas&Destination.Icon=nmbs-logo&Destination.ExtId=8894508&Destination.IsBelgian=true&Destination.Reference=A%3D1%40O%3DSint-Niklaas%40X%3D4142966%40Y%3D51171472%40U%3D80%40L%3D8894508%40B%3D1%40p%3D1776549195%40&DatePicker=20%2F04%2F2026&TimePicker=1950&BoardType=DepartureBoard&Language=Dutch&ExtraOptions=%5B%7B%22name%22%3A%22Minimum+transfer+time%22%2C%22value%22%3A%220%22%7D%2C%7B%22name%22%3A%22First+mile%22%2C%22value%22%3A%22By+foot%22%7D%2C%7B%22name%22%3A%22Last+mile%22%2C%22value%22%3A%22By+foot%22%7D%2C%7B%22name%22%3A%22Transport+means%22%2C%22value%22%3A%22Train%22%7D%2C%7B%22name%22%3A%22Slower+trains%22%2C%22value%22%3A%22exclude%22%7D%5D&SaveExtraOptions=False&IsInternationalTrip=False"

	var data = strings.NewReader(tokenString)
	req, err = http.NewRequest("POST", "https://www.belgiantrain.be/api/routeplanner/GetJourneySearchResult/", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "nl-BE,nl-NL;q=0.9,nl;q=0.8,en-US;q=0.7,en;q=0.6")
	req.Header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("origin", "https://www.belgiantrain.be")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://www.belgiantrain.be/nl/")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="147", "Not.A/Brand";v="8", "Chromium";v="147"`)
	req.Header.Set("sec-ch-ua-arch", `"x86"`)
	req.Header.Set("sec-ch-ua-bitness", `"64"`)
	req.Header.Set("sec-ch-ua-full-version-list", `"Google Chrome";v="147.0.7727.102", "Not.A/Brand";v="8.0.0.0", "Chromium";v="147.0.7727.102"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-model", `""`)
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-ch-ua-platform-version", `"19.0.0"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/147.0.0.0 Safari/537.36")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	addSimpleCookies(req, `BRailWebLang=NL; BrailWidgetsStationsHistory=[[],[{"uic":"8894672","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"TEMSE","value":"4672","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null}]]; BrailDelayCertificateStationsHistory=[[{"uic":"8894672","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"TEMSE","value":"4672","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null}]]; belgiantrain#lang=nl; shell#lang=en; ASP.NET_SessionId=1mlnc3ofdghy4finfoa2mdw5; BrailRPStationsHistory=[[{"uic":"8894672","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Temse","value":"A%3d1%40O%3dTemse%40X%3d4221352%40Y%3d51126086%40U%3d80%40L%3d8894672%40B%3d1%40p%3d1776635490%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8822004","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Mechelen","value":"A%3d1%40O%3dMechelen%40X%3d4482786%40Y%3d51017649%40U%3d80%40L%3d8822004%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8894508","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Sint-Niklaas","value":"A%3d1%40O%3dSint-Niklaas%40X%3d4142966%40Y%3d51171472%40U%3d80%40L%3d8894508%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8894672","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Temse","value":"A%3d1%40O%3dTemse%40X%3d4221352%40Y%3d51126086%40U%3d80%40L%3d8894672%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8894672","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Temse","value":"A%3d1%40O%3dTemse%40X%3d4221352%40Y%3d51126086%40U%3d80%40L%3d8894672%40B%3d1%40p%3d1775684432%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null}],[{"uic":"8894508","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Sint-Niklaas","value":"A%3d1%40O%3dSint-Niklaas%40X%3d4142966%40Y%3d51171472%40U%3d80%40L%3d8894508%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8833001","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Leuven","value":"A%3d1%40O%3dLeuven%40X%3d4715867%40Y%3d50882280%40U%3d80%40L%3d8833001%40B%3d1%40p%3d1776635490%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8822004","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Mechelen","value":"A%3d1%40O%3dMechelen%40X%3d4482786%40Y%3d51017649%40U%3d80%40L%3d8822004%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8822715","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Puurs","value":"A%3d1%40O%3dPuurs%40X%3d4282704%40Y%3d51077220%40U%3d80%40L%3d8822715%40B%3d1%40p%3d1776549195%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null},{"uic":"8821121","isbelgian":true,"isbusstop":false,"lat":0.0,"long":0.0,"text":"Antwerpen-Berchem","value":"A%3d1%40O%3dAntwerpen-Berchem%40X%3d4432221%40Y%3d51199231%40U%3d80%40L%3d8821121%40B%3d1%40p%3d1775684432%40","icon":"nmbs-logo","multistation":false,"stringdifference":0,"stationnameformobile":null}]]`)

	resp, err = client.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(bodyBytes))
}

func addSimpleCookies(req *http.Request, cookieHeader string) {
	parts := strings.Split(cookieHeader, ";")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		name := kv[0]
		value := kv[1]

		// ❌ skip complexe cookies
		if strings.Contains(value, "{") ||
			strings.Contains(value, "[") ||
			strings.Contains(value, "\"") {
			continue
		}

		req.AddCookie(&http.Cookie{
			Name:  name,
			Value: value,
		})
	}
}