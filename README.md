# Go-Logger-MDC

A Go library for logging with MDC. This library is a wrapper around the [zap](https://godoc.org/go.uber.org/zap) library.
Using context in Go, it supports MDC (Mapped Diagnostic Context) and provides a simple way to add MDC to your logging. 

# How To Use It

Here is an example of a simple program that uses the library to log requestID and userID when a request is sent by the user

This program makes use of the context field of http.Request struct in Go. We will save the info that we want to log in the context field of the http.request in the middleware function.


```
func Verify(r *http.Request) error {
        tokenString := Extract(r)
	requestID, _ := uuid.NewUUID()
	ctx := r.Context() // Get the context from the request
	//Adding info that we want to log throughout the request
	ctx = context.WithValue(ctx, "requestID", requestID.String())
	ctx = context.WithValue(ctx, "userID", tokenString)
	ctx = context.WithValue(ctx, "Accept", r.Header.Get("Accept"))
	ctx = context.WithValue(ctx, "Accept-Encoding", r.Header.Get("Accept-Encoding"))
	ctx = context.WithValue(ctx, "Connection", r.Header.Get("Connection"))
	ctx = context.WithValue(ctx, "Content-Length", r.Header.Get("Content-Length"))
	ctx = context.WithValue(ctx, "Content-Type", r.Header.Get("Content-Type"))
	req := r.WithContext(ctx) // Set the context to the request
	*r = *req // Set the request to the new request with the context

```

Then we initialize the logger using the SugarLog(ctx) function in the libary and then use the logger to log to console in JSON/console format.
Passing in the http.request in the functions where you will need to log and you can log out the info that you want to see throughout the entire request.
````
func CheckJwt(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		err := jwt.Verify(r)

		logger := logger.SugarLog(r.Context()) \\initialize the logger

		if err != nil {
			logger.Errorf("JWT is invalid")
			res.ERROR(w, 401, errors.New("Unauthorized"))
			return
		}
		logger.Infow("JWT Valid") \\logging the info that we want

		next(w, r, ps)

	}
}
```
