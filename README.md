# larder

A CLI for Larder - for people who don't like browser extensions!

### Creating a Larder Client

Before installing the Larder CLI, you need to go to the [Larder app management page](https://larder.io/apps/clients/). You might have to log in first.

You need to create a new Client. It's pretty standard, nothing that will catch you off-guard, just fill out the form to create the Oauth2 Client. 

### Larder Config

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

### Build

Once you've stored your client information in `~/.larder/config.yml`, you are ready to build the CLI. To build, run `make install` in the top level of this repository. Requires sudo or root permissions.

### Run

You can run the Larder CLI from anywhere by running the `larder` command

You can checkout the available commands with `larder --help`.
