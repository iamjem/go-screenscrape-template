This repository contains a simple screen scrape CLI written in Go. You can read the full blog post [here](http://www.pixeldonor.com/2015/mar/01/golang-screen-scrape-cli/).

# Environment Variables
The CLI uses the following environment variables:

``` text
DEBUG=true
DATABASE_URL=user=youruser dbname=yourdbname sslmode=disable
EMAIL_FROM=fromuser@gmail.com
EMAIL_HOST=smtp.gmail.com
EMAIL_HOST_PASSWORD=emailpassword
EMAIL_HOST_USER=hostuser@gmail.com
EMAIL_PORT=587
EMAIL_TO=touser@gmail.com

```

# Database Setup
This skeleton is configured to use Postgres. Change `lib/db.go` if you prefer something else.

To run migrations with [goose](https://bitbucket.org/liamstask/goose) run `goose -env development up`.

# Running the CLI
This project uses [godep](https://github.com/tools/godep) for dependency management.

1. Clone the repository and `cd` into it.
2. Run `godep go install` to build and install binary.
3. Run `./run_screenscrape`
