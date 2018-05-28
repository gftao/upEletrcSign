package handle

import (
	"net/http"
	"fmt"
	"upEletrcSign/logr"
)

func DoHandle(w http.ResponseWriter, r *http.Request) {
	logr.Debug("天龙人")
	logr.Debug("海贼王")

	fmt.Fprintf(w, "Welcome to the home page!")
	return
}
