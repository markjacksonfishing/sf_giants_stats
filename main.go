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
	"sync"

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

	headers, data := extractData(table)

	// Check that the data slice is not empty
	if len(data) == 0 {
		fmt.Fprintln(output, "No data found in HTML")
		return
	}

	// Create a dense matrix from the data slice
	if len(headers) < 2 {
		fmt.Fprintln(output, "Not enough headers found in HTML")
		return
	}

	X, Y := createMatrix(headers, data)

	// Compute the means of X and Y
	meanX, meanY := computeMeans(X, Y)

	// Compute the variances and covariance of X and Y
	varX, varY, covXY := computeVariancesAndCovariance(X, Y, meanX, meanY)

	// Compute the regression coefficients
	beta, alpha := computeRegressionCoefficients(varX, covXY, meanX, meanY)

	// Use the regression coefficients to predict the team's wins
	wins := predictWins(headers, beta, alpha)

	// Print the predicted wins for the upcoming season
	fmt.Fprintf(output, "\nPredicted wins for the upcoming season: %.2f\n", wins)
}

func extractData(table *goquery.Selection) ([]string, [][]float64) {
	headers := make([]string, 0)
	data := make([][]float64, 0)

	table.Find("thead tr th").Each(func(i int, header *goquery.Selection) {
		headers = append(headers, strings.TrimSpace(header.Text()))
	})

	var wg sync.WaitGroup
	table.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		wg.Add(1)
		go func(row *goquery.Selection) {
			defer wg.Done()
			vals := make([]float64, len(headers))
			row.Find("td").Each(func(j int, cell *goquery.Selection) {
				if j == 0 {
					if f, err := strconv.ParseFloat(strings.TrimSpace(cell.Text()), 64); err == nil {
						vals[0] = f
					}
				} else {
					if f, err := strconv.ParseFloat(strings.TrimSpace(cell.Text()), 64); err == nil {
						vals[j] = f
					}
				}
			})
			data = append(data, vals)
		}(row)
	})
	wg.Wait()

	return headers, data
}

func createMatrix(headers []string, data [][]float64) (*mat.Dense, *mat.Dense) {
	X := mat.NewDense(len(data), 1, nil)
	Y := mat.NewDense(len(data), 1, nil)

	for i, row := range data {
		for j, val := range row {
			if headers[j] == "W" {
				Y.Set(i, 0, val)
			} else {
				X.Set(i, 0, val)
			}
		}
	}

	return X, Y
}

func computeMeans(X, Y mat.Matrix) (float64, float64) {
	var meanX, meanY float64

	r, _ := X.Dims()
	for i := 0; i < r; i++ {
		meanX += X.At(i, 0)
		meanY += Y.At(i, 0)
	}

	meanX /= float64(r)
	meanY /= float64(r)

	return meanX, meanY
}
