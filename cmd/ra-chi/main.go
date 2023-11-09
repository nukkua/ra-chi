package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nukkua/ra-chi/internal/app/router"
)

func main (){

	r:= router.SetupRouter();

	fmt.Println("Serving at port :8080");
	log.Fatal(http.ListenAndServe(":8080", r));
}
