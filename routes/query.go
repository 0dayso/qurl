// Copyright 2017 The qurl Authors. All rights reserved.

// Package routes implements all the HTTP entry points for this microservice.
package routes

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo"
	"github.com/repejota/qurl"
)

// Query fetch an URL and returns JSON with the data obtained.
func Query(c echo.Context) error {
	queryParams := c.QueryParams()
	u := queryParams.Get("url")

	result := qurl.NewResponse()
	result.URL = u
	result.Status = http.StatusOK

	// Validate URL
	_, err := url.ParseRequestURI(u)
	if err != nil {
		result.Status = http.StatusBadRequest
		return c.JSON(result.Status, result)
	}

	// Fetch URL content
	response, err := http.Get(u)
	if err != nil {
		result.Status = http.StatusInternalServerError
		return c.JSON(result.Status, result)
	}
	defer response.Body.Close()

	// Process headers
	for _, v := range queryParams["header"] {
		result.Headers[v] = response.Header[v]
	}

	// Process selectors
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		result.Status = http.StatusInternalServerError
		return c.JSON(result.Status, result)
	}
	for _, v := range queryParams["selector"] {
		result.Selectors[v] = append(result.Selectors[v], doc.Find(v).Text())
	}

	return c.JSON(result.Status, result)
}