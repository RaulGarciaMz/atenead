package tests

import (
	"net/http"
	"testing"

	"github.com/RaulGarciaMz/atenead/colector"
	"github.com/RaulGarciaMz/atenead/modelos"
	"github.com/RaulGarciaMz/go-httpclient/gohttp"
)

func TestColector_ColectaInformacion(t *testing.T) {

	f := colector.NewColector(&http.Client{})

	eq := modelos.Equipo{
		Id:            36,
		Nombre:        "Dispositivo I",
		Ip:            "195.159.183.43",
		Descripcion:   "Equipo con autenticación 1",
		Puerto:        "5959",
		Usuario:       "admin",
		Password:      "lyngsat",
		Autenticacion: true,
	}

	data, err := f.ColectaInformacion(eq, "http://195.159.183.43:5959/automation/service/v1")
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
