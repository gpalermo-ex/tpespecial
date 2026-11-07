package main
import (
	"log"
	"net/http"
	"tpespecial/servidor-musica/views"
	"tpespecial/servidor-musica/db/sqlc"
)
func main() {
	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal("Error al abrir la base de datos: %v", err)
	}
	defer conn.Close()
	queries := db.New(conn)
	
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path != "/" || r.Method != http.MethodGet) {
			http.NotFound(w, r)
			return
		}
		registros, err := queries.ListRegisters(r.Context())
		if err != nil {
			http.Error(w, "Error al obtener la lista de registros", http.StatusInternalServerError)
			log.Printf("Error en ListRegisters: %v", err)
			return
		}

		pagina := views.IndexPage("Catálogo de música", registros)

		err = pagina.Render(r.Context(), w)
		if err != nil {
			http.Error(w, "Error al renderizar la página.", http.StatusInternalServerError)
			log.Printf("Error al renderizar la vista: %v", err)
		}
	})

	http.HandleFunc("/albums", func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path != "/albums" || r.Method != http.MethodPost) {
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
			return
		}
		if err := r.ParseForm(); err != nil {
			log.Printf("Error al ingresar datos: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		id_album := r.FormValue("id_album")
		id_usuario := r.FormValue("id_usuario")
		titulo := r.FormValue("titulo")
		artista := r.FormValue("artista")
		descripcion := r.FormValue("descripcion")
		genero := r.FormValue("genero")
		duracion_min := r.FormValue("duracion_min")
		rating := r.FormValue("rating")

		id_album, err1 := strconv.Atoi(id_album)
		duracion_min, err2 := strconv.Atoi(duracion_min)
		rating, err3 := strconv.Atoi(rating)
		if (err1 != nil || err2 != nil || err3 != nil){
			log.Printf("Valores inválidos: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		arg := db.CreateRegisterParams{
			IDAlbum: id_album,
			IDUsuario: id_usuario,
			Titulo: titulo,
			Artista: artista,
			Descripcion: descripcion,
			Genero: genero,
			DuracionMin: duracion_min,
			Rating: rating,
		}
		err := queries.CreateRegister(r.Context(), arg)
		if err != nil {
			log.Printf("Error al crear el registro en la base de datos: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	
	
		
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
