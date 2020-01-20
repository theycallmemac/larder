package cmd

import (
    "fmt"

    "github.com/spf13/cobra"
)

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

func getBookmarkInfo(cmd *cobra.Command, bookmark Bookmark, folder string, folders []Folders) Bookmark {
        bookmark.Title, _ = cmd.Flags().GetString("name")
        bookmark.URL, _ = cmd.Flags().GetString("url")
        bookmark.Parent = getFolderID(folder, folders)
        return bookmark
}

func addBookmark(cmd *cobra.Command) {
        var bookmark Bookmark
        response := makeGetRequest("https://larder.io/api/1/@me/folders/")
        json := getFolders(response)
        folder, _ := cmd.Flags().GetString("folder")
        bookmark = getBookmarkInfo(cmd, bookmark, folder, json.Results)
        makePostRequest("https://larder.io/api/1/@me/links/add/", setPostData(bookmark))
        fmt.Println("Bookmark added!")
}
