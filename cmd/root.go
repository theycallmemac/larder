package cmd

import (
	"bytes"
	"encoding/json"
        "fmt"
	"io"
        "log"
	"net/http"
	"net/url"
        "os"
        "os/user"
        "strings"

        "gopkg.in/yaml.v2"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "larder-cli",
	Short: "a cli to your bookmarks",
	Long: `a cli to your bookmarks`,
}


func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


// ------------------ ~/.larder/config.yml ------------------ //
// ---------------------------------------------------------- //
type Config struct {
    AccessToken string `yaml:"access_token"`
    RefreshToken string `yaml:"refresh_token"`
    ClientID string `yaml:"client_id"`
    ClientSecret string `yaml:"client_secret"`
}

func readFile(filename string) *os.File {
    file, err := os.Open(filename)
    if err != nil {
        processError(err)
        var tmp *os.File
        return tmp
    }
    return file
}

func parseYaml(cfg *Config, file *os.File) bool {
    decoder:= yaml.NewDecoder(file)
    err := decoder.Decode(cfg)
    if err != nil {
        processError(err)
        return false
    }
    return true
}

func getYaml(filename string) Config {
    var cfg Config
    yamlFile := readFile(filename)
    success := parseYaml(&cfg, yamlFile)
    if success == false {
        log.Fatal("Failed to read file.")
    }
    return cfg
}

func pathToConfig() string {
    usr, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }
    return usr.HomeDir + "/.larder/config.yml"
}

func processError(err error) {
    fmt.Println(err)
    os.Exit(2)
}

// ------------------ SHARED STRUCTS ------------------ //
// ---------------------------------------------------- //
type Bookmark struct {
        Parent string `json:"parent"`
        Title string `json:"title"`
        URL string `json:"url"`
        Tags []Tags `json:"tags"`
}

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

type FolderContentAPIResponse struct {
    Count int `json:"count"`
    Next string `json:"next"`
    Previous string `json:"previous"`
    Results []FolderContents `json:"results"`
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

type Tags struct {
    ID string `json:"id"`
    Name string `json:"name"`
    Color string `json:"color"`
    Created string `json:"created"`
    Modified string `json:"modified"`
}


// ------------------ SHARED GET METHODS ------------------ //
// -------------------------------------------------------- //
func makeGetRequest(link string) io.ReadCloser {
            client := &http.Client{}
            req, _ := http.NewRequest("GET", link, nil)
            token := getAccessToken()
            req.Header.Set("Authorization", "Bearer " + token)
            res, _ := client.Do(req)
            return res.Body
}

func getFolders(blob io.ReadCloser) FolderAPIResponse {
            decoder := json.NewDecoder(blob)
            var f FolderAPIResponse
            err := decoder.Decode(&f)
            if err != nil {
                processError(err)
            }
            return f
}

func getFolderContents(blob io.ReadCloser) FolderContentAPIResponse {
            decoder := json.NewDecoder(blob)
            var fc FolderContentAPIResponse
            err := decoder.Decode(&fc)
            if err != nil {
                processError(err)
            }
            return fc
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


// ------------------ SHARED POST METHODS ------------------ //
// --------------------------------------------------------- //
func makePostRequest(link string, data *bytes.Buffer) {
        client := &http.Client{}
        req, _ := http.NewRequest("POST", link, data)
        token := getAccessToken()
        req.Header.Set("Authorization", "Bearer " + token)
        req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
        _, clientErr := client.Do(req)
        if clientErr != nil {
            processError(clientErr)
        }
}

func setPostData(bookmark Bookmark) *bytes.Buffer {
        data := url.Values{}
        data.Set("parent", bookmark.Parent)
        data.Set("title", bookmark.Title)
        data.Set("url", bookmark.URL)
        data.Set("tags", "[]")
        return bytes.NewBufferString(data.Encode())
}


// ------------------ SHARED ACCESS METHODS ------------------ //
// ----------------------------------------------------------- //
func getAccessToken() string {
    var cfg = getYaml(pathToConfig())
    return cfg.AccessToken
}

func getRefreshToken() string {
    var cfg = getYaml(pathToConfig())
    return cfg.RefreshToken
}

func getClientID() string {
    var cfg = getYaml(pathToConfig())
    return cfg.ClientID
}

func getClientSecret() string {
    var cfg = getYaml(pathToConfig())
    return cfg.ClientSecret
}
