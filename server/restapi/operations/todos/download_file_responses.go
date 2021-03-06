// Code generated by go-swagger; DO NOT EDIT.

package todos

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"io"
	"net/http"

	"github.com/go-openapi/runtime"

	"flightAPI/server/models"
)

// DownloadFileOKCode is the HTTP code returned for type DownloadFileOK
const DownloadFileOKCode int = 200

/*DownloadFileOK Downloading

swagger:response downloadFileOK
*/
type DownloadFileOK struct {

	/*
	  In: Body
	*/
	Payload io.ReadCloser `json:"body,omitempty"`
}

// NewDownloadFileOK creates DownloadFileOK with default headers values
func NewDownloadFileOK() *DownloadFileOK {

	return &DownloadFileOK{}
}

// WithPayload adds the payload to the download file o k response
func (o *DownloadFileOK) WithPayload(payload io.ReadCloser) *DownloadFileOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the download file o k response
func (o *DownloadFileOK) SetPayload(payload io.ReadCloser) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DownloadFileOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	payload := o.Payload
	if err := producer.Produce(rw, payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}
}

/*DownloadFileDefault error

swagger:response downloadFileDefault
*/
type DownloadFileDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.Error `json:"body,omitempty"`
}

// NewDownloadFileDefault creates DownloadFileDefault with default headers values
func NewDownloadFileDefault(code int) *DownloadFileDefault {
	if code <= 0 {
		code = 500
	}

	return &DownloadFileDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the download file default response
func (o *DownloadFileDefault) WithStatusCode(code int) *DownloadFileDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the download file default response
func (o *DownloadFileDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the download file default response
func (o *DownloadFileDefault) WithPayload(payload *models.Error) *DownloadFileDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the download file default response
func (o *DownloadFileDefault) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *DownloadFileDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
