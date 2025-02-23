package api

import (
	"net/http"

	"backend-pajak/handlers"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Buat server G4Vercel
	mux := http.NewServeMux()

	// Gunakan metode yang tepat untuk operasi
	mux.HandleFunc("/api/insertdata", handlers.InsertData)
	mux.HandleFunc("/api/getdata", handlers.GetData)
	mux.HandleFunc("/api/getlaporan", handlers.GetLaporan)
	mux.HandleFunc("/api/getdatabyid", handlers.GetDataById)
	mux.HandleFunc("/api/updatedata", handlers.UpdateData)
	mux.HandleFunc("/api/deletedata", handlers.DeleteData)
	mux.HandleFunc("/api/deletelaporan", handlers.DeleteLaporan)
	mux.HandleFunc("/api/detectplate", handlers.DetectPlate)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	corsWrappedHandler := corsHandler.Handler(mux)
	corsWrappedHandler.ServeHTTP(w, r)
}
