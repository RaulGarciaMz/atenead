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

	//Falta configuraciónde IP de Bdd de atenea y puerto

	colInfo := colector.NewColector(&http.Client{})
	clHttp := gohttp.NewBuilder().SetHttpClient(&http.Client{}).DisableTimeouts(true).Build()

	/* 	ticker := time.NewTicker(60 * time.Second)
	   	for {
	   		select {
	   		case <-ticker.C: */

	//Obtiene la lista de los equipos
	urlEq := "http://184.72.110.87:8085/atenea/admin/equipo"

	responseEq, errcli := clHttp.Get(urlEq)
	if errcli != nil {
		return
	}

	if responseEq.StatusCode != 200 {
		return
	}

	equipos := []modelos.Equipo{}

	err := responseEq.UnmarshallJson(&equipos)
	if err != nil {
		return
	}

	// Colecta la información de alarmas de cada equipo
	for _, eq := range equipos {

		// Aquí debe haber goroutines

		sb := strings.Builder{}
		sb.WriteString("http://")
		sb.WriteString(eq.Ip)
		sb.WriteString(":5959/automation/service/v1")
		datos, err := colInfo.ColectaInformacion(eq.Id, eq.Nombre, sb.String())
		//datos, err := colInfo.ColectaInformacion(eq.Id, eq.Nombre, "http://195.159.183.43:5959/automation/service/v1")
		//

		// si no fue posible obtener datos del equipo, debe mandar mensaje y pasar al siguiente equipo
		if err != nil {
			continue
		}

		/* 			if datos.HayFiltroAlarmas {
			//Marcar al equipos con bandera que indique que tiene filtros configurados
		} */

		//Generar estructura para procesar alarmas y enviarlas a la API correspondiente
		numAlarmas := len(datos.Alarmas)

		msgIds := make([]int64, numAlarmas)
		msgSlots := make([]int32, numAlarmas)
		msgPorts := make([]*int64, numAlarmas)
		msgTexts := make([]string, numAlarmas)
		msgSourcesNames := make([]string, numAlarmas)
		msgSeveryties := make([]string, numAlarmas)
		msgInstances := make([]int32, numAlarmas)
		msgSetTimes := make([]string, numAlarmas)
		msgCardIds := make([]int32, numAlarmas)
		timestampInicios := make([]string, numAlarmas)

		count := 0
		for _, v := range datos.Alarmas {
			msgIds[count] = v.MsgId
			msgSlots[count] = v.MsgSlot
			msgPorts[count] = v.MsgPort
			msgTexts[count] = v.MsgText
			msgSourcesNames[count] = v.MsgSourceName
			msgSeveryties[count] = v.MsgSeverity
			msgInstances[count] = v.MsgInstance
			msgSetTimes[count] = v.MsgSetTime.Format("2006-01-02 15:04:05")
			msgCardIds[count] = v.MsgCardId
			timestampInicios[count] = v.SetTimeAdjusted.Format("2006-01-02 15:04:05")
			count++
		}

		pa := modelos.ListaAlarmasParam{
			Id:               eq.Id,
			MsgIds:           msgIds,
			MsgSlots:         msgSlots,
			MsgPorts:         msgPorts,
			MsgTexts:         msgTexts,
			MsgSourcesNames:  msgSourcesNames,
			MsgSeveryties:    msgSeveryties,
			MsgInstances:     msgInstances,
			MsgSetTimes:      msgSetTimes,
			MsgCardIds:       msgCardIds,
			TimestampInicios: timestampInicios,
			DateServer:       datos.EquipoTimestamp.Format("2006-01-02"),
			TimeServer:       datos.EquipoTimestamp.Format("15:04:05"),
		}

		//Envía al procesamiento de alarmas
		urlEq := "http://184.72.110.87:8085/atenea/admin/alarma/procesa"
		responseEq, errcli := clHttp.Post(urlEq, pa)
		if errcli != nil {
			return
		}

		if responseEq.StatusCode != 200 {
			return
		}

	}

	/*
		 		}
			}
	*/
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
