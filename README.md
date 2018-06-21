# go-servicebus-api

### Usage
The following environment variables will need to be set

* SB_URL: url for the service bus endpoint
* SB_KEY: shared key for service bus
* SB_KEYTYPE: type of shared key
* AUTH_URL: authentication endpoint for validating JWT tokens
* SB_PORT: port to run the application (default is 8001)

Main endpoint is <server_name>/message

The Authorization header needs to be set with a valid token from the authentication source specified in the AUTH_URL environment variable

Example header: Authorization: Bearer token
  
Messages must be sent via POST method and the message body must follow the format shown for the SbMessage struct [here](https://github.com/waustinlynn/go-servicebus)
