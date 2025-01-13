# httpmock

[![Build](https://github.com/NdoleStudio/httpmock/actions/workflows/ci.yml/badge.svg)](https://github.com/NdoleStudio/httpmock/actions/workflows/ci.yml)
[![GitHub contributors](https://img.shields.io/github/contributors/NdoleStudio/httpmock)](https://github.com/NdoleStudio/httpmock/graphs/contributors)
[![GitHub license](https://img.shields.io/github/license/NdoleStudio/httpmock?color=brightgreen)](https://github.com/NdoleStudio/httpmock/blob/master/LICENSE)
![Docker Pulls](https://img.shields.io/docker/pulls/ndolestudio/httpmock)
[![Netlify Status](https://api.netlify.com/api/v1/badges/6a751c80-ac38-4fa0-a470-3d2a69f98dfc/deploy-status)](https://app.netlify.com/sites/httpmock/deploys)

This is a mock http server which can be used to test HTTP requests and responses when building an HTTP client.
You can also use it to mock a backend API in your frontend app e.g in in situations where you're still waiting for the
backend API to be ready.

## Server

You can use the mock server for free at httpmock.dev. The server will return the data which you specify using custom http request headers.
The server will use these headers to generate the response to your http request.

- `response-body`: This can any string which will be returned as the response body by the server
- `response-headers`: This should be a JSON array of headers that will be returned by the server for the request
- `response-status`: This is the HTTP status code of the response e.g `500`, `200`, `404`
- `response-delay`: This is the time in milliseconds that the server will wait before returning the response e.g `1000` for 1 second. The max delay you can set is 10 seconds. If you provide a larger delay, it will capped at 10 seconds.

```bash
curl -X GET https://httpmock.dev/server \
  -H 'response-body: {"id": 12334, "name": "e.g John Doe"}' \
  -H 'response-headers: [{"Content-Type":"application/json"}]' \
  -H 'response-status: 200'
```

or this is an example with javascript

```js
fetch("https://httpmock.dev/server", {
  headers: {
    "response-body": '{"id": 12334, "name": "John Doe"}',
    "response-headers":
      '[{"Content-Type":"application/json"}, {"x-request-id":"dea576ed-ba18-4dd3-baa7-7c865c14b444"}]',
    "response-status": 200,
    "response-delay": 1000,
  },
  method: "GET",
});
```

## Credits

- Color Palette: https://coolors.co/palette/606c38-283618-fefae0-dda15e-bc6c25
- Icon: https://www.svgrepo.com/svg/361904/json-ld

## License

This project is licensed under the GNU AFFERO GENERAL PUBLIC LICENSE Version 3 - see the [LICENSE](LICENSE) file for details
