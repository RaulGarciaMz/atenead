package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RaulGarciaMz/atenead/modelos"
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

	colInfo := NewColector(&http.Client{})

	ticker := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-ticker.C:

			// Obtener la lista de equipos de la BDD, ya sea por repositorio o por llamado a API
			var equipos []Equipo
			// Obtener la lista de alarmas en BDD, ya sea por repositorio o por llamado a API
			var alarmas []modelos.ListaAlarma
			// crear map de alarmas en BDD para busquedas
			alarmasEnBdd := make(map[string]modelos.ListaAlarma, len(alarmas))
			for _, ala := range alarmas {

				var mp string
				if ala.MsgPort != nil {
					mp = strconv.Itoa(*ala.MsgPort)
				}

				clave := fmt.Sprintf("%d%s%d%d%d%d", ala.IdEquipo, ala.Equipo, ala.MsgId, ala.MsgSlot, ala.MsgInstance, &mp)

				alarmasEnBdd[clave] = ala
			}

			for _, eq := range equipos {

				sb := strings.Builder{}
				sb.WriteString("http://")
				sb.WriteString(*eq.Ip)
				sb.WriteString(":5959/automation/service/v1")
				datos, err := colInfo.ColectaInformacion(*eq.Id, *eq.Nombre, sb.String())
				if err != nil {
				}

				if datos.HayFiltroAlarmas {
					//Marcar al equipos con bandera que indique que tiene filtros configurados
				}

				// Obtener las alarmas que no están registradas
				var alNoRegistradas []modelos.ListaAlarma
				for k, a := range datos.Alarmas {

					if _, ok := alarmasEnBdd[k]; !ok {
						alNoRegistradas = append(alNoRegistradas, a)
					}
				}
				//Obtener las alarmas que no están en BDD para clarearlas
				var alNoEnBdd []modelos.ListaAlarma
				for k, a := range alarmasEnBdd {

					if _, ok := datos.Alarmas[k]; !ok {
						alNoEnBdd = append(alNoEnBdd, a)
					}
				}

			}

		}
	}

	//
	// //con getAlarmaList -------
	//
	//Obtener lista de alarmas registradas en la BDD - ALReg (¿TODAS las alarmas? ... faltaría una función que traiga todas las alarmas... o sea sin clasificar)
	//Algo así como alarmasEnBdd := Repo.GetAlarmas
	//De la lista de alarmas, filtrar aquellas que no están registradas aún

	// // las filtradas se registran en la BDD con fecha modificada con la diferencia calculada y se marca como activa
	//
	// //Discriminar las que ya están registradas en BDD
	// //O sea de una lista de alarmas enviada a una consulta, obtener las
	// //.... Cómo obtener alarmas registradas o de una lista de alarmas, regresar las que NO están registradas
	//
	// //Al contrario, De la lista de alarmas registradas en la BDD, determinar las que ya no están en la Lista obtenida
	// //egistrar hora del clareo y eliminarla de la bdd

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
