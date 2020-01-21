package cmd

import (
    "bytes"
    "fmt"
    "net/url"
    "os"
    "strings"

    "github.com/spf13/cobra"
)


// ------------------ INIT COBRA FOLDER CMD ------------------ //
// ----------------------------------------------------------- //
func init() {
    RootCmd.AddCommand(FolderCmd)
    var List bool
    FolderCmd.Flags().StringP("add", "a", "", "add a new folder")
    FolderCmd.Flags().StringP("delete", "d", "", "delete a folder")
    FolderCmd.Flags().BoolVarP(&List, "list", "l", false, "list current folders")
    FolderCmd.Flags().StringP("show", "s", "", "show contents of a folder")
}

var FolderCmd = &cobra.Command {
    Use:   "folder",
    Short: "interact with folders",
    Long:  `This subcommand will allow the user to interact with folders`,
    Run: func(cmd *cobra.Command, args []string) {
	response := makeGetRequest("https://larder.io/api/1/@me/folders/")
        json := getFolders(response)

        folderAddCmd, _ := cmd.Flags().GetString("add")
        folderDeleteCmd, _ := cmd.Flags().GetString("delete")
        folderListCmd, _ := cmd.Flags().GetBool("list")
        folderShowCmd, _ := cmd.Flags().GetString("show")

        var f Folder
        if folderListCmd {
            listFolders(json.Results)
            os.Exit(0)
        }
        if folderAddCmd != "" {
            addFolder(folderAddCmd, f)
            os.Exit(0)
        }
        if folderDeleteCmd != "" {
            deleteFolder(folderDeleteCmd, f)
            os.Exit(0)
        }
        if folderShowCmd != "" {
            showFolder(folderShowCmd, f, json.Results)
            os.Exit(0)
        }
        cmd.Help()
    },
}


// ------------------ COMMAND STRUCTURES ------------------ //
// -------------------------------------------------------- //
type FolderAPIResponse struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []Folders `json:"results"`
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

type Folder struct {
    Name string `json:"name"`
    Parent string `json:"parent"`
}

type EmptyFolder struct {
    EMPTY_TO string `json:"empty_to"`
}


// ------------------ COMMAND OPTION FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func addFolder(folder string, f Folder) {
    response := makeGetRequest("https://larder.io/api/1/@me/folders/")
    json := getFolders(response)
    if exists(folder, json.Results) {
        fmt.Println("Folder '" + folder + "' already exists.")
        os.Exit(2)
    }
    f.Name = folder
    f.Parent = getFolderID(folder, json.Results)
    code := makePostRequest("https://larder.io/api/1/@me/folders/add/", setFolderPostData(f))
    if checkSuccess(code, 201) {
        fmt.Println("Folder added!")
    }
}

func deleteFolder(folder string, f Folder) {
    var ef EmptyFolder
    response := makeGetRequest("https://larder.io/api/1/@me/folders/")
    json := getFolders(response)
    if exists(folder, json.Results) {
        f.Name = folder
        f.Parent = getFolderID(folder, json.Results)
        ef.EMPTY_TO = ""
        code := makePostRequest("https://larder.io/api/1/@me/folders/" + f.Parent + "/delete/", emptyFolderPostData(ef))
        if checkSuccess(code, 204) {
            fmt.Println("Folder deleted!")
        }
        os.Exit(0)
    }
    fmt.Println("Folder '" + folder + "' does not exist.")
}

func listFolders(folders []Folders) {
    for _, folder := range folders {
        fmt.Println(folder.Name, "-", folder.Links)
    }
}

func showFolder(folder string, f Folder, folders []Folders) {
    if exists(folder, folders) {
        f.Name = folder
        f.Parent = getFolderID(f.Name, folders)
        response := makeGetRequest("https://larder.io/api/1/@me/folders/" + f.Parent + "?limit=200")
        json := getBookmarks(response)
        for _, contents := range json.Results {
            fmt.Println(contents.Title, "(" + contents.URL + ")")
        }
        os.Exit(0)
    }
    fmt.Println("Folder '" + folder + "' does not exist.")
}

// ------------------ COMMAND HELPER FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func exists(folder string, folders []Folders) bool {
    folderName := strings.ToLower(folder)
    for _, folder := range folders {
        currentFolder := strings.ToLower(folder.Name)
        if currentFolder == folderName {
            return true
        }
    }
    return false
}

func setFolderPostData(folder Folder) *bytes.Buffer {
    data := url.Values{}
    data.Set("name", folder.Name)
    data.Set("parent", folder.Parent)
    return bytes.NewBufferString(data.Encode())
}

func setFoldersPostData(fs Folders) *bytes.Buffer {
    data := url.Values{}
    data.Set("name", fs.Name)
    return bytes.NewBufferString(data.Encode())
}

func emptyFolderPostData(folder EmptyFolder) *bytes.Buffer {
    data := url.Values{}
    data.Set("empty_to", folder.EMPTY_TO)
    return bytes.NewBufferString(data.Encode())
}
