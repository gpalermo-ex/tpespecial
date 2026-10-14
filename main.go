package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "tpespecial/servidor-musica/db/sqlc"
)

// Creamos una variable global para el "Queries" de sqlc, así lo pueden usar todos los handlers
var queries *sqlc.Queries

func main() {
	// 1. Establecemos la conexión con la base de datos PostgreSQL que corre en Docker
	// Usamos los mismos datos configurados en tu docker-compose.yml
	connStr := "host=localhost port=5432 user=admin password=123 dbname=base_catalogo"
	
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal("Error al configurar la base de datos:", err)
	}

	// 2. Verificamos que la conexión realmente funcione haciendo un "Ping"
	if err := db.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la base de datos. ¿Está prendido Docker?:", err)
	}
	fmt.Println("¡Conexión a la base de datos exitosa!")

	// Importante: Cerramos la conexión al final del programa para no dejar conexiones colgadas
	defer db.Close()

	// 3. Creamos la instancia del repositorio generado por sqlc
	// Esto nos va a permitir usar las funciones CreateAlbum, GetAlbum, etc. en los handlers
	queries = sqlc.New(db)

	// Configurar rutas
	http.HandleFunc("/albums", albumsHandler)
	http.HandleFunc("/albums/", albumHandler)
	
	// Iniciar servidor
	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Manejador para /albums
func albumsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getAlbums(w, r)
	case http.MethodPost:
		createAlbum(w, r)
	default:
		http.Error(w, "Method not allowed",
      http.StatusMethodNotAllowed)
	}
}

func albumHandler(w http.ResponseWriter, r *http.Request) {
	// Extraer ID del path
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "URL inválida", http.StatusBadRequest)
		return
	}
	_, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "ID de álbum inválida", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodGet:
		getAlbum(w, r)
	case http.MethodPut:
		updateAlbum(w, r)
	case http.MethodDelete:
		deleteAlbum(w, r)
	default:
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
	}
}

// GET /albums - Listar todos los albums
func getAlbums(w http.ResponseWriter, r *http.Request) {
	albums, err := queries.ListAlbums(r.Context())
	if err != nil {
		log.Println("Error de la BD al obtener álbumes:", err)
		http.Error(w, "Error al obtener los álbumes", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(albums)
}

// POST /albums - Crear nuevo album
func createAlbum(w http.ResponseWriter, r *http.Request) {
	var params sqlc.CreateAlbumParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(params.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	newAlbum, err := queries.CreateAlbum(r.Context(), params)
	if err != nil {
		log.Println("Error de la BD al crear álbum:", err)
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(newAlbum)
}

// Listar UN álbum
func getAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/albums/")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	album, err := queries.GetAlbum(r.Context(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Álbum no encontrado", http.StatusNotFound) // 404
			return
		}
		log.Println("Error de la BD al obtener álbum:", err)
		http.Error(w, "Error al buscar el álbum", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(album)
}

func updateAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/albums/")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	var params sqlc.UpdateAlbumParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(params.Titulo) == "" {
		http.Error(w, "El título no puede estar vacío", http.StatusBadRequest)
		return
	}

	params.IDAlbum = int32(id)

	err = queries.UpdateAlbum(r.Context(), params)
	if err != nil {
		log.Println("Error de la BD al actualizar el álbum:", err)
		http.Error(w, "Error al actualizar el álbum", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/albums/")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = queries.DeleteAlbum(r.Context(), int32(id))
	if err != nil {

		log.Println("Error de la BD al borrar álbum:", err)
		http.Error(w, "Error al eliminar el álbum", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent) 
}
