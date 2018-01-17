package translate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const DefaultURL = "https://translate.yandex.net/api/v1.5/tr.json"

const codeOK = 200

var ErrNoKey = errors.New("translate: emtpy key")

type Config struct {
	Key string
	URL string
}

// New initializes a Client with Config.
func New(cfg Config) (*Client, error) {
	if cfg.Key == "" {
		return nil, ErrNoKey
	}
	if cfg.URL == "" {
		cfg.URL = DefaultURL
	}
	c := &Client{key: cfg.Key, url: cfg.URL}
	return c, nil
}

type Client struct {
	url string
	key string
}

// GetLangs returns a list of supported directions and languages.
//
// https://tech.yandex.ru/translate/doc/dg/reference/getLangs-docpage/
func (c *Client) GetLangs() (directions []string, langs map[string]string, err error) {
	resp, err := c.postForm("/getLangs", url.Values{})
	if err != nil {
		return
	}
	defer resp.Body.Close()

	v, err := unmarshal(resp.Body)
	if err != nil {
		return
	}
	return v.Directions, v.Langs, nil
}

// Detect returns a lang of the text.
//
// https://tech.yandex.ru/translate/doc/dg/reference/detect-docpage/
func (c *Client) Detect(text string, hints ...string) (string, error) {
	hint := strings.Join(hints, ",")

	form := url.Values{"text": []string{text}}
	if hint != "" {
		form.Set("hint", hint)
	}

	resp, err := c.postForm("/detect", form)
	if err != nil {
		return Unknown, err
	}
	defer resp.Body.Close()

	v, err := unmarshal(resp.Body)
	if err != nil {
		return Unknown, err
	}
	return v.Lang, nil
}

// Translate translates the text to the lang.
//
// https://tech.yandex.ru/translate/doc/dg/reference/translate-docpage/
func (c *Client) Translate(text string, lang string, opts ...Opt) ([]string, error) {
	r := &request{Text: text, To: lang}
	for _, opt := range opts {
		opt(r)
	}

	resp, err := c.postForm("/translate", r.Values())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v, err := unmarshal(resp.Body)
	if err != nil {
		return nil, err
	}
	return v.Text, nil
}

func (c *Client) postForm(path string, data url.Values) (*http.Response, error) {
	data.Set("key", c.key)
	return http.DefaultClient.PostForm(c.url+path, data)
}

func unmarshal(r io.Reader) (*response, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	var v *response
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, fmt.Errorf("translate: %v: %q", err, string(b))
	}
	if v.Code != nil && *v.Code != codeOK {
		return nil, fmt.Errorf("translate: %v: %v", *v.Code, v.Message)
	}
	return v, nil
}

type request struct {
	From string
	To   string
	Text string
}

func (r *request) Values() url.Values {
	var lang = r.To
	if r.From != Unknown {
		lang = fmt.Sprintf("%s-%s", r.From, r.To)
	}
	return url.Values{
		"text": []string{r.Text},
		"lang": []string{lang},
	}
}

// Opt allows to make a request more concrete.
type Opt func(*request)

// From defines the source language of a text to translate.
func From(lang string) Opt {
	return func(r *request) {
		r.From = lang
	}
}

type response struct {
	Code       *int              `json:"code"`
	Message    string            `json:"message"`
	Directions []string          `json:"dirs"`
	Langs      map[string]string `json:"langs"`
	Lang       string            `json:"lang"`
	Text       []string          `json:"text"`
}
