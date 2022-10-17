package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/anaskhan96/soup"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func getGameNerdzDeal() (string, string) {
	resp, err := soup.Get("https://www.gamenerdz.com/deal-of-the-day")
	if err != nil {
		log.Fatalf("failed to get gamenerdz info, %v", err)
	}
	doc := soup.HTMLParse(resp)

	root := doc.FindAll("article", "class", "card")
	tmp := root[:0]
	for _, e := range root {
		if !strings.Contains(e.FullText(), "Crisis") {
			tmp = append(tmp, e)
		}
	}
	card := tmp[0]
	log.Println(card)
	title := card.Find("h4", "class", "card-title").Find("a").Text()
	link := card.Find("h4", "class", "card-title").Find("a").Attrs()["href"]
	price := card.Find("span", "class", "price--withoutTax").Text()
	log.Println(title)
	log.Println(link)
	log.Println(price)

	showPrice := true
	if strings.Contains(title, "see price") {
		showPrice = false
	}

	title = strings.Replace(title, "(Add to cart to see price)", "", -1)
	title = strings.Replace(title, "(Deal of the Day)", "", -1)

	oldtitle := title
	if showPrice {
		title = fmt.Sprintf("[GN] [DOTD] - %v - %v", title, price)
	} else {
		title = fmt.Sprintf("[GN] [DOTD] - %v", title)
	}
	title += "\n\n" + link
	return oldtitle, title
}

func getCardHausDeal() (string, string) {
	resp, err := soup.Get("https://www.cardhaus.com/")
	if err != nil {
		log.Fatalf("failed to get cardhaus info: %v\n", err)
	}

	doc := soup.HTMLParse(resp)
	img := doc.Find("img", "title", "daily-deal-generic-rectangle.png")
	link := img.Pointer.Parent.Attr[0].Val

	resp, err = soup.Get(link)
	if err != nil {
		log.Fatalf("failed to get cardhaus game info: %v\n", err)
	}

	doc = soup.HTMLParse(resp)
	title := doc.Find("h1", "class", "productView-title").Text()
	price := doc.Find("span", "class", "price--withoutTax").Text()

	data := fmt.Sprintf("[Cardhaus] [DOTD] - %v - %v\n\n%v", title, price, link)

	return title, data
}

func getTabletopMerchantDotd() (string, string) {
	resp, err := soup.Get("https://tabletopmerchant.com/collections/deal-of-the-day")
	if err != nil {
		log.Fatalf("failed to get tabletop merchant data: %v", err)
	}

	doc := soup.HTMLParse(resp)
	card := doc.Find("a", "class", "product-item__title")
	title := card.Text()
	link := card.Attrs()["href"]
	link = "https://tabletopmerchant.com" + link
	log.Println(title)
	log.Println(link)

	variantId := doc.Find("form", "class", "button-stack").Find("input", "name", "id").Attrs()["value"]
	print(variantId)
	price := getPrice(variantId)
	title = strings.Replace(title, "(DEAL OF THE DAY)", "", -1)
	title = strings.Replace(title, "(SEE LOW PRICE AT CHECKOUT)", "", -1)

	var data string
	if price == 0 {
		data = fmt.Sprintf("[Tabletop Merchant] [DoTD] - %v - %v\n\n%v", title, " PRICE UNAVAILABLE ", link)
	} else {
		data = fmt.Sprintf("[Tabletop Merchant] [DoTD] - %v - %v\n\n%v", title, price, link)
	}
	return title, data
}

