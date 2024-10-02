package market

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func PrintMyString(writer http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)

	s := vars["string"]

	fmt.Fprintf(writer, "Here is a string: %s\n", s)

}
