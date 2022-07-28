module.exports = {
      apps: [
            {
                  name: "STB_marion",
                  script: "./main",
                  env: {
                        file_path: "./assets/SMV_test_fichier.csv",
                        output_file_name: "SMV_Output_file.csv",
                        
                  },
            }
      ]
}