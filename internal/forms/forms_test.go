package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	field := "a"

	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)
	hasField := form.Has(field, r)
	if hasField {
		t.Error("form got the field when the field missing")
	}

	postedData := url.Values{}
	postedData.Add(field, "a")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.Form = postedData
	form = New(r.Form)
	hasField = form.Has(field, r)
	if !hasField {
		t.Error("does not have the field when it does")
	}
}

func TestForm_MinLength(t *testing.T) {
	field := "a"
	length := 3

	values := url.Values{}
	values.Add(field, "aa")
	form := New(values)
	hasLengthOverMin := form.MinLength(field, length)
	if hasLengthOverMin {
		t.Error("form got a longer field value than min when should have got a less")
	}

	isError := form.Errors.Get(field)
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	values = url.Values{}
	values.Add(field, "aaa")
	form = New(values)
	hasLengthOverMin = form.MinLength(field, length)
	if !hasLengthOverMin {
		t.Error("field value is not longer than min when it is")
	}

	isError = form.Errors.Get(field)
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	field := "email"
	
	values := url.Values{}
	values.Add(field, "invalid email address")
	form := New(values)
	form.IsEmail(field)

	if form.Valid() {
		t.Error("email was valid when should have been invalid")
	}

	values = url.Values{}
	values.Add(field, "valid@email.com")
	form = New(values)
	form.IsEmail(field)

	if !form.Valid() {
		t.Error("email was invalid when should have been valid")
	}


}