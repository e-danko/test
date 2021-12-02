package worker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"parsers/database"
	"strconv"
	"strings"
	"time"
)

type Application struct {
	Id              string
	Slug            string
	CompanyName     string
	LogoDetail      string
	LogoList        string
	CategoryCode    string
	CategoryName    string
	Rate            float32
	UsersRated      int32
	UsersDownloaded int32
	Website         string
	Link            string
	IsOffline       bool
	IsGDrive        bool `gorm:"column:is_gdrive"`
	IsGoogle        bool
	IsAndroid       bool
	IsFree          bool
	FaqPage         string
	Version         string
	LastUpdated     string
	CreatedAt       string
	LastModified    string
	Type            string
	Size            string
	ContactEmail    string
	ContactAddress  string
	PrivacyPolicy   string
	Price           string
	AndroidLink     string
	DefaultLang     string
	CountLangs      int
	CountReviews    int
}

type ApplicationContainer struct {
	App      *Application
	Request  *AppRequest
}

var lastUpdatedMask = "January 02, 2006"

func saveApplication(appContainer *ApplicationContainer) {
	database.Connection.Create(appContainer.App)
	fmt.Println("Application saved: " + appContainer.App.Id)
}

func getHTTPClient(ip string) (client *http.Client) {
	var transport http.RoundTripper
	if ip != "" {
		ipData := strings.Split(ip, ":")
		proxyUrl := &url.URL{
			Scheme: "http",
			Host: ipData[0] + ":" + ipData[1],
			User: url.UserPassword(ipData[2], ipData[3]),
		}
		transport = &http.Transport{
			DisableKeepAlives: true,
			Proxy:             http.ProxyURL(proxyUrl),
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, 6*time.Second)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(6 * time.Second))
				conn.SetWriteDeadline(time.Now().Add(6 * time.Second))
				conn.SetReadDeadline(time.Now().Add(6 * time.Second))
				return conn, nil
			},
		}
	} else {
		transport = &http.Transport{
			DisableKeepAlives:  true,
			DisableCompression: true,
		}
	}

	client = &http.Client{Transport: transport}
	return client
}

func retryApp(app *AppRequest) {
	return
}

