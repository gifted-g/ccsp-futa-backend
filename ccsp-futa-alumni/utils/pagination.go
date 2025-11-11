package utils

import "strconv"

type Page struct { Page, Limit int }

func GetPage(limitStr, pageStr string) Page {
	p := Page{Page:1, Limit:20}
	if v, err := strconv.Atoi(limitStr); err==nil && v>0 && v<=100 { p.Limit=v }
	if v, err := strconv.Atoi(pageStr); err==nil && v>0 { p.Page=v }
	return p
}

func Offset(p Page) int { return (p.Page-1)*p.Limit }