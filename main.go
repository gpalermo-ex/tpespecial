package main
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "tpespecial/servidor-musica/db/sqlc"
	"tpespecial/servidor-musica/views"
)
func main() {
	connStr := "host=localhost port=5432 user=admin password=123 dbname=base_catalogo"
	
	// Establecemos la conexión con la bd de PostgreSQL que corre en Docker
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error al configurar la base de datos:", err)
	}

	// Verificamos la conexión a la bd haciendo un ping
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos. ¿Está prendido Docker?:", err)
	}
	fmt.Println("Conexión a la base de datos exitosa")
	defer db.Close()

	// Instancia de queries de sqlc
	queries := sqlc.New(db)
	
	// Handler para "/" con GET
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

		id_album_text := r.FormValue("id_album")
		id_usuario := r.FormValue("id_usuario")
		titulo := r.FormValue("titulo")
		artista := r.FormValue("artista")
		descripcion := r.FormValue("descripcion")
		genero := r.FormValue("genero")
		duracion_min_text := r.FormValue("duracion_min")
		rating_text := r.FormValue("rating")

		id_album, err1 := strconv.Atoi(id_album_text)
		duracion_min, err2 := strconv.Atoi(duracion_min_text)
		rating, err3 := strconv.Atoi(rating_text)
		if (err1 != nil || err2 != nil || err3 != nil){
			log.Printf("Valores inválidos: %v", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		arg := sqlc.CreateRegisterParams{
			IDAlbum: int32(id_album),
			IDUsuario: id_usuario,
			Titulo: titulo,
			Artista: artista,
			Descripcion: descripcion,
			Genero: genero,
			DuracionMin: int32(duracion_min),
			Rating: int32(rating),
		}
		_, err := queries.CreateRegister(r.Context(), arg)
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
