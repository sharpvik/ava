# Server Docs

## Routing

```http
/api    => see API
GET /@  => index.html
GET /   => FileServer
```

## API

```http
/api
--> POST /signup    => see signup
--> POST /login     => see login
--> GET /chats      => see chats
--> GET /profile    => see profile
```

### `/signup`

Request:

```json
{
  "handle": "...",
  "name": "...",
  "password": "***"
}
```

Response is a simple `string` message. It is simply to be displayed to the user.
Here's the table of possible status codes and their semantics:

| Status                      | Description                            |
| --------------------------- | -------------------------------------- |
| `500` Internal server error | Usually due to failed password hashing |
| `409` Conflict              | Username is taken                      |
| `200` OK                    | Success                                |

### `/login`

Request:

```json
{
  "handle": "...",
  "password": "***"
}
```

Response is a simple `string` message. It is simply to be displayed to the user.
In case of a successful login, response will also contain a cookie with a JWT
token for further communication. Here's the table of possible status codes and
their semantics:

| Status                      | Description                             |
| --------------------------- | --------------------------------------- |
| `404` Not found             | Username does not exist in the database |
| `401` Unauthorized          | Invalid password                        |
| `500` Internal server error | Failed to generate a JWT token          |
| `200` OK                    | Success                                 |

### `/chats`

Request to chats is a plain `GET` with no data whatsoever. User's identity is
established from their JWT token that is passed as an HTTP-Only cookie on
successful login.

Response:

```json
["COZY CHAT", "Kinks and jinks", "Crushampton (uncensored)", "..."]
```

Here's the table of possible status codes and their semantics:

| Status             | Description                              |
| ------------------ | ---------------------------------------- |
| `401` Unauthorized | JWT cookie not found or token is invalid |
| `200` OK           | Success                                  |

### `/profile`

Request to profile is a plain `GET` with no data whatsoever. User's identity is
established from their JWT token that is passed as an HTTP-Only cookie on
successful login.

Response:

```json
{
  "handle": "sharpvik",
  "name": "Viktor A. Rozenko Voitenko"
}
```

Here's the table of possible status codes and their semantics:

| Status             | Description                              |
| ------------------ | ---------------------------------------- |
| `401` Unauthorized | JWT cookie not found or token is invalid |
| `200` OK           | Success                                  |

### `/avatar?handle=sharpvik`

Request to avatar is a plain `GET` with no data whatsoever. User's identity is
either established from their JWT token or passed explicitly through the URL
parameter (optional) called `handle`.

Response is an image or `404 Not Found` if something went wrong.
