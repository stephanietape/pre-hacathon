package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("Go app...")

	type User struct {
		Id         string `json:"id"`
		Nom        string `json:"nom"`
		Prenom     string `json:"prenom"`
		Age        string `json:"age"`
		Contact    string `json:"conatct"`
		Email      string `json:"email"`
		Competence string `json:"competence"`
		Profession string `json:"prefession"`
		Fichier    string `json:"fichier"`
		Video      string `json:"video"`
		Photo      string `json:"Photo"`
	}

	// PageData représente les données transmises au modèle HTML
	type PageData struct {
		Personnes []User
	}

	type DetailData struct {
		Personne User
	}
	//fonction pour la page acceuil
	home := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" && r.URL.Path == "/" {

			fileContent, err := ioutil.ReadFile("users.json")
			if err != nil {
				fmt.Println("Erreur lors de la lecture du fichier JSON:", err)
				return
			}
			var listuser []User
			err = json.Unmarshal(fileContent, &listuser)
			if err != nil {
				fmt.Println("Erreur lors de la désérialisation JSON:", err)
				return
			}

			fmt.Println(listuser)
			// Crée la structure de données pour le modèle HTML
			data := PageData{
				Personnes: listuser,
			}

			tmpl := template.Must(template.ParseFiles("Templates/index.html"))
			error := tmpl.Execute(w, data)

			if error != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		} else {
			http.NotFound(w, r)

		}
	}
	//fonction du formulaire
	formulaire := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "GET" && r.URL.Path == "/form" {
			tmpl := template.Must(template.ParseFiles("Templates/formulaire.html"))
			err := tmpl.Execute(w, tmpl)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.NotFound(w, r)

		}
	}
	//fonction de l'inscription
	inscription := func(w http.ResponseWriter, r *http.Request) {

		if r.Method == "POST" && r.URL.Path == "/enregistrer" {
			id := uuid.New()
			var mystruct []User
			var user_one User
			var listFiles = uploadHandler(w, r)

			// Use json.Unmarshal to decode the JSON string into the []string variable

			fmt.Println("listFiles", listFiles)
			user_one.Id = id.String()
			user_one.Nom = r.PostFormValue("nom")
			user_one.Prenom = r.PostFormValue("prenom")
			user_one.Age = r.PostFormValue("age")
			user_one.Email = r.PostFormValue("email")
			user_one.Contact = r.PostFormValue("contact")
			user_one.Competence = r.PostFormValue("competence")
			user_one.Profession = r.PostFormValue("profession")
			user_one.Fichier = listFiles[0]
			user_one.Video = "video1.mp4"
			user_one.Photo = "fem1.png"
			mystruct = append(mystruct, user_one)
			fmt.Println(mystruct)

			fileContent, err := ioutil.ReadFile("users.json")
			if err != nil {
				fmt.Println("Erreur lors de la lecture du fichier JSON:", err)
				return
			}
			var listuser []User
			err = json.Unmarshal(fileContent, &listuser)
			if err != nil {
				fmt.Println("Erreur lors de la désérialisation JSON:", err)
				return
			}

			listuser = append(listuser, mystruct[0])
			fmt.Println("final user", listuser)
			// Sérialise la structure en JSON avec une indentation lisible
			jsonData, err := json.MarshalIndent(listuser, "", "  ")
			if err != nil {
				fmt.Println("Erreur lors de la sérialisation JSON:", err)
				return
			}
			// Écrit le JSON dans un fichier
			err = ioutil.WriteFile("users.json", jsonData, 0644)
			if err != nil {
				fmt.Println("Erreur lors de l'écriture dans le fichier JSON:", err)
				return
			}
			fmt.Println("Données stockées avec succès dans le fichier donnees.json")
			// Redirige vers la page detail information
			http.Redirect(w, r, "/detail?id="+id.String(), http.StatusSeeOther)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}
	//page detail des information
	apropos := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/detail" {
			userID := r.URL.Query().Get("id")
			fmt.Println(userID)

			fileContent, err := ioutil.ReadFile("users.json")
			if err != nil {
				http.Error(w, "Erreur lors de la lecture du fichier JSON", http.StatusInternalServerError)
				return
			}

			// Désérialise le JSON dans une slice de Person
			var personnes []User
			err = json.Unmarshal(fileContent, &personnes)
			if err != nil {
				http.Error(w, "Erreur lors de la désérialisation JSON", http.StatusInternalServerError)
				return
			}

			// Recherche de l'utilisateur en fonction de l'ID
			var utilisateurTrouve User
			for _, personne := range personnes {
				if personne.Id == userID {
					utilisateurTrouve = personne
					break
				}
			}

			// Vérifie si l'utilisateur a été trouvé
			if utilisateurTrouve.Id != "" {
				// Affiche les détails de l'utilisateur ou effectue d'autres actions
				fmt.Println("user trouver", utilisateurTrouve)
				data := DetailData{
					Personne: utilisateurTrouve,
				}
				tmpl := template.Must(template.ParseFiles("Templates/aPropos.html"))
				error := tmpl.Execute(w, data)

				if error != nil {
					http.Error(w, error.Error(), http.StatusInternalServerError)
					return
				}

			} else {
				fmt.Fprint(w, "Utilisateur non trouvé.")
			}

		}
	}
	//afficher les pdf
	pdfHandler := func(w http.ResponseWriter, r *http.Request) {
		// Récupère le nom du fichier PDF depuis l'URL
		fileName := r.URL.Path[len("/pdf/"):]

		// Lit le contenu du fichier PDF
		pdfContent, err := ioutil.ReadFile("public/" + fileName)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du fichier PDF", http.StatusInternalServerError)
			return
		}

		// Définit le type de contenu à "application/pdf"
		w.Header().Set("Content-Type", "application/pdf")
		// Écrit le contenu du fichier PDF dans la réponse
		w.Write(pdfContent)
	}
	//afficher les video
	videoHandler := func(w http.ResponseWriter, r *http.Request) {
		// Récupère le nom du fichier vidéo depuis l'URL
		fileName := r.URL.Path[len("/video/"):]

		// Lit le contenu du fichier vidéo
		videoContent, err := ioutil.ReadFile("public/" + fileName)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du fichier vidéo", http.StatusInternalServerError)
			return
		}

		// Définit le type de contenu à "video/mp4"
		w.Header().Set("Content-Type", "video/mp4")
		// Écrit le contenu du fichier vidéo dans la réponse
		w.Write(videoContent)
	}

	// define handlers
	http.HandleFunc("/form", formulaire)
	http.HandleFunc("/", home)
	http.HandleFunc("/detail", apropos)
	http.HandleFunc("/enregistrer", inscription)
	http.HandleFunc("/pdf/", pdfHandler)
	http.HandleFunc("/video/", videoHandler)
	http.Handle("/ressources/", http.StripPrefix("/ressources/", http.FileServer(http.Dir("ressources"))))
	log.Fatal(http.ListenAndServe(":8000", nil))

}

