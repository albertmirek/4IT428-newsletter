# Api specification for mailing-api
Implemented was only v1

This api handles sending of Newsletter posts to subscribers

### Send newsletter post to subscribers
expectes valid **JWT token** in header (Bearer)

POST:`/v1/{newsletterId}/send/{postId}`


##### Expected response:
HTTP 204