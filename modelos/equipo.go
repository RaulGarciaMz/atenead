package modelos

type Equipo struct {
	Id          int32  `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
	Ip          string `json:"ip"`
	Puerto      int32  `json:"puerto"`
}
