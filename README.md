# About this project

This project was made within 2 hours as part of [GitHub Copilot 1-Day Build Challenge](https://dev.to/devteam/join-us-for-the-github-copilot-1-day-build-challenge-3000-in-prizes-3o2i?bb=202755).

# EOL CLI

This CLI application lists items from the End of Life API.

[![YouTube Video](https://img.youtube.com/vi/zcXKJ1KPn28/0.jpg)](https://www.youtube.com/watch?v=zcXKJ1KPn28)



## How to build


```sh
go build -o eol cmd/main.go
export PATH=$PATH:.

```


## How to run

To run the compiled application, use the following command:

```sh
eol -h

```

```sh
eol -l
```

```sh
eol --list
```

This command will list all items from the End of Life API.

```sh
# End of life of my favorite java stack
eol maven java quarkus
```
