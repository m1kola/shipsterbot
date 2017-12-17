# Shipster bot

Shipster is a stupid bot that helps you maintain your shopping lists.

This is deliberately over-engineered project:
I wanted to play with the Go language, its standard library and with some popular libraries.
But I also wanted to produce something useful,
something I'll, probably, use.

## Compilation instructions

We currently support only MacOS and Linux platforms. To compile binary on
your machine you need to take the following steps.


1. Clone the project using `git`:

    ```
    git clone git@github.com:m1kola/shipsterbot.git
    ```

    or download using the `go get` tool:

    ```
    go get github.com/m1kola/shipsterbot
    ```

2. Compile a binary. From the project root directory run the build process:

    ```
    make build
    ```

We use `make` to automate some actions. `make build` will take care of
dependencies for your. You would be able to find a binary `shipsterbot` for your
platform in the root directory of the project.


## Other `make` commands

Compile into a custom location:

```
make build OUTPUT_BIN=/your/custom/location/app
```

Run tests:
```
make test
```

Clean up the working directory:

```
make clean
```
