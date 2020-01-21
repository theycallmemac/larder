package cmd

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"

    "gopkg.in/yaml.v2"
    "github.com/spf13/cobra"
    "github.com/tidwall/gjson"
)


// ------------------ INIT COBRA REFRESH TOKEN CMD ------------------ //
// ------------------------------------------------------------------ //
func init() {
    RootCmd.AddCommand(RefreshTokenCmd)
}

var RefreshTokenCmd = &cobra.Command {
    Use:   "refresh-token",
    Short: "refresh your access token",
    Long:  `This subcommand will automate help refresh your access token`,
    Run: func(cmd *cobra.Command, args []string) {
        refreshAccessToken()
    },
}


// ------------------ COMMAND OPTION FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func refreshAccessToken() {
    var newCfg Config = getNewValues(makeRefreshPost())
    MarshalledCfg, _ := yaml.Marshal(&newCfg)
    err := ioutil.WriteFile(pathToConfig(), MarshalledCfg, 0755)
    if err != nil {
        fmt.Println("An error occured during writing the new file, here are your new tokens:\n  - Access Token = " + newCfg.AccessToken + "\n - Refresh Token = " + newCfg.RefreshToken)
        os.Exit(2)
    }
    fmt.Println("New token generated and stored in ~/.larder/config.yml!")
}


// ------------------ COMMAND HELPER FUNCTIONS ------------------ //
// -------------------------------------------------------------- //
func makeRefreshPost() string {
    client := &http.Client{}
    data := setRefreshPostData()
    req, _ := http.NewRequest("POST", "https://larder.io/oauth/token/", data)
    token := getAccessToken()
    req.Header.Set("Authorization", "Bearer " + token)
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    response, clientErr := client.Do(req)
    if clientErr != nil {
        processError(clientErr)
    }
    body, _ := ioutil.ReadAll(response.Body)
    return string(body)
}

func setRefreshPostData() *bytes.Buffer {
    data := url.Values{}
    data.Set("refresh_token", getRefreshToken())
    data.Set("grant_type", "refresh_token")
    data.Set("client_id", getClientID())
    data.Set("client_secret", getClientSecret())
    return bytes.NewBufferString(data.Encode())
}

func getNewValues(body string) Config {
    var cfg Config
    cfg = getYaml(pathToConfig())
    cfg.AccessToken = gjson.Get(string(body), "access_token").String()
    cfg.RefreshToken = gjson.Get(string(body), "refresh_token").String()
    return cfg
}
