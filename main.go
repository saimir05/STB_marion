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
	//fichier par défaut , fichier de test
	if FilePath == "" {
		FilePath = "./assets/SMV_test_fichier.csv"
	}

	if OutputFileName == "" {
		OutputFileName = "SMV_Output_file.csv"
	}
}

// ReadFile : Lit le fichier ligne par ligne et renvoie le contenu de chaque ligne dans un tableau a 2 dimension
func ReadFile(filePath string) [][]string {
	file_content, err := os.Open(filePath)
	var fileLines [][]string
	if err != nil {
		log.Fatal(err)
	}
	// Lecture ligne par ligne du fichier ouvert
	fileScanner := bufio.NewScanner(file_content)
	fileScanner.Split(bufio.ScanLines)
	//ajout de chaque ligne sous forme de tableau de string dans un autre tableau pour faire un tableau a 2 dimensions
	for fileScanner.Scan() {
		fileLines = append(fileLines, strings.Split(fileScanner.Text(), ";"))
	}
	// fermeture du fichier
	file_content.Close()

	return fileLines
}

// Initialise l'objet MapFromFile depuis la ligne lu
// func initMapFromLine(line []string) MapFromLine {
// 	var map MapFromFile
// 	map.srvid = fileLines[line][0]
// 	map.techno = fileLines[line][3]
// 	map.reseau = fileLines[line][4]
// 	map.imei_stb, map.nb_stb = numAndImeiStb(fileLines[line])
// 	map.date_souscription = fileLines[line][2]
// 	map.date_ouverture_commerical = fileLines[line][1]
// 	map.code_insee = fileLines[line][6]
// 	map.adresse_ip = getIP(fileLines[line])
// 	arrayOfMap = append(arrayOfMap, mapToArrayOfMap)
// 	fmt.Println("ligne : %s", line+1)
// 	return map
// }

// recupére les infos du clients pour créer un objet de type MapFromFile , definit plus haut
func arrayToMap(fileLines [][]string) []MapFromFile {
	//Tableau de MapFromFIle qui contiendra les lignes a insérer dans le fichier CSV de sortie
	var arrayOfMap []MapFromFile
	// defintion d'une ligne dans le fichier CSV de sortie
	var mapOfLine MapFromFile
	var line = 0
	for line < len(fileLines) {
		if line == 0 {
			mapOfLine.srvid = fileLines[line][0]
			mapOfLine.techno = fileLines[line][3]
			mapOfLine.reseau = fileLines[line][4]
			mapOfLine.imei_stb, mapOfLine.nb_stb = numAndImeiStb(fileLines[line])
			mapOfLine.date_souscription = fileLines[line][2]
			mapOfLine.date_ouverture_commerical = fileLines[line][1]
			mapOfLine.code_insee = fileLines[line][6]
			mapOfLine.adresse_ip = getIP(fileLines[line])
			mapOfLine.options = append(mapOfLine.options, getOptions(fileLines[line]))
			fmt.Println("ligne : %s", line+1)
			line++
		} else if line > 0 {
			if fileLines[line][0] == fileLines[line-1][0] {
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
				fmt.Println("ligne : %s", line+1)
				line++
			}
		} else {
			fmt.Println("ERROR : OUT OF RANGE , curseur negative !")
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
		if re.FindAllString(value, -1) != nil {
			nbImei++
			listImei = append(listImei, value)
		}
	}
	return listImei, nbImei
}

// via une expression reguliére , on récupere l'option de la ligne actuelle
func getOptions(line []string) string {
	var option string

	for i := 4; i < len(line); i++ {
		for _, value := range optionsList {
			if value == line[i] {
				option = line[i]
			}
		}
	}
	return option
}

// retourne l'adresse IP a aprtir de la ligne de fichier lue actuellement
func getIP(line []string) string {
	var re = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
	var ip string
	for _, value := range line {
		if re.FindAllString(value, -1) != nil {
			ip = value
		}
	}
	return ip
}

//Trie les données du tableau à deux dimensions
func sortData(arrayOfMap []MapFromFile) [][]string {
	var dataToCSV [][]string
	header := []string{"SRVID", "Date d'ouverture co", "Date de souscription", "Techno", "Reseau", "Code INSEE", "Nb STB", "STB1", "STB2", "STB3", "STB4", "STB5", "STB6", "STB7", "STB8", "STB9", "STB10", "Adresse IP", "OPTIONS"}
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
		line = append(line, strconv.FormatInt(int64(arrayOfMap[i].nb_stb), 10))
		var cursor = 7
		for index := 0; index <= 9; index++ {

			if index < len(arrayOfMap[i].imei_stb) {
				line = append(line, arrayOfMap[i].imei_stb[index])
			} else {
				line = append(line, "")
			}
			cursor++
		}

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
		dataToCSV = append(dataToCSV, line)
	}
	return dataToCSV
}

// Creation d'un fichier CSV a partir d'un tableau a 2 dimensions et un header
func writeCSV(dataToCSV [][]string, outputFileName string) {
	csvFile, err := os.Create(outputFileName)

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)
	csvwriter.Comma = ';'

	for _, empRow := range dataToCSV {
		_ = csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	csvFile.Close()
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
	// configuration depuis les variables environement du fichier ecosystem.config.js
	Assign_env()
	// lecture du fichier
	fmt.Println("lecture du fichier")
	// Récuperation de toutes les lignes du fichier
	fileLines := ReadFile(FilePath)
	fmt.Println("reception de l'array de Map ")
	// Création d'un tableau de map
	arrayOfMapFromFile := arrayToMap(fileLines)
	for _, value := range arrayOfMapFromFile {
		fmt.Println("")
		fmt.Println(value)
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("trie des données ")
	fmt.Println("---------------------------------------------")
	// Création du tableau a deux dimensions nécessaire pour l'ecriture du fichier 
	dataToCSVOut := sortData(arrayOfMapFromFile)
	for _, value := range dataToCSVOut {
		fmt.Println("")
		fmt.Println(value)
	}
	// Ecriture du fichier csv temporaire
	writeCSV(dataToCSVOut, "temp.csv")
	// Réorganisation des colonnés pour avoir les imei et le nombre d'STB a la fin 
	awkCommand("temp.csv", OutputFileName)
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
