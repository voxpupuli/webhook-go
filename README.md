# Webhook Go

Webhook Go is a port of the [puppet_webhook](https://github.com/voxpupuli/puppet_webhook) Sinatra API server to Go.
This is designed to be more streamlined, performant, and easier to ship for users than the Sinatra/Ruby API server.

This server is a REST API server designed to accept Webhooks from version control systems, such as GitHub or GitLab, and execute actions based on those webhooks. Specifically, the following tasks:

* Trigger r10k environment and module deploys onto Puppet Servers
* Send notifications to ChatOps systems, such as Slack and RocketChat

## Prerequisites

While there are no prerequisites for running the webhook server itself, for it to be useful, you will need the following installed on the same server or another server for this tool to be useful:

* Puppet Server
* [r10k](https://github.com/puppetlabs/r10k)
* Puppet Bolt (optional)
* Windows or Linux server to run the server on. MacOS is not supported.

## Installation

Download a Pre-release Binary from the [Releases](https://github.com/voxpupuli/webhook-go/releases) page, make it executable, and run the server.

## Configuration

The Webhook API server uses a configuration file called `webhook.yml` to configure the server. Several of the required options have defaults pre-defined so that a configuration file isn't needed for basic function.

`webhook.yaml.example`:
```yaml
server:
  protected: false
  user: puppet
  password: puppet
  port: 4000
  tls:
    enabled: false
    certificate: "/path/to/tls/certificate"
    key: "/path/to/tls/key"
chatops:
  enabled: false
  service: slack
  channel: "#general"
  user: r10kbot
  auth_token: 12345
  server_uri: "https://rocketchat.local"
orchestration:
  enabled: true
  type: bolt
  user: webhook
  password: password
  bolt:
    transport: local
    targets:
      - localhost
r10k:
  config_path: /etc/puppetlabs/r10k/r10k.yaml
  default_branch: main
  allow_uppercase: false
  verbose: true

```

### Server options

#### `protected`

Type: bool
Description: Enforces authetication via basic Authentication
Default: `false`

#### `user`

Type: string
Description: Username to use for Basic Authentication. Optional.
Default: `nil`

#### `password`

Type: string
Description: Password to use for Basic Authentication. Optional.
Default: `nil`

#### `port`

Type: int64
Description: Port to run the server on. Optional.
Default: `4000`

#### `tls`

Type: struct
Description: Struct containing server TLS options

##### `enabled`

Type: bool
Description: Enforces TLS with http server
Default: `false`

##### `certificate`

Type: string
Description: Full path to certificate file. Optional.
Default: `nil`

##### `key`

Type: string
Description: Full path to key file. Optional.
Default: `nil`

### ChatOps options

#### `enabled`

Type: boolean
Description: Enable/Disable chatops support
Default: false

#### `service`

Type: string
Description: Which service to use. Supported options: [`slack`, `rocketchat`]
Default: nil

#### `channel`

Type: string
Description: ChatOps communication channel to post to.
Default: nil

#### `user`

Type: string
Description: ChatOps user to post as
Default: nil

#### `auth_token`

Type: string
Description: The authentication token needed to post as the ChatOps user in the chosen, supported ChatOps service
Default: nil

#### `server_uri`

Type: string
Description: The ChatOps service API URI to send the message to.
Default: nil

### Orchestration

#### `enabled`

Type: boolean
Description: Enable/Disable orchestration support
Default: false

#### `type`

Type: string
Description: Which orchestration tool to use. Currently only supports `bolt`.
Default: nil

#### `user`

Type: string
Description: User to authenticate to the target as.
Default: nil

#### `password`

Type: string
Description: Password for authentication use.
Default: nil

#### `bolt`

Type: hash
Description: Hash of Puppet Bolt specific settings.
Default: `nil`

##### `transport`

Type: string
Description: What kind of bolt network transport protocol to use. Supported types: [`local`, `ssh`] with more coming soon.
Default: `nil`

##### `targets`

Type: array
Description: A list of target nodes that Bolt will attempt to run `r10k` on.
Default: []

##### `concurrency`

Type: integer
Description: Maximum number of simultaneous connections.
Default: 100

##### `run_as`

Type: string
Description: User to run the `r10k` command as.
Default: nil

##### `sudo_password`

Type: string
Description: The password to use when using `sudo` to run `r10k` as the `run_as` user.
Default: nil

##### `HostKeyCheck`

Type: boolean
Description: Enable/Disable SSH host key checking
Default: false

### r10k options

#### `config_path`

Type: string
Description: Full path to the r10k configuration file. Optional.
Default: `/etc/puppetlabs/r10k/r10k.yaml`

#### `default_branch`

Type: string
Description: Name of the default branch for r10k to pull from. Optional.
Default: `main`

#### `prefix`

Type: string
Description: An r10k prefix to apply to the module or environment being deployed. Optional.
Default: `nil`

#### `allow_uppercase`

Type: bool
Description: Allow Uppercase letters in the module, branch, or environment name. Optional.
Default: `false`

#### `verbose`

Type: bool
Description: Log verbose output when running the r10k command
Default: `true`

### `deploy_modules`

Type: bool
Description: Deploy modules in environments.
Default: `true`

### `generate_types`

Type: bool
Description: Run `puppet generate types` after updating an environment
Default: `true`

