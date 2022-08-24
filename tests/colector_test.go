package tests

import (
	"net/http"
	"testing"

	"github.com/RaulGarciaMz/atenead/colector"
)

func TestColector_ColectaInformacion(t *testing.T) {

	f := colector.NewColector(&http.Client{})

	data, err := f.ColectaInformacion(1, "Equipo 1", "http://195.159.183.43:5959/automation/service/v1")
	if err != nil {
		t.Fatal("No pudo colectar información")
	}

	_ = data
}
