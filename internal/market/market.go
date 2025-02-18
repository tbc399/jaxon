package market

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func PrintMyString(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	s := vars["string"]

	fmt.Fprintf(writer, "Here is a string: %s\n", s)
}
