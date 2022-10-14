package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Items struct {
	XMLName xml.Name `xml:"items"`
	Items   []Item   `xml:"item"`
	Total   string   `xml:"total,attr"`
}

type Item struct {
	XMLName xml.Name `xml:"item"`
	//Type    string   `xml:"type,attr"`
	Id   string `xml:"id,attr"`
	Name Name   `xml:"name"`
}

type Name struct {
	Name string `xml:"value,attr"`
}

func searchBGG(name string) string {
	name = strings.Replace(strings.ToLower(name), "expansion", "", -1)
	name = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(name, "")
	name = strings.Replace(name, "  ", " ", -1)
	name = strings.Replace(name, " ", "+", -1)
	resp, err := http.Get("https://boardgamegeek.com/xmlapi2/search?query=" + name + "&type=boardgame")
	checkErr(err)
	bytes, err := io.ReadAll(resp.Body)
	checkErr(err)
	var items Items
	err = xml.Unmarshal(bytes, &items)
	checkErr(err)
	if i, _ := strconv.Atoi(items.Total); i < 1 {
		return "no bgg link"
	}
	return "https://boardgamegeek.com/boardgame/" + items.Items[0].Id
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
