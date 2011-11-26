package bra

/*
Objet exportable en json

*/

type Vue struct {
	Z            int16 // la profondeur de la couche, 0 pour la surface
	Time         int64 // secondes depuis 1970. Une date à 0 signifie que l'objet est vide ou invalide
	Voyeur       uint  // id du braldun. Un id à 0 signifie que l'objet est vide ou invalide
	PrénomVoyeur string
	XMin         int16
	XMax         int16
	YMin         int16
	YMax         int16
	Bralduns     []*Braldun
	Cadavres     []*VueCadavre
	Monstres     []*VueMonstre
	Objets       []*VueObjet
}

func NewVue() *Vue {
	vue := new(Vue)
	vue.Bralduns = make([]*Braldun, 0, 5)
	vue.Cadavres = make([]*VueCadavre, 0, 5)
	vue.Monstres = make([]*VueMonstre, 0, 0)
	vue.Objets = make([]*VueObjet, 0, 0)
	return vue
}