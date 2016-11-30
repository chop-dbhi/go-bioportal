package bioportal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-querystring/query"
)

var (
	DefaultTimeout = 10 * time.Second
	BaseURL        = "https://data.bioontology.org"
)

type APIError struct {
	Status int      `json:"status"`
	Errors []string `json:"errors"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, strings.Join(e.Errors, ", "))
}

type BaseOptions struct {
	// Format string `url:"format"`
	Include        string `url:"include,omitempty"`
	Pagesize       int    `url:"pagesize,omitempty"`
	Page           int    `url:"page,omitempty"`
	IncludeViews   bool   `url:"include_views"`
	DisplayContext bool   `url:"display_context"`
	DisplayLinks   bool   `url:"display_links"`
}

func DefaultBaseOptions() *BaseOptions {
	return &BaseOptions{
		// Format:         "json",
		Pagesize:       0,
		Page:           1,
		IncludeViews:   false,
		DisplayContext: true,
		DisplayLinks:   true,
	}
}

type Client struct {
	APIKey string
	HTTP   *http.Client
}

func (c *Client) request(path string) (*http.Request, error) {
	var u string

	if strings.HasPrefix(path, "/") {
		x, err := url.Parse(BaseURL)
		if err != nil {
			panic(err)
		}
		x.Path = path
		u = x.String()
	} else {
		u = path
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("apikey token=%s", c.APIKey))
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) Send(path string, params interface{}) (io.ReadCloser, error) {
	req, err := c.request(path)
	if err != nil {
		return nil, err
	}

	if params != nil {
		v, err := query.Values(params)
		if err != nil {
			panic(err)
		}

		req.URL.RawQuery = v.Encode()
	}

	var resp *http.Response

	for {
	MAKE_REQ:
		resp, err = c.HTTP.Do(req)
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode {
		case http.StatusTooManyRequests:
			time.Sleep(100 * time.Millisecond)
			goto MAKE_REQ

		case http.StatusRequestURITooLong:
			return nil, fmt.Errorf("request URI too long:\n%s", path)
		}

		if resp.StatusCode != 200 {
			b, _ := ioutil.ReadAll(resp.Body)
			var apiErr APIError
			if err := json.Unmarshal(b, &apiErr); err != nil {
				panic(fmt.Sprintf("[%s] %s:\n%s", resp.Status, err, string(b)))
			}
			return nil, &apiErr
		}

		break
	}

	return resp.Body, nil

}

func (c *Client) search(opts *SearchOptions) (io.ReadCloser, error) {
	if opts.Query == "" {
		return nil, errors.New("query cannot be empty")
	}

	return c.Send("/search", opts)
}

func (c *Client) SearchRead(w io.Writer, opts SearchOptions) (int64, error) {
	rc, err := c.search(&opts)
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	return io.Copy(w, rc)
}

func (c *Client) Search(opts SearchOptions) (*SearchResult, error) {
	rc, err := c.search(&opts)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var res SearchResult
	if err := json.NewDecoder(rc).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) recommend(opts *RecommendOptions) (io.ReadCloser, error) {
	if len(opts.Terms) == 0 {
		return nil, errors.New("at least one term is required")
	}

	return c.Send("/recommender", opts)
}

func (c *Client) RecommendRead(w io.Writer, opts RecommendOptions) (int64, error) {
	rc, err := c.recommend(&opts)
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	return io.Copy(w, rc)
}

func (c *Client) Recommend(opts RecommendOptions) ([]*RecommendResult, error) {
	rc, err := c.recommend(&opts)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var res []*RecommendResult
	if err := json.NewDecoder(rc).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) annotate(opts *AnnotateOptions) (io.ReadCloser, error) {
	if opts.Text == "" {
		return nil, errors.New("text cannot be empty")
	}

	return c.Send("/annotator", opts)
}

func (c *Client) AnnotateRead(w io.Writer, opts AnnotateOptions) (int64, error) {
	rc, err := c.annotate(&opts)
	if err != nil {
		return 0, err
	}
	defer rc.Close()

	return io.Copy(w, rc)
}

func (c *Client) Annotate(opts AnnotateOptions) ([]*AnnotationResult, error) {
	rc, err := c.annotate(&opts)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var res []*AnnotationResult
	if err := json.NewDecoder(rc).Decode(&res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) Class(ontology, class string) (*Class, error) {
	path := fmt.Sprintf("/ontologies/%s/classes/%s", ontology, url.QueryEscape(class))
	rc, err := c.Send(path, nil)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var cl Class
	if err := json.NewDecoder(rc).Decode(&cl); err != nil {
		return nil, err
	}

	return &cl, nil
}

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTP: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}
