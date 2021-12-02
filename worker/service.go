package worker

import (
	"bufio"
	"io"
	"log"
	"os"
	"parsers/common"
	"parsers/database"
	"regexp"
	"time"

	"fmt"
	"net/http"
	"strings"
)

var currentSitemap = ""
var sitemapUrl = "https://chrome.google.com/webstore/sitemap"
var applicationLangPattern = regexp.MustCompile("hreflang=\"(.*?)\"")

type AppRequest struct {
	Id           string
	LastModified string
	Langs        []string
	Exist        bool
}

func downloadSitemapFile(name string) bool {
	if _, err := os.Stat(common.GetAbsPath("data/sitemap" + name + ".xml")); !os.IsNotExist(err) {
		return true
	}

	out, _ := os.Create(common.GetAbsPath("data/sitemap" + name + ".xml"))
	defer out.Close()

	resp, _ := http.Get(sitemapUrl + "?" + name)
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		fmt.Errorf("bad status: %s", resp.Status)
		return false
	}
	io.Copy(out, resp.Body)
	return true
}

// Parse current sitemap file with extensions
func parseSitemap(name string) {
	file, err := os.Open(common.GetAbsPath("data/sitemap" + name + ".xml"))
	if err != nil {
		fmt.Errorf("Invalid sitemap", err)
		return
	}
	defer file.Close()
	currentSitemap = name

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	applicactionId := ""
	applicactionLastModified := ""
	applicationLangs := []string{}
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(s, "</url>") {
			if len(applicationLangs) == 0 {
				applicationLangs = []string{"en"}
			}
			addApplication(applicactionId, applicactionLastModified, applicationLangs)
			applicactionId = ""
			applicactionLastModified = ""
			applicationLangs = []string{}
			continue
		}

		if strings.HasPrefix(s, "<loc>") {
			s = replaceTags(s, "loc")
			applicactionId = s[len(s)-32:]
			continue
		}

		if strings.HasPrefix(s, "<lastmod>") {
			applicactionLastModified = replaceTags(s, "lastmod")
			continue
		}

		if strings.HasPrefix(s, "<xhtml:link") {
			ls := applicationLangPattern.FindStringSubmatch(s)
			if ls != nil && len(ls) > 1 {
				applicationLangs = append(applicationLangs, ls[1])
			}
		}
	}
}

func addApplication(id string, lastModified string, langs []string) {
	var exist uint8
	database.Connection.Raw("select count(id) from `applications` WHERE id=?", id).Row().Scan(&exist)
	if exist != 0 {
		return
	}
	time.Sleep(300 * time.Millisecond)
	app := &AppRequest{id, lastModified, langs, exist != 0}
	getRawApplication(app)
}

func replaceTags(s string, tag string) string {
	s = strings.Replace(s, "<"+tag+">", "", -1)
	s = strings.Replace(s, "</"+tag+">", "", -1)
	return s
}

func getApplicationsFromSitemap() {
	file, err := os.Open(common.GetAbsPath("data/sitemap.xml"))
	if err != nil {
		fmt.Errorf("Invalid sitemap", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	i := 0
	findSitemap := currentSitemap
	currentSitemap = ""

	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if isLoc := strings.HasPrefix(s, "<loc>"); !isLoc {
			continue
		}

		if isLang := strings.Contains(s, "hl="); isLang {
			break
		}

		name := strings.Replace(s, "<loc>"+sitemapUrl+"?", "", -1)
		name = strings.Replace(name, "</loc>", "", -1)
		name = strings.Replace(name, "&amp;", "&", -1)
		i++

		if sitemapDownloaded := downloadSitemapFile(name); !sitemapDownloaded {
			continue
		}

		if len(findSitemap) > 0 {
			if findSitemap != name {
				continue
			} else {
				// clear
				findSitemap = ""
			}
		}

		parseSitemap(name)
		println("Sitemap ", name, " was parsed")
		time.Sleep(10 * time.Second)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func Start() {

	// download first file with sitemaps list
	downloadSitemapFile("")

	database.Connection.DB().SetMaxIdleConns(50)
	database.Connection.DB().SetMaxOpenConns(50)
	database.Connection.DB().SetConnMaxLifetime(5 * time.Minute)

	getApplicationsFromSitemap()
}
