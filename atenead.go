package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RaulGarciaMz/atenead/colector"
	"github.com/RaulGarciaMz/atenead/modelos"
	"github.com/RaulGarciaMz/go-httpclient/gohttp"
)

var (
	Sha1ver        string // clave sha1 del commit usado para compilar el programa
	Branch         string //nombre de la rama usada para compilar el programa
	BuildTime      string // when the executable was built
	LastCommitTime string // when the last commit was
	Tag            string // nombre del último tag registrado en la rama
	flgVersion     bool
)

type Equipo struct {
	Id          *int32  `json:"id"`
	Nombre      *string `json:"nombre"`
	Ip          *string `json:"ip"`
	Descripcion *string `json:"descripcion"`
}

func main() {

	version := generaVersion()
	parseCmdLineFlags(version)

	colInfo := colector.NewColector(&http.Client{})
	clHttp := gohttp.NewBuilder().SetHttpClient(&http.Client{}).Build()

	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:

			//Obtiene la lista de los equipos
			urlEq := "https://184.172.110.87:8085/atenea/admin/equipo"
			responseEq, errcli := clHttp.Get(urlEq)
			if errcli != nil {
				return
			}

			if responseEq.StatusCode != 200 {
				return
			}

			var equipos []modelos.Equipo

			err := responseEq.UnmarshallJson(equipos)
			if err != nil {
				return
			}

			// Colecta la información de alarmas de cada equipo
			for _, eq := range equipos {

				sb := strings.Builder{}
				sb.WriteString("http://")
				sb.WriteString(eq.Ip)
				sb.WriteString(":5959/automation/service/v1")
				datos, err := colInfo.ColectaInformacion(eq.Id, eq.Nombre, sb.String())
				if err != nil {
				}

				if datos.HayFiltroAlarmas {
					//Marcar al equipos con bandera que indique que tiene filtros configurados
				}

			}

		}
	}
}

func parseCmdLineFlags(version string) {
	flag.BoolVar(&flgVersion, "version", false, "si true, imprime la versión y termina el programa")
	flag.Parse()

	if flgVersion {
		fmt.Println(version)
		os.Exit(0)
	}
}

func generaVersion() string {
	var version string

	tUnix, err := strconv.ParseInt(LastCommitTime, 10, 64)
	if err != nil {
		version = "Compilado el: " + BuildTime + " Rama: " + Branch + " Commit (sha1): " + Sha1ver + " Tag: " + Tag
	}
	timeT := time.Unix(tUnix, 0)
	version = "Compilado el: " + BuildTime +
		" Rama: " + Branch +
		" Commit (sha1): " + Sha1ver +
		" Fecha Commit: " + timeT.Format("02-01-2006 15:04:00") +
		" Tag: " + Tag

	return version
}
