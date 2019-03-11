# pw-gen-api
A very simple password generation API.

## API

**URL** : `/passwords`

**Method** : `GET`

**Query Params** : 
  * `numCandidates`: number of password candidates (*`int`*, *optional*, *default=4*)
  * `length` of password candidates (*`int`*, *optional*, *default=8*)
  * `numDigits` the number of digits (*`int`*, *optional*, *default=0*)
  * `numSymbols` the number of symbols (*`int`*, *optional*, *default=0*)
  * `replaceVowels` replace (german) vowels of password candidates randomly (*`bool`*, *optional*, , *default=false*)

### Success Responses

**Code** : `200 OK`

**Content** :
```
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