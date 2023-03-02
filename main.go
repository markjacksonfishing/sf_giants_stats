package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"gonum.org/v1/gonum/mat"
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
	var data [][]float64
	table.Find("thead tr th, tbody tr:not(.thead)").Each(func(i int, row *goquery.Selection) {
		cells := make([]string, 0)
		row.Find("th, td").Each(func(j int, cell *goquery.Selection) {
			cells = append(cells, strings.TrimSpace(cell.Text()))
		})
		if len(cells) > 0 {
			if i == 0 {
				headers = cells[1:]
			} else {
				vals := make([]float64, len(cells)-1)
				for k, v := range cells[1:] {
					if f, err := strconv.ParseFloat(v, 64); err == nil {
						vals[k] = f
					} else {
						// Skip non-numeric values
						continue
					}
				}
				data = append(data, vals)
			}
		}
	})

	// Check that the data slice is not empty
	if len(data) == 0 {
		fmt.Fprintln(output, "No data found in HTML")
		return
	}

	// Create a dense matrix from the data slice
	X := mat.NewDense(len(data), len(headers)-1, nil)
	Y := mat.NewVecDense(len(data), nil)
	for i, row := range data {
		if len(row) != len(headers) {
			log.Fatalf("Error parsing data: row %d has length %d instead of %d", i, len(row), len(headers))
		}
		for j, val := range row[1:] {
			X.Set(i, j, val)
		}
		Y.SetVec(i, row[0])
	}

	// Compute the means of X and Y
	meanX := mat.Sum(X.ColView(0)) / float64(X.RawMatrix().Rows)
	meanY := mat.Sum(Y) / float64(Y.Len())

	// Compute the variances and covariance of X and Y
	varX := 0.0
	varY := 0.0
	covXY := 0.0
	for i := 0; i < X.RawMatrix().Rows; i++ {
		x := X.At(i, 0)
		y := Y.AtVec(i)
		devX := x - meanX
		devY := y - meanY
		varX += devX * devX
		varY += devY * devY
		covXY += devX * devY
	}
	varX /= float64(X.RawMatrix().Rows)
	varY /= float64(Y.Len())
	covXY /= float64(X.RawMatrix().Rows)

	// Compute the regression coefficients
	beta := covXY / varX
	alpha := meanY - beta*meanX

	// Use the regression coefficients to predict the team's wins
	wins := 0.0
	for i := 0; i < len(headers)-1; i++ {
		wins += beta * X.At(0, i)
	}
	wins += alpha

	// Print the predicted wins for the upcoming season
	fmt.Fprintf(output, "\nPredicted wins for the upcoming season: %.2f\n", wins)
}
