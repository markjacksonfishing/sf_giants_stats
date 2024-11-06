package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
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
	varX, covXY := computeVariancesAndCovariance(X, Y, meanX, meanY)

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
					if cell.HasClass("left") {
						vals[j] = 1
					}
				} else {
					val, _ := strconv.ParseFloat(strings.TrimSpace(cell.Text()), 64)
					vals[j] = val
				}
			})
			data = append(data, vals)
		}(row)
	})
	wg.Wait()

	return headers, data
}

func createMatrix(headers []string, data [][]float64) (*mat.Dense, *mat.Dense) {
	// Create a dense matrix from the data slice
	X := mat.NewDense(len(data), len(headers)-1, nil)
	Y := mat.NewDense(len(data), 1, nil)

	for i, row := range data {
		for j, val := range row {
			if j == 0 {
				Y.Set(i, 0, val)
			} else {
				X.Set(i, j-1, val)
			}
		}
	}

	return X, Y
}

func computeMeans(X, Y mat.Matrix) (float64, float64) {
	// Compute the means of X and Y
	rows, _ := X.Dims()
	meanX := mat.Sum(X) / float64(rows)
	meanY := mat.Sum(Y) / float64(rows)

	return meanX, meanY
}

func computeVariancesAndCovariance(X mat.Matrix, Y mat.Matrix, meanX float64, meanY float64) (float64, float64) {
	var varX, covXY float64

	r, _ := X.Dims()
	for i := 0; i < r; i++ {
		varX += math.Pow(X.At(i, 0)-meanX, 2)
		covXY += (X.At(i, 0) - meanX) * (Y.At(i, 0) - meanY)
	}

	varX /= float64(r)
	covXY /= float64(r)

	return varX, covXY
}

func computeRegressionCoefficients(varX, covXY, meanX, meanY float64) (float64, float64) {
	// Compute the regression coefficients
	beta := covXY / varX
	alpha := meanY - beta*meanX

	return beta, alpha
}

func predictWins(headers []string, beta, alpha float64) float64 {
	// Use the regression coefficients to predict the team's wins
	wins := alpha

	for i, header := range headers {
		if i == 0 {
			continue
		}
		if header == "Year" {
			continue
		}
		fmt.Printf("Enter the team's %s: ", header)
		reader := bufio.NewReader(os.Stdin)
		val, _ := reader.ReadString('\n')
		val = strings.TrimSpace(val)
		v, _ := strconv.ParseFloat(val, 64)
		wins += beta * v
	}

	// Use raw HTTP request to OpenAI API for prediction
	winsGPT := makeOpenAIPrediction(headers)

	// Use the average of the two predictions
	return (wins + winsGPT) / 2
}

// makeOpenAIPrediction makes a prediction by calling the OpenAI API directly
func makeOpenAIPrediction(headers []string) float64 {
	apiURL := "https://api.openai.com/v1/completions"
	apiKey := "<your-api-key>" // Replace with your actual API key

	// Construct the request payload
	payload := map[string]interface{}{
		"model":       "text-davinci-002",
		"prompt":      fmt.Sprintf("Predict the number of wins for the %s based on the following stats:\n%s", headers[0], strings.Join(headers[1:], "\n")),
		"max_tokens":  100,
		"temperature": 0.5,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Create the request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		log.Fatalf("Error creating OpenAI request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error making OpenAI API request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading OpenAI response: %v", err)
	}

	// Parse the response JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Fatalf("Error parsing OpenAI response JSON: %v", err)
	}

	// Extract the generated text and parse it as a float
	choices := result["choices"].([]interface{})
	text := choices[0].(map[string]interface{})["text"].(string)
	winsGPT, err := strconv.ParseFloat(strings.TrimSpace(text), 64)
	if err != nil {
		log.Fatalf("Error parsing OpenAI prediction: %v", err)
	}

	return winsGPT
}
