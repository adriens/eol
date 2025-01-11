package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "net/http"
    "os"
    "text/tabwriter"
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
            fmt.Printf("Product: %s\n", product)
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

    var productInfo []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&productInfo); err != nil {
        fmt.Println("Error decoding JSON:", err)
        os.Exit(1)
    }

    writer := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
    fmt.Fprintln(writer, "Cycle\tLTS\tRelease Date\tSupport\tEOL\tLatest\tLink\t")
    fmt.Fprintln(writer, "-----\t---\t-------------\t-------\t---\t------\t----\t")
    for _, info := range productInfo {
        link := ""
        if info["link"] != nil {
            link = info["link"].(string)
        }
        support := ""
        if info["support"] != nil {
            support = info["support"].(string)
        }
        eol := ""
        switch v := info["eol"].(type) {
        case bool:
            if v {
                eol = "true"
            } else {
                eol = "false"
            }
        case string:
            eol = v
        case nil:
            eol = ""
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

        fmt.Fprintf(writer, "%s\t%s\t%s\t%s\t%s%s%s\t%s\t%s\t\n",
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
    writer.Flush()
}

func printHelp() {
    fmt.Println("CLI to show end-of-life dates for a number of products, from https://endoflife.date")
    fmt.Println()
    fmt.Println("For example:")
    fmt.Println()
    fmt.Println("* `eol python` to see Python EOLs")
    fmt.Println("* `eol ubuntu` to see Ubuntu EOLs")
    fmt.Println("* `eol centos fedora` to see CentOS and Fedora EOLs")
    fmt.Println("* `eol all` or `eol` to list all available products")
}