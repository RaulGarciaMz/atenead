package tests

import (
	"net/http"
	"testing"

	"github.com/RaulGarciaMz/atenead/colector"
	"github.com/RaulGarciaMz/go-httpclient/gohttp"
)

func TestColector_ColectaInformacion(t *testing.T) {

	f := colector.NewColector(&http.Client{})

	data, err := f.ColectaInformacion(1, "Equipo 1", "http://195.159.183.43:5959/automation/service/v1")
	if err != nil {
		t.Fatal("No pudo colectar información")
	}

	_ = data
}

func TestColector_HttpClient(t *testing.T) {

	clHttp := gohttp.NewBuilder().SetHttpClient(&http.Client{}).DisableTimeouts(true).Build()

	//Obtiene la lista de los equipos
	urlEq := "http://184.172.110.87:8085/atenea/admin/equipo"
	responseEq, errcli := clHttp.Get(urlEq)
	if errcli != nil {
		return
	}

	if responseEq.StatusCode != 200 {
		return
	}

	/*
		 	var equipos []modelos.Equipo

			err := responseEq.UnmarshallJson(equipos)
			if err != nil {
				return
			}
	*/
}
