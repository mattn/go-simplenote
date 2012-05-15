package simplenote

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	email string
	token string
}

func getUrl(path string, c *Client, values url.Values) string {
	u := fmt.Sprintf("https://simple-note.appspot.com%s?auth=%s&email=%s",
		path, url.QueryEscape(c.token), url.QueryEscape(c.email))
	if values != nil {
		u += "&" + values.Encode()
	}
	return u
}

func NewClient(email, password string) (*Client, error) {
	s := fmt.Sprintf("email=%s&password=%s",
			url.QueryEscape(email), url.QueryEscape(password))
	sr := strings.NewReader(base64.StdEncoding.EncodeToString([]byte(s)))
	uri := "https://simple-note.appspot.com/api/login"
	r, err := http.Post(uri, "application/x-www-form-urlencoded", sr)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return &Client {email, string(b)}, nil
}

type notes struct {
	Count int
	Data []Note
}

type Note struct {
	Tags []string
	SystemTags []string
	Key string
	ModifyDate string
	CreateDate string
	Deleted int
	Content string
	MinVersion int
}

func (c *Client) GetNotes() ([]Note, error) {
	u := getUrl("/api2/index", c, url.Values {
		"length": []string{"10"},
		"mark": []string{""},
	})
	r, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var ns notes
	d := json.NewDecoder(r.Body)
	err = d.Decode(&ns)
	if err != nil {
		return nil, err
	}
	return ns.Data, nil
}

func (c *Client) GetNote(note *Note) error {
	u := getUrl("/api2/data/" + note.Key, c, nil)
	r, err := http.Get(u)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	err = d.Decode(note)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) UpdateNote(note *Note) error {
	u := getUrl("/api2/data/" + note.Key, c, nil)
	s := fmt.Sprintf("content=%s&tags=%s",
			url.QueryEscape(note.Content),
			url.QueryEscape(strings.Join(note.Tags, ",")))
	sr := strings.NewReader(s)
	r, err := http.Post(u, "application/x-www-form-urlencoded", sr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	d := json.NewDecoder(r.Body)
	err = d.Decode(note)
	if err != nil {
		return err
	}
	return nil
}
