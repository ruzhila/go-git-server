# A git-http-server with 100 lines of Golang code

By [ruzhila.cn](http://ruzhila.cn/?from=github_git_server), a campus for learning backend development through practice.

This is a tutorial code demonstrating how to use Golang write git server. Pull requests are welcome. 👏

## Features
- 🚀 Simple and easy to use, not need nginx or fcgiwrap
- 👏 Support create repository
- 📦 Only 100 lines of code

## Build from source
```shell
$ git clone https://github.com/ruzhila/go-git-server.git
$ cd go-git-server
$ go build
```

## Install from source
```shell
$ go install github.com/ruzhila/go-git-server@latest
# Run the server
$ go-git-server
```
## Usage
```shell
$ ./go-git-server -h
Usage of ./go-git-server:
  -addr string
        address to listen on (default "127.0.0.1:8080")
  -create
        create repository only without serving
  -prefix string
        prefix path for git server (default "/")
  -root string
        root repository path (default "/tmp/git")

$ ./go-git-server -addr 127.0.0.1:8080 -root /var/git

# Create a new repository
$ ./go-git-server -root /var/git -create hellorepos

# Visit the repository
$ git clone http://127.0.0.1:8080/hellorepos
$ cd hellorepos
$ echo "Hello, Git!" > README.md
$ git add README.md
$ git commit -m "Initial commit"
$ git push origin master

```