package soapclient

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	//"net/url"
	"strings"
	"sync"

	"github.com/RaulGarciaMz/atenead/modelos"
)

type AteneaSoapHttp struct {
	cliente *http.Client
	//url     string //http://195.159.183.43:5959/automation/service/v1"

	clientOnce sync.Once
}

func NewAteneaSoapHttp(c *http.Client) (*AteneaSoapHttp, error) {

	s := AteneaSoapHttp{}

	s.clientOnce.Do(func() {
		if c != nil {
			s.cliente = c
		} else {
			s.cliente = &http.Client{}
		}
	})

	return &s, nil
}

func (s *AteneaSoapHttp) GetDateTime(url string) (*modelos.GetTimeDateResponse, error) {

	payload := []byte(strings.TrimSpace(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://www.appeartv.com/automation/v1">
	   <soapenv:Header/>
	   <soapenv:Body>
		  <v1:getTimeDate>?</v1:getTimeDate>
	   </soapenv:Body>
	</soapenv:Envelope>`,
	))

	// prepare the request
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	soapAction := "getTimeDate"
	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	//client := &http.Client{}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		return nil, errors.New("error no controlado. código 500 - Internal server error o superior")
	}

	if res.StatusCode == http.StatusOK {
		result := new(modelos.GetTimeDateResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("error no considerado para controlar")
}

func (s *AteneaSoapHttp) GetAlarmFilter(url string) (*modelos.GetAlarmFilterResponse, error) {

	payload := []byte(strings.TrimSpace(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://www.appeartv.com/automation/v1">
	   <soapenv:Header/>
	   <soapenv:Body>
		  <v1:getAlarmFilter>?</v1:getAlarmFilter>
	   </soapenv:Body>
	</soapenv:Envelope>`,
	))

	// prepare the request
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	soapAction := "getAlarmFilter"
	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	//client := &http.Client{}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		return nil, errors.New("error no controlado. código 500 - Internal server error o superior")
	}

	if res.StatusCode == http.StatusOK {
		result := new(modelos.GetAlarmFilterResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("error no considerado para controlar")
}

func (s *AteneaSoapHttp) GetAlarmList(url string) (*modelos.GetAlarmListResponse, error) {

	// payload
	payload := []byte(strings.TrimSpace(`
	<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:v1="http://www.appeartv.com/automation/v1">
	   <soapenv:Header/>
	   <soapenv:Body>
		  <v1:getAlarmList>
		  </v1:getAlarmList>
	   </soapenv:Body>
	</soapenv:Envelope>`,
	))

	// prepare the request
	req, err := http.NewRequest("POST", url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}

	soapAction := "getAlarmList"
	// set the content type header, as well as the oter required headers
	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	//client := &http.Client{}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode > http.StatusOK {
		return nil, errors.New("error no controlado. código 500 - Internal server error o superior")
	}

	if res.StatusCode == http.StatusOK {
		result := new(modelos.GetAlarmListResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil
	}

	return nil, errors.New("error no considerado para controlar")
}
