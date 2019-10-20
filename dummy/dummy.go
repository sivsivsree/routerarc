package dummy

import (
	"fmt"
	"net/http"
)

func Run() {

	fmt.Println("Test Server running on port http://localhost:8000  for test")
	http.HandleFunc("/", HelloServer)
	_ = http.ListenAndServe(":8000", nil)
}

func HelloServer(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
