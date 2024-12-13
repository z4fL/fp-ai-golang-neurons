package utility

import (
	"log"

	"github.com/gorilla/mux"
)

func ListRoutes(router *mux.Router) {
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil {
			log.Printf("Route: %s (no methods)", pathTemplate)
			return nil
		}

		log.Printf("Route: %s %s", methods, pathTemplate)
		return nil
	})
	if err != nil {
		log.Fatalf("Error listing routes: %v", err)
	}
}
