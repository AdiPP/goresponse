package goresponse

import (
	"encoding/json"
	"github.com/lvqingan/gopager"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlusher_CommonResponse(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	statusCode := 200
	response := NewCommon(&CommonOpt{
		Message:    "common response success",
		Data:       "",
		Error:      nil,
		StatusCode: statusCode,
	})

	flusher := NewFlusher(w, req)
	flusher.JSON(response)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.Errorf("got error %v", err)
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("got error %v", err)
	}

	var jsonCommon JSONCommon
	err = json.Unmarshal(resBody, &jsonCommon)

	if jsonCommon.Message != response.Message() {
		t.Errorf("expected %v got %v", response.Message(), jsonCommon.Message)
	}

	if res.StatusCode != statusCode {
		t.Errorf("expected %v got %v", statusCode, res.Status)
	}

	if jsonCommon.Status != true {
		t.Errorf("expected %v got %v", true, jsonCommon.Status)
	}
}

func TestFlusher_PaginationFromPaginator(t *testing.T) {
	w := httptest.NewRecorder()

	target := "/resources"
	req := httptest.NewRequest(http.MethodGet, target, nil)

	statusCode := 200
	result := make([]int, 100)
	perPage := 5
	currentPage := 1
	paginator := gopager.NewLengthAwarePaginator(
		result[currentPage-1:perPage*currentPage],
		len(result),
		perPage,
		currentPage,
		map[string]string{
			"path": target,
		},
	)

	response := NewPaginationFromLengthAwarePaginator(&PaginationFromLengthAwarePaginatorOpt{
		Message:    "pagination response success",
		Paginator:  paginator,
		Error:      nil,
		StatusCode: statusCode,
	})

	flusher := NewFlusher(w, req)
	flusher.JSON(response)

	res := w.Result()
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.Errorf("got error %v", err)
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		t.Errorf("got error %v", err)
	}

	var jsonPagination JSONPagination
	err = json.Unmarshal(resBody, &jsonPagination)

	if jsonPagination.Message != response.Message() {
		t.Errorf("expecyed %v got %v", response.Message(), jsonPagination.Message)
	}

	if res.StatusCode != statusCode {
		t.Errorf("expected %v got %v", statusCode, res.Status)
	}

	if jsonPagination.Status != true {
		t.Errorf("expected %v got %v", true, jsonPagination.Status)
	}
}
