package main

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

var baseURL = "https://www.imdb.com"

func main() {
	movieID := getMovieID("titanic")
	getAwardsPage(movieID)
}

func handleErr(err error) {

	if err != nil {
		fmt.Println(err)
	}

}

func getMovieID(movieName string) string {
	searchURL := baseURL + "/find?q="
	resp, err := soup.Get(searchURL + movieName)
	handleErr(err)
	rootNode := soup.HTMLParse(resp)
	links := rootNode.FindAll("a")
	var movieID string

	for _, link := range links {

		if strings.HasPrefix(link.Attrs()["href"], "/title") {
			movieID = strings.Split(link.Attrs()["href"], "/?ref")[0]
			break
		}

	}

	return movieID
}

func getAwardsPage(movieID string) soup.Root {
	awardsURL := baseURL + movieID + "/awards"
	resp, err := soup.Get(awardsURL)
	handleErr(err)
	rootNode := soup.HTMLParse(resp)
	return rootNode
}

func getEvents(awardsPage soup.Root) {
	events := awardsPage.FindAll("table", "class", "awards")

	for _, event := range events {

		for _, award := range event.FindAll("td", "class", "award_description") {

			awardName := strings.TrimSpace(award.Text())
			personName := award.Find("a").Text()

			fmt.Printf("%s - %s\n", personName, awardName)

		}

	}

}

func getCategories() {

}

func getAwards() {

}

func getNominations() {

}
