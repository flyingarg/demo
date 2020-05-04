# Image Dump

The application hosts a URL that allows clients to dump images onto the server.
The clients are to be authorized using access tokens. Clients have quotas on the maximum number of 
dumps that are allowed within an hour. 

## Description
The following are the set of expectations from the service. 
* Expose a http end point that takes the following xml and its equivalent json as input

```buildoutcfg
<items>
   <item>
     <name>Tomato Soup</name>
  <manufacturer>Cafe world</manufacturer>
  <brand>Cafe</brand>
  <category>Food</category>
  <images>
    <url>https://dummyimages/test.png</url>
 <url>https://dummyimages/test1.png</url>
  ....
  </images>
</item>
<item> ...... </item>
<item> ...... </item>
</items>
```

The number of items in a request should be more than 50.
The response must have a status field that returns the processed status of each item.
Images need to be stored in file system.


## Implementation
* Authentication Mechanism : 
    * User authenticates with the system by passing f(username, hash(password))
    * A JWT token gets generated for the user. This token gets signed using the secret coded into the demoservice.go
    file. This is a symmetric key implementation of signature.
    * The client would be using this token in their consecutive requests to authorize its requests.
    * Tokens have an expiration duration of 1 hour. Can me overridden using env variable TOKEN_EXPIRY (seconds).
* Authorization Mechanism :
    * Users who have successfully authenticated with the system, will get a token.
    * This token has an expiry time. The user will be asked to authenticate again if the token expires.
    * When the user does this, the column ```tokens.user_secret``` is regenerated.
* Customer Rate Limits
    * Gold : 100 requests within an hour
    * Silver : 50 requests within an hour
    * Bronze: 25 requests within an hour
* Rate Limit Mechanism:
    * Every user request gets logged to the table ```user_requests```.
    * Each request has the following attributes
        * user
        * requestID
        * request_data_format - json/xml
        * data blob
        * items need to have additional itemID and processStatus
        * images need to have additional imageID, URL, Status
    * The user gets a sync response when all the images have been downloaded.
    * The saving of a request to the system is performed as a transaction. This transaction rolls back if the rate limit
    is found to have been exceeded.
        
## Tables
refer testdata/mysql.init.sql

## Build
```buildoutcfg
    $ make clean
    $ make build
```

## Test
Integration test cases are added to testdata. 

* Database setup
Reset/Create a fresh mysql database called demo with data and tables initialized
```buildoutcfg
    $ bash mysql.init.sql
```

* Test the application
Run the application using 
```buildoutcfg
    $ ./bin/demoservice
```
Once the application is running, use the test_json.py and test_xml.py to perform integration "tests".
These scripts, hit the server, perform a login, which provides the access tokens. Then these application use the data 
files to perform requests on the server. These can be run using the commands
```buildoutcfg
    $ cd testdata
    $ python test_json.py
```

