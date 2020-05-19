package cmd

import (
    "bufio"
    "bytes"
    "fmt"
    "net/url"
    "os"
    "strings"

    "github.com/spf13/cobra"
)

// ------------------ INIT COBRA BOOKMARK CMD ------------------ //
// ------------------------------------------------------------- //
func init() {
    var ID bool
    RootCmd.AddCommand(BookmarkCmd)
    BookmarkCmd.Flags().StringP("add", "a", "", "add a new bookmark via url")
    BookmarkCmd.Flags().StringP("delete", "d", "", "delete a bookmark via id")
    BookmarkCmd.Flags().StringP("folder", "f", "", "set folder to store a bookmark in - (required)")
    BookmarkCmd.Flags().BoolVarP(&ID, "id", "i", false, "get the id's of each bookmark in a folder")
    BookmarkCmd.MarkFlagRequired("folder")
}

var BookmarkCmd = &cobra.Command {
    Use:   "bookmark",
    Short: "interact with bookmarks",
    Long:  `This subcommand will allow the user to interact with bookmarks`,
    Run: func(cmd *cobra.Command, args []string) {
        response := makeGetRequest("https://larder.io/api/1/@me/folders/")
        json := getFolders(response)

        bookmarkAddCmd, _ := cmd.Flags().GetString("add")
        bookmarkDeleteCmd, _ := cmd.Flags().GetString("delete")
        bookmarkFolderCmd, _ := cmd.Flags().GetString("folder")
        bookmarkIDCmd, _ := cmd.Flags().GetBool("id")

        if bookmarkAddCmd != "" {
            addBookmark(bookmarkAddCmd, bookmarkFolderCmd)
            os.Exit(0)
        }
        if bookmarkDeleteCmd != "" {
            deleteBookmark(bookmarkDeleteCmd, bookmarkFolderCmd, json.Results)
            os.Exit(0)
        }
        if bookmarkIDCmd {
            getIDFrom(bookmarkFolderCmd, json.Results)
            os.Exit(0)
        }
        cmd.Help()
    },
}


// ------------------ COMMAND STRUCTURES ------------------ //
// -------------------------------------------------------- //
type BookmarkAPIResponse struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []Bookmark `json:"results"`
}

type Bookmark struct {
    ID string `json:"id"`
    Parent Tags `json:"parent"`
    Tags []Tags `json:"tags"`
    Title string `json:"title"`
    Description string `json:"description"`
    URL string `json:"url"`
    Domain string `json:"domain"`
    Created string `json:"created"`
    Modified string `json:"modified"`
    Meta interface{} `json:"meta"`
}

type Tags struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Color string `json:"color"`
    Created string `json:"created"`
    Modified string `json:"modified"`
}


// ------------------ COMMAND OPTION FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func addBookmark(link string, folder string) {
    var b Bookmark
    response := makeGetRequest("https://larder.io/api/1/@me/folders/")
    json := getFolders(response)
    b = getBookmarkInfo(b, link, folder, json.Results)
    code := makePostRequest("https://larder.io/api/1/@me/links/add/", setBookmarkPostData(b))
    if checkSuccess(code, 201) {
        fmt.Println("Bookmark added!")
        os.Exit(0)
    }
    os.Exit(2)
}

func deleteBookmark(id string, folder string, folders []Folders) {
    var b Bookmark
    folderID := getFolderID(folder, folders)
    response := makeGetRequest("https://larder.io/api/1/@me/folders/" + folderID)
    json := getBookmarks(response)
    json.Count= 0
    code := makePostRequest("https://larder.io/api/1/@me/links/" + id + "/delete/", setBookmarkPostData(b))
    if checkSuccess(code, 204) {
        fmt.Println("Bookmark deleted!")
        os.Exit(0)
    }
    os.Exit(2)
}

// ------------------ COMMAND HELPER FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func getBookmarkInfo(b Bookmark, link string, folder string, folders []Folders) Bookmark {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter title of bookmark: ")
    b.Title, _ = reader.ReadString('\n')
    fmt.Print("Enter tags for bookmark: ")
    slice, _ := reader.ReadString('\n')
    b = setTags(b, slice)
    b.Parent.ID = getFolderID(folder, folders)
    b.URL = link
    if b.URL == "" {
        fmt.Println("No link was supplied.")
        os.Exit(2)
    }
    return b
}

func setTags(b Bookmark, slice string) Bookmark{
    var ts []Tags
    var t Tags
    for _, s := range strings.Split(slice, ",") {
        t.Name = strings.TrimSuffix(s, "\n")
        ts = append(ts, t)
    }
    b.Tags = ts
    return b
}

func getIDFrom(folder string, folders []Folders) {
    folderID := getFolderID(folder, folders)
    response := makeGetRequest("https://larder.io/api/1/@me/folders/" + folderID)
    json := getBookmarks(response)
    for _, item := range json.Results {
        fmt.Println(item.Title + " (" + item.ID + ")")
    }
}

func setBookmarkPostData(bookmark Bookmark) *bytes.Buffer {
    data := url.Values{}
    data.Set("title", bookmark.Title)
    data.Set("url", bookmark.URL)
    data.Set("parent", bookmark.Parent.ID)
    for _, t := range bookmark.Tags {
        data.Add("tags", t.Name)
    }
    return bytes.NewBufferString(data.Encode())
}
