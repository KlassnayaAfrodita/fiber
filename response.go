package fiber

import (
	"bytes"
	"net/http"

	"github.com/gofiber/fiber/v2"

	httpcontract "github.com/goravel/framework/contracts/http"
)

type Response struct {
	instance *fiber.Ctx
	origin   httpcontract.ResponseOrigin
}

func NewResponse(instance *fiber.Ctx, origin httpcontract.ResponseOrigin) *Response {
	return &Response{instance, origin}
}

func (r *Response) Data(code int, contentType string, data []byte) {
	_ = r.instance.Status(code).Type(contentType).Send(data)
}

func (r *Response) Download(filepath, filename string) {
	_ = r.instance.Download(filepath, filename)
}

func (r *Response) File(filepath string) {
	_ = r.instance.SendFile(filepath)
}

func (r *Response) Header(key, value string) httpcontract.Response {
	r.instance.Set(key, value)

	return r
}

func (r *Response) Json(code int, obj any) {
	_ = r.instance.Status(code).JSON(obj)
}

func (r *Response) Origin() httpcontract.ResponseOrigin {
	return r.origin
}

func (r *Response) Redirect(code int, location string) {
	_ = r.instance.Redirect(location, code)
}

func (r *Response) String(code int, format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(code).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(code).Type(format).SendString(values[0].(string))
}

func (r *Response) Success() httpcontract.ResponseSuccess {
	return NewFiberSuccess(r.instance)
}

func (r *Response) Status(code int) httpcontract.ResponseStatus {
	return NewFiberStatus(r.instance, code)
}

func (r *Response) Writer() http.ResponseWriter {
	// Fiber doesn't support this
	return nil
}

func (r *Response) Flush() {
	r.instance.Fresh()
}

type FiberSuccess struct {
	instance *fiber.Ctx
}

func NewFiberSuccess(instance *fiber.Ctx) httpcontract.ResponseSuccess {
	return &FiberSuccess{instance}
}

func (r *FiberSuccess) Data(contentType string, data []byte) {
	_ = r.instance.Type(contentType).Send(data)
}

func (r *FiberSuccess) Json(obj any) {
	_ = r.instance.Status(http.StatusOK).JSON(obj)
}

func (r *FiberSuccess) String(format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(http.StatusOK).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(http.StatusOK).Type(format).SendString(values[0].(string))
}

type FiberStatus struct {
	instance *fiber.Ctx
	status   int
}

func NewFiberStatus(instance *fiber.Ctx, code int) httpcontract.ResponseSuccess {
	return &FiberStatus{instance, code}
}

func (r *FiberStatus) Data(contentType string, data []byte) {
	_ = r.instance.Status(r.status).Type(contentType).Send(data)
}

func (r *FiberStatus) Json(obj any) {
	_ = r.instance.Status(r.status).JSON(obj)
}

func (r *FiberStatus) String(format string, values ...any) {
	if len(values) == 0 {
		_ = r.instance.Status(http.StatusOK).Type(format).SendString(format)
		return
	}

	_ = r.instance.Status(http.StatusOK).Type(format).SendString(values[0].(string))
}

func ResponseMiddleware() httpcontract.Middleware {
	return func(ctx httpcontract.Context) {
		switch ctx := ctx.(type) {
		case *Context:
			ctx.Instance().Response().ResetBody()
		}

		ctx.Request().Next()
	}
}

type ResponseOrigin struct {
	*fiber.Ctx
}

func (w *ResponseOrigin) Body() *bytes.Buffer {
	return bytes.NewBuffer(w.Ctx.Response().Body())
}

func (w *ResponseOrigin) Header() http.Header {
	result := http.Header{}
	w.Ctx.Response().Header.VisitAll(func(key, value []byte) {
		result.Add(string(key), string(value))
	})

	return result
}

func (w *ResponseOrigin) Size() int {
	return len(w.Ctx.Response().Body())
}

func (w *ResponseOrigin) Status() int {
	return w.Ctx.Response().StatusCode()
}
