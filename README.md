# ![My Love for the San Francisco Giants](/images/my-picture.jpg)

Hello there, fellow baseball fans! Are you ready to hear all about my love for the San Francisco Giants?

The Giants are a team with a rich history, devoted fans, and a winning spirit. As a software factory chief engineer at Booz Allen Hamilton, I appreciate the dedication and precision that the Giants bring to the field.

And let's not forget about the Giants' incredible fan base. There's nothing quite like being surrounded by thousands of fellow Giants fans, all chanting and cheering in unison. It's enough to give anyone goosebumps.

So why am I telling you all this? Well, I just want to share my love for the Giants with the world. Whether you're a die-hard fan or a casual observer, I invite you to join me in rooting for this amazing team. Who knows? Maybe we'll even get to watch them win another World Series championship one day.

And hey, if you're ever in San Francisco and want to catch a Giants game with me, just hit me up. I'm always down for some good old-fashioned baseball and camaraderie. Let's go Giants!

## A Giant's Code Tour

## Disclaimer

Hey there, before we dive into the code, I want to make something clear - this is just a simple program I made to scratch an itch. I know there are other more elegant and visually-pleasing ways to get baseball stats, but I just wanted to have some fun and share it with fellow baseball fans.

That being said, I welcome any feedback or suggestions you may have for how to improve this program or make it more useful for baseball fans. But above all, let's remember to be nice! We're all here because we love baseball and want to share that love with others.

So if you have any comments or critiques, please share them respectfully and constructively. Let's keep the spirit of the game alive and well in everything we do.

And of course, if you find any bugs or issues with the program, please let me know so I can fix them.

So with that out of the way, let's dive into the code and see what kind of baseball stats we can scrape up!

This code is a Go program that prompts the user to enter a Major League Baseball team's abbreviation (e.g. SF for San Francisco Giants), then uses web scraping techniques to fetch and parse the team's 2022 season batting statistics from Baseball Reference. The program then outputs the table headers and the team's stats for the 2022 season to the console.

Let's take a closer look at how the code works:

```go
func main() {
	run(os.Stdin, os.Stdout)
}
```

The main function simply calls the run function with standard input and output as the input and output parameters.

```go
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
```

The run function takes an input and output parameter which allow for testing and customization of the input and output. It first prompts the user to enter a team abbreviation and then constructs a URL based on the entered abbreviation. The URL is used to make an HTTP request to Baseball Reference, and the HTML response is parsed using the goquery package.

The program then searches the HTML for a table with the ID team_batting, and extracts the table headers and rows. It stores the headers in a slice of strings called headers, and the rows in a slice of slices of strings called rows.

Finally, the program prints the headers and the team's stats for the 2022 season to the console. It first joins the headers with tabs and prints them, then loops through the rows until it finds the row with the year "2022". It then joins the cells of that row

## A Giant's Unit Test Code Tour

This code is a Go test for the run function in the previous code example. It sets up a fake input string, captures the output of the run function when it's called with the fake input, and compares the actual output to an expected output.

Let's take a closer look at how the code works:

```go
func TestRun(t *testing.T) {
	// Set up a fake input string and a buffer to capture output
	inputStr := "SF\n"
	expectedOutput := "Year\tTm\t#Bat\tBatAge\tR/G\tG\tPA\tAB\tR\tH\t2B\t3B\tHR\tRBI\tSB\tCS\tBB\tSO\tBA\tOBP\tSLG\tOPS\tOPS+\n2022\tSFG\t27\t28.2\t4.4\t162\t6428\t5557\t708\t1431\t282\t43\t209\t694\t57\t20\t621\t1426\t.258\t.330\t.432\t.762\t106\n"
	outputBuf := bytes.NewBuffer(nil)

	// Run the function with the fake input and capture the output
	run(strings.NewReader(inputStr), outputBuf)
	actualOutput := outputBuf.String()

	// Compare the actual and expected output
	if actualOutput != expectedOutput {
		t.Errorf("Unexpected output:\nExpected: %s\nActual: %s", expectedOutput, actualOutput)
	}
}
```
The TestRun function is a Go test function that takes a testing.T parameter. It first sets up a fake input string and an expected output string. It then creates a buffer to capture the output of the run function.

The run function is called with a strings.NewReader that reads the fake input string and the output is written to the buffer created earlier.

Finally, the actual output captured in the buffer is compared to the expected output. If the actual output does not match the expected output, an error message is printed. This test ensures that the run function works as expected and produces the correct output for a given input.

## How to Contribute

Hello there, fellow baseball fans! Do you want to help improve this Go program and make it even better? Great, because we could use your help!

Here are some ways you can contribute:

1. Add More Stats

Are you a stat-head who can't get enough baseball data? Then we need your expertise! Add more stats to the table to make it even more informative. Maybe you could even come up with your own new stat - something like the "slug-tastic index" or the "batting-rama score". The sky's the limit!

2. Improve the UI

Are you a design whiz who loves making things look pretty? Then we need you too! The current program output is pretty basic, so let's spruce it up a bit. Maybe you could add some fancy colors or maybe some ASCII art of your favorite player. The possibilities are endless!

3. Fix Bugs

Are you a master bug-squasher who loves nothing more than fixing pesky issues? Then we could definitely use your skills. If you find any bugs or issues with the program, please let us know so we can fix them.

4. Play Ball!

Are you a baseball player or coach with insights on how to make the program even more useful? We want to hear from you! Maybe you have some ideas on how to incorporate defensive statistics or how to make the program more relevant to different positions. We're all ears!

So there you have it - four ways you can contribute to this amazing baseball program. Whether you're a stat-head, a design whiz, a bug-squasher, or a baseball pro, we welcome your contributions with open arms. Let's make this program the best it can be!