func getPrice(id string) float64 {

	cookies := []http.Cookie{{Name: "keep_alive", Value: "16c1d5e7-7d42-48d7-b356-b44ea7abf8a6"},
		{Name: "secure_customer_sig", Value: ""},
		{Name: "localization", Value: "US"},
		{Name: "cart_currency", Value: "USD"},
		{Name: "_orig_referrer", Value: ""},
		{Name: "_landing_page", Value: "%2F"},
		{Name: "_y", Value: "89d13d6b-98c8-402f-b88b-d450bad6e894"},
		{Name: "_s", Value: "16c1d5e7-7d42-48d7-b356-b44ea7abf8a6"},
		{Name: "_shopify_y", Value: "89d13d6b-98c8-402f-b88b-d450bad6e894"},
		{Name: "_shopify_s", Value: "16c1d5e7-7d42-48d7-b356-b44ea7abf8a6"},
		{Name: "_shopify_sa_p", Value: ""},
		{Name: "shopify_pay_redirect", Value: "pending"},
		{Name: "swym-session-id", Value: "97dqo9rgjet0nrk5aoh9ws7owsuxeg12wcgm4mwf1exa9mm031hh4ukits9wn2jj"},
		{Name: "swym-pid", Value: "/9c/ZJJ3ILSOjuo6VuLmhIBAdJyXBzKEliHleA8T0I0="},
		{Name: "swym-o_s", Value: "true"},
		{Name: "swym-swymRegid", Value: "EQ2CI9WG-fL5CSuS7QUOKObL2WQbIPxSDxIYToqmlSRroMsCVzPpt97T5ooLP-dq7Z0wPZIHpWZH2rEM-35QPEB7-ay9DOmp826vyxsCj_JXu0FFPQIcLNr13cvmSlhaI3tkwLg1cmtCgiy5M5DIhw8k6-H01t3zihxyW5rgV_k"},
		{Name: "swym-email", Value: "null"},
		{Name: "swym-cu_ct", Value: "undefined"},
		{Name: "_shopify_sa_t", Value: "2022-08-12T19%3A19%3A28.172Z"},
		{Name: "swym-instrumentMap", Value: "{}"}}

	headers := map[string]string{
		"authority":       "tabletopmerchant.com",
		"accept":          "*/*",
		"accept-language": "en-US,en;q=0.9",
		"cache-control":   "no-cache",
		//# Already added when you pass json=
		"content-type": "application/json; charset=UTF-8	",
		//# Requests sorts cookies= alphabetically
		//# "cookie": "keep_alive=16c1d5e7-7d42-48d7-b356-b44ea7abf8a6; secure_customer_sig=; localization=US; cart_currency=USD; _orig_referrer=; _landing_page=%2F; _y=89d13d6b-98c8-402f-b88b-d450bad6e894; _s=16c1d5e7-7d42-48d7-b356-b44ea7abf8a6; _shopify_y=89d13d6b-98c8-402f-b88b-d450bad6e894; _shopify_s=16c1d5e7-7d42-48d7-b356-b44ea7abf8a6; _shopify_sa_p=; shopify_pay_redirect=pending; swym-session-id="97dqo9rgjet0nrk5aoh9ws7owsuxeg12wcgm4mwf1exa9mm031hh4ukits9wn2jj"; swym-pid="/9c/ZJJ3ILSOjuo6VuLmhIBAdJyXBzKEliHleA8T0I0="; swym-o_s=true; swym-swymRegid="EQ2CI9WG-fL5CSuS7QUOKObL2WQbIPxSDxIYToqmlSRroMsCVzPpt97T5ooLP-dq7Z0wPZIHpWZH2rEM-35QPEB7-ay9DOmp826vyxsCj_JXu0FFPQIcLNr13cvmSlhaI3tkwLg1cmtCgiy5M5DIhw8k6-H01t3zihxyW5rgV_k"; swym-email=null; swym-cu_ct=undefined; _shopify_sa_t=2022-08-12T19%3A19%3A28.172Z; swym-instrumentMap={}",
		"origin":             "https://tabletopmerchant.com",
		"pragma":             "no-cache",
		"referer":            "https://tabletopmerchant.com/collections/featured-products/products/parks-nightfall",
		"sec-ch-ua":          "\"Chromium\";v=\"104\", \" Not A;Brand\";v=\"99\", \"Microsoft Edge\";v=\"104\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "Windows",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.81 Safari/537.36 Edg/104.0.1293.47",
		"x-requested-with":   "XMLHttpRequest",
	}

	var jsonData = []byte(`{"form_type": "product", "utf8": "✓", "id": ` + id + `, "quantity": "1"}`)

	//postdata := strings.NewReader("")

	req, err := http.NewRequest("POST", "https://tabletopmerchant.com/cart/add.js", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	for _, cookie := range cookies {
		req.AddCookie(&cookie)
	}

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to get price: %v", err)
	}

	if resp.StatusCode != 200 {
		log.Println(resp.StatusCode)
		log.Println(resp.Status)
		//log.Fatalf("Error response from tabletop merchant")
		return 0
	}

	// Read response
	data, err := io.ReadAll(resp.Body)

	// error handle
	if err != nil {
		log.Fatalf("error reading data: %v", err)
	}

	// Print response
	fmt.Printf("Response = %s", string(data))

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatalf("failed to parse json: %v", err)
	}

	return result["final_price"].(float64) / 100
}

func diff(site, title string) bool {
	dat, err := os.ReadFile("dotdData/" + site + ".dat")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalln(err)
		}
	}
	if title == string(dat) {
		return false
	}
	dat = []byte(title)
	err = os.WriteFile("dotdData/"+site+".dat", dat, 0666)
	checkErr(err)
	return true

}
