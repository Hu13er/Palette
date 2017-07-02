package smsVerification

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"gitlab.com/NagByte/Palette/common"
)

func init() {
	apiKey = common.ConfigString("SMS_API_KEY")
	if apiKey == "" {
		log.Fatalln("The variable SMS_API_KEY is not present")
	}
}

const (
	smsAPIHost = "http://api.kavenegar.com"
	smsAPIPath = "/v1/%s/%s.json?%s"

	smsSimpleMethod = "sms/send"
	smsLookupMethod = "verify/lookup"

	smsVerificationTemplate = "Palette: %token"
)

var (
	ErrUnsuccessfulRequest = errors.New("unsuccessful request")
	ErrEmptyFields         = errors.New("empty fields")
	apiKey                 string
)

func getAPIPath(method, query string) string {
	return smsAPIHost + "/" + fmt.Sprintf(smsAPIPath, apiKey, method, query)
}

func sendVerification(phoneNumber, code string) (*Response, error) {
	//message := strings.Replace(smsVerificationTemplate, "%token", code, 1)
	return lookupMethod(phoneNumber, code, "kababverify")
}

func lookupMethod(phoneNumber, code, template string) (*Response, error) {
	return sendRequest(&lookUpRequest{Receptor: phoneNumber, Token: code, Template: template})
}

func simpleMethod(phoneNumber, message string) (*Response, error) {
	return sendRequest(&simpleRequest{Receptor: phoneNumber, Message: message})
}

func sendRequest(req request) (*Response, error) {
	query, err := req.ToQuery()
	if err != nil {
		return nil, err
	}

	url := getAPIPath(req.Method(), query)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	jsonResp := &Response{}
	err = json.NewDecoder(resp.Body).Decode(jsonResp)
	if err != nil {
		return nil, err
	}

	if jsonResp.Return.Status != 200 {
		err = ErrUnsuccessfulRequest
	}

	return jsonResp, err
}

type Return struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Entry struct {
	MessageId  int64  `json:"messageid"`
	Message    string `json:"message"`
	Status     int    `json:"status"`
	StatusText string `json:"statustext"`
	Sender     string `json:"sender"`
	Receptor   string `json:"receptor"`
	Date       int    `json:"date"`
	Cost       int    `json:"cost"`
}

type Response struct {
	Return  Return  `json:"return"`
	Entries []Entry `json:"entries"`
}

type request interface {
	ToQuery() (string, error)
	Method() string
}

type lookUpRequest struct {
	Receptor string
	Token    string
	Template string
}

func (lur *lookUpRequest) ToQuery() (string, error) {
	if lur.Receptor == "" || lur.Token == "" || lur.Template == "" {
		return "", ErrEmptyFields
	}
	return fmt.Sprintf("receptor=%s&token=%s&template=%s", url.QueryEscape(lur.Receptor),
		url.QueryEscape(lur.Token), url.QueryEscape(lur.Template)), nil
}

func (lur *lookUpRequest) Method() string {
	return smsLookupMethod
}

type simpleRequest struct {
	Receptor string
	Message  string
}

func (sr *simpleRequest) ToQuery() (string, error) {
	if sr.Receptor == "" || sr.Message == "" {
		return "", ErrEmptyFields
	}
	return fmt.Sprintf("receptor=%s&message=%s", url.QueryEscape(sr.Receptor), url.QueryEscape(sr.Message)), nil
}

func (sr *simpleRequest) Method() string {
	return smsSimpleMethod
}
