package cmd

import (
    "fmt"
    "os"
    "strings"
    "github.com/spf13/cobra"
)

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "list the user's folders",
    Long:  `This subcommand will list the user's folders`,
    Run: func(cmd *cobra.Command, args []string) {
            response := getResponse("https://larder.io/api/1/@me/folders/")
            json := getFolders(response)
            if cmd.Flags().Changed("folder") {
                listFolderContents(cmd, json.Results)
                os.Exit(0)
            }
            listAllFolders(json.Results)
    },
}

func init() {
        RootCmd.AddCommand(ListCmd)
        ListCmd.Flags().StringP("folder", "f", "", "list the links in a folder")
}

func listAllFolders(folders []Folders) {
        for _, folder := range folders {
            fmt.Println(folder.Name, "-", folder.Links)
        }
}

func getFolderID(folder string, folders []Folders) string {
        folderName := strings.ToLower(folder)
        var id string
        for _, folder := range folders {
            currentFolder := strings.ToLower(folder.Name)
            if currentFolder == folderName {
                id = folder.ID
                break
            }
        }
        return id
}

func listFolderContents(cmd *cobra.Command, folders []Folders) {
        folder, _ := cmd.Flags().GetString("folder")
        id := getFolderID(folder, folders)
        response := getResponse("https://larder.io/api/1/@me/folders/" + id)
        json := getFolderContents(response)
        for _, contents := range json.Results {
            fmt.Println(contents.Title, "(" + contents.URL + ")")
        }
}
