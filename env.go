package main;

import (
	"os"
)

// chemin du fichier
var FilePath = os.Getenv("file_path")

func Assign_env() {
	//fichier par défaut , fichier de test
	if FilePath == "" {
		FilePath = "fichier_test"
	}
}