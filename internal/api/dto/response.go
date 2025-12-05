package dto

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    PaginatedData `json:"data"`
}

type PaginatedData struct {
	Items []any       `json:"items"`
	Meta  *Pagination `json:"meta"`
}

type Pagination struct {
	Total  int64 `json:"total"`
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Count  int   `json:"count"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func Success(data any) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func SuccessWithMessage(message string, data any) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

func SuccessPaginated(data []any, total int64, page, pageSize int) PaginatedResponse {
	return PaginatedResponse{
		Code:    200,
		Message: "success",
		Data: PaginatedData{
			Items: data,
			Meta: &Pagination{
				Total:  total,
				Limit:  pageSize,
				Offset: (page - 1) * pageSize,
				Count:  len(data),
			},
		},
	}
}

func Error(code int, message string, err error) ErrorResponse {
	resp := ErrorResponse{
		Code:    code,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	return resp
}

func BadRequest(message string, err error) ErrorResponse {
	return Error(400, message, err)
}

func NotFound(message string) ErrorResponse {
	return Error(404, message, nil)
}

func InternalServerError(message string, err error) ErrorResponse {
	return Error(500, message, err)
}
