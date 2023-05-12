package api

type Client interface {
	Get(event, category string) (response *Response, err error)
}
