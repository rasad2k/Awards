package main

import (
	"fmt"
	"strings"

	"github.com/anaskhan96/soup"
)

var baseURL = "https://www.imdb.com"

func main() {
	movieID := getMovieID("parasite")
	awardsPage := getAwardsPage(movieID)
	eventNames := getEvents(awardsPage)
	getAwards(awardsPage, eventNames)
}

func handleErr(err error) {

	if err != nil {
		fmt.Println(err)
	}

}

func getRootNode(URL string) soup.Root {
	resp, err := soup.Get(URL)
	handleErr(err)
	rootNode := soup.HTMLParse(resp)
	return rootNode
}

func getMovieID(movieName string) string {
	searchURL := baseURL + "/find?q="
	links := getRootNode(searchURL + movieName).FindAll("a")
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
	return getRootNode(awardsURL)
}

func getEvents(awardsPage soup.Root) []string {
	links := awardsPage.FindAll("a")
	var eventNames []string

	for _, link := range links {

		if strings.HasPrefix(link.Attrs()["href"], "/event/ev") {
			eventNames = append(eventNames, strings.TrimSpace(link.FindPrevSibling().NodeValue)+" "+link.Text())
		}

	}

	return eventNames
}

func getAwards(awardsPage soup.Root, eventNames []string) {
	eventTables := awardsPage.FindAll("table", "class", "awards")

	for i, table := range eventTables {
		fmt.Println(eventNames[i])
		var tdList []soup.Root
		tdList = append(tdList, table.FindAll("td")...)

		for _, td := range table.FindAll("td", "class", "title_award_outcome") {
			category := td.Children()[4].Text()
			outcome := td.Children()[1].Text()
			fmt.Printf("%s (%s):\n", category, outcome)

			check := false
			for _, tableData := range tdList {

				if tableData.Attrs()["class"] == "title_award_outcome" && tableData.Children()[1].Text() == outcome {
					check = true
					continue
				}

				if check {

					if tableData.Attrs()["class"] == "title_award_outcome" {
						break
					}

					fmt.Println(extractAward(tableData))
				}
			}

		}

		break
	}

}

func extractAward(td soup.Root) string {
	category := strings.TrimSpace(td.Text())
	var people []string

	for _, personLink := range td.FindAll("a") {
		person := strings.TrimSpace(personLink.Text())
		people = append(people, person)
	}

	if category == "" {
		return ""
	}

	if len(people) == 0 {
		return fmt.Sprintf("%s\n", category)
	}

	return fmt.Sprintf("%s: %s\n", category, strings.Join(people, ", "))

}
