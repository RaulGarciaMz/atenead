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

func NewColector(client *http.Client) *Colector {

	s, err := soapclient.NewAteneaSoapHttp(client)
	if err != nil {
	}

	return &Colector{
		sc: s,
	}
}

func (c Colector) ColectaInformacion(idEquipo int32, nombreEquipo string, myurl string) (*modelos.ColectadoServidor, error) {

	_, err := url.ParseRequestURI(myurl)
	if err != nil {
		return nil, err
	}

	var resp *modelos.ColectadoServidor
	hayFiltro := false
	ahorita := time.Now()
	gtd, err := c.sc.GetDateTime(myurl)
	if err != nil {
		return resp, err
	}

	//calcula diferencia entre hora local y hora obtenida
	timeServer, err := time.Parse("2006-01-02 15:04:05", fmt.Sprintf("%s %s", gtd.SoapBody.Resp.TimeAndDate.Date, gtd.SoapBody.Resp.TimeAndDate.Time))
	if err != nil {
		return resp, err
	}

	diff := ahorita.Sub(timeServer)

	filter, err := c.sc.GetAlarmFilter(myurl)
	if err != nil {
		return resp, err
	}
	if len(filter.SoapBody.Resp.AlarmFilterList) > 0 {
		hayFiltro = true
	}

	alarmasEnEquipo, err := c.sc.GetAlarmList(myurl)
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
		var MsgPortVal *int
		if alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgPort != nil {
			mp = *alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgPort
			intMsgPortVal, err := strconv.Atoi(mp)
			if err != nil {
				return resp, err
			}
			MsgPortVal = &intMsgPortVal
		}

		clave := fmt.Sprintf("%d%s%s%s%s%s", idEquipo, nombreEquipo, mid, ms, mins, mp)

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

		alEnEquipo[clave] = modelos.ListaAlarma{
			IdEquipo:      idEquipo,
			Equipo:        nombreEquipo,
			MsgId:         int32(MsgIdVal),
			MsgSlot:       MsgSlotVal,
			MsgPort:       MsgPortVal,
			MsgText:       alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgText,
			MsgSourceName: alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSourceName,
			MsgSeverity:   alarmasEnEquipo.SoapBody.Resp.AlarmList[i].MsgSeverity,
			MsgInstance:   MsgInstanceVal,
			MsgSetTime:    MsgSetTimeVal,
			MsgCardId:     int32(MsgCardIdVal),
		}
	}

	resp = &modelos.ColectadoServidor{
		HayFiltroAlarmas: hayFiltro,
		DiferenciaTiempo: diff,
		Alarmas:          alEnEquipo,
	}

	return resp, nil
}
