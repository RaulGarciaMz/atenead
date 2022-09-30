package soapclient

import (
	"bytes"
	"encoding/xml"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/RaulGarciaMz/atenead/modelos"
)

// AteneaSoapHttp define el cliente para servicios SOAP
type AteneaSoapHttp struct {
	cliente    *http.Client
	clientOnce sync.Once
}

// NewAteneaSoapHttp crea un cliente para comunicación con servicios SOAP
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

// GetDateTime consulta el servicio SOAP para obtener la fecha del equipo
func (s *AteneaSoapHttp) GetDateTime(url, user, password string, autenticacion bool) (*modelos.GetTimeDateResponse, error) {

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

	if autenticacion {
		req.SetBasicAuth(user, password)
	}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		result := new(modelos.GetTimeDateResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil

	case http.StatusUnauthorized:
		return nil, errors.New("error en la autenticación credenciales no válidas")

	case http.StatusInternalServerError:
		return nil, errors.New("error 500 interno en el servidor")

	default:
		return nil, errors.New("error no considerado para controlar")
	}
}

// GetAlarmFilter consulta el servicio SOAP para obtener las alarmas con filtro del equipo
func (s *AteneaSoapHttp) GetAlarmFilter(url, user, password string, autenticacion bool) (*modelos.GetAlarmFilterResponse, error) {

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

	if autenticacion {
		req.SetBasicAuth(user, password)
	}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		result := new(modelos.GetAlarmFilterResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil

	case http.StatusUnauthorized:
		return nil, errors.New("error en la autenticación credenciales no válidas")

	case http.StatusInternalServerError:
		return nil, errors.New("error 500 interno en el servidor")

	default:
		return nil, errors.New("error no considerado para controlar")
	}
}

// GetAlarmList consulta el servicio SOAP para obtener las alasrmas del equipo
func (s *AteneaSoapHttp) GetAlarmList(url, user, password string, autenticacion bool) (*modelos.GetAlarmListResponse, error) {

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

	if autenticacion {
		req.SetBasicAuth(user, password)
	}

	// dispatch the request
	res, err := s.cliente.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusOK:
		result := new(modelos.GetAlarmListResponse)
		err = xml.NewDecoder(res.Body).Decode(result)
		if err != nil {

			return nil, err
		}
		return result, nil

	case http.StatusUnauthorized:
		return nil, errors.New("error en la autenticación credenciales no válidas")

	case http.StatusInternalServerError:
		return nil, errors.New("error 500 interno en el servidor")

	default:
		return nil, errors.New("error no considerado para controlar")
	}
}
