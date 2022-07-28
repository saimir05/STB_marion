package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type mapFromCsv struct {
	Srvid                     string
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

// envoyer chaque ligne dans un objet
// Regrouper les meme stb dans une ligne sur le fichier de sortie

// ReadFile : Lit le fichier ligne par ligne et renvoie le contenu de chaque ligne dans un tableau a 2 dimension
func ReadFile(filePath string) [][]string {
	file_content, err := os.Open("fichier_test")
	var fileLines [][]string
	//var fileLinesSliced []string
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file_content)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		fileLines = append(fileLines, strings.Split(fileScanner.Text(), ";"))
	}

	file_content.Close()

	for _, line := range fileLines {
		fmt.Println(line)
		fmt.Println("nombre d'element : %s", len(line))
		//fileLinesSliced = append(fileLinesSliced,slicedLigne)

	}

	return fileLines
}

func main() {
	ReadFile("fichier_test")
}
