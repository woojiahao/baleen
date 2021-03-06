# baleen

Custom tool for migrating from Trello to Notion.

## Installation

```bash
git clone https://github.com/woojiahao/baleen.git
cd baleen/
```

Create a file `.env` at the root of the project directory with the following fields filled in:

```bash
# .env
# Create Trello API key here: https://trello.com/app-key
TRELLO_API_KEY=<>
# Create Trello token by selecting "Token" in the page above
TRELLO_TOKEN=<>
# Create Notion Integration key here: https://www.notion.so/my-integrations
NOTION_INTEGRATION_KEY=<>
```

Run the CLI for the available commands.

```bash
go run cmd/main.go

NAME:
   main - A new cli application

USAGE:
   main [global options] command [command options] [arguments...]

COMMANDS:
   export   exports a Trello board and creates a save file (to import, use "baleen import <save path>"
   import   imports saved cards into Notion (cards saved from "baleen migrate" or "baleen export")
   migrate  exports a Trello board into the integrated Notion page (full flow from saving exports to importing to Notion)
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --board value, -b value   specify name of Trello board (default: "Programming Bucket")
   --config value, -c value  specify the configuration JSON for the migration (default: "configs/conf.json")
   --env value, -e value     specify the environment file hodling the API keys (default: ".env")
   --help, -h                show help (default: false)
```

## Motivation

I started `baleen` as a personal project to migrate my evergrowing Trello board to a custom Notion workspace to host all
of the old content.

I was inspired to work on this when I had saw this error: ![Motivation](./res/motivation.png). I had somehow managed to
nearly max out Trello's card limit - that's quite amazing!

`baleen` is named after the [Baleen Whales](https://en.wikipedia.org/wiki/Baleen_whale) which embark on large migrations
and their mouths contain baleen plates to act as filters for planktonic creatures.

## Why Go?

I chose Go for this project as I wanted to access both the Trello and Notion APIs hassle-free while building a CLI. I
did not choose to use Javascript (even though there are pre-built libraries for each API) because I want to experiment
with the latest Go features.

## TODO

### Development

- [ ] Migrate any image attachments to be hosted on Google Drive and embed them into the page
- [ ] Ignore any attachments taht start with https://docs.google.com/viewer?embedded=true
- [ ] Add support for image attachments

### Documentation

- [ ] Talk about how the Notion API is incredibly nested and how certain values need to be explicitly filled in as the marshalled JSON does not include them by default
