import requests
import json

def get_token():
    data = requests.get("http://127.0.0.1:8080/login?username=mohit&password=asdas234-9nhj324brf9f834")
    extracted = json.loads(data.content)
    return extracted["token"]

def readPayload(format="xml"):
    response = ""
    filename = "data"
    if format == "json":
        filename += ".json"
    else:
        filename += ".xml"
    with open(filename, "r") as f:
        for l in f.readlines():
            response += l
    return response

def post_data_to_server(data, token):
    url = "http://127.0.0.1:8080/fetch"
    headers = {"Authorization":"Bearer " + token,
               "Content-Type":"application/json"}
    r = requests.post(url, data=data, headers=headers)
    return r.content

token = get_token()
payload = readPayload("json")
print "pushing payload \n {0}to server {1}".format(payload, "127.0.0.1:8080")
response = post_data_to_server(payload, token)
print response


