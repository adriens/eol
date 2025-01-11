# EOL CLI

This CLI application lists items from the End of Life API.

## How to build

To build the application, follow these steps:

1. Open a terminal and navigate to the root directory of your project.

2. Run the following command to compile the application:

    ```sh
    go build -o eol cmd/main.go
    ```

    This command will create an executable file named [eol](http://_vscodecontentref_/1) in the current directory.

## How to run

To run the compiled application, use the following command:

```sh
./eol -l
```

```sh
./eol --list
```

This command will list all items from the End of Life API.

```sh
# End of life of my favorite stack
./eol maven java quarkus
```
