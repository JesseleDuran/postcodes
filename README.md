# **Postcodes API** #
---

[![Go version](https://img.shields.io/badge/Go-1.15-blue.svg?logo=go)]()


This service is in charge of giving information about the postcodes given a
 lat and lng.


## Basic architecture

This microservice, when starting, loads all the zones (or geographic divisions)
of postal codes, taken from this repository: https://github.com/missinglink/uk-postcode-polygons

Once the polygons are loaded and indexed, we proceed to "hydrate" them, this
means that with the centroid of each polygon, a request is made to the
 postcodes API (https://postcodes.io/), and relates it to the polygon. As a
fallback, if the API does not return any result at the centroid of the polygon, 
the polygon is tesselated into level 11 cells (4km x 4km) and with the
 centroid of each cell, we call the postcode API until it returns a valid
 result.

A total of 2728 polygons with their respective zip code are loaded into memory,
and it is only 28Mb in memory.

Every time the main endpoint is called with a set of coordinates, with each 
coordinate it searches in its index of areas in which area those coordinates
 are contained (with the ray tracing algorithm), and returns it and its
  corresponding postcode.

Since the memory access is very fast and the Point in polygon algorithm,
 the requests have a p99 of 300ms.
 
This microservice is faithful to the "Memory Access Pattern".

## Code Pattern

The code pattern used in the structure of this service is DDD (Domain Driven
 Design).

## How to run locally

Just run

```
docker-compose up -d
```

## API
---

### Get postcodes given a set of lat/lng
```http
POST /postcodes/v1/postcodes
```

- **EXAMPLE BODY:**

```javascript
// POST /postcodes/v1/postcodes

{
    "coordinates": [
        {
            "lat": 53.024182,
            "lon": -2.210694
        },
        {
            "lat": 51.417427,
            "lon": -0.080494
        }
    ]
}
```

- **EXAMPLE RESPONSE:**

```javascript
[
    {
        "lat": 53.024182,
        "lon": -2.210694,
        "postcode": "ST5 0JS"
    },
    {
        "lat": 51.417427,
        "lon": -0.080494,
        "postcode": "SE19 3BA"
    }
]
```

## Git

The branching model used for Git is a very simple GitFlow, with its
 conrresponding main, develop and feature branches (the only ones necessaries).

---