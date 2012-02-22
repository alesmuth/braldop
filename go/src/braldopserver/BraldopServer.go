package main

import (
	"bra"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime/pprof"
	"strconv"
	"time"
)

const (
	port = 8001
)

var versionActuelleExtension Version

func init() {
	versionActuelleExtension = MakeVersion(3, 1)
}

type MapServer struct {
	répertoireCartes *string // répertoire racine dans lequel on trouve les répertoires des utilisateurs
	bd               *bra.BaseMysql
	fv               FusionneurVue
}

func getFormValue(hr *http.Request, name string) string {
	values := hr.Form[name]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

func envoieRéponse(w http.ResponseWriter, out *MessageOut) {
	bout, err := json.Marshal(out)
	if err != nil {
		log.Println("Erreur encodage réponse : ", err)
		return
	}
	fmt.Fprint(w, "receiveFromMapServer(")
	w.Write(bout)
	fmt.Fprint(w, ")")
}

func vérifieVersion(vs string) (html string) {
	if version, err := ParseVersion(vs); err != nil {
		log.Println("Version utilisateur incomprise : " + vs)
	} else if CompareVersions(&version, &versionActuelleExtension) == -1 {
		log.Println("Version utilisateur obsolète : " + vs)
		html = "L'extension Braldop n'est pas à jour.<br>Vous devriez installer <a href=http://canop.org/braldop/index.html>la nouvelle version</a>."
	}
	return
}

func (ms *MapServer) répertoireCartesBraldun(idBraldun uint, mdpr string) string {
	return fmt.Sprintf("%s/%d-%s", *ms.répertoireCartes, idBraldun, mdpr)
}

func (ms *MapServer) ServeHTTP(w http.ResponseWriter, hr *http.Request) {
	startTime := time.Now().UnixNano()
	defer func() {
		log.Printf(" durée totale traitement requête : %d ms\n", (time.Now().UnixNano()-startTime)/1e6)
	}()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Request-Method", "GET")

	//> analyse et vérification de la requête
	hr.ParseForm()
	in := new(MessageIn)
	out := new(MessageOut)
	defer envoieRéponse(w, out)
	bin := ([]byte)(getFormValue(hr, "in"))
	err := json.Unmarshal(bin, in)
	if err != nil {
		out.Erreur = "Erreur décodage : " + err.Error()
		log.Println("Erreur décodage : ", err.Error())
		return
	}
	//fmt.Printf("Message reçu : %+v\n", in)
	out.Text = vérifieVersion(in.Version)
	if in.IdBraldun == 0 || len(in.Mdpr) != 64 {
		log.Println("IdBraldun ou Mot de passe restreint invalide")
		return
	}
	dirBase := ms.répertoireCartesBraldun(in.IdBraldun, in.Mdpr)
	log.Println("Requête Braldun ", in.IdBraldun)
	var couche *bra.Couche
	if in.Vue == nil || len(in.Vue.Couches) == 0 {
		log.Println(" Pas de données de vue")
	} else {
		couche = in.Vue.Couches[0]
	}

	//> récupération du compte braldop authentifié et du tableau des amis
	var cb *bra.CompteBraldop
	var amis []*bra.CompteBraldop
	con, err := ms.bd.DB()
	if err != nil {
		log.Println("Erreur à la connexion bd", err)
		out.Text = "Erreur Braldop : connexion BD"
		return
	}
	cb, errmess := con.AuthentifieCompte(in.IdBraldun, in.Mdpr, true)
	if errmess != "" {
		out.Text = "Une erreur s'est produite durant l'authentification sur Braldop : <i>" + errmess + "</i><br>Contactez Canopée du Haut-Rac pour plus d'informations."
		return // on a besoin du compte
	}
	log.Println(" compte authentifié")
	amis, err = con.Amis(cb.IdBraldun)
	if err != nil {
		log.Println(" erreur durant récupération amis")
	}

	log.Println(" commande :", in.Cmd)

	//> stockage des données de vue
	if couche != nil {
		hasher := sha1.New()
		hasher.Write(bin)
		sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
		dir := dirBase + "/" + time.Now().Format("2006/01/02")
		path := dir + "/carte-" + sha + ".json"
		if _, err = os.Stat(path); err != nil { // le fichier n'existe pas, ce sont des données intéressantes
			log.Println(" Carte à modifier")
			//> on sauvegarde le fichier json
			os.MkdirAll(dir, 0777)
			f, _ := os.Create(path)
			defer f.Close()
			f.Write(bin)
			//> on crée ou enrichit l'image png correspondant à la couche
			bra.EnrichitCouchePNG(dirBase, couche, 20)
			//> et on s'occupe aussi des amis
			if amis != nil {
				for _, ami := range amis {
					log.Println(" enrichissement carte ami ", ami.IdBraldun)
					bra.EnrichitCouchePNG(ms.répertoireCartesBraldun(ami.IdBraldun, ami.Mdpr), couche, 20)
				}
			}
		} else {
			log.Println(" Carte inchangée")
		}
	}

	if in.Cmd == "carte" || in.Cmd == "" { // pour la compatibilité ascendante, la commande est provisoirement optionnelle
		
		//> renseignements sur les couches disponibles
		out.ZConnus, err = bra.CouchesPNGDisponibles(dirBase)
		if err != nil {
			log.Println(" erreur durant la détermination des couches disponibles")
		}

		//> renvoi des données de vues provenant des amis
		if amis != nil && couche != nil {
			ms.fv.Reçoit(in.Vue.Vues[0])
			vues := ms.fv.Complète(in.Vue.Vues[0], amis)
			if len(vues) > 0 {
				out.DV = new(bra.DonnéesVue)
				out.DV.Vues = vues
			}
		}

		//> renvoi de la carte en png
		log.Println(" ZRequis : ", in.ZRequis)
		out.Z = in.ZRequis
		cheminLocalImage := fmt.Sprintf("%s/%d-%s/couche%d.png", *ms.répertoireCartes, in.IdBraldun, in.Mdpr, in.ZRequis)
		if f, err := os.Open(cheminLocalImage); err == nil {
			defer f.Close()
			bytes, _ := ioutil.ReadAll(f)
			out.PngCouche = "data:image/png;base64," + base64.StdEncoding.EncodeToString(bytes)
		}
	} else if in.Cmd == "partages" {
		if in.Action != "" && in.Cible > 0 {
			con.ModifiePartage(in.IdBraldun, in.Cible, in.Action)
		}
		out.Partages, err = con.Partages(in.IdBraldun)
		if err != nil {
			log.Println(" erreur à la récupération des partages :", err)
		}
	}
}

func (server *MapServer) Start() {
	http.Handle("/", server)
	log.Printf("mapserver démarre sur le port %d\n", port)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Println("Erreur au lancement : ", err)
	}
}

