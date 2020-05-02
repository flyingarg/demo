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
    * Assume that the users with the following details have been created on the server. A sample text file users.txt 
    contains 3 users belonging to 3 user classes.
        * username
        * password
        * plan
    * User authenticates with the system by passing f(username, hash(password))
    * A JWT token gets generated for the user. This token is generated using signed using the secret 
    saved in a table ```tokens``` with last updated timestamp.
    * The client would be using this token in their consecutive requests to authorize its requests.
    * Tokens have an expiration duration of 20 minutes. Can me overridden using env variable TOKEN_EXPIRY (seconds).
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
* Image Upload status
    a. SUCCESS
    b. FAILURE
        i. HTTP ERROR CODE
        ii. INTERNAL ERRORS (io/network errors)
        
## Tables
Users table
```
create table users(
    userid int, 
    username varchar[50],
    plan varchar[10],
    PRIMARY KEY (userid))
```

Authentication keys table
```
create table authentication(
    userid int,
    passwordhash varchar[64]
    PRIMARY KEY (userid)
)
```

Authorization token
```
create table tokens(
    userid int,
    user_secret blob,
    last_updated timestamp
    PRIMARY KEY (userid)
)
```

User Requests
```
create table requests(
    id MEDIUMINT NOT NULL AUTO_INCREMENT,
    userid int,
    dataformat varchar[4],
    data blob,
    status bool,
    PRIMARY KEY (userid)
)
```