package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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
	option                    string
}

// varible environement pour le nom du fichier a traiter
var FilePath = os.Getenv("file_path")

// header du fichier csv de sortie
//header := ["SRVID","Date d\'ouverture co","Date de souscription","Techno","Reseau","Code INSEE","Nb STB","STB1","STB2","STB3","STB4","STB5","STB6","STB7","STB8","STB9","STB10","Adresse IP","OPTION"]
// initialise les variables avec les variables env
func Assign_env() {
	//fichier par défaut , fichier de test
	if FilePath == "" {
		FilePath = "./assets/SMV_test_fichier.csv"
	}
}

// ReadFile : Lit le fichier ligne par ligne et renvoie le contenu de chaque ligne dans un tableau a 2 dimension
func ReadFile(filePath string) [][]string {
	file_content, err := os.Open(filePath)
	var fileLines [][]string
	if err != nil {
		log.Fatal(err)
	}

	fileScanner := bufio.NewScanner(file_content)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileLines = append(fileLines, strings.Split(fileScanner.Text(), ";"))
	}

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

func arrayToMap(fileLines [][]string) []MapFromFile {
	var arrayOfMap []MapFromFile
	var mapToArrayOfMap MapFromFile
	var line = 0
	for line < len(fileLines) {
		if line == 0 {
			mapToArrayOfMap.srvid = fileLines[line][0]
			mapToArrayOfMap.techno = fileLines[line][3]
			mapToArrayOfMap.reseau = fileLines[line][4]
			mapToArrayOfMap.imei_stb, mapToArrayOfMap.nb_stb = numAndImeiStb(fileLines[line])
			mapToArrayOfMap.date_souscription = fileLines[line][2]
			mapToArrayOfMap.date_ouverture_commerical = fileLines[line][1]
			mapToArrayOfMap.code_insee = fileLines[line][6]
			mapToArrayOfMap.adresse_ip = getIP(fileLines[line])
			mapToArrayOfMap.option = "MULTI"
			arrayOfMap = append(arrayOfMap, mapToArrayOfMap)
			fmt.Println("ligne : %s", line+1)
			line++
		} else if line > 0 {
			if fileLines[line][0] == fileLines[line-1][0] {
				line++
			} else {
				mapToArrayOfMap.srvid = fileLines[line][0]
				mapToArrayOfMap.techno = fileLines[line][3]
				mapToArrayOfMap.reseau = fileLines[line][4]
				mapToArrayOfMap.imei_stb, mapToArrayOfMap.nb_stb = numAndImeiStb(fileLines[line])
				mapToArrayOfMap.date_souscription = fileLines[line][2]
				mapToArrayOfMap.date_ouverture_commerical = fileLines[line][1]
				mapToArrayOfMap.code_insee = fileLines[line][6]
				mapToArrayOfMap.adresse_ip = getIP(fileLines[line])
				mapToArrayOfMap.option = "MULTI"
				arrayOfMap = append(arrayOfMap, mapToArrayOfMap)
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

func main() {
	// configuration depuis les variables environement du fichier ecosystem.config.js
	Assign_env()
	// lecture du fichier
	fileLines := ReadFile("./assets/SMV_test_fichier.csv")
	arrayOfMapFromFile := arrayToMap(fileLines)
	for _, value := range arrayOfMapFromFile {
		fmt.Println(value)
	}
}
