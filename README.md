# Webhook Go

Webhook Go is a port of the [puppet_webhook](https://github.com/voxpupuli/puppet_webhook) Sinatra API server to Go.
This is designed to be more streamlined, performant, and easier to ship for users than the Sinatra/Ruby API server.

This server is a REST API server designed to accept Webhooks from version control systems, such as GitHub or GitLab, and execute actions based on those webhooks. Specifically, the following tasks:

* Trigger r10k environment and module deploys onto Puppet Servers
* Send notifications to ChatOps systems, such as Slack, RocketChat, and Microsoft Teams

## Prerequisites

While there are no prerequisites for running the webhook server itself, for it to be useful, you will need the following installed on the same server or another server for this tool to be useful:

* Puppet Server
* [r10k](https://github.com/puppetlabs/r10k) >= 3.9.0 (or [g10k](https://github.com/xorpaul/g10k))
* Windows or Linux server to run the server on. MacOS is not supported.

## Installation

### Packages and binaries

Released artifacts are published on the [Releases](https://github.com/voxpupuli/webhook-go/releases) page, including:

* `deb` and `rpm` packages (linux amd64, arm, and arm64)
* Archive files for linux (amd64, arm, arm64), windows (amd64), and freebsd (amd64, arm64)

The `deb` and `rpm` packages install the `webhook-go` binary to `/usr/bin/webhook-go`, a configuration file to `/etc/voxpupuli/webhook.yml` (marked as `config|noreplace`, so it is not overwritten on upgrades), and a systemd unit (`webhook-go.service`) that enables and starts the service on install.

For a manual binary installation, download the archive for your platform, extract the `webhook-go` binary, make it executable, and place it somewhere in your `PATH`.

### Docker

Versioned container images are published to [ghcr.io/voxpupuli/webhook-go](https://github.com/voxpupuli/webhook-go/pkgs/container/webhook-go):

```sh
docker run -d -p 4000:4000 ghcr.io/voxpupuli/webhook-go:2.15.1
```

The image is built `FROM scratch` and contains only the `webhook-go` binary and a default configuration file at `/webhook.yml` (which is in the default config search path). To use your own configuration, mount it over the baked-in file, or pass the `--config` flag (arguments after the image name are appended to the `webhook-go server` entrypoint):

```sh
docker run -d -p 4000:4000 \
  -v /etc/voxpupuli/webhook.yml:/webhook.yml:ro \
  ghcr.io/voxpupuli/webhook-go:2.15.1
```

Note that the server executes `r10k` (or `g10k`) as a local process, and the published image does not contain one. To deploy from a container, derive your own image that includes an `r10k` binary, or mount one in and point `r10k.command_path` at it.

### From source

Building from source requires [Go](https://go.dev/dl/) >= 1.25:

```sh
git clone https://github.com/voxpupuli/webhook-go.git
cd webhook-go
go build -o webhook-go .
./webhook-go server --config ./build/webhook.yml
```

## Configuration

The Webhook API server uses a configuration file called `webhook.yml` to configure the server. Several of the required options have defaults pre-defined so that a configuration file isn't needed for basic function.

When the `--config` flag is not given, the server looks for a `webhook.yml` in the following locations, in order:

* `./webhook.yml` (current working directory)
* `/etc/voxpupuli/webhook/webhook.yml`
* `../config/webhook.yml`
* `config/webhook.yml`

Example `webhook.yml`:

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
    max_history_items: 50
chatops:
  enabled: false
  service: slack
  channel: "#general"
  user: r10kbot
  auth_token: 12345
  server_uri: "https://rocketchat.local"
r10k:
  config_path: /etc/puppetlabs/r10k/r10k.yaml
  default_branch: master
  allow_uppercase: false
  verbose: true
mappings:
  long-repo-name: lrp
```

### Server options

#### `protected`

Type: bool
Description: Enforces authentication via Basic Authentication on all `/api/v1/*` routes. The `/health` route always stays unauthenticated.
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

Type: int
Description: Port to run the server on. Optional.
Default: `4000`

#### `deploy_on_success_only`

Type: bool
Description: Only run a deploy if the pipeline/workflow associated with the webhook (e.g. a GitHub `workflow_run` or GitLab `pipeline` event) completed successfully. Failed pipelines are rejected with `424 Failed Dependency`.
Default: `false`

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
Description: Struct containing Queue options. When enabled, deploy requests are queued and processed by a background worker instead of running synchronously; the API responds immediately with the queued job.

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
yellow: Warning

Press the `Details` button to get more information.

If the queue is enabled in the `server` part of `webhook.yaml`, then two notifications are emitted: First, when the request is added to the queue and second, when the request was processed.

### r10k options

#### `config_path`

Type: string
Description: Full path to the r10k configuration file. Optional.
Default: `/etc/puppetlabs/r10k/r10k.yaml`

#### `default_branch`

Type: string
Description: Name of the default branch for r10k to pull from, used when the webhook payload doesn't include a branch. Optional.
Default: `master`

#### `prefix`

Type: string
Description: An r10k prefix to apply to the environment being deployed, producing `prefix_branch` environment names. Optional. Set to the special value `mapping` to look up the prefix per repository in the `mappings` section instead.
Default: `nil`

#### `allow_uppercase`

Type: bool
Description: Allow Uppercase letters in the module, branch, or environment name. Optional.
Default: `false`

#### `verbose`

Type: bool
Description: Log verbose output when running the r10k command
Default: `true`

#### `deploy_modules`

Type: bool
Description: Deploy modules in environments.
Default: `true`

#### `use_legacy_puppetfile_flag`

Type: bool
Description: Use the legacy `--puppetfile` flag instead of `--modules`. This should only be used when your version of r10k doesn't support the newer flag.
Default: `false`

#### `generate_types`

Type: bool
Description: Run `puppet generate types` after updating an environment
Default: `true`

#### `env_incremental`

Type: bool
Description: Use `--incremental` flag when updating an environment (only used with the `--modules` deploy mode, not the legacy `--puppetfile` flag)
Default: `false`

#### `command_path`

Type: `string`
Description: Allow overriding the default path to r10k.
Default: `/opt/puppetlabs/puppetserver/bin/r10k`

#### `blocked_branches`

Type: `array of strings`
Description: A list of branches to not allow deployments to.
Default: `[]`

#### `ignore_branch_prefixes`

Type: `array of strings`
Description: A list of branch name prefixes to skip. A branch matching one of these
is answered with `200 OK` and no r10k run, which is what you want for branches that
are not environments at all. Set it to the same value as `ignore_branch_prefixes` in
your `r10k.yaml`: r10k refuses to deploy those branches, so without this the webhook
runs r10k, r10k fails with `Environment(s) '<name>' cannot be found in any source`, and
the sender records a failed delivery. Unlike `blocked_branches`, which matches whole
branch names and answers `403`, this matches by prefix.
Default: `[]`

Example, for a repository that has Renovate and Dependabot enabled:

```yaml
r10k:
  ignore_branch_prefixes:
    - renovate/
    - dependabot/
```

#### `use_g10k_commands`

Type: `boolean`
Description: Use g10k commands instead of r10k (parameters are a little bit different)
Default: `false`

### `mappings`

Type: `map`
Description: A map of long repository names to short names. This is useful for repositories that have long names that are not suitable for use in the URL. This is useful for multi tenant environments where you also want to use a shorter prefix for the environment. Used when `r10k.prefix` is set to `mapping`.
Default: `{}`

## Usage

Start the server with:

```sh
webhook-go server --config /path/to/webhook.yml
```

The binary is a small CLI (run `webhook-go --help` for details) with a single subcommand, `server`, plus the usual `help` and `completion` commands.

Webhook API provides following paths

### GET /health

Get health assessment about the Webhook API server

```sh
curl http://localhost:4000/health
```

```json
{"message":"running"}
```

### GET /api/v1/queue

Get current queue status of the Webhook API server. Returns the list of queued jobs, including each job's id, name, command, state (`added`, `success`, `failed`), timestamps, and the command output once processed. This endpoint requires authentication when `server.protected` is enabled.

```sh
curl -u puppet:puppet http://localhost:4000/api/v1/queue
```

### POST /api/v1/r10k/environment

Updates a given puppet environment, ie. `r10k deploy environment`. This only updates a specific environment governed by the branch name.

Available URL arguments (`?argument=value`):

* `no_mods=(true|false)` - If set, this will only update an environment with no modules (no flags `--modules` or `--puppetfile`). This option is usefull when you need to update only files (for example [separate hiera data](https://github.com/puppetlabs/r10k/blob/main/doc/dynamic-environments/configuration.mkd#separate-hiera-data))

Example, sending a GitHub `push` event for the `production` branch:

```sh
curl -u puppet:puppet -X POST http://localhost:4000/api/v1/r10k/environment \
  -H "X-GitHub-Event: push" \
  -H "Content-Type: application/json" \
  -d '{"ref":"refs/heads/production","deleted":false,"repository":{"name":"webhook-go","full_name":"voxpupuli/webhook-go","owner":{"name":"voxpupuli"}}}'
```

On success the server responds with `202 Accepted`: the raw r10k output when the queue is disabled, or the queued job when it is enabled. Deleting a ref or pushing a branch matching an `ignore_branch_prefixes` entry is answered with `200 OK` and nothing is deployed; a `blocked_branches` entry is answered with `403 Forbidden`.

### POST /api/v1/r10k/module

Updates a puppet module, ie. `r10k deploy module`. The default behavior of r10k is to update the module in all environments that have it. Module name defaults to the git repository name.

Available URL arguments (`?argument=value`):

* `branch_only=(true|false)` - If set to any non-empty value, this will only update the module in an environment set by the branch, as opposed to all environments. This is equivalent to the `--environment` r10k option. DEFAULT: `false`
* `module_name=name` - Sometimes git repository and module name cannot have the same name due to arbitrary naming restrictions. This option forces the module name to be the given value instead of repository name. Must match `^[a-z][a-z0-9_]*$`.

Example, deploying the `stdlib` module only in the environment of the pushed branch:

```sh
curl -u puppet:puppet -X POST "http://localhost:4000/api/v1/r10k/module?branch_only=true&module_name=stdlib" \
  -H "X-GitHub-Event: push" \
  -H "Content-Type: application/json" \
  -d '{"ref":"refs/heads/main","deleted":false,"repository":{"name":"puppetlabs-stdlib","full_name":"puppetlabs/puppetlabs-stdlib","owner":{"name":"puppetlabs"}}}'
```

### Supported webhook providers

The server detects the provider from the event headers each request carries and parses the payload accordingly. Supported providers and events:

| Provider | Detected via | Supported events |
|---|---|---|
| GitHub | `X-GitHub-Event` | `push`, `workflow_run` |
| GitLab | `X-Gitlab-Event` | `Push Hook`, `Pipeline Hook` |
| Bitbucket Cloud | `X-Event-Key` + `X-Hook-UUID` | `repo:push` |
| Bitbucket Server | `X-Event-Key` + `X-Request-Id` | `repo:refs_changed` |
| Azure DevOps | `X-Azure-DevOps` | `git.push` |
| Gitea | `X-Gitea-Event` | `push`, `delete` |

Requests from other providers, or with unsupported event types, fail with `500` and an explanatory message.

## Development

The [Makefile](Makefile) provides the usual targets; run `make help` to list them:

* `make run` - run the server from source with `./build/webhook.yml`
* `make test` - run the test suite (`go test ./...`)
* `make binary` - build a local binary into `bin/` with GoReleaser
* `make compile` - build for all supported OS/architectures (snapshot)
* `make snapshot` - build all release artifacts without publishing (snapshot)
* `make release` - build and publish a release (requires a `v*` tag)
* `make clean` - clean up build and dependency artifacts

Building release artifacts requires [GoReleaser](https://goreleaser.com/). CI runs on every push/PR to `master` and builds and tests with Go 1.25.

## License

This project is licensed under the Apache 2.0 license. See [LICENSE](LICENSE) for details.
