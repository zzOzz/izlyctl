package crous

import (
	"github.com/spf13/viper"
	"fmt"
	"encoding/base64"
	"net/http"
	"strings"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"time"
)
//[
//{
//"idTransmitter":65
//"idMapping":1
//"idZdc":36882379
//"zdcCreationDate":"2015-09-23T02:00:00.000+02"
//"pixSs":"03"
//"pixNn":"0000"
//"appl":"D1"
//"uid":"0471928A083C80"
//"rid":"0471928A083C80"
//"reason":null
//"cancelDate":null
//"revalidationDate":null
//"deliveryDate":null
//"origin":null
//}
//]

type Card struct {
	RightHolder RightHolder `json:"rightholder,omitempty"`
	SmartCards []SmartCard `json:"smartcards,omitempty"`
}

type SmartCard struct {
	IdTransmitter int `json:"idTransmitter,omitempty"`
	IdMapping int `json:"idMapping,omitempty"`
	IdZdc int `json:"idZdc,omitempty"`
	ZdcCreationDate CustomTime `json:"zdcCreationDate,omitempty"`
	PixSs string `json:"pixSs,omitempty"`
	PixNn string `json:"pixNn,omitempty"`
	Appl string `json:"appl,omitempty"`
	UID string `json:"uid,omitempty"`
	RID string `json:"rid,omitempty"`
	Reason string `json:"reason,omitempty"`
	CancelDate CustomTime `json:"cancelDate,omitempty"`
	RevalidateDate CustomTime `json:"revalidateDate,omitempty"`
	DeliveryDate CustomTime `json:"deliveryDate,omitempty"`
	Origin string `json:"origin,omitempty"`
}

type RightHolder struct {
	Identifier string `json:"identifier,omitempty"`
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	Email string `json:"email,omitempty"`
	DueDate CustomTime `json:"dueDate,omitempty"`
	IdCompanyRate int `json:"idCompanyRate,omitempty"`
	IdRate int `json:"idRate,omitempty"`
	BirthDate CustomTime `json:"birthDate,omitempty"`
	RneOrgCode string `json:"rneOrgCode,omitempty"`
	RneDepCode string `json:"rneDepCode,omitempty"`
	InternalId string `json:"internalId,omitempty"`
	Ine string `json:"ine,omitempty"`
	CellNumber string `json:"cellNumber,omitempty"`
	Address1 string `json:"address1,omitempty"`
	ZipCode string `json:"zipCode,omitempty"`
	City string `json:"city,omitempty"`
	Country string `json:"country,omitempty"`
	Other1 string `json:"other1,omitempty"`
	Other2 string `json:"other2,omitempty"`
	Other3 string `json:"other3,omitempty"`
	Other4 string `json:"other4,omitempty"`
	Other5 string `json:"other5,omitempty"`
	Student bool `json:"student,omitempty"`
	CreatedAt int `json:"createdAt,omitempty"`
	UpdatedAt int `json:"updatedAt,omitempty"`
	IdCrous int `json:"idCrous,omitempty"`
	IdCompanyRateForced int `json:"idCompanyRateForced,omitempty"`
	IdRateForced int `json:"idRateForced,omitempty"`
	ApplyRateForced bool `json:"applyRateForced,omitempty"`
	Pic int `json:"pic,omitempty"`
}

type Auth struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
}

type Api struct {
	Auth Auth `json:"auth,omitempty"`
	Url string `json:"url,omitempty"`
}

type Config struct {
	Api Api
}

func GetAuth() Auth {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	return config.Api.Auth
}

func New() *Client {
	var config Config

	err := viper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("unable to decode into config struct, %v", err)
	}
	client := Client{
		config,
		&http.Client{},
		"",
	}
	return &client
}

// Client is the API client
type Client struct {
	config Config
	client *http.Client
	token string
}

