# mockservice

Mock HTTP Service that can register mock endpoints and respond with respective responses.

## Features

- Registering mock endpoints by request by sending a `POST` request to an URL path of your choice
- Registering the mock registration endpoint and bulk loading pre-determined requests of your choice

## Upcoming Features

- Registering endpoints that can send callbacks hooks back to your service
- Registering endpoints that return specific responses based on timing; used to mock APIs that require polling to keep on top of statuses

## Examples

### Adding a mock service with `/mocks` as the registration endpoint

```go
mockService, err := mockservice.New("/mocks")
if err != nil {
    log.Fatalf("Failed to created mock service %s", err)
}
http.ListenAndServe(":8080", mockService)
```

To register HTTP responses, send `POST` requests to `/mocks` with a JSON request body like the following

```json
{
    "method": "GET",
    "endpoint": " ",
    "httpStatusCode": 204,
    "responseBody": "",
    "responseHeaders": { "foo": "bar" }
}
```

### Adding a mock service by bulk loading via a configuration

Note: JSON and XML configuration files can also be loaded by unmarshaling into `mockservice.Conf`. See [here](https://github.com/wchan2/mock_service/blob/master/mock_service.go#L25-L28) and [here](https://github.com/wchan2/mock_service/blob/master/endpoints.go#L27-L35) for details.

```go
conf := mockservice.Conf{
    RegistrationEndpoint: "/mocks",
    Endpoints: []*mockservice.MockEndpoint{
        {
            Method:          http.MethodPost,
            Endpoint:        "/mock/test",
            StatusCode:      http.StatusCreated,
            ResponseBody:    `hello world`,
            ResponseHeaders: map[string]string{"Foo": "Bar"},
        },
    },
}
service, err := mockservice.NewWithConf(&conf)
```

## Contributing

In order to contribute, please:

1. Check if an issue that you'd like to contribute a fix for doesn't already exist
2. Create an issue for something you'd like fixed
3. Fork the repository
4. Make a pull request

Please also be sure to follow the following guidelines.

- `go fmt` all your code
- `golint` all your code

## Contributors

- [enriqueChen](https://github.com/enriqueChen)
