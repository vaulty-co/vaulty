package main

import (
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

func errResponse(r *http.Request, message string) *http.Response {
	return goproxy.NewResponse(r,
		goproxy.ContentTypeText,
		http.StatusBadGateway,
		message)
}

func transformRequestBody(req *http.Request) error {
	err := transformBody(req)
	if err != nil {
		return err
	}
	return nil

	// Sidekiq
	//
	// 1. Generate request ID
	// 2. Subscribed for request ID channel

	// 3. Create Sidekiq payload
	// 4. Encode payload to JSON
	// 5. Push to default redis queue

	// 6. Read reply from request ID channel
	// 7. Get body from redis requests ID key
	// 8. Decode body from JSON
	// 9. Continue

	// send body to sidekiq
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     os.Getenv("REDIS_URL"),
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// pubsub, err := client.Subscribe("mychannel")

	// wait for sidekiq response
	// read modified body
	// update request body with new one

	// client.LPush("queue:default", "{\"queue\":\"default\",\"class\":\"ProxyWorker::Worker\",\"args\":\"{}\",\"jid\":\"ecbdb1e0927b71c5228839e\",\"enqueued_at\":1584172882.16241,\"at\":1584172882.16241}")

	// body, err := ioutil.ReadAll(req.Body)
	// if err != nil {
	// 	return err
	// }

	// oldBody := string(body)
	// newBody := " + new body is here!!!"

	// buf2 := bytes.NewBufferString(oldBody + newBody)
	// req.Header.Del("Content-Length")
	// req.ContentLength = int64(buf2.Len())
	// req.Body = ioutil.NopCloser(bufio.NewReader(buf2))
	// return nil
}

func main() {

	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		err := transformRequestBody(req)
		if err != nil {
			return nil, errResponse(req, err.Error())
		}

		return req, nil
	})

	proxy.Verbose = true

	log.Fatal(http.ListenAndServe(":8080", proxy))
}
