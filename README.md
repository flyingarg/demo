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

* Images need to stored in file system or database blobs. Better implemented with an interface.
* Access-Token based authentication and authorization per customer. Authorization categories
    * Gold : 100 requests within an hour
    * Silver : 50 requests within an hour
    * Bronze: 25 requests within an hour
    
* Image Upload status
    a. SUCCESS
    b. FAILURE
        i. HTTP ERROR CODE
        ii. INTERNAL ERRORS (io/network errors)