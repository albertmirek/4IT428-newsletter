# Api specification for User-api
Implemented was only v1


### User Registration
**POST:**`/v1`

##### Expected request body:
```json
{
  "email": "string",
  "password": "string"
}
```



### User Login

**POST**: `/v1/login`

##### Expected request body:
```json
{
  "email": "string",
  "password": "string"
}
```

##### Expected response body:
```json
{
  "token": "string"
}
```


##### Expected response:
HTTP 201

### Update User Password
expectes valid **JWT token** in header (Bearer)

**PUT:**`/v1`

##### Expected request body:
```json
{
  "newPassword": "string"
}
```