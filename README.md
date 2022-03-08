# httpclient

http request client with in build request and response logging.

## Integration

either you can import default config and edit it out according to your needs or you can use `Config` struct to create your own config. for example:

```
// passing the config as nil.
// internally http client uses
// default config
client, err := httpclient.New()
if err != nil {
	log.Println(err)
}
```

or 

```
// initialize the config
hConfig := httpclient.Config{
	Timeout: 10 * time.Second,
}

// passing the http config while creating new http client
client, err := httpclient.New(hConfig...)
if err != nil {
	log.Println(err)
}
```

then you can use this client inside your code.

## Example

below are the examples which gonna help you to get started with the the integration.

[example](example/)