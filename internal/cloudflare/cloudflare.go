package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const BaseUrl = "https://api.cloudflare.com/client/v4"

type CloudflareClient struct {
	host       string
	token      string
	email      string
	httpClient *http.Client
}

func NewClient(token string) *CloudflareClient {
	client := &http.Client{}
	return &CloudflareClient{
		host:       BaseUrl,
		httpClient: client,
		token:      token,
	}
}

func (c *CloudflareClient) do(method, endpoint string, params map[string]string, body []byte) (*http.Response, error) {
	baseUrl := fmt.Sprintf("%s/%s", c.host, endpoint)
	req, err := http.NewRequest(method, baseUrl, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("authorization", "Bearer "+c.token)

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()

	return c.httpClient.Do(req)
}

type CloudflareError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type DnsRecord struct {
	Id      string `json:"id"`
	Content string `json:"content"`
	Name    string `json:"name"`
}
type DnsRecordResponse struct {
	Success bool              `json:"success"`
	Result  []DnsRecord       `json:"result"`
	Errors  []CloudflareError `json:"errors"`
}

func (c *CloudflareClient) GetDnsRecords(zoneId string) (resp DnsRecordResponse, err error) {
	endpoint := fmt.Sprintf("zones/%s/dns_records", zoneId)
	params := map[string]string{
		"type": "A",
	}
	res, err := c.do(http.MethodGet, endpoint, params, nil)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	if !resp.Success {
		return resp, errors.New("Cloudflare API call failed")
	}

	return resp, err
}

type UpdateRecordResponse struct {
	Success bool              `json:"success"`
	Errors  []CloudflareError `json:"errors"`
}

func (c *CloudflareClient) UpdateDnsRecord(zoneId, dnsId, name, content string) (resp UpdateRecordResponse, err error) {
	endpoint := fmt.Sprintf("zones/%s/dns_records/%s", zoneId, dnsId)
	reqBody := []byte(fmt.Sprintf(`{"content":"%s","name":"%s","type":"A"}`, content, name))

	res, err := c.do(http.MethodPut, endpoint, map[string]string{}, reqBody)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	if !resp.Success {
		return resp, errors.New("Cloudflare API call failed")
	}

	return resp, err
}
