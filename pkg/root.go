package cmd

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
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

func getBookmarks(blob io.ReadCloser) BookmarkAPIResponse {
    decoder := json.NewDecoder(blob)
    var f BookmarkAPIResponse
    err := decoder.Decode(&f)
    if err != nil {
        processError(err)
    }
    return f
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
func makePostRequest(link string, data *bytes.Buffer) int {
    client := &http.Client{}
    req, _ := http.NewRequest("POST", link, data)
    token := getAccessToken()
    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    response, clientErr := client.Do(req)
    if clientErr != nil {
        processError(clientErr)
    }
    return response.StatusCode
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

func checkSuccess(code int, expected int) bool {
    return code == expected
}
