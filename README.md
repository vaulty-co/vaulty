https://github.com/vaulty-co/vaulty/workflows/Go/badge.svg
    
# Vaulty

Vaulty is a reverse and forward proxy that modifies (encrypt, decrypt, tokenize, etc.) data on the fly and securely stores it in a safe. Vaulty can be used for the following:

- Anonymize data before it reaches your APIs and backends
- Get encryption / decryption for your APIs without changing a line of code
- Provide filtered data for specific group of users like support stuff, etc.
- Tokenize credit cards, SSNs, etc. for companies that are PCI compliant
- Encrypt your customers data when you import this data from 3d party services (access tokens, PII, etc.)

Currently you can play with Vaulty, think how you would like to use it and share your ideas and feedback so we can make it work. It's not ready for production yet.

## Try it now!

Let's create simple routes.json file with transformation rules just to try Vaulty:

```json
{
    "vault":{
        "upstream":"https://postman-echo.com/"
    },
    "routes":{
        "inbound":[
            {
                "method":"POST",
                "path":"/post",
                "request_transformations":[
                    {
                        "type":"json",
                        "expression":"card.number",
                        "action":{
                            "type":"encrypt"
                        }
                    }
                ]
            }
        ]
    }
}
```

In short, all requests to Vaulty will be transformed and then send to [http://postman-echo.com](http://postman-echo.com) which is echo server and will display all data it receives.

Now, let's run Vaulty as a proxy:

```bash
docker run -p 8080:8080 -v ${PWD}:/vaulty/.vaulty/ vaulty 
```

You should see something like this:

```
==> Vaulty proxy server started on port 8080! in development environment
```

Let's make a request with card number:

```bash
curl http://127.0.0.1:8080/post \
  -d '{ "card": { "number": "4242424242424242", "exp": "10/22" } }' \
  -H "Content-Type: application/json"
```

In response you will see that Vauly has encrypted card number and our upstream has received encrypted data instead of plain card number.

```
{"args":{},"data":{"card":{"number":"NDI0MjQyNDI0MjQyNDI0Mg(demo encryption)","exp":"10/22"}},"files":{},"form":{},"headers":{"x-forwarded-proto":"https","x-forwarded-port":"443","host":"127.0.0.1","x-amzn-trace-id":"Root=1-5ec1412f-6ab8d3f28110822b8a425e81","content-length":"83","user-agent":"curl/7.64.1","accept":"*/*","content-type":"application/json","accept-encoding":"gzip"},"json":{"card":{"number":"NDI0MjQyNDI0MjQyNDI0Mg(demo encryption)","exp":"10/22"}},"url":"https://127.0.0.1/post"}%
```

More information about Vaulty can be found here: https://vaulty.co
