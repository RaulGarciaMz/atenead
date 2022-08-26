package modelos

import "time"

type ListaAlarma struct {
	Equipo          string    `json:"equipo"`
	IdEquipo        int32     `json:"id_equipo"`
	MsgId           int64     `json:"msg_id"`
	MsgSlot         int32     `json:"msg_slot"`
	MsgPort         *int64    `json:"msg_port"`
	MsgText         string    `json:"msg_text"`
	MsgSourceName   string    `json:"msg_source_name"`
	MsgSeverity     string    `json:"msg_severity"`
	MsgInstance     int32     `json:"msg_instance"`
	MsgSetTime      time.Time `json:"msg_set_time"`
	MsgCardId       int32     `json:"msg_card_id"`
	SetTimeAdjusted time.Time `json:"set_time_adjusted"`
}

type ColectadoServidor struct {
	HayFiltroAlarmas bool
	Alarmas          map[string]ListaAlarma
	EquipoTimestamp  time.Time
}

type ListaAlarmasParam struct {
	Id               int32
	MsgIds           []int64
	MsgSlots         []int32
	MsgPorts         []*int64
	MsgTexts         []string
	MsgSourcesNames  []string
	MsgSeveryties    []string
	MsgInstances     []int32
	MsgSetTimes      []string
	MsgCardIds       []int32
	TimestampInicios []string
	DateServer       string
	TimeServer       string
}
