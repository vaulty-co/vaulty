# Example

Create card.json with the following content:

```json
{
	"amount": 100,
	"currency": "USD",
	"card": {
		"number": "424242424242424242",
		"exp": "20/23",
		"holder": "John Doe",
		"cvv": "123"
	}
}
```

Make tokenization request to vaulty:

```
curl vltgW1X5uISWemT.proxy.vaulty.dev:8080/credit-cards -d @card.json
```

# Setup

## Redis

```redis-server /usr/local/etc/redis.conf````
