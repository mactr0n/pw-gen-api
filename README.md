# pw-gen-api
A very simple password generation API.

## API

**URL** : `/passwords`

**Method** : `GET`

### Success Responses

**Code** : `200 OK`

**Content** :
```json
[
    "candidates": [
        "password1",
        "password2",
        ...
    ]
]
```

## Build
```
$ chmod +x .build.bash
$ ./build.bash [VERSION]
``` 