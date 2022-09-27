package modelos

import "time"

type Equipo struct {
	Id          int32  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Ip          string `json:"ip"`
	Puerto      string `json:"puerto"`
}

type EquipoAlcanzable struct {
	Id         int32     `json:"id"`
	Alcanzable string    `json:"alcanzable"`
	Fecha      time.Time `json:"fecha"`
}
