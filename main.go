package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
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
	var re = regexp.MustCompile(`^[A-Z_]*$`)
	for i := 4; i < len(line); i++ {
		if re.FindString(line[i]) != "" && line[i] != "BYT" && line[i] != "FT" || line[i] == "VOD_24" || line[i] == "FIBRE_BEIN_CRYP1" || line[i] == "FIBRE_BEIN_CRYP2" || line[i] == "FIBRE_BEIN_CRYP3" || line[i] == "FIBRE_BEIN_CRYP4" || line[i] == "FIBRE_BEIN_CRYP5" {
			option = line[i]
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

func main() {
	// configuration depuis les variables environement du fichier ecosystem.config.js
	Assign_env()
	// lecture du fichier
	fmt.Println("lecture du fichier")
	fileLines := ReadFile(FilePath)
	fmt.Println("reception de l'array de Map ")
	arrayOfMapFromFile := arrayToMap(fileLines)
	for _, value := range arrayOfMapFromFile {
		fmt.Println("")
		fmt.Println(value)
	}
	fmt.Println("---------------------------------------------")
	fmt.Println("trie des données ")
	fmt.Println("---------------------------------------------")
	dataToCSVOut := sortData(arrayOfMapFromFile)
	for _, value := range dataToCSVOut {
		fmt.Println("")
		fmt.Println(value)
	}
	writeCSV(dataToCSVOut, OutputFileName)
}
