# Api specification for Newsletter-api
Implemented was only v1



### Create Newsletter
**POST:**`/v1/newsletter`

expectes valid **JWT token** in header (Bearer)

##### Expected request body:
```json
{
  "name": "string",
}
```

##### Expected response:
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "name": "string"
}
```

### UpdateNewsletter
expectes valid **JWT token** in header (Bearer)

**PUT:**`/v1/newsletter/{id}`
##### Expected request body:
```json
{
  "name": "string"
}
```

##### Expected response:
HTTP 204


### Delete Newsletter
expectes valid **JWT token** in header (Bearer)

**DELETE:**`/v1/newsletter/{id}`

##### Expected response:
HTTP 204




### Post for Newsletter
expectes valid **JWT token** in header (Bearer)

**POST:**`/v1/newsletter/{id}/post`

##### Expected request body:
```json
{
  "heading": "string",
  "body": "string"
}
```

##### Expected response:
HTTP 201


### Subscribe to Newsletter
expects valid **JWT token** in header (Bearer)

**POST:**`/v1/newsletter/{id}/subscribe`

#### Expected response:
HTTP 200


### Unsubscribe from Newsletter

**GET** `/v1/newsletter/{token}/unsubscribe`

Data from token will determine which newsletter to unsubscribe from


#### Expected response:
HTTP 200
