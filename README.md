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
# Get target Trello board ID, follow instructions here: https://community.atlassian.com/t5/Trello-questions/How-to-get-Trello-Board-ID/qaq-p/1347525
TRELLO_BOARD_ID=<>
# Create Notion Integration key here: https://www.notion.so/my-integrations
NOTION_INTEGRATION_KEY=<>
```

Run the CLI.

```bash
go run cmd/main.go
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
