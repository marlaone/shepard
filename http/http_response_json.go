package http

import (
	"encoding/json"

	"github.com/marlaone/shepard/collections/slice"
)

func JsonResponse(body any) Response[Body] {

	builder := NewResponseBuilder(NewHttpResponseBytes())
	potentialRes := builder.Header("Content-Type", "application/json").Status(StatusCodeOk).Body(NewBytesBody())
	if potentialRes.IsErr() {
		errRes := NewHttpResponseBytes()
		errRes.SetStatusCode(StatusCodeInternalServerError)
		errRes.SetBody(NewBytesBody())
		errRes.Body().Write(slice.Init[byte]([]byte(potentialRes.Err().Unwrap().Error())...))
		errRes.Finish()
		return errRes
	}
	res := potentialRes.Unwrap()

	bytes, err := json.Marshal(body)
	if err != nil {
		res.SetStatusCode(StatusCodeInternalServerError)
		res.Body().Write(slice.Init[byte]([]byte(err.Error())...))
		res.Finish()
		return res
	}

	res.SetStatusCode(StatusCodeOk)
	res.Body().Write(slice.Init[byte](bytes...))
	res.Finish()

	return res
}
