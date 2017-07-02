package common

type ErrorJSONResponse struct {
	ErrorDescription string `json:"error_description"`
}

var (
	StatusOK                  = 200
	StatusBadRequestError     = 400
	StatusInternalServerError = 500

	ResponseInternalServerError = ErrorJSONResponse{"internalServerError"}
	ResponseBadRequest          = ErrorJSONResponse{"badRequestError"}
)
