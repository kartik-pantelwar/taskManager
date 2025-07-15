package errorhandling

import (
	"net/http"
	pkgresponse "notificationservice/src/pkg/response"
)

func HandleError(w http.ResponseWriter, msg string, statusCode int) {
	response := pkgresponse.StandardResponse{
		Status:  "FAILURE",
		Message: msg,
	}
	pkgresponse.WriteResponse(w, statusCode, response)

}

// func HandleError(w http.ResponseWriter, statusCode int, err error) {
// 	response := pkgresponse.StandardResponse{
// 		Status: "FAILURE",
// 		Error:  err,
// 	}
// 	pkgresponse.WriteResponse(w, statusCode, response)

// }
