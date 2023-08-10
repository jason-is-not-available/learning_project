# Liquor Inventory

## Description

In this project we will be making an HTTP API to track inventory of liquor in your house. It will be unauthenticated, and will use JSON for response payloads and requests. It will operate on port 8090. Use a routing framework (like gin) to do this.


### Endpoints

#### GET `/liquors`
This endpoint should return the entire structure of liquors and their amounts in JSON format. An example response is below.

```
GET /liquors

Response (200):

[
  {"type": "bourbon", "amount": 4},
  {"type": "vodka", "amount": 1},
  ...
]
```

#### GET `/liquors/[type]`
This endpoint should return the amount filtered by liquor. An example response is below. The filter should allow for substring matches. i.e. a request with `schnapps` should also return `peach schnapps`

```
GET /liquors/bourbon

Response (200):

{"type": "bourbon", "amount": 4}
```

#### POST `/liquors/add`
This endpoint should receive a json object and add the amount to the existing amount (or create the new entry). An example POST is below. The response should be the corresponsing current total amount.

```
POST /liquors/add {"type": "bourbon", "amount": 4}

Response (200):

{"type": "bourbon", "amount": 8}
```

#### POST `/liquors/remove`
This endpoint should receive a json object and remove the amount from the existing amount. If the number requested is more than the current total, a 500 error should be thrown.

An example POST is below. The response should be the corresponsing current total amount.

```
POST /liquors/remove {"type": "bourbon", "amount": 4}

Response (200):

{"type": "bourbon", "amount": 4}

OR

Response (500):

"Not Enough Liquor"
```

### Additional Notes:

All types in json and in queries should be case INSENSITIVE.