func (c *Client) GetResource(resource string, authenticated bool) ([]byte, http.Header, error) {
	var url string
	url = fmt.Sprintf("%s%s", strings.TrimRight(c.config.Api.Url, "/"), resource)
	logrus.Debugf(">>> GET  %q", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("User-Agent", "izlyctl")
	if authenticated {
		req.Header.Set("Authorization", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	header:= resp.Header
	logrus.Debugf("<<< %s", body)

	if resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Status code: %d ", resp.StatusCode)
	}

	return body, header, err
}

func (c *Client) PostResource(resource string, data interface{}, authenticated bool) ([]byte, http.Header, error) {
	var url string
	url = fmt.Sprintf("%s%s", strings.TrimRight(c.config.Api.Url, "/"), resource)
	logrus.Debugf(">>> POST  %q", url)

	payload := new(bytes.Buffer)
	encoder := json.NewEncoder(payload)
	if err := encoder.Encode(data); err != nil {
		return nil, nil, err
	}

	payloadString := strings.TrimSpace(fmt.Sprintf("%s", payload))
	logrus.Debugf(">>> POST %s payload=%s", url, payloadString)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("User-Agent", "izlyctl")
	if authenticated {
		req.Header.Set("Authorization", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	header:= resp.Header
	logrus.Debugf("headers <<< %s", header)
	logrus.Debugf("body <<< %s", body)
	if resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Status code: %d ", resp.StatusCode)
	}

	return body, header, err
}

func (c *Client) PostAuth(auth Auth) ([]byte, http.Header, error) {
	var url string
	var creds string
	var form string

	url = fmt.Sprintf("%s/v1/token", c.config.Api.Url)
	logrus.Debugf(">>> POST  %q", url)

	creds = fmt.Sprintf("%s:%s", auth.Login, auth.Password)
	creds = base64.StdEncoding.EncodeToString([]byte(creds))
	creds = fmt.Sprintf("Basic %s", creds)
	logrus.Debugf(">>> Basic  %q", creds)

	form = "grant_type=client_credentials&env=PRD"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(form)))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Authorization", creds)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "izlyctl")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	header:= resp.Header
	logrus.Debugf("headers <<< %s", header)
	logrus.Debugf("body <<< %s", body)
	if resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Status code: %d ", resp.StatusCode)
	}

	return body, header, err
}

func (c *Client) PutResource(resource string, data interface{}, authenticated bool) ([]byte, http.Header, error) {
	var url string
	url = fmt.Sprintf("%s%s", strings.TrimRight(c.config.Api.Url, "/"), resource)
	logrus.Debugf(">>> PUT  %q", url)

	//payload := new(bytes.Buffer)
	//encoder := json.NewEncoder(payload)
	//if err := encoder.Encode(data); err != nil {
	//	return nil, nil, err
	//}
	payload, err := json.Marshal(data)
	if err != nil {
		return nil, nil, err
		fmt.Println("error:", err)
	}

	//payloadString := strings.TrimSpace(fmt.Sprintf("%s", payload))
	logrus.Debugf(">>> PUT %s payload=%s", url, string(payload))

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(payload)))
	//req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("User-Agent", "izlyctl")
	if authenticated {
		req.Header.Set("Authorization", c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	header:= resp.Header
	logrus.Debugf("headers <<< %s", header)
	logrus.Debugf("body <<< %s", body)
	if resp.StatusCode > 299 {
		return nil, nil, fmt.Errorf("Status code: %d ", resp.StatusCode)
	}

	return body, header, err
}

func (c *Client) Authenticate() {
	body, header, err := c.PostAuth(c.config.Api.Auth)
	c.token = header.Get("Authorization")
	logrus.Debugf("auth <<< %s", c.token)
	if err != nil {
		logrus.Errorf("err <<< %s", err)
		logrus.Errorf("body <<< %s", body)
		logrus.Errorf("header <<< %s", header)
	}
}

func (c *Client) GetRightHolder(identifier string) (RightHolder){
	body, header, err := c.GetResource("/v1/rightholders/" + identifier, true)
	if err != nil {
		logrus.Errorf("err <<< %s", err)
		logrus.Errorf("body <<< %s", body)
		logrus.Errorf("header <<< %s", header)
	}
	// fmt.Printf("%s", body)
	var response RightHolder
	err = json.Unmarshal(body, &response)
	return response
}

func (c *Client) ChangeRightHolderRate(identifier string, idRate int, dueDate string) (RightHolder){
	const inputCtLayout = "2006-01-02"
	var rightHolder RightHolder = c.GetRightHolder(identifier)
	rightHolder.IdRate = idRate
	var timeError error
	rightHolder.DueDate.Time, timeError = time.Parse(inputCtLayout, dueDate)
	if timeError != nil {
		logrus.Errorf("timeError <<< %s", timeError)
	}
	//rightHolder.BirthDate = nil
	//rightHolder.DueDate = nil
	body, header, err := c.PutResource("/v1/rightholders/" + rightHolder.Identifier, rightHolder, true)
	if err != nil {
		logrus.Errorf("err <<< %s", err)
		logrus.Errorf("body <<< %s", body)
		logrus.Errorf("header <<< %s", header)
	}
	var response RightHolder
	err = json.Unmarshal(body, &response)
	return response
}

func (c *Client) GetSmartCard(identifier string) ([]SmartCard){
	body, header, err := c.GetResource("/v1/rightholders/" + identifier + "/smartcard", true)
	if err != nil {
		logrus.Errorf("err <<< %s", err)
		logrus.Errorf("body <<< %s", body)
		logrus.Errorf("header <<< %s", header)
	}
	// fmt.Printf("%s", body)
	var response []SmartCard
	err = json.Unmarshal(body, &response)
	return response
}

func (c *Client) GetCard(identifier string) (Card){
	return	Card{
		RightHolder: c.GetRightHolder(identifier),
		SmartCards: c.GetSmartCard(identifier),
	}
}
