package modelos

import "encoding/xml"

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
}

type GetTimeDateResponse struct {
	XMLName  xml.Name `xml:"Envelope"`
	SoapBody *SOAPBodyGetTimeDateResponse
}

type SOAPBodyGetTimeDateResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *GetTimeResponseBody
	FaultDetails *Fault
}

type GetTimeResponseBody struct {
	XMLName     xml.Name `xml:"getTimeDateResponse"`
	Code        string   `xml:"code"`
	Message     string   `xml:"message"`
	TimeAndDate *TimeDateBody
	Status      string `xml:"Status"`
}

type TimeDateBody struct {
	XMLName           xml.Name `xml:"timeAndDate"`
	UseGpsTime        string   `xml:"useGpsTime"`
	LocalTimezonePath string   `xml:"localTimezonePath"`
	Date              string   `xml:"date"`
	Time              string   `xml:"time"`
}

type GetAlarmFilterResponse struct {
	XMLName  xml.Name `xml:"Envelope"`
	SoapBody *SOAPBodyGetAlarmFilterResponse
}

type SOAPBodyGetAlarmFilterResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *GetAlarmFilterResponseBody
	FaultDetails *Fault
}

type GetAlarmFilterResponseBody struct {
	XMLName         xml.Name               `xml:"getAlarmFilterResponse"`
	Code            string                 `xml:"code"`
	Message         string                 `xml:"message"`
	AlarmFilterList []AlarmFilterListEntry `xml:"alarmFilterList>alarmFilterListEntry"`
}

type AlarmFilterListEntry struct {
	MsgId       string `xml:"msgId"`
	MsgCardType string `xml:"msgCardType"`
	MsgSlot     string `xml:"msgSlot"`
	MsgModule   string `xml:"msgModule"`
	MsgPort     string `xml:"msgPort"`
	MsgSeverity string `xml:"msgSeverity"`
}

type GetAlarmListResponse struct {
	XMLName  xml.Name `xml:"Envelope"`
	SoapBody *SOAPBodyGetAlarmListResponse
}

type SOAPBodyGetAlarmListResponse struct {
	XMLName      xml.Name `xml:"Body"`
	Resp         *GetAlarmListResponseBody
	FaultDetails *Fault
}

type GetAlarmListResponseBody struct {
	XMLName   xml.Name          `xml:"getAlarmListResponse"`
	Code      string            `xml:"code"`
	Message   string            `xml:"message"`
	AlarmList []*AlarmListEntry `xml:"alarmList>alarmListEntry"`
}

type AlarmListEntry struct {
	MsgId         string  `xml:"msgId"`
	MsgSlot       string  `xml:"msgSlot"`
	MsgPort       *string `xml:"msgPort"`
	MsgSeverity   string  `xml:"msgSeverity"`
	MsgText       string  `xml:"msgText"`
	MsgSourceName string  `xml:"msgSourceName"`
	MsgInstance   string  `xml:"msgInstance"`
	MsgSetTime    string  `xml:"msgSetTime"`
	MsgCardId     string  `xml:"msgCardId"`
}
