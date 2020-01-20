package cmd

import (
        "fmt"
        "log"
        "os"
        "os/user"

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

type Config struct {
    AccessToken string `yaml:"access_token"`
    RefreshToken string `yaml:"refresh_token"`
    ClientID string `yaml:"client_id"`
    ClientSecret string `yaml:"client_secret"`
}

func processError(err error) {
    fmt.Println(err)
    os.Exit(2)
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
