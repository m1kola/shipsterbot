# Shipster bot

Shipster is a stupid bot that helps you maintain your shopping lists.

This is deliberately over-engineered project:
I wanted to play with the Go language, its standard library
and with some popular libraries.
But I also wanted to produce something useful,
something I'll, probably, use.

## Compilation instructions

To build the application you need

* MacOS or Linux OS (amd64)
* Go 1.9
* GNU `make` utility


To compile an application on your machine you need to take the following steps:


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

## How to run the application

Currently Shipster bot is only implemented for the Telegram messenger.
To run the bot you need the following:

* PostgreSQL (9.6 or higher is recommended)
* Telegram [bot token](https://core.telegram.org/bots#6-botfather)
* [Register a webhook](https://core.telegram.org/bots/api#setwebhook)
  using the Telegram Bot API which involve
  using a [verified TLS certificate](https://core.telegram.org/bots/webhooks#a-verified-supported-certificate)
  or generating a [self-signed certificate](https://core.telegram.org/bots/webhooks#a-self-signed-certificate).
  Note that in case of self-signed certificate you must upload a public key
  while registering your webhook.

You would need to define the following required environment variables to run the bot:

* `DATABASE_URL` - connection string for your PostgreSQL
  database. `postgres://localhost/shipster?sslmode=disable`, for example
* `TELEGRAM_API_TOKEN` - Telegram bot token

There are some optional environment arguments:

* `TELEGRAM_TLS_CERT_PATH` and `TELEGRAM_TLS_KEY_PATH` - Path to your TLS
  certificate and key files. Note that to successfully run the Shipster bot
  you **must** use either a verified TLS certificate
  or a self-signed one. These env vars are optional,
  because you might want to run the bot behind
  a web server (nginx, for example) that terminates
  secure connections and proxies requests to your application
  in a private network
* `TELEGRAM_WEBHOOK_PORT` - Possible values: `443`, `80`, `88`, `8443`.
  Default is `8443`.

  Telegram requires us to run a bot on any of the listed ports,
  in order to be able to deliver webhooks.
  `8443` doesn't require root privileges, so it seems like a sensible default
* `DEBUG` - Possible values: `true` and `false`. Default is `false`.

  Enables the debug mode. In the debug mode bot produces
  more verbose logs

So to run the application you need to do the following:

1. Create a database
2. Set env vars:

    ```
    # Required env vars
    export DATABASE_URL=postgres://localhost/shipster?sslmode=disable
    export TELEGRAM_API_TOKEN=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11

    # Optional
    export TELEGRAM_TLS_CERT_PATH=/path/to/your/telegram_public.pem
    export TELEGRAM_TLS_KEY_PATH=/path/to/your/telegram_private.key
    export DEBUG=trueÂ
    ```
3. Run migrations database migrations:

    ```
    ./shipsterbot migrate up
    ```

4. Run the application:

    ```
    ./shipsterbot startbot telegram
    ```
