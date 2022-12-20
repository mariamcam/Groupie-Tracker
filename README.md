
# Groupie Tracker

Groupie Trackers consists on receiving a given API and manipulate the data contained in it, in order to create a site, displaying the information.



## Usage/Examples

Clone the repository and start the server

```bash
  go run .
```

go to http://127.0.0.1:8080/


## Implementation details

`main.go` Creates multiplexer and starts server on 8080 port. When the request reaches the server, a multiplexer will inspect the URL being requested and redirect the request to the correct handler fucntion

`handlers.go`

Handlers fucntions process requests. When the processing is complete, the handler passes the data to the template engine, which will use templates to generate HTML to be returned to the client.

##
Will be updated

## Author

- [@Tlep](https://www.github.com)/Tlepkali

