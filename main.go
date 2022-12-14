package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MapFromFile struct {
	srvid                     string
	date_ouverture_commerical string
	date_souscription         string
	techno                    string
	reseau                    string
	code_insee                string
	nb_stb                    int
	imei_stb                  []string
	adresse_ip                string
	options                   []string
}

// varible environement pour le nom du fichier a traiter
var FilePath = os.Getenv("file_path")

// varibale environement pour le nom du fichier de sortie
var OutputFileName = os.Getenv("output_file_name")

// header du fichier csv de sortie

// initialise les variables avec les variables env
func Assign_env() {
	//fichier par défaut a traiter par defaut , si aucun argument n'est passé en paramétre  : /assets/extract.csv
	if FilePath == "" {
		fmt.Println(" args : ")
		fmt.Println(len(os.Args))
		fmt.Println(os.Args)
		if len(os.Args) >= 2 {
			FilePath = "./assets/" + os.Args[1]
		} else {
			FilePath = "./assets/extract.csv"
		}
	}

	if OutputFileName == "" {

		OutputFileName = "Output_file.csv"
	}
}

// timeTrack : calcul du temps d'execution d'une fonction , a utilisé avec "defer"
func timeTrack(start time.Time, name string) {
	println("-----------------------")
	elapsed := time.Since(start)
	log.Printf("%s : %s", name, elapsed)
}

// ReadFile : Lit le fichier ligne par ligne et renvoie le contenu de chaque ligne dans un tableau a 2 dimension
func ReadFile(filePath string) [][]string {

	defer timeTrack(time.Now(), "Read file")

	fmt.Println("Ouverture du fichier...")
	file_content, err := os.Open(filePath)
	// fermeture du fichier a la fin de l'execution de la fonction
	defer file_content.Close()
	var fileLines [][]string
	if err != nil {
		log.Fatal(err)
	}
	// Lecture ligne par ligne du fichier ouvert
	fileScanner := bufio.NewScanner(file_content)
	println("-file scanner split ...")
	fileScanner.Split(bufio.ScanLines)
	//ajout de chaque ligne sous forme de tableau de string dans un autre tableau pour faire un tableau a 2 dimensions
	fmt.Println("-Lecture du fichier ligne par ligne ...")
	for fileScanner.Scan() {
		var lineScanner = fileScanner.Text()
		fileLines = append(fileLines, strings.Split(lineScanner, ";"))
	}
	fmt.Println("-Tableau a 2 dimensions créé !")
	fileScanner = nil
	return fileLines
}

// recupére les infos du clients pour créer un objet de type MapFromFile , definit plus haut
func arrayToMap(fileLines [][]string) []MapFromFile {
	defer timeTrack(time.Now(), "Array to map")
	//Tableau de MapFromFIle qui contiendra les lignes a insérer dans le fichier CSV de sortie
	var arrayOfMap []MapFromFile
	// defintion d'une ligne dans le fichier CSV de sortie
	var mapOfLine MapFromFile
	var line = 0
	println("Récuperation des différents champs grace avec une RegEx ...")
	for line < len(fileLines) {
		if line > 0 && fileLines[line][0] == fileLines[line-1][0] {
			mapOfLine.options = append(mapOfLine.options, getOptions(fileLines[line]))
			if line+1 >= len(fileLines) {
				arrayOfMap = append(arrayOfMap, mapOfLine)
			}
			line++
		} else {
			arrayOfMap = append(arrayOfMap, mapOfLine)
			mapOfLine.options = nil
			mapOfLine.srvid = fileLines[line][0]
			mapOfLine.techno = fileLines[line][3]
			mapOfLine.reseau = fileLines[line][4]
			mapOfLine.imei_stb, mapOfLine.nb_stb = numAndImeiStb(fileLines[line])
			mapOfLine.date_souscription = fileLines[line][2]
			mapOfLine.date_ouverture_commerical = fileLines[line][1]
			mapOfLine.code_insee = fileLines[line][6]
			mapOfLine.adresse_ip = getIP(fileLines[line])
			mapOfLine.options = append(mapOfLine.options, getOptions(fileLines[line]))
			line++
		}

	}
	return arrayOfMap
}

//retourne la liste des imei , et le nombre de stb a aprtir de la ligne de fichier lue actuellement
func numAndImeiStb(line []string) ([]string, int) {
	var re = regexp.MustCompile(`^2[0-9]{14}`)
	var listImei []string
	var nbImei = 0
	for _, value := range line {
		if re.FindString(value) != "" {
			nbImei++
			listImei = append(listImei, value)
		}
	}
	return listImei, nbImei
}

// via une expression reguliére , on récupere l'option de la ligne actuelle
func getOptions(line []string) string {
	var option = line[26]
	return option
}

// retourne l'adresse IP a aprtir de la ligne de fichier lue actuellement
func getIP(line []string) string {
	var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	var ip string
	for _, value := range line {
		if re.FindString(value) != "" {
			ip = value
		}
	}
	return ip
}

