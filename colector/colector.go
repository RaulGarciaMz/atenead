package colector

import (
	"fmt"
	"github.com/RaulGarciaMz/atenead/modelos"
	"github.com/RaulGarciaMz/atenead/soapclient"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Colector struct {
	sc *soapclient.AteneaSoapHttp
}

// NewColector crea una instancia del colector de datos para servicios SOAP
func NewColector(client *http.Client) *Colector {

	s, err := soapclient.NewAteneaSoapHttp(client)
	if err != nil {
	}

	return &Colector{
		sc: s,
	}
}

// ColectaInformacion Realiza el llamado al servicio SOAP del equipo indicado y procesa el resultado para obtener las alarmas con timestamp ajustado
func (c Colector) ColectaInformacion(equipo modelos.Equipo, myurl string) (*modelos.ColectadoServidor, error) {

	_, err := url.ParseRequestURI(myurl)
	if err != nil {
		return nil, err
	}

	var resp *modelos.ColectadoServidor
	hayFiltro := false

	gtd, err := c.sc.GetDateTime(myurl, equipo.Usuario, equipo.Password, equipo.Autenticacion)
	if err != nil {
		return resp, err
	}

	//calcula diferencia entre hora local y hora obtenida
	timeServer, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", gtd.SoapBody.Resp.TimeAndDate.Date, gtd.SoapBody.Resp.TimeAndDate.Time))
	if err != nil {
		return resp, err
	}

	filter, err := c.sc.GetAlarmFilter(myurl, equipo.Usuario, equipo.Password, equipo.Autenticacion)
	if err != nil {
		return resp, err
	}
	if len(filter.SoapBody.Resp.AlarmFilterList) > 0 {
		hayFiltro = true
	}

	alarmasEnEquipo, err := c.sc.GetAlarmList(myurl, equipo.Usuario, equipo.Password, equipo.Autenticacion)
	if err != nil {
		return resp, err
	}

	totAlarmas := len(alarmasEnEquipo.SoapBody.Resp.AlarmList)
	alEnEquipo := make(map[string]modelos.ListaAlarma, totAlarmas)

	for i := range alarmasEnEquipo.SoapBody.Resp.AlarmList {

		mid := alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgId
		ms := alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSlot
		mins := alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgInstance

		var mp string
		var MsgPortVal *int64
		if alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgPort != nil {
			mp = *alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgPort
			intMsgPortVal, err := strconv.ParseInt(mp, 10, 64)
			if err != nil {
				return resp, err
			}
			MsgPortVal = &intMsgPortVal
		}

		clave := fmt.Sprintf("%d%s%s%s%s%s", equipo.Id, equipo.Nombre, mid, ms, mins, mp)

		MsgIdVal, err := strconv.Atoi(alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgId)
		if err != nil {
			return resp, err
		}
		MsgSlotVal, err := strconv.Atoi(alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSlot)
		if err != nil {
			return resp, err
		}

		MsgInstanceVal, err := strconv.Atoi(alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgInstance)
		if err != nil {
			return resp, err
		}

		MsgSetTimeVal, err := time.Parse("2006-01-02 15:04:05", alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSetTime)
		if err != nil {
			return resp, err
		}
		MsgCardIdVal, err := strconv.Atoi(alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgCardId)
		if err != nil {
			return resp, err
		}

		diff := MsgSetTimeVal.Sub(timeServer)

		alEnEquipo[clave] = modelos.ListaAlarma{
			IdEquipo:        equipo.Id,
			Equipo:          equipo.Nombre,
			MsgId:           int64(MsgIdVal),
			MsgSlot:         int32(MsgSlotVal),
			MsgPort:         MsgPortVal,
			MsgText:         alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgText,
			MsgSourceName:   alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSourceName,
			MsgSeverity:     alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSeverity,
			MsgInstance:     int32(MsgInstanceVal),
			MsgSetTime:      MsgSetTimeVal,
			MsgCardId:       int32(MsgCardIdVal),
			SetTimeAdjusted: MsgSetTimeVal.Add(diff),
		}
	}

	resp = &modelos.ColectadoServidor{
		HayFiltroAlarmas: hayFiltro,
		Alarmas:          alEnEquipo,
		EquipoTimestamp:  timeServer,
	}

	return resp, nil
}
