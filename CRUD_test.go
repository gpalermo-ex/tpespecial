package main

import (
	"testing"
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	sqlc "tpespecial/servidor-musica/db/sqlc"
)

func TestCrud(t *testing.T){
	connStr := "host=localhost port=5432 user=admin password=123 dbname=base_catalogo"
	db, err := sql.Open("pgx",connStr)
	if err != nil {
		t.Fatalf("Error al conectarse a la BD: %v", err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		t.Fatalf("Error ping de BD: %v", err)
	}
	t.Logf("Conexión exitosa a la base de datos.")
	
	queries := sqlc.New(db)
	ctx := context.Background()
	
	var createdID int32

	t.Run("CreateAlbum", func(t *testing.T){
		createdAlbum, err := queries.CreateAlbum(ctx,
			sqlc.CreateAlbumParams{
				IDUsuario: 1,
				Titulo: "Ya es mañana",
				Artista: "Morat",
				Descripcion: "Habla de que hoy es mañana",
				Genero: "Pop-rock",
				DuracionS: 40000,
				Ranking: 4,
			})
		if err != nil {
			t.Fatalf("Error al crear álbum: %v", err)
		}

		createdID = createdAlbum.IDAlbum
		t.Logf("Álbum creado: %+v\n", createdAlbum)
	})

	t.Run("GetAlbum", func(t *testing.T){
		album, err := queries.GetAlbum(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al obtener álbum: %v", err)
		}
		t.Logf("Álbum obtenido: %+v\n", album)
	})

	t.Run("ListAlbums", func(t *testing.T){
		albums, err := queries.ListAlbums(ctx)
		if err != nil {
			t.Fatalf("Error al listar albumes: %v", err)
		}
		t.Logf("Todos los albumes: %+v\n", albums)
	})

	t.Run("UpdateAlbum", func(t *testing.T){
		err = queries.UpdateAlbum(ctx, sqlc.UpdateAlbumParams{
			IDAlbum: createdID,
			IDUsuario: 2,
			Titulo: "Detour",
			Artista: "Kim Petras",
			Descripcion: "Álbum de música pop",
			Genero: "Pop",
			DuracionS: 360,
			Ranking: 4,
		})
		if err != nil {
			t.Fatalf("Error al actualizar álbum: %v", err)
		}
		t.Logf("Álbum actualizado correctamente.")
		
		
		updatedAlbum, err := queries.GetAlbum(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al obtener álbum actualizado: %v", err)
		}
		t.Logf("Álbum actualizado: %+v\n", updatedAlbum)
	})
		
	t.Run("UpdateAlbum", func(t *testing.T){
			err = queries.UpdateAlbum(ctx, sqlc.UpdateAlbumParams{
			IDAlbum: createdID,
			IDUsuario: 2,
			Titulo: "Detour Deluxe",
			Artista: "Kim Petras 2.0",
			Descripcion: "Álbum de música pop 2.0",
			Genero: "Pop-Rock",
			DuracionS: 360,
			Ranking: 0,
		})
		if err != nil {
			t.Logf("Error al actualizar álbum con ranking nulo: %v", err)
		} else {
			t.Fatalf("Álbum con ranking nulo actualizado correctamente")
		}
	})


	t.Run("DeleteAlbum", func(t *testing.T){
		err = queries.DeleteAlbum(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al borrar álbum: %v", err)
		}
		t.Logf("Álbum borrado correctamente.")

		_, err = queries.GetAlbum(ctx, createdID)
		if err == sql.ErrNoRows {
			t.Logf("Álbum borrado no encontrado.")
		} else if err != nil {
			t.Fatalf("Error al obtener álbum despues de borrarlo: %v", err)
		}
	})
}
