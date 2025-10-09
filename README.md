# Webhook Go

Webhook Go is a port of the [puppet_webhook](https://github.com/voxpupuli/puppet_webhook) Sinatra API server to Go.
This is designed to be more streamlined, performant, and easier to ship for users than the Sinatra/Ruby API server.

This server is a REST API server designed to accept Webhooks from version control systems, such as GitHub or GitLab, and execute actions based on those webhooks. Specifically, the following tasks:

* Trigger r10k environment and module deploys onto Puppet Servers
* Send notifications to ChatOps systems, such as Slack and RocketChat

## Prerequisites

While there are no prerequisites for running the webhook server itself, for it to be useful, you will need the following installed on the same server or another server for this tool to be useful:

* Puppet Server
* [r10k](https://github.com/puppetlabs/r10k) >= 3.9.0
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
  queue:
    enabled: true
    max_concurrent_jobs: 10
    max_history_items: 20
chatops:
  enabled: false
  service: slack
  channel: "#general"
  user: r10kbot
  auth_token: 12345
  server_uri: "https://rocketchat.local"
r10k:
  config_path: /etc/puppetlabs/r10k/r10k.yaml
  default_branch: main
  allow_uppercase: false
  verbose: true
mappings:
  long-repo-name: lrp
```

### Microsoft Teams notifications

Create an "Incoming Webhook" connector in Teams at the designated channel as described in the documentation: [Create Incoming Webhooks at learn.microsoft.com](https://learn.microsoft.com/en-us/microsoftteams/platform/webhooks-and-connectors/how-to/add-incoming-webhook). Keep the URL confidential!

Configure the service in the `webhook.yaml`:
```yaml
chatops:
  enabled: true
  service: teams
  server_uri: "<Teams Webhook URI>"
```

Notifications are colored, according to their status.
green: Success
red: Failure
orange: Warning

Press the `Details` button to get more information.

If the queue is enabled in the `server` part of `webhook.yaml`, then two notifications are emitted: First, when the request is added to the queue and second, when the request was processed.

### Bolt authentication

Due to the inherent security risk associated with passing plain text passwords to the Bolt CLI tool, all ability to set it within the application have been removed.

Instead, it is recommended to instead utilize the Bolt [Transport configuration options](https://puppet.com/docs/bolt/latest/bolt_transports_reference.html) and place them within the `bolt-defaults.yaml` file.

If you want to utilize an `inventory.yaml` and place the targets and auth config within that file, you can. Just be sure to remember to add the target name containing the nodes you need to the `webhook.yml` file

### Server options

#### `protected`

Type: bool
Description: Enforces authentication via basic Authentication
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

#### `queue`

Type: struct
Description: Struct containing Queue options

##### `enabled`

Type: bool
Description: Should queuing be used
Default: `false`

##### `max_concurrent_jobs`

Type: int
Description: How many jobs could be stored in queue
Default: `10`

##### `max_history_items`
Type: int
Description: How many queue items should be stored in the history
Default: `50`

### ChatOps options

#### `enabled`

Type: boolean
Description: Enable/Disable chatops support
Default: false

#### `service`

Type: string
Description: Which service to use. Supported options: [`slack`, `rocketchat`, `teams`]
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
Description: The ChatOps service API URI to send the message to. For MS Teams, this is the Webhook URL created at the channel connectors.
Default: nil

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

### `use_legacy_puppetfile_flag`

Type: bool
Description: Use the legacy `--puppetfile` flag instead of `--modules`. This should only be used when your version of r10k doesn't support the newer flag.
Default: `false`

### `generate_types`

Type: bool
Description: Run `puppet generate types` after updating an environment
Default: `true`

### `env_incremental`

Type: bool
Description: Use `--incremental` flag when updating an environment
Default: `false`

### `command_path`

Type: `string`
Description: Allow overriding the default path to r10k.
Default: `/opt/puppetlabs/puppetserver/bin/r10k`

### `blocked_branches`

Type: `array of strings`
Description: A list of branches to not allow deployments to.
Default: `[]`

### `mappings`

Type: `map`
Description: A map of long repository names to short names. This is useful for repositories that have long names that are not suitable for use in the URL. This is useful for multi tenant environments where you also want to use a shorter prefix for the environment.
Default: `{}`


## Usage

Webhook API provides following paths

### GET /health

Get health assessment about the Webhook API server

### GET /api/v1/queue

Get current queue status of the Webhook API server

### POST /api/v1/r10k/environment

Updates a given puppet environment, ie. `r10k deploy environment`. This only updates a specific environment governed by the branch name.

### POST /api/v1/r10k/module

Updates a puppet module, ie. `r10k deploy module`. The default behavior of r10k is to update the module in all environments that have it. Module name defaults to the git repository name.

Available URL arguments (`?argument=value`):

* `branch_only=(true|false)` - If set, this will only update the module in an environment set by the branch, as opposed to all environments. This is equivalent to the `--environment` r10k option. DEFAULT: `false`
* `module_name=name` - Sometimes git repository and module name cannot have the same name due to arbitrary naming restrictions. This option forces the module name to be the given value instead of repository name.
