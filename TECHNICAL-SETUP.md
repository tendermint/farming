# Technical Setup

To ensure you have a successful experience working with our farming module, we recommend this technical setup.

## Github Integration

Click the GitHub icon in the sidebar for GitHub integration and follow the prompts.

Clone the repos you work in

- Fork or clone the https://github.com/tendermint/farming repository.

Internal Tendermint users have different permissions, if you're not sure, fork the repo.

## Software Requirement

To build the project:

- [Golang](https://golang.org/dl/) v1.16 or higher
- [make](https://www.gnu.org/software/make/) to use `Makefile` targets

## Development Environment Setup

Setup git hooks for conventional commit. 

1. Install [`pre-commit`](https://pre-commit.com/)

2. Run the following command:
    ```bash
    pre-commit install --hook-type commit-msg
    ```

3. (Optional for macOS users) Install GNU `grep`:

4. Run the following command
    ```bash
    brew install grep
    ```

5. Add the line to your shell profile:
    ```bash
    export PATH="/usr/local/opt/grep/libexec/gnubin:$PATH"
    ```

Now, whenever you make a commit, the `pre-commit` hook will be run to check if the commit message conforms [Conventional Commit](https://www.conventionalcommits.org/) rule.

## Building

To build the farming module node and command line client, run the `make build` command from the project's root folder. The output of the build will be generated in the `build` folder.

For cross-builds use the standard `GOOS` and `GOARCH` env vars. i.e. to build for windows:

```bash
GOOS=windows GOARCH=amd64 make build
```

## Installation

To install the node client on your machine, run `make install` command from the project's root folder. 

> 💡 you can also use the default `go` command to build the project, check the content of the [Makefile](https://github.com/tendermint/farming/blob/main/Makefile#L90) for reference

## Testing

Run `make test-all` command to run tests.

> 💡 you can also use the default `go` command to build the project, check the content of the [Makefile](https://github.com/tendermint/farming/blob/main/Makefile#L145) for reference

## Getting Started

The farming module uses [Starport](https://docs.starport.network/guide/). You can install it by simply running the following command. If you want to know more about the installation, you can go to the [official docs](https://docs.starport.network/guide/install.html).

```bash
curl https://get.starport.network/starport! | bash
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

```bash
starport chain serve
```
