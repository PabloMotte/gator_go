# gator_go

RSS feed aggregator in Go

## Requirements

The following is required to run the aggregator:

- [Go](https://go.dev/) - Cilck the link if you don't already have it installed
- [PostgreSQL](https://www.postgresql.org/) - same here

If you are using Windows, use [WSL](https://learn.microsoft.com/en-us/windows/wsl/install)
  
## Installation

Clone the repo, then use `go install` and type `gator_go` to get started

### List Users

`gator_go users`  

### Add Users

`gator_go register <user>`  

### Login a User

`gator_go register <user>`  

### List Available Feeds

`gator_go feeds`  

### Add Feeds

`gator_go addfeed <url>`  
Add and follow a feed not already in gator_go. You must be a logged in user.

### Follow Feeds

`gator_go follow <url>`  
Follow the feed already in gator_go. You can only use this command as a logged in user.

### Unfollow Feeds

`gator_go unfollow <url>`  
Stop following the specified feed. You can only use this command as a logged in user.

### Aggregate Feeds

`gator_go agg time`  
where time is a string representing the duration between feed downloads, such as 1m0s for "one minute and 0 seconds". It is usually best to run this in a separate process.

### Browse

`gator_go browse [limit]`  
Browse the latest feed entries. Defaults to listing only the latest 2. You can only browse feeds if you are a logged in user.
  
## Notes on Development

Followed a tutorial on [boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang)
  
Developed using [Go](https://go.dev/) on [WSL](https://en.wikipedia.org/wiki/Windows_Subsystem_for_Linux), and [PostgreSQL](https://www.postgresql.org/):

- go version go1.27.3 linux/amd64
- PostgreSQL 16.15 (Ubuntu 16.15-0ubuntu0.24.04.1) on x86_64-pc-linux-gnu, compiled by gcc (Ubuntu 13.3.0-6ubuntu2~24.04.1) 13.3.0, 64-bit
- DB work through psql, [SQLC](https://sqlc.dev/) and [Goose](https://github.com/pressly/goose)
- - goose version: v3.29.0
- - SQLC version: v1.31.1
