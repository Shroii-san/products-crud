package handler

import (
	"crud-lengkap/model"
	"crud-lengkap/repository"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductHandler struct {
	Repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{Repo: repo}
}

func sendJSON(w http.ResponseWriter, status int, response model.Response) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) HandleProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		idStr := r.URL.Query().Get("id")
		if idStr != "" {
			id, err := strconv.Atoi(idStr)
			if err != nil || id <= 0 {
				sendJSON(w, http.StatusBadRequest, model.Response{
					Status:  "error",
					Message: "Query parameter 'id' tidak valid",
				})
				return
			}

			product, err := h.Repo.GetProductById(id)
			if err != nil {
				sendJSON(w, http.StatusInternalServerError, model.Response{
					Status:  "error",
					Message: err.Error(),
				})
				return
			}
			sendJSON(w, http.StatusOK, model.Response{
				Status:  "success",
				Message: "Berhasil mengambil produk",
				Data:    product,
			})
			return
		}

		products, err := h.Repo.GetAllProduct()
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, model.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		sendJSON(w, http.StatusOK, model.Response{
			Status:  "success",
			Message: "Berhasil mengambil daftar produk",
			Data:    products,
		})

	case http.MethodPost:
		var p model.Product

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			sendJSON(w, http.StatusBadRequest, model.Response{
				Status:  "error",
				Message: "Payload JSON tidak valid",
			})
			return
		}

		if p.Name == "" {
			sendJSON(w, http.StatusBadRequest, model.Response{
				Status:  "error",
				Message: "Nama produk tidak boleh kosong",
			})
			return
		}
		if p.Price <= 0 {
			sendJSON(w, http.StatusBadRequest, model.Response{
				Status:  "error",
				Message: "Harga produk harus lebih dari 0",
			})
			return
		}

		newProduct, err := h.Repo.AddProduct(p)
		if err != nil {
			sendJSON(w, http.StatusInternalServerError, model.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		sendJSON(w, http.StatusCreated, model.Response{
			Status:  "success",
			Message: "Produk berhasil ditambahkan",
			Data:    newProduct,
		})

	case http.MethodPut:
		var p model.Product
		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil || p.ID <= 0 {
			sendJSON(w, http.StatusBadRequest, model.Response{
				Status:  "error",
				Message: "Data produk atau ID tidak valid",
			})
			return
		}

		err = h.Repo.UpdateProduct(p)
		if err != nil {
			sendJSON(w, http.StatusNotFound, model.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		sendJSON(w, http.StatusOK, model.Response{
			Status:  "success",
			Message: "Produk berhasil diperbarui",
			Data:    p,
		})

	case http.MethodDelete:
		idStr := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			sendJSON(w, http.StatusBadRequest, model.Response{
				Status:  "error",
				Message: "Query parameter 'id' tidak valid",
			})
			return
		}

		err = h.Repo.DeleteProduct(id)
		if err != nil {
			sendJSON(w, http.StatusNotFound, model.Response{
				Status:  "error",
				Message: err.Error(),
			})
			return
		}

		sendJSON(w, http.StatusOK, model.Response{
			Status:  "success",
			Message: "Produk berhasil dihapus",
		})

	default:
		sendJSON(w, http.StatusMethodNotAllowed, model.Response{
			Status:  "error",
			Message: "Method tidak diizinkan",
		})
	}
	return
}
