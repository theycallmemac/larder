# 📚 larder

A [Larder.io](https://larder.io) command line interface - for all us people who don't like browser extensions!

## 🔧 Build

#### Creating a Larder Client

Before installing the Larder CLI, you need to go to the [Larder app management page](https://larder.io/apps/clients/). You might have to log in first.

You need to create a new Client. It's pretty standard, nothing that will catch you off-guard, just fill out the form to create the Oauth2 Client. 


#### Larder Config

Once you've created your client, you can head to that clients information page and fetch the following tokens:
- `access_token`
- `refresh_token`
- `client_id`
- `client_secret`

Once you have these four tokens, store them in `~/.larder/config.yml` like so:

``` yaml
access_token: $YOUR_ACCESS_TOKEN 
refresh_token: $YOUR_REFRESH_TOKEN 
client_id: $YOUR_CLIENT_ID
client_secret: $YOUR_CLIENT_SECRET
```


Once you've stored your client information in `~/.larder/config.yml`, you are ready to build the CLI. To build, run `make install` in the top level of this repository. Requires sudo or root permissions.



## 🚀 Usage

You can run the Larder CLI from anywhere by running the `larder` command

You can checkout the available commands with `larder --help`.

#### Help

```
A CLI to your bookmarks

Usage:
  larder [command]

Available Commands:
  bookmark      interact with bookmarks
  folder        interact with folders
  help          Help about any command
  refresh-token refresh your access token
  search        search through bookmarks

Flags:
  -h, --help   help for larder
```

#### Folders

```
This subcommand will allow the user to interact with folders

Usage:
  larder folder [flags]

Flags:
  -a, --add string      add a new folder
  -d, --delete string   delete a folder
  -h, --help            help for folder
  -l, --list            list current folders
  -s, --show string     show contents of a folder
```

#### Bookmarks

```
This subcommand will allow the user to interact with bookmarks

Usage:
  larder bookmark [flags]

Flags:
  -a, --add string      add a new bookmark via url
  -d, --delete string   delete a bookmark via id
  -f, --folder string   set folder to store a bookmark in - (required)
  -h, --help            help for bookmark
  -i, --id              get the id's of each bookmark in a folder

```

As seen above, the `-f` / `--folder` paramter is always required. This means that when adding, removing or dispalying bookmark id's a destination folder must always be provided.

#### Searching

```
This subcommand will allow the user to search through bookmarks

Usage:
  larder search [flags]

Flags:
  -h, --help            help for search
  -p, --params string   search by given parameters
```

The `-p` / `--params` flag is a bt vague here. The user provides a string of search terms delimited by commas. For example: "texas,bbq". This will search through all folders, for names or tags containing those terms.

#### Refreshing Access Tokens

Tokens expire in a month and can be refreshed for a new access token at any time, invalidating the original access and refresh tokens.

Tokens can be automatically refreshed by running `larder refresh-token`. 

## 👤 Author

**James McDermott**

- Email: <james.mcdermott7@mail.dcu.ie>
- Twitter: [@theycallmemac_](https://twitter.com/theycallmemac_)
- Github: [@theycallmemac](https://github.com/theycallmemac)

## ⭐️ Show your support

Give a ⭐️ if this project helped you!
