package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
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

func main() {

	version := generaVersion()
	parseCmdLineFlags(version)

	myConf, err := leeConfiguracion()
	if err != nil {
		panic(err)
	}

	preUrl, err := generaUrlRestService(myConf)
	if err != nil {
		panic(err)
	}

	soapClient := colector.NewColector(&http.Client{})
	clHttp := gohttp.NewBuilder().SetHttpClient(&http.Client{}).DisableTimeouts(true).Build()
	ticker := time.NewTicker(time.Duration(myConf.Tick) * time.Minute)
	defer ticker.Stop()

	for ; true; <-ticker.C {
		equipos, err := obtieneEquiposFromRestService(preUrl, clHttp)
		if err != nil {
			//log.Fatal("no fue posible obtener la lista de equipos")
			continue
		}

		for _, eq := range equipos {
			// Aquí puede haber goroutines
			urlSoap := generaUrlSoapService(eq.Ip, "5959")
			datos, err := soapClient.ColectaInformacion(eq.Id, eq.Nombre, urlSoap)
			if err != nil {
				// si no fue posible obtener datos del equipo, debe mandar mensaje y pasar al siguiente equipo
				//Escribe en el log el error de lectura en el equipo
				//log.Fatal("no fue posible obtener datos del equipo tal en la IP tal")
				//fmt.Println("No se obtuvieron datos")
				continue
			}

			fltro := modelos.FiltroEquipoParam{
				Id:     eq.Id,
				Filtro: datos.HayFiltroAlarmas,
			}
			chFiltro, err := filtroEquipoToRestService(preUrl, clHttp, &fltro)
			if err != nil {
				//Escribe en el log el error de modificación del filtro
				//log.Fatal("error al intentar modificar el valor del campo filtro del equipo")
				//fmt.Println("Error en FiltroEquipo")
				continue
			}

			if chFiltro == "Corrupto" {
				//crear error por no poder cambiar el filtro
				//Escribe en el log el error al ide modificación del filtro
				//log.Fatal("el intento de modificar el valor del campo filtro del equipo no fue exitoso")
				//fmt.Println("Filtro no cambiado - corrupto")
				continue
			}

			pa := generaAlarmasDelEquipo(eq.Id, datos)
			res, err := procesaAlarmasToRestService(preUrl, clHttp, pa)
			if err != nil {
				//fmt.Println("Error en ProcesaAlarmasAtenea")
				//log.Fatal("error al intentar procesar las alarmas del sistema")
				continue
			}
			if res == "Corrupto" {
				//fmt.Println("Error en ProcesaAlarmasAtenea - Corrupto")
				//Registrar error de procesamiento de alarmas
				//log.Fatal("el intento de procesar las alarmas no fue exitoso")
			}
		}
	}

}

// parseCmdLineFlags imprime la información de la versión del programa
func parseCmdLineFlags(version string) {
	flag.BoolVar(&flgVersion, "version", false, "si true, imprime la versión y termina el programa")
	flag.Parse()

	if flgVersion {
		fmt.Println(version)
		os.Exit(0)
	}
}

// generaVersion genera información de la versión del programa con datos de GIT
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

// leeConfiguracion obtiene los datos de configuraciṕon para el sistema desde un archivo JSON
func leeConfiguracion() (*configuration.ServerConfig, error) {

	confFile, err := os.Open("atenea_dconf.json")
	if err != nil {
		return nil, err
	}
	defer confFile.Close()
	conf, err := ioutil.ReadAll(confFile)
	if err != nil {
		return nil, err
	}

	myConf := configuration.ServerConfig{}
	err = json.Unmarshal(conf, &myConf)
	if err != nil {
		return nil, err
	}

	if myConf.Tick < 1 {
		myConf.Tick = 1
	} else if myConf.Tick > 2 {
		myConf.Tick = 2
	}

	return &myConf, nil
}

// generaUrlRestService genera la sección inicial del URL del servicio REST
func generaUrlRestService(myConf *configuration.ServerConfig) (string, error) {
	sb := strings.Builder{}

	if !IsValidIp(myConf.Ip) {
		return "", fmt.Errorf("error en archivo de configuración. IP inválida")
	}

	sb.WriteString("http://")
	sb.WriteString(myConf.Ip)
	sb.WriteString(":")
	sb.WriteString(myConf.Port)
	return sb.String(), nil
}

// generaUrlSoapService genera el URL del servicio SOAP de un equipo
func generaUrlSoapService(ip, puerto string) string {
	sb := strings.Builder{}
	sb.WriteString("http://")
	sb.WriteString(ip)
	sb.WriteString(":")
	sb.WriteString(puerto)
	sb.WriteString("/automation/service/v1")
	return sb.String()
}

// obtieneEquiposFromRestService obtiene la lista de los equipos desde el servicio REST
func obtieneEquiposFromRestService(preUrl string, clHttp gohttp.Client) ([]modelos.Equipo, error) {
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/equipo")

	responseEq, errcli := clHttp.Get(sb.String())
	if errcli != nil {
		return nil, errcli
	}

	if responseEq.StatusCode != 200 {
		return nil, fmt.Errorf("llamado a API %s con código %d ", sb.String(), responseEq.StatusCode)
	}

	equipos := []modelos.Equipo{}

	err := responseEq.UnmarshallJson(&equipos)
	if err != nil {
		return nil, err
	}

	return equipos, nil
}

// procesaAlarmasToRestService envía las alarmas obtenidas al servicio REST
func procesaAlarmasToRestService(preUrl string, clHttp gohttp.Client, pa *modelos.ListaAlarmasParam) (string, error) {
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/alarma/procesa")

	responseEq, errcli := clHttp.Post(sb.String(), pa)
	if errcli != nil {
		return "", errcli
	}

	if responseEq.StatusCode != 200 {
		return "", fmt.Errorf("llamado a API %s con código %d ", sb.String(), responseEq.StatusCode)
	}

	return responseEq.String(), nil
}

// filtroEquipoToRestService envía al servicio REST el valor de la bandera de filtro para el equipo
func filtroEquipoToRestService(preUrl string, clHttp gohttp.Client, pa *modelos.FiltroEquipoParam) (string, error) {
	sb := strings.Builder{}
	sb.WriteString(preUrl)
	sb.WriteString("/atenea/admin/equipo/filtro")

	responseEq, errcli := clHttp.Post(sb.String(), pa)
	if errcli != nil {
		return "", errcli
	}

	if responseEq.StatusCode != 200 {
		return "", fmt.Errorf("llamado a API %s con código %d ", sb.String(), responseEq.StatusCode)
	}

	return responseEq.String(), nil
}

// IsValidIp verifica si el valor de ip es un dirección IP válida
func IsValidIp(ip string) bool {
	return net.ParseIP(ip) != nil
}

// generaAlarmasDelEquipo genera la estructura de alarmas requerida para enviarla al servicio REST
func generaAlarmasDelEquipo(idEqu int32, datos *modelos.ColectadoServidor) *modelos.ListaAlarmasParam {
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

	return &modelos.ListaAlarmasParam{
		Id:               idEqu,
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
}
