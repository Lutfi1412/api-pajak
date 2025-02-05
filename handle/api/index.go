package api

import (
	"fmt"
	"net/http"

	"backend-pajak/handlers"

	_ "github.com/lib/pq"
	"github.com/rs/cors"
	. "github.com/tbxark/g4vercel"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	// Buat server G4Vercel
	server := New()
	server.Use(Recovery(func(err interface{}, c *Context) {
		if httpError, ok := err.(HttpError); ok {
			c.JSON(httpError.Status, H{
				"message": httpError.Error(),
			})
		} else {
			message := fmt.Sprintf("%s", err)
			c.JSON(500, H{
				"message": message,
			})
		}
	}))

	// Gunakan metode yang tepat untuk operasi
	server.POST("/api/insertdata", handlers.InsertData)
	server.GET("/api/getdata", handlers.GetData)
	server.GET("/api/getlaporan", handlers.GetLaporan)
	server.POST("/api/getdatabyid", handlers.GetDataById)
	server.POST("/api/updatedata", handlers.UpdateData)
	server.POST("/api/deletedata", handlers.DeleteData)
	server.POST("/api/deletelaporan", handlers.DeleteLaporan)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                       // Domain React Anda
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Metode yang diizinkan
		AllowedHeaders:   []string{"Content-Type", "Authorization"},           // Header yang diizinkan
		AllowCredentials: true,                                                // Izinkan kredensial seperti cookies
	})

	// Bungkus handler utama dengan CORS
	corsWrappedHandler := corsHandler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Handle(w, r)
	}))
	corsWrappedHandler.ServeHTTP(w, r)
}
