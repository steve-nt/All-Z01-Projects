# forum

This project can be used to host a forum on the web. It serves via the HTTP
protocol and uses a sqlite3 database to store information.

It supports:

- categorized posts,
- comments on posts,
- likes and dislikes per post/comment,
- user login and registration,
- attaching an image to a post,
- user can login with Google and Github OAuth,
- notifications per comment to post and/or reactions to post/comment,
- log of all user activities,
- edit or remove posts/comments.

## Build

To build the project you can use the command either directly:

```
export CGO_ENABLED=1; go build -o forum cmd/forum/main.go
```

or by using the make tool:

```
make
```

We advise you to use `make`. Resulted binary executables will be at `./bin`
directory.

## Docker

Build the image with:

```bash
docker build -t forum-app .
```

Start it with:

```bash
docker run
    -e GOOGLE_CLIENT_ID="<YOUR_GOOGLE_CLIENT_ID_HERE>" \
    -e GOOGLE_CLIENT_SECRET="<YOUR_GOOGLE_CLIENT_SECRET_HERE>" \
    -e GITHUB_CLIENT_ID="<YOUR_GITHUB_CLIENT_ID_HERE>" \
    -e GITHUB_CLIENT_SECRET="<YOUR_GITHUB_CLIENT_SECRET_HERE>" \
    -e BASE_URL="http://localhost:<PORT>/" \
    -p <PORT>:8080 forum-app
```

You can alternatively use a `.env` file. An example is provided with the
repository at `./.env_example`.

```bash
docker run --env-file .env -p <PORT>:8080 forum-app
```

Note:
`BASE_URL` must always end with a trailing slash!

## Usage

### `forum`

You can start it with:

```
./bin/forum [--db-path <DB_PATH>] [IP] [PORT]
./bin/forum [IP] [PORT]
./bin/forum [PORT]
./bin/forum
```
and navigate to `http://[IP]:[PORT]`.

Make sure you add your API keys and `BASE_URL` when executing:

```bash
GOOGLE_CLIENT_ID="YOUR_GOOGLE_CLIENT_ID_HERE" \
GOOGLE_CLIENT_SECRET="YOUR_GOOGLE_CLIENT_SECRET_HERE" \
GITHUB_CLIENT_ID="YOUR_GITHUB_CLIENT_ID_HERE" \
GITHUB_CLIENT_SECRET="YOUR_GITHUB_CLIENT_SECRET_HERE" \
BASE_URL="http://localhost:8080/" \
./bin/forum
```

This would result to start the server listening at: http://localhost:8080/

#### Defaults

- `DB_PATH`: ./db.db
- `IP`: 127.0.0.1
- `PORT`: 8080

### `mockgen`

`mockgen` can be used to generate some sample data for your DB. Use it by
running:

```
./bin/mockgen [DB_PATH]
```

#### Defaults

- `DB_PATH`: ./db.db

### `hashgen`

`hashgen` can be used to generate hashes compatible with the ones stored in the
database. If you want to manually add a user inside the database, you can use:

```
./bin/hashgen <PASSWORD>
```

This will return you hash, which you can enter in your custom query. After
creating the user, you'd be able to login with the password you used to generate
the hash.

## Run tests

To run the tests you can use the command:
```
go test -v ./...
```

or

```
make test
```
