package cmd

import (
    "fmt"
    "os"
    "strings"

    "github.com/spf13/cobra"
)


// ------------------ INIT COBRA SEARCH CMD ------------------ //
// ----------------------------------------------------------- //
func init() {
    RootCmd.AddCommand(SearchCmd)
    SearchCmd.Flags().StringP("params", "p", "", "search by given parameters")
}

var SearchCmd = &cobra.Command {
    Use:   "search",
    Aliases: []string{"se"},
    Short: "search through bookmarks",
    Long:  `This subcommand will allow the user to search through bookmarks`,
    Run: func(cmd *cobra.Command, args []string) {
        params, _ := cmd.Flags().GetString("params")
        if params != "" {
            params := strings.Split(params, ",")
            searchString := buildSearchString(params)
            response := makeGetRequest("https://larder.io/api/1/@me/search/" + searchString)
            json := getBookmarks(response)
            getSearchResults(json.Results)
            os.Exit(0)
        }
        cmd.Help()
    },
}


// ------------------ COMMAND HELPER FUNCTIONS ------------------ //
// ------------------------------------------------------------- //
func buildSearchString(params []string) string {
    searchString := "?q="
    for i, p := range params {
        if i +1 == len(params) {
            searchString = searchString + p
            break
        } else {
            searchString = searchString + p + "&q="
        }
    }
    return searchString
}

func getSearchResults(results []Bookmark) {
    for _, r := range results {
        fmt.Println(r.Title + " (" + r.URL + ")")
    }
}
