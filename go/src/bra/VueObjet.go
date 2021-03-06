package bra

/*
 * correspond aux
 *  éléments
 *  charettes
 *  runes
 *  buissons
 *  munition
 */

import (
	"errors"
	"fmt"
	"strings"
)

type VueObjet struct {
	X        int16
	Y        int16
	Type     string
	Quantité uint
	Label    string
	IdType   uint
}

func (o *VueObjet) ReadCsvSimple(cells []string, Type string, Label string) (err error) {
	if len(cells) < 5 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = Type
	o.Label = Label
	return
}

func (o *VueObjet) ReadCsvSimpleLabel(cells []string, Type string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = Type
	o.Label = cells[5]
	return
}

// quantité, label, butin
func (o *VueObjet) ReadCsvQLB(cells []string, Type string) (err error) {
	if len(cells) < 7 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = Type
	o.Label = cells[5]
	if o.IdType, err = Atoui(cells[6]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvGraine(cells []string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "graine"
	o.Quantité, _ = Atoui(cells[4])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, "graine")
	if o.Quantité > 1 {
		o.Label += "s de "
	} else {
		o.Label += " de "
	}
	o.Label += cells[5]
	if o.IdType, err = Atoui(cells[6]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvTabac(cells []string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "tabac"
	o.Quantité, _ = Atoui(cells[4])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, "feuille")
	if o.Quantité > 1 {
		o.Label += "s "
	} else {
		o.Label += " "
	}
	o.Label += cells[5]
	if o.IdType, err = Atoui(cells[6]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvLingot(cells []string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "lingot"
	o.Quantité, _ = Atoui(cells[4])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, "lingot")
	if o.Quantité > 1 {
		o.Label += "s de "
	} else {
		o.Label += " de "
	}
	o.Label += cells[5]
	if o.IdType, err = Atoui(cells[6]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvMinerai(cells []string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "minerai"
	o.Quantité, _ = Atoui(cells[4])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, "minerai")
	if o.Quantité > 1 {
		o.Label += "s de "
	} else {
		o.Label += " de "
	}
	o.Label += cells[5]
	if o.IdType, err = Atoui(cells[6]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvElement(cells []string) (err error) {
	if len(cells) < 6 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = strings.ToLower(cells[4])
	o.Quantité, _ = Atoui(cells[5])

	o.Label = fmt.Sprintf("%d %s", o.Quantité, o.Type)
	// lignes suivantes provisoires
	if o.Quantité > 1 {
		if len(cells) > 6 {
			o.Label += cells[6]
		} else {
			o.Label += "s"
		}
	}
	return
}
func (o *VueObjet) ReadCsvMunition(cells []string) (err error) {
	if len(cells) < 8 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "munition"
	o.Quantité, _ = Atoui(cells[6])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, cells[5])
	o.IdType, _ = Atoui(cells[7])
	return
}

func (o *VueObjet) ReadCsvRune(cells []string) (err error) {
	if len(cells) < 3 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "rune"
	o.Label = "rune"
	return
}

func (o *VueObjet) ReadCsvPlante(cells []string, brut bool) (err error) {
	if len(cells) < 8 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "plante"
	o.Quantité, _ = Atoui(cells[4])
	o.Label = fmt.Sprintf("%d %s", o.Quantité, cells[5])
	if o.Quantité > 1 {
		o.Label += "s de "
	} else {
		o.Label += " de "
	}
	o.Label += cells[6]
	if brut {
		o.Label += " (plante brute)"
	} else {
		o.Label += " (plante préparée)"
	}
	if o.IdType, err = Atoui(cells[8]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvPotion(cells []string) (err error) {
	if len(cells) < 8 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "potion"
	o.Label = fmt.Sprintf("%s %s de qualité %s", cells[5], cells[6], cells[7])
	if o.IdType, err = Atoui(cells[4]); err != nil {
		return
	}
	return
}

func (o *VueObjet) ReadCsvAliment(cells []string) (err error) {
	if len(cells) < 8 {
		err = errors.New(fmt.Sprintf("pas assez de champs (%d)", len(cells)))
		return
	}
	if o.X, err = Atoi16(cells[1]); err != nil {
		return
	}
	if o.Y, err = Atoi16(cells[2]); err != nil {
		return
	}
	o.Type = "aliment"
	o.Label = fmt.Sprintf("%s de qualité %s", cells[5], cells[6])
	if o.IdType, err = Atoui(cells[7]); err != nil {
		return
	}
	return
}
