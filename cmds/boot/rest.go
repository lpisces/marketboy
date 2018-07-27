package boot

import "gopkg.in/resty.v1"

type (
	Endpoint struct {
		Scheme string
		Host   string
		Prefix string
		Path   string
		Verb   string
	}

	Body interface{}

	Header map[string]string

	Request struct {
		*Endpoint
		*Body
		*Header
	}
)

func (r *Request) Do() (resp resty.Response) {
	// endpoint
	ep := fmt.Sprintf("%s://%s%s%s", r.Endpoint.Scheme, r.Endpoint.Host, r.Endpoint.Prefix, r.Endpoint.Path)
	return
}
