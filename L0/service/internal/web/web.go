package web

import (
	"encoding/json"
	"log"
	"net/http"
	"service/service/internal/cache"
)

func StartServer(port string, cache *cache.Cache) {
	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		uid := r.URL.Path[len("/order/"):]
		if uid == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order UID is required"})
			return
		}

		order, ok := cache.Get(uid)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "Order not found"})
			return
		}

		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Printf("JSON encode error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Internal server error"})
		}
	})

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
