package errorhandling

import (
	"net/http"
	pkgresponse "user_service/src/pkg/response"
)

func HandleError(w http.ResponseWriter, statusCode int, err error) {
	response := pkgresponse.StandardResponse{
		Status: "FAILURE",
		Error:  err,
	}
	pkgresponse.WriteResponse(w, statusCode, response)

}
