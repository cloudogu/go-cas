package cas

import (
	"net/url"
	"net/http"
	"github.com/golang/glog"
	"fmt"
	"path"
	"io/ioutil"
)

// https://apereo.github.io/cas/4.2.x/protocol/REST-Protocol.html

// TicketGrantingTicket represents a SSO session for a user, also known as TGT
type TicketGrantingTicket string

// ServiceTicket stands for the access granted by the CAS server to an application for a specific user, also known as ST
type ServiceTicket string


// RestOptions provide options for the RestClient
type RestOptions struct {
	CasURL     *url.URL
	Client     *http.Client
}

// RestClient uses the rest protocol provided by cas
type RestClient struct {
	casUrl     *url.URL
	client     *http.Client
}

// NewRestClient creates a new client for the cas rest protocol with the provided options
func NewRestClient(options *RestOptions) *RestClient {
	if glog.V(2) {
		glog.Infof("cas: new rest client with options %v", options)
	}

	var client *http.Client
	if options.Client != nil {
		client = options.Client
	} else {
		client = &http.Client{}
	}

	return &RestClient{
		casUrl:      options.CasURL,
		client:      client,
	}
}

// RequestGrantingTicket returns a new TGT, if the username and password authentication was successful
func (c *RestClient) RequestGrantingTicket(username string, password string) (TicketGrantingTicket, error) {
	// request:
	// POST /cas/v1/tickets HTTP/1.0
	// username=battags&password=password&additionalParam1=paramvalue

	endpoint, err := c.casUrl.Parse(path.Join(c.casUrl.Path, "v1", "tickets"))
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("username", username)
	values.Set("password", password)

	resp, err := c.client.PostForm(endpoint.String(), values)
	if err != nil {
		return "", err
	}

	// response:
	// 201 Created
	// Location: http://www.whatever.com/cas/v1/tickets/{TGT id}

	if resp.StatusCode != 201 {
		return "", fmt.Errorf("ticket endoint returned status code %v", resp.StatusCode)
	}

	tgt := path.Base(resp.Header.Get("Location"))
	if tgt == "" {
		return "", fmt.Errorf("does not return a valid location header")
	}

	return TicketGrantingTicket(tgt), nil
}

// RequestServiceTicket requests a service ticket with the TGT for the given service url
func (c *RestClient) RequestServiceTicket(tgt TicketGrantingTicket, serviceUrl *url.URL) (ServiceTicket, error) {
	// request:
	// POST /cas/v1/tickets/{TGT id} HTTP/1.0
	// service={form encoded parameter for the service url}
	endpoint, err := c.casUrl.Parse(path.Join(c.casUrl.Path, "v1", "tickets", string(tgt)))
	if err != nil {
		return "", err
	}

	values := url.Values{}
	values.Set("service", serviceUrl.String())

	resp, err := c.client.PostForm(endpoint.String(), values)
	if err != nil {
		return "", err
	}

	// response:
	// 200 OK
	// ST-1-FFDFHDSJKHSDFJKSDHFJKRUEYREWUIFSD2132

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("service ticket endoint returned status code %v", resp.StatusCode)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return ServiceTicket(data), nil
}

// Logout destroys the given granting ticket
func (c *RestClient) Logout(tgt TicketGrantingTicket) error {
	// DELETE /cas/v1/tickets/TGT-fdsjfsdfjkalfewrihfdhfaie HTTP/1.0
	endpoint, err := c.casUrl.Parse(path.Join(c.casUrl.Path, "v1", "tickets", string(tgt)))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", endpoint.String(), nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("could not destroy granting ticket %v, server returned %v", tgt, resp.StatusCode)
	}

	return nil
}