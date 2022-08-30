package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/RaulGarciaMz/atenead/colector"
	"github.com/RaulGarciaMz/atenead/configuration"
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
	confFile, err := os.Open("atenea_dconf.json")
	if err != nil {
		log.Fatal("no pudo abrir el archivo de configuración")
		panic(err)
	}
	defer confFile.Close()
	conf, err := ioutil.ReadAll(confFile)
	if err != nil {
		log.Fatal("no pudo leer el archivo de configuración")
		panic(err)
	}

	myConf := configuration.ServerConfig{}
	err = json.Unmarshal(conf, &myConf)
	if err != nil {
		log.Fatal("no pudo interpretar el archivo de configuración")
		panic(err)
	}

	if !IsValidIp(myConf.Ip) {
		log.Fatal("error en archivo de configuración. IP inválida")
		panic(fmt.Errorf("error en archivo de configuración. IP inválida"))
	}

	if myConf.Tick < 60 {
		myConf.Tick = 60
	} else if myConf.Tick > 120 {
		myConf.Tick = 120
	}

	sb := strings.Builder{}
	sb.WriteString("http://")
	sb.WriteString(myConf.Ip)
	sb.WriteString(":")
	sb.WriteString(myConf.Port)
	preUrl := sb.String()

	colInfo := colector.NewColector(&http.Client{})
	clHttp := gohttp.NewBuilder().SetHttpClient(&http.Client{}).DisableTimeouts(true).Build()

	/* 	ticker := time.NewTicker(myConf.Tick * time.Second)
	   	for {
	   		select {
	   		case <-ticker.C:


			}
		}*/

	equipos, err := ObtieneEquiposFromAtenea(preUrl, clHttp)
	if err != nil {

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
			//Escribe en el log el error de lectura en el equipo
			continue
		}

		//Cambia el valor del filtro
		fltro := modelos.FiltroEquipoParam{
			Id:     eq.Id,
			Filtro: datos.HayFiltroAlarmas,
		}
		chFiltro, err := FiltroEquipo(preUrl, clHttp, &fltro)
		if err != nil {
			//Escribe en el log el error de modificación del filtro
			continue
		}

		if chFiltro != "" {
			//crear error por no poder cambiar el filtro
			//Escribe en el log el error al ide modificación del filtro
			continue
		}

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

		res, err := ProcesaAlarmasAtenea(preUrl, clHttp, &pa)
		if err != nil {
		}

		if res == "Corrupto" {
			//Registrar error de procesamiento de alarmas
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

func ObtieneEquiposFromAtenea(preUrl string, clHttp gohttp.Client) ([]modelos.Equipo, error) {
	//urlEq := "http://184.72.110.87:8085/atenea/admin/equipo"
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/equipo")

	responseEq, errcli := clHttp.Get(sb.String())
	if errcli != nil {
		return nil, errcli
	}

	if responseEq.StatusCode != 200 {
		return nil, fmt.Errorf("")
	}

	equipos := []modelos.Equipo{}

	err := responseEq.UnmarshallJson(&equipos)
	if err != nil {
		return nil, fmt.Errorf("")
	}

	return equipos, nil
}

func ProcesaAlarmasAtenea(preUrl string, clHttp gohttp.Client, pa *modelos.ListaAlarmasParam) (string, error) {
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/alarma/procesa")
	responseEq, errcli := clHttp.Post(sb.String(), pa)
	if errcli != nil {
		return "", errcli
	}

	if responseEq.StatusCode != 200 {
		return "", fmt.Errorf("")
	}

	return responseEq.String(), nil
}

func FiltroEquipo(preUrl string, clHttp gohttp.Client, pa *modelos.FiltroEquipoParam) (string, error) {
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/equipo/filtro")
	responseEq, errcli := clHttp.Post(sb.String(), pa)
	if errcli != nil {
		return "", errcli
	}

	if responseEq.StatusCode != 200 {
		return "", fmt.Errorf("")
	}

	return responseEq.String(), nil
}

func IsValidIp(ip string) bool {
	return net.ParseIP(ip) != nil
}
