package html

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"
)

type MyForm struct {
	Firstname string
	Age       int
	Height    float32
	Cool      bool
	Blob      []byte
	Form      *multipart.Form
}

func TestUnmarshall(t *testing.T) {
	var buf bytes.Buffer
	mp := multipart.NewWriter(&buf)
	must(mp.WriteField("Firstname", "Torben"))
	must(mp.WriteField("Age", "12"))
	must(mp.WriteField("Height", "1.73"))
	must(mp.WriteField("Cool", "true"))
	tmp, err := mp.CreateFormFile("Blob", "abc.pdf")
	must(err)
	tmp.Write([]byte("12345"))
	must(mp.Close())

	req, err := http.NewRequest(http.MethodPost, "/test", bytes.NewReader(buf.Bytes()))
	must(err)
	req.Header.Set("content-type", mp.FormDataContentType())

	if err := req.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		t.Fatal(err)
	}

	var form MyForm
	must(UnmarshallForm(&form, req))

	fmt.Printf("%+v\n", form)

	if form.Firstname != "Torben" {
		t.Fatal(form)
	}

	if form.Age != 12 {
		t.Fatal(form)
	}

	if form.Height != 1.73 {
		t.Fatal(form)
	}

	if form.Cool != true {
		t.Fatal(form)
	}

	if !reflect.DeepEqual(form.Blob, []byte("12345")) {
		t.Fatal(form)
	}

	if form.Form.Value["Firstname"][0] != "Torben" {
		t.Fatal(form)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
