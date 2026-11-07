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

	t.Run("CreateRegister", func(t *testing.T){
		createdRegister, err := queries.CreateRegister(ctx,
			sqlc.CreateRegisterParams{
				IDAlbum: 1,
				IDUsuario: "vale_jaurena",
				Titulo: "Ya es mañana",
				Artista: "Morat",
				Descripcion: "Habla de que hoy es mañana",
				Genero: "Pop-rock",
				DuracionMin: 40,
				Rating: 4,
			})
		if err != nil {
			t.Fatalf("Error al crear registro: %v", err)
		}

		createdID = createdRegister.IDRegistro
		t.Logf("Registro creado: %+v\n", createdRegister)
	})

	t.Run("GetRegister", func(t *testing.T){
		registro, err := queries.GetRegister(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al obtener registro: %v", err)
		}
		t.Logf("Registro obtenido: %+v\n", registro)
	})

	t.Run("ListRegisters", func(t *testing.T){
		registros, err := queries.ListRegisters(ctx)
		if err != nil {
			t.Fatalf("Error al listar registros: %v", err)
		}
		t.Logf("Todos los registros: %+v\n", registros)
	})

	t.Run("UpdateRegister", func(t *testing.T){
		err = queries.UpdateRegister(ctx, sqlc.UpdateRegisterParams{
			IDRegistro: createdID,
			IDAlbum: 4,
			IDUsuario: "gonza_pal",
			Titulo: "Detour",
			Artista: "Kim Petras",
			Descripcion: "Álbum de música pop",
			Genero: "Pop",
			DuracionMin: 36,
			Rating: 4,
		})
		if err != nil {
			t.Fatalf("Error al actualizar registro: %v", err)
		}
		t.Logf("Registro actualizado correctamente.")
		
		
		updatedRegister, err := queries.GetRegister(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al obtener registro actualizado: %v", err)
		}
		t.Logf("Registro actualizado: %+v\n", updatedRegister)
	})
		
	t.Run("UpdateRegister", func(t *testing.T){
		err = queries.UpdateRegister(ctx, sqlc.UpdateRegisterParams{
			IDRegister: createdID,
			IDAlbum: 6,
			IDUsuario: "gonza_pal",
			Titulo: "Detour Deluxe",
			Artista: "Kim Petras 2.0",
			Descripcion: "Álbum de música pop 2.0",
			Genero: "Pop-Rock",
			DuracionMin: 30,
			Rating: 0,
		})
		if err != nil {
			t.Logf("Error al actualizar registro con rating nulo: %v", err)
		} else {
			t.Fatalf("Registro con rating nulo actualizado correctamente")
		}
	})


	t.Run("DeleteRegister", func(t *testing.T){
		err = queries.DeleteRegister(ctx, createdID)
		if err != nil {
			t.Fatalf("Error al borrar registro: %v", err)
		}
		t.Logf("Registro borrado correctamente.")

		_, err = queries.GetRegister(ctx, createdID)
		if err == sql.ErrNoRows {
			t.Logf("Registro borrado no encontrado.")
		} else if err != nil {
			t.Fatalf("Error al obtener registro despues de borrarlo: %v", err)
		}
	})
}