func uploadHandler(w http.ResponseWriter, r *http.Request) []string {
	// Parse le formulaire avec une limite de taille pour éviter les attaques DOS
	err := r.ParseMultipartForm(10 << 20) // 10 MB limite
	if err != nil {
		http.Error(w, "Impossible de parser le formulaire", http.StatusInternalServerError)
	}
	// Récupère tous les fichiers uploadés
	files := r.MultipartForm.File["files"]

	// Dossier où stocker les fichiers
	uploadDir := "./public/"
	// Crée le dossier s'il n'existe pas
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Impossible de créer le dossier d'upload", http.StatusInternalServerError)
	}

	var listFiles []string

	// Boucle à travers les fichiers et les stocke dans le dossier d'upload
	for _, file := range files {
		// Ouvre le fichier uploadé
		src, err := file.Open()
		if err != nil {
			http.Error(w, "Erreur lors de l'ouverture du fichier uploadé", http.StatusInternalServerError)
		}
		defer src.Close()

		// Crée le fichier destination dans le dossier d'upload
		listFiles = append(listFiles, file.Filename)
		dst, err := os.Create(filepath.Join(uploadDir, file.Filename))
		if err != nil {
			http.Error(w, "Erreur lors de la création du fichier destination", http.StatusInternalServerError)
		}
		defer dst.Close()

		// Copie le contenu du fichier uploadé vers le fichier destination
		if _, err := io.Copy(dst, src); err != nil {
			http.Error(w, "Erreur lors de la copie du contenu du fichier", http.StatusInternalServerError)

		}
	}

	// Redirige vers la liste des fichiers après l'upload
	//http.Redirect(w, r, "/", http.StatusSeeOther)
	return listFiles
}