//Trie les données du tableau à deux dimensions
func sortData(arrayOfMap []MapFromFile) [][]string {
	defer timeTrack(time.Now(), "Sort Data")
	fmt.Println("trie des données dans les bonnes colonnes...")
	var dataToCSV [][]string
	header := []string{"SRVID", "Date d'ouverture co", "Date de souscription", "Techno", "Reseau", "Code INSEE", "Adresse IP", "OPTIONS", "Nb STB", "STB1", "STB2", "STB3", "STB4", "STB5", "STB6", "STB7", "STB8", "STB9", "STB10"}
	//ajout des titres des colonnes
	dataToCSV = append(dataToCSV, header)
	for i := 0; i < len(arrayOfMap); i++ {
		var line []string
		line = append(line, arrayOfMap[i].srvid)
		line = append(line, arrayOfMap[i].date_ouverture_commerical)
		line = append(line, arrayOfMap[i].date_souscription)
		line = append(line, arrayOfMap[i].techno)
		line = append(line, arrayOfMap[i].reseau)
		line = append(line, arrayOfMap[i].code_insee)

		line = append(line, arrayOfMap[i].adresse_ip)
		var strOption = ""
		for index, value := range arrayOfMap[i].options {
			if index == 0 {
				strOption = strOption + value
			} else {
				strOption = strOption + " | " + value
			}
		}
		line = append(line, strOption)
		// Ce curseur se deplace de colonne en colonne , il permet de sauter le colonne vide pour les imei
		var cursor = 8
		line = append(line, strconv.FormatInt(int64(arrayOfMap[i].nb_stb), 10))
		for index := 0; index <= 9; index++ {

			if index < len(arrayOfMap[i].imei_stb) {
				line = append(line, arrayOfMap[i].imei_stb[index])
			} else {
				line = append(line, "")
			}
			cursor++
		}
		dataToCSV = append(dataToCSV, line)
	}
	fmt.Println("- Fin du trie !")
	return dataToCSV
}

// Creation d'un fichier CSV a partir d'un tableau a 2 dimensions et un header
func writeCSV(dataToCSV [][]string, outputFileName string) {
	fmt.Println("Ecriture du fichier de sortie...")
	csvFile, err := os.Create(outputFileName)

	if err != nil {
		log.Fatalf("Erreur dans la création di fichier: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ';'

	for _, empRow := range dataToCSV {
		_ = csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	csvFile.Close()
	fmt.Println("fichier de sortie \"Output_file.csv\" créé !")
}

func awkCommand(fileName string, outPutfileName string) {
	command := "awk -F\\; '{print $1,$2,$3,$4,$5,$6,$18,$19,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17}' OFS=';' \"./" + fileName + "\" > " + outPutfileName

	cmd := exec.Command(command)

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	defer timeTrack(time.Now(), "FONCTION MAIN")
	// configuration depuis les variables environement du fichier ecosystem.config.js
	Assign_env()
	// lecture du fichier

	// Récuperation de toutes les lignes du fichier

	// Création du tableau a deux dimensions nécessaire pour l'ecriture du fichier

	writeCSV(sortData(arrayToMap(ReadFile(FilePath))), OutputFileName)

	fmt.Println("FIN !!!!")
	//awkCommand("temp.csv", OutputFileName)
}

// liste des options
var optionsList = []string{
	"BBOX_FIBRE",
	"RADIO_FIBRE",
	"EXT_OTT_ACCESS",
	"EXT_OTT_TVBASIC",
	"FIBRE_BASIC",
	"FIBRE_CHAINE_HUSTLER",
	"HD",
	"IPTV_MYTF1MAX_INCLUS",
	"IPTV_STO",
	"LIVE",
	"SVOD",
	"VOD",
	"VOD_24",
	"IPTV_NPVR_BASIC_PARC",
	"NPVR_TYPE",
	"BBOX_BYT",
	"RADIO_BYT",
	"IPTV_BYT_BASIC",
	"IPTV_BYT_BASIC_ADVAN",
	"OTT_LIVE_DRM_STB",
	"TIME",
	"FIBRE_SALTO",
	"FIBRE_BBOXCINE",
	"RISK_2",
	"BBOX_FT",
	"RADIO_FT",
	"IPTV_FT_BASIC",
	"IPTV_FT_BASIC_ADVANC",
	"FIBRE_OL_TV",
	"PVR_40",
	"FIBRE_BBOXCINE_ES",
	"IPTV_BYT_OCS",
	"IPTV_SVOD_OCS",
	"FIBRE_BEIN_CRYP1",
	"FIBRE_BEIN_CRYP2",
	"FIBRE_BEIN_CRYP3",
	"FIBRE_BEIN_SPORT",
	"FIBRE_DIVERTI",
	"FIBRE_OCS",
	"MULTI",
	"IPTV_MYTF1MAX",
	"PVR",
	"FIBRE_PACK_X",
	"IPTV_BYT_PACK_LUSOPH",
	"IPTV_PBCANALB",
	"IPTV_BYT_JEUNESSEB",
	"FIBRE_PACK_3",
	"IPTV_BYT_BBOXCINE",
	"IPTV_BYT_GRAND_ANGLE",
	"IPTV_BYT_SALTO",
	"FIBRE_JEUNESSE",
	"IPTV_PLAYZERB",
	"FIBRE_DISNEY",
	"IPTV_FT_JEUNESSEB",
	"IPTV_FT_BBOXCINE",
	"IPTV_FT_DIVERTI",
	"IPTV_FT_BBOXCINE_ES",
	"IPTV_FT_SALTO",
	"FIBRE_PCK_SPORT",
	"FIBRE_GRAND_ANGLE",
	"IPTV_FT_BLACKPILLS",
	"IPTV_FT_CHAINE_STAR",
	"IPTV_FT_PACK_MUSIC",
	"IPTV_BYT_ARABE_MAXI",
	"USAGE_REPORT",
	"FIBRE_JEUNESSEB",
	"FIBRE_EUROSPORTB",
	"FIBRE_MAGHREB",
	"RISK_1"}
