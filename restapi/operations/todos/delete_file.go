// Code generated by go-swagger; DO NOT EDIT.

package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	"github.com/go-openapi/runtime/middleware"
)

// DeleteFileHandlerFunc turns a function with the right signature into a delete file handler
type DeleteFileHandlerFunc func(DeleteFileParams) middleware.Responder

// Handle executing the request and returning a response
func (fn DeleteFileHandlerFunc) Handle(params DeleteFileParams) middleware.Responder {
	return fn(params)
}

// DeleteFileHandler interface for that can handle valid delete file params
type DeleteFileHandler interface {
	Handle(DeleteFileParams) middleware.Responder
}

// NewDeleteFile creates a new http.Handler for the delete file operation
func NewDeleteFile(ctx *middleware.Context, handler DeleteFileHandler) *DeleteFile {
	return &DeleteFile{Context: ctx, Handler: handler}
}

/* DeleteFile swagger:route DELETE /file/{uuid} todos deleteFile

Delete a file.

*/
type DeleteFile struct {
	Context *middleware.Context
	Handler DeleteFileHandler
}

func (o *DeleteFile) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, rCtx, _ := o.Context.RouteInfo(r)
	if rCtx != nil {
		r = rCtx
	}
	var Params = NewDeleteFileParams()
	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request
	o.Context.Respond(rw, r, route.Produces, route, res)

}
