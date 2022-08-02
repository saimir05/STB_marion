version 1.2

Build:
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build main.go
-----------------------------------------------------------
Lancement du programme :

./main <NOM_DU_FICHIER> <NOM_DU_FICHIER_DE_SORTIE.CSV> // je preconise de mettre le fichier a traiter dans le dossier <ASSETS>


pm2 start ecosystem_gpark.config.js