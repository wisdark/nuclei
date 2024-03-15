package component

import (
	"context"

	"github.com/pkg/errors"
	"github.com/projectdiscovery/nuclei/v3/pkg/fuzz/dataformat"
	"github.com/projectdiscovery/retryablehttp-go"
	urlutil "github.com/projectdiscovery/utils/url"
)

// Query is a component for a request query
type Query struct {
	value *Value

	req *retryablehttp.Request
}

var _ Component = &Query{}

// NewQuery creates a new query component
func NewQuery() *Query {
	return &Query{}
}

// Name returns the name of the component
func (q *Query) Name() string {
	return RequestQueryComponent
}

// Parse parses the component and returns the
// parsed component
func (q *Query) Parse(req *retryablehttp.Request) (bool, error) {
	if req.URL.Query().IsEmpty() {
		return false, nil
	}
	q.req = req

	q.value = NewValue(req.URL.Query().Encode())

	parsed, err := dataformat.Get(dataformat.FormDataFormat).Decode(q.value.String())
	if err != nil {
		return false, err
	}
	q.value.SetParsed(parsed, dataformat.FormDataFormat)
	return true, nil
}

// Iterate iterates through the component
func (q *Query) Iterate(callback func(key string, value interface{}) error) error {
	for key, value := range q.value.Parsed() {
		if err := callback(key, value); err != nil {
			return err
		}
	}
	return nil
}

// SetValue sets a value in the component
// for a key
func (q *Query) SetValue(key string, value string) error {
	if !q.value.SetParsedValue(key, value) {
		return ErrSetValue
	}
	return nil
}

// Delete deletes a key from the component
func (q *Query) Delete(key string) error {
	if !q.value.Delete(key) {
		return ErrKeyNotFound
	}
	return nil
}

// Rebuild returns a new request with the
// component rebuilt
func (q *Query) Rebuild() (*retryablehttp.Request, error) {
	encoded, err := q.value.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "could not encode query")
	}
	cloned := q.req.Clone(context.Background())
	cloned.URL.RawQuery = encoded

	// Clear the query parameters and re-add them
	cloned.Params = nil
	cloned.Params = urlutil.NewOrderedParams()
	cloned.Params.Decode(encoded)
	cloned.Update()
	return cloned, nil
}
