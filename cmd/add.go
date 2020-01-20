package cmd

import (
    "bytes"
    "fmt"
    "net/http"
    "net/url"

    "github.com/spf13/cobra"
)


type Bookmark struct {
	Parent string `json:"parent"`
	Title string `json:"title"`
	URL string `json:"url"`
	Tags []Tags `json:"tags"`
}

var AddCmd = &cobra.Command{
    Use:   "add",
    Short: "add a bookmark to larder",
    Long:  `This subcommand will add a bookmark to larder`,
    Run: func(cmd *cobra.Command, args []string) {
            addBookmark(cmd)
    },
}


func init() {
        RootCmd.AddCommand(AddCmd)
        AddCmd.Flags().StringP("folder", "f", "", "set folder to store bookmark (required)")
        AddCmd.Flags().StringP("url", "u", "", "set url to store be stored (required)")
        AddCmd.Flags().StringP("name", "n", "", "set name of the bookmark (required)")
	AddCmd.MarkFlagRequired("folder")
	AddCmd.MarkFlagRequired("url")
	AddCmd.MarkFlagRequired("name")
}

func getBookmarkData(cmd *cobra.Command, bookmark Bookmark, folder string, folders []Folders) Bookmark {
        bookmark.Title, _ = cmd.Flags().GetString("name")
        bookmark.URL, _ = cmd.Flags().GetString("url")
        bookmark.Parent = getFolderID(folder, folders)
        return bookmark
}


func getPostData(bookmark Bookmark) *bytes.Buffer {
        data := url.Values{}
        data.Set("parent", bookmark.Parent)
        data.Set("title", bookmark.Title)
        data.Set("url", bookmark.URL)
        data.Set("tags", "[]")
        return bytes.NewBufferString(data.Encode())
}

func addBookmark(cmd *cobra.Command) {
        var bookmark Bookmark
        response := getResponse("https://larder.io/api/1/@me/folders/")
        json := getFolders(response)
        folder, _ := cmd.Flags().GetString("folder")
        bookmark = getBookmarkData(cmd, bookmark, folder, json.Results)
        data := getPostData(bookmark)
        client := &http.Client{}
        req, err := http.NewRequest("POST", "https://larder.io/api/1/@me/links/add/", data)
        if err != nil {
            processError(err)
        }
        token := getAccessToken()
        req.Header.Set("Authorization", "Bearer " + token)
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        _, clientErr := client.Do(req)
        if clientErr != nil {
            processError(clientErr)
        }
        fmt.Println("Bookmark added!")
}
