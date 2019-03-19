# pw-gen-api
A very simple password generation API.

## Run
1. Download latest binary from [the releases page](https://github.com/mactr0n/pw-gen-api/releases/latest) or build binary:
    ```
    $ chmod +x .build.bash
    $ ./build.bash [VERSION]
    ```
2. Executing the binary will start a REST API server on port `3334`.

## API

**URL** : `/passwords`

**Method** : `GET`

**Query Params** : 
  * `numCandidates`: number of password candidates (*`int`*, *optional*, *default=4*)
  * `length` of password candidates (*`int`*, *optional*, *default=8*)
  * `numDigits` the number of digits (*`int`*, *optional*, *default=0*)
  * `numSymbols` the number of symbols (*`int`*, *optional*, *default=0*)
  * `replaceVowels` replace (german) vowels of password candidates randomly (*`bool`*, *optional*, *default=false*)

### Success Responses

**Code** : `200 OK`

**Content** :
```
[
    "candidates": [
        "password_candidate_1",
        "password_candidate_1",
        ...
        "password_candidate_n"
    ]
]
```

## Curl Example:
```
$ curl -X GET "http://localhost:3334/passwords?length=32&numDigits=4&numSymbols=4&replaceVowels=true"
```