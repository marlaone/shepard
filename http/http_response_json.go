package http

import (
	"encoding/json"
)

func JsonResponse(body any) Response[Body] {

	builder := NewResponseBuilder(NewHttpResponseBytes())
	potentialRes := builder.Header("Content-Type", "application/json").Status(StatusCodeOk).Body(NewBytesBody())
	if potentialRes.IsErr() {
		errRes := NewHttpResponseBytes()
		errRes.SetStatusCode(StatusCodeInternalServerError)
		errRes.SetBody(NewBytesBody())
		errRes.Body().Write() <- []byte(potentialRes.Err().Unwrap().Error())
		errRes.Body().Finish()
		return errRes
	}
	res := potentialRes.Unwrap()

	res.SetHeader("Content-Type", "application/json")

	bytes, err := json.Marshal(body)
	if err != nil {
		res.SetStatusCode(StatusCodeInternalServerError)
		res.Body().Write() <- []byte(err.Error())
		return res
	}

	res.SetStatusCode(StatusCodeOk)
	res.Body().Write() <- bytes
	res.Body().Finish()

	return res
}
