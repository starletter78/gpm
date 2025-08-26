package casbin_util

import (
	"fmt"
	"strings"
)

type Sub struct {
	Type string
	Id   string
}

func NewSub() *Sub {
	return &Sub{}
}

func (s *Sub) encode() string {
	return fmt.Sprintf("%s:%s", s.Type, s.Id)
}

func (s *Sub) EncodeUserId(id string) string {
	s.Id = id
	s.Type = "user"
	return s.encode()
}
func (s *Sub) EncodeRoleId(id string) string {
	s.Id = id
	s.Type = "role"
	return s.encode()
}

func (s *Sub) DecodeStr(str string) *Sub {
	decodeStr := strings.Split(str, ":")
	if len(decodeStr) == 2 {
		s.Type = decodeStr[0]
		s.Id = decodeStr[1]
	}
	return s
}
