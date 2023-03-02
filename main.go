package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	run(os.Stdin, os.Stdout)
}

func run(input io.Reader, output io.Writer) {
	reader := bufio.NewReader(input)

	// Prompt the user to enter the team's abbreviation
	fmt.Fprint(output, "Enter the team's abbreviation (e.g. SF for San Francisco Giants): ")
	teamAbbr, _ := reader.ReadString('\n')
	teamAbbr = strings.TrimSpace(teamAbbr)

	// Construct the URL based on the team's abbreviation
	url := fmt.Sprintf("https://www.baseball-reference.com/teams/%s/2022.shtml", teamAbbr)

	// Make the HTTP request and parse the HTML response
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error making HTTP request: %v", err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML response: %v", err)
	}

	// Find the table and extract the headers and rows
	table := doc.Find("#team_batting").First()
	if table.Length() == 0 {
		fmt.Fprintln(output, "Table not found in HTML")
		return
	}

	headers := make([]string, 0)
	rows := make([][]string, 0)
	table.Find("thead tr th, tbody tr:not(.thead) td").Each(func(i int, header *goquery.Selection) {
		headers = append(headers, strings.TrimSpace(header.Text()))
	})
	table.Find("tbody tr:not(.thead)").Each(func(i int, row *goquery.Selection) {
		cells := make([]string, 0)
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			cells = append(cells, strings.TrimSpace(cell.Text()))
		})
		rows = append(rows, cells)
	})

	// Print the headers and the team's stats for the 2022 season
	fmt.Fprintln(output, strings.Join(headers, "\t"))
	for _, row := range rows {
		if len(row) > 0 && row[0] == "2022" {
			fmt.Fprintln(output, strings.Join(row, "\t"))
			break
		}
	}
}
