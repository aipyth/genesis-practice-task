# genesis-practice-task
Web API for current BTC exchange rate with user authentication.

## API Description
### Auth
Authentication is implemented with Token in body of the request.
Example:
```json
{
  "token": "4eae45b75a0f72b12628"
}
```
**The following queries are available:**
- `POST` `/user/create` Create user record.
```json
{
  "email": [string],
  "password": [string]
}
```

Possible response variants:

|  Code  |  Body      |  Explanation  |
| :------: | :------: | :--------------: |
| 201    |    -       | The user is successfully created |
| 500    |    -       | Server Error |
| 400    | email or password is not provided | |
| 409    |    -       | This user already exists |

- `POST` `/user/login` Get the token required for other API requests.
```json
{
  "email": [string],
  "password": [string]
}
```
Possible response variants:

| Code  |  Body    |  Explanation  |
| :---: | :------: | :--------------: |
| 200   | `{ "token": [string] }` | Authentication successful. Use this for following requests. |
| 500   |    -     | Server Error |
| 403   |    -     | No such user with this email or password |

### BTC Exchange Rate

- `GET` `/btcRate` Get current BTC exchange rate in UAH.
```json
{
  "token": [string]
}
```

Possible response variants:

|  Code  |  Body  |  Explanation  |
| :------: | :------: | :--------------: |
| 200    | `[float]` | Current BTC exchange rate in UAH. |
| 500    |   -    | Server Error |
| 503    |   -    | The BTC information providing service responded with an error or is currently unavailable. |


## Project Structure

`cmd` - available main applications. Currently, there is only server.  
`internal` - private applications: storage, server handlers.



1. `internal/storage`  
`common.go` - contains needed interface for any storage in the project. The 
   API may use any implementation of these: csv, json (not implemented), etc.  
`csv.go` - `CSVStorage` is file storage, that stores information in a form 
   of csv. It reads into memory and maps all data while calling `Connect()` 
   method. The information is stored in memory till it's written into the 
   file by calling `Save()`.
   
2. `internal/api` The server itself. To run it we simply call `Run()`.  
`api.go` - contains `Server` struct definition and methods to run it. Here 
   we specify the usage of `CSVStorage`.  
`auth.go` - handlers to provide authentication and authorization.     
`authMiddleware.go` - middleware to provide proper user auth check.  
`btc.go` - contains only `GetBTCRateInUAH` handler.
