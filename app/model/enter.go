package model

type OptionsRes struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

type IdListReq struct {
	IdList []uint `json:"idList"`
}
type IdReq struct {
	Id uint `json:"id" form:"id"`
}
