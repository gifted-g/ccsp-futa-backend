package util

import (
	"errors"
	"strings"
)

func Require(fields map[string]string) error {
	for k,v := range fields { if strings.TrimSpace(v)=="" { return errors.New(k+" is required") } }
	return nil
}