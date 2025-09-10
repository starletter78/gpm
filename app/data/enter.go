package data

import "gpm/app/data/api_data"

type Data struct {
	ApiData api_data.ApiData
}

func NewData() *Data {
	return &Data{}
}