func requestData(client *http.Client, urlLang string) (statusCode int, respBody []byte) {
	data := url.Values{
		"login": {""},
	}

	resp, err := client.PostForm(urlLang, data)
	if err != nil || resp.StatusCode != 200 {
		fmt.Printf("Errror %s", resp)
		if err == nil {
			statusCode = resp.StatusCode
			resp.Body.Close()
		}
		return
	}

	respBody, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func getRawApplication(app *AppRequest) {
	client := getHTTPClient("")
	defer client.CloseIdleConnections()
	println("Receiving Application info: " + app.Id)

	query := map[string]string{
		"gl":        "RU",
		"pv":        "20210820",
		"mce":       "atf,pii,rtr,rlb,gtc,hcn,svp,wtd,hap,nma,dpb,utb,hbh,ebo,c3d,ncr,hns,ctm,ac,hot,hsf,mac,epb,fcf,rma",
		"container": "CHROME",
		"id":        app.Id,
		"rt":        "j",
	}
	queryList := []string{}
	for k, v := range query {
		queryList = append(queryList, k+"="+v)
	}

	baseUrl := "https://chrome.google.com/webstore/ajax/detail?"
	var application *Application

	urlApp := baseUrl + strings.Join(queryList, "&")

	statusCode, respBody := requestData(client, urlApp)
	time.Sleep(500 * time.Millisecond)
	if statusCode == 404 {
		fmt.Printf("Application not found: %s\n", app.Id)
		return
	}

	if len(respBody) < 6 {
		fmt.Printf("Application request failed: %s\n", app.Id)
		retryApp(app)
		return
	}
	respBody = respBody[6:]

	var appList [][][]interface{}
	if err := json.Unmarshal(respBody, &appList); err != nil {
		fmt.Println("APP list JSON error: ", err)
		retryApp(app)
		return
	}

	if appList[0][1][0] != "getitemdetailresponse" {
		fmt.Printf("No valid response: %s\n", app.Id)
		retryApp(app)
		return
	}

	valuesAll := appList[0][1][1].([]interface{})
	valuesMain := valuesAll[0].([]interface{})
	application = &Application{
		Id:           fmt.Sprint(valuesMain[0]),
		Slug:         fmt.Sprint(valuesMain[61]),
		LogoDetail:   fmt.Sprint(valuesMain[25]),
		CategoryCode: fmt.Sprint(valuesMain[9]),
		CategoryName: fmt.Sprint(valuesMain[10]),
		CompanyName:  fmt.Sprint(valuesMain[2]),
		Price:        fmt.Sprint(valuesMain[30]),
		Website:      fmt.Sprint(valuesMain[35]),
		Link:         fmt.Sprint(valuesMain[37]),
		IsOffline:    (fmt.Sprint(valuesMain[53]) == "1"),
		IsGDrive:     (fmt.Sprint(valuesMain[54]) == "1"),
		IsGoogle:     (fmt.Sprint(valuesMain[56]) == "1"),
		IsFree:       (fmt.Sprint(valuesMain[75]) == "Free"),
		CreatedAt:    time.Now().Format(database.DateTimeFormat),
		FaqPage:      fmt.Sprint(valuesAll[5]),
		Version:      fmt.Sprint(valuesAll[6]),
		Type:         fmt.Sprint(valuesAll[10]),
		Size:         fmt.Sprint(valuesAll[25]),
	}

	if valuesMain[69] != nil {
		application.AndroidLink = fmt.Sprint(valuesMain[69])
		application.IsAndroid = true
	}

	if valuesMain[4] != nil {
		application.LogoList = fmt.Sprint(valuesMain[4])
	}

	if lastUpdated, err := time.Parse(lastUpdatedMask, fmt.Sprint(valuesAll[7])); err == nil {
		application.LastUpdated = lastUpdated.Format(database.DateFormat)
	} else {
		application.LastUpdated = time.Now().Format(database.DateFormat)
	}

	if lastModified, err := time.Parse(time.RFC3339Nano, app.LastModified); err == nil {
		application.LastModified = lastModified.Format(database.DateTimeFormat)
	} else {
		application.LastModified = time.Now().Format(database.DateTimeFormat)
	}

	if valuesAll[35] != nil {
		valuesContact := valuesAll[35].([]interface{})
		if l := len(valuesContact); l > 0 {
			application.ContactEmail = fmt.Sprint(valuesContact[0])
			if l > 1 {
				application.ContactAddress = fmt.Sprint(valuesContact[1])
				if l > 2 {
					application.PrivacyPolicy = fmt.Sprint(valuesContact[2])
				}
			}
		}
	}
	rate, _ := strconv.ParseFloat(fmt.Sprint(valuesMain[12]), 32)
	usersRated, _ := strconv.Atoi(fmt.Sprint(valuesMain[22]))
	usersDownloaded, _ := strconv.Atoi(fmt.Sprint(valuesMain[23]))
	application.Rate = float32(rate)
	application.UsersRated = int32(usersRated)

	usersDownloadedStr := fmt.Sprint(valuesMain[23])
	if strings.HasSuffix(usersDownloadedStr, "+") {
		usersDownloadedStr = strings.Replace(usersDownloadedStr, "+", "", 1)
		usersDownloadedStr = strings.Replace(usersDownloadedStr, ",", "", 5)
		d, _ := strconv.Atoi(usersDownloadedStr)
		application.UsersDownloaded = int32(d)
	} else if usersDownloadedStr != "" && usersDownloadedStr != "0" {
		d, _ := strconv.Atoi(usersDownloadedStr)
		application.UsersDownloaded = int32(d)
	} else {
		application.UsersDownloaded = int32(usersDownloaded)
	}

	saveApplication(&ApplicationContainer{
		App: application,
	})
	time.Sleep(1 * time.Second)
}
