package cmd

import (
    "fmt"
    "encoding/json"
    "net/http"
    "os"
    "strings"
    "github.com/spf13/cobra"
)


type Tags struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Color string `json:"color"`
    Created string `json:"created"`
    Modified string `json:"modified"`
}

type FolderContents struct {
    ID string `json:"id"`
    Tags []Tags `json:"tags"`
    Title string `json:"title"`
    Description string `json:"description"`
    URL string `json:"url"`
    Domain string `json:"Domain"`
    Created string `json:"created"`
    Modified string `json:"modified"`
    Meta interface{} `json:"meta"`
}

type Folders struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Color string `json:"color"`
    Icon string `json:"icon"`
    Created string `json:"created"`
    Modified string `json:"modified"`
    Parent string `json:"parent"`
    Folders []string `json:"folders"`
    Links int `json:"Links"`
}

type FolderContentAPIResponse struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []FolderContents `json:"results"`
}

type FolderAPIResponse struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []Folders `json:"results"`
}

var ListCmd = &cobra.Command{
    Use:   "list",
    Short: "list the user's folders",
    Long:  `This subcommand will list the user's folders`,
    Run: func(cmd *cobra.Command, args []string) {

            client := &http.Client{}
            req, _ := http.NewRequest("GET", "https://larder.io/api/1/@me/folders/", nil)
            token := getAccessToken()
            req.Header.Set("Authorization", "Bearer " + token)
            res, _ := client.Do(req)
            decoder := json.NewDecoder(res.Body)
            var f FolderAPIResponse
            err := decoder.Decode(&f)
            if err != nil {
                fmt.Println(err)
            }
            if cmd.Flags().Changed("folder") {
                listFolderContent(cmd, f.Results)
                os.Exit(0)
            }
            listAllFolders(f.Results)
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

func listFolderContent(cmd *cobra.Command, folders []Folders) {
        folderName, _ := cmd.Flags().GetString("folder")
        folderName = strings.ToLower(folderName)
        var id string
        for _, folder := range folders {
            currentFolder := strings.ToLower(folder.Name)
            if currentFolder == folderName {
                id = folder.ID
            }
        }
        client := &http.Client{}
        req, _ := http.NewRequest("GET", "https://larder.io/api/1/@me/folders/" + id, nil)
        token := getAccessToken()
        req.Header.Set("Authorization", "Bearer " + token)
        res, _  := client.Do(req)
        decoder := json.NewDecoder(res.Body)
        var fc FolderContentAPIResponse
        err := decoder.Decode(&fc)
        if err != nil {
            fmt.Println(err)
        }
        for _, contents:= range fc.Results {
            fmt.Println(contents.Title, "(" + contents.URL + ")")
        }
}
