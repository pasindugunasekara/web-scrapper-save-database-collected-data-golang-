package main

import (
	"database/sql"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gocolly/colly"
)

func getData(url *url.URL) {

	collection := colly.NewCollector()

	//////

	collection.OnHTML(".amount--3NTpl", func(e *colly.HTMLElement) {
		price := e.Attr("class")
		fmt.Println(e.Text)
		collection.Visit(e.Request.AbsoluteURL(price))
	})
	collection.OnHTML(".word-break--2nyVq", func(e *colly.HTMLElement) {
		information := e.Attr("class")
		fmt.Println(e.Text)
		collection.Visit(e.Request.AbsoluteURL(information))
	})

	collection.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	collection.Visit(url.String())

}

func main() {
	fmt.Print("Please enter here catagory Type ->>> ")
	var category string
	fmt.Scanln(&category)

	fmt.Print("Please enter here District ->>> ")
	var District string
	fmt.Scanln(&District)

	collectionn := colly.NewCollector()

	dataBase, _ := sql.Open("mysql", "root:ijse@tcp(127.0.0.1:3306)/ikman")

	var Contact string
	var Description string
	var model string
	var price string
	var des string

	collectionn.OnHTML(".gtm-normal-ad", func(element *colly.HTMLElement) {
		model := element.ChildText(".heading--2eONR")
		des := element.ChildText(".description--2-ez3")
		price := element.ChildText(".price--3SnqI")
		url := element.ChildAttr(".card-link--3ssYv", "href")
		fmt.Println("\n")
		fmt.Println("\tmodel : ", model)
		fmt.Println("\tprice : ", price)
		fmt.Println("\tdes : ", des)
		fmt.Println("\turl: ", url)

		//insert, err := dataBase.Query("INSERT INTO ikman (category, District, model, price, descr) VALUES (?, ?, ?, ?, ?)", category, District, model, price, des)
		e := element.Request.Visit(url)
		check(e)

	})

	//get contact
	collectionn.OnHTML(".contact-name--m97Sb", func(element *colly.HTMLElement) {
		Contact = element.Text
		fmt.Println("\tContact: ", Contact)
	})

	// collectionn.OnRequest(func(t *colly.Request) {
	// 	url := (t.URL)
	// 	getData(url)

	// })
	//get add full discription

	collectionn.OnHTML(".description-section--oR57b > div > .description--1nRbz", func(element *colly.HTMLElement) {
		Description = element.Text
		fmt.Println("\tfull Description: ", Description)
		fmt.Println("}")
	})

	collectionn.OnScraped(func(r *colly.Response) {
		insert, e := dataBase.Query("INSERT INTO ikman (District, category, model, price, des, Contact, Description) VALUES (?, ?, ?, ?, ?, ?, ?)",
			District, category, model, price, des, Contact, Description)
		check(e)
		defer insert.Close()
	})

	_ = collectionn.Visit("https://ikman.lk/en/ads/" + District + "/" + category)

}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
}