func main() {
	ms := new(MapServer)
	ms.répertoireCartes = flag.String("cartes", "", "répertoire des cartes")
	cpuprofile := flag.String("cpuprofile", "", "fichier dans lequel écrire un bilan de profiling cpu")
	memprofile := flag.String("memprofile", "", "fichier dans lequel écrire un bilan de profiling mémoire (lors de l'ordre d'arrêt)")
	mysqluser := flag.String("mysqluser", "", "user pour l'accès mysql")
	mysqlmdp := flag.String("mysqlmdp", "", "mdp pour l'accès mysql")
	mysqldb := flag.String("mysqldb", "braldop", "base mysql")
	flag.Parse()

	ms.bd = bra.NewBaseMysql(*mysqluser, *mysqlmdp, *mysqldb)

	if *ms.répertoireCartes == "" {
		log.Println("Chemin des cartes non fourni")
	} else {
		log.Println("Répertoire des cartes : " + *ms.répertoireCartes)
	}
	if *cpuprofile != "" {
		log.Println("Profiling CPU actif, résultats dans le fichier ", *cpuprofile)
		fp, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(fp)
	}
	ms.fv.Charge(*ms.répertoireCartes)
	go func() {
		sigchan := make(chan os.Signal)
		signal.Notify(sigchan, os.Interrupt, os.Kill)
		for {
			sig := <-sigchan
			log.Println("Signal : %+v", sig)
			log.Println("Mapserver tué ! (", sig, ")")
			if *memprofile != "" {
				log.Println("Ecriture heap dans le fichier ", *memprofile)
				fp, err := os.Create(*memprofile)
				if err != nil {
					log.Fatal(err)
				}
				pprof.WriteHeapProfile(fp)
			}
			bra.BloqueEcrituresPNG()
			if *cpuprofile != "" {
				pprof.StopCPUProfile()
			}
			os.Exit(0)
		}
	}()
	ms.Start()
}
