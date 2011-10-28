package main

import (
	"fmt"
	"os"
	"strconv"
)

type VueRoute struct {
	Point
	IdRoute   uint
	TypeRoute string
}

func (o *VueRoute) readCsv(cells []string) (err os.Error) {
	if len(cells) < 6 {
		err = os.NewError(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	o.readCsvPoint(cells)
	o.IdRoute, _ = strconv.Atoui(cells[4])
	o.TypeRoute = cells[5]
	return
}

func (o *VueRoute) store(mm *MemMap) {
	mm.StoreRoute(o)
}
