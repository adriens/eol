package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const baseURL = "https://endoflife.date/api"

func main() {
	listFlag := flag.Bool("list", false, "List all items from the API")
	listFlagShort := flag.Bool("l", false, "List all items from the API (short)")
	helpFlag := flag.Bool("help", false, "Show help")
	helpFlagShort := flag.Bool("h", false, "Show help")

	flag.Parse()

	if *helpFlag || *helpFlagShort {
		printHelp()
		os.Exit(0)
	}

	if !*listFlag && !*listFlagShort && len(flag.Args()) == 0 {
		printHelp()
		os.Exit(1)
	}

	if *listFlag || *listFlagShort {
		listItems()
	} else if len(flag.Args()) > 0 {
		for _, product := range flag.Args() {
			fmt.Printf("## %s\n\n", strings.Title(product))
			getProductInfo(product)
			fmt.Println()
		}
	}
}

func listItems() {
	resp, err := http.Get(baseURL + "/all.json")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	var items []string
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	for _, item := range items {
		fmt.Println(item)
	}
}

func getProductInfo(product string) {
	resp, err := http.Get(fmt.Sprintf("%s/%s.json", baseURL, product))
	if err != nil {
		fmt.Println("Error fetching data:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Print the error message in orange color : Error: The product '%s' could not be found. Please run `eol -l` to find the supported products.\n
		fmt.Printf("\033[33mError: The product '%s' could not be found. Please run `eol -l` to find the supported products.\033[0m\n", product)
		fmt.Print()
		return
	}

	var productInfo []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&productInfo); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}

	fmt.Printf("| Cycle | LTS | Release Date | Support | EOL | Latest | Link |\n")
	fmt.Printf("|-------|-----|--------------|---------|-----|--------|------|\n")
	for _, info := range productInfo {
		link := ""
		if info["link"] != nil {
			link = info["link"].(string)
		}

		support := ""
		if info["support"] != nil {
			support = fmt.Sprintf("%v", info["support"])
		}

		eol := ""
		if info["eol"] != nil {
			// print support as string
			eol = fmt.Sprintf("%v", info["eol"])

		}

		lts := "false"
		if info["lts"].(bool) {
			lts = "true"
		}

		eolColor := ""
		resetColor := "\033[0m"
		if eol != "" && eol != "false" && eol != "true" {
			eolDate, err := time.Parse("2006-01-02", eol)
			if err == nil {
				now := time.Now()
				if eolDate.Before(now) {
					eolColor = "\033[31m" // Red color
				} else if eolDate.Before(now.AddDate(0, 3, 0)) {
					eolColor = "\033[33m" // Orange color
				} else {
					eolColor = "\033[32m" // Green color
				}
			}
		}

		fmt.Printf("| %s | %s | %s | %s | %s%s%s | %s | %s |\n",
			info["cycle"],
			lts,
			info["releaseDate"],
			support,
			eolColor,
			eol,
			resetColor,
			info["latest"],
			link)
	}
}

func printHelp() {
	fmt.Println("CLI to show end-of-life (eol) dates for a number of products, from https://endoflife.date. See https://bit.ly/4jd3lbH for more.")
	fmt.Println()
	fmt.Println("For example:")
	fmt.Println()
	fmt.Println("* `eol python` to see Python EOLs")
	fmt.Println("* `eol ubuntu` to see Ubuntu EOLs")
	fmt.Println("* `eol centos fedora` to see CentOS and Fedora EOLs")
	fmt.Println("* `eol java quarkus` to see Quarkus and Java EOLs")
	fmt.Println("* `eol -l` or `eol --list` to list all available products")
}
