# ava

The microservice for storing your users' avatars (or any kind of images for that
matter).

Start it up via Docker or `docker-compose`. It will listen on port `42069`.
You could create a `.env` file where you _may or may not_ specify an
`AVA_API_KEY` variable like so:

```bash
# .env
AVA_API_KEY=mySuperSecretAPIKey
```

If you don't specify any `AVA_API_KEY` it is assumed to be `""` and any
request with a missing `Ava-API-Key` header will be authorized to proceed.

## API

```
POST /upload/jpg
```

Put your image file into the request body and enjoy. This will respond with a
simple string (e.g. `166ebdf66f4174b5.jpg`) that you'll be able to use to
retrieve the image.

```
GET /download/166ebdf66f4174b5.jpg
```

Given the string, send this request to get your image back. It will be attached
to the response body with a proper MIME type header set.
