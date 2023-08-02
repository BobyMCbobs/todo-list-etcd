package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/BobyMCbobs/todo-list-etcd/pkg/types"
	"github.com/gorilla/mux"
)

// Logging ...
// log the HTTP requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v %v %v %v", r.Method, r.URL, r.Proto, r.Response, r.Header)
		next.ServeHTTP(w, r)
	})
}

func (h *HTTPServer) apiListLists() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		ctx := context.TODO()
		lists, err := h.todolistManager.Lists().List(ctx)
		if err != nil {
			log.Printf("error getting lists: %v\n", err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		listsBytes, err := json.Marshal(lists)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(listsBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiGetList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		ctx := context.TODO()
		list, err := h.todolistManager.Lists().Get(ctx, id)
		if err != nil {
			log.Printf("error getting list by id '%v': %v\n", id, err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if list.ID == "" {
			log.Println("error: list ID not found")
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		listBytes, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(listBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiPostOrPutList(method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }

		vars := mux.Vars(r)
		var id string
		if method == http.MethodPut {
			existingID, ok := vars["id"]
			if !ok {
				log.Println("Failed to find name param")
				response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
				return
			}
			existingList, err := h.todolistManager.Lists().Get(context.TODO(), existingID)
			if err != nil || existingList.ID == "" {
				fmt.Println(err)
				response(w, "NOT_FOUND", http.StatusNotFound)
				return
			}
			fmt.Println("found list with id", id)
			id = existingID
		}
		fmt.Println("list id", id)
		var list *types.List
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &list); err != nil {
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		if method == http.MethodPut {
			list.ID = id
		}
		listCreated, err := h.todolistManager.Lists().Put(context.TODO(), list)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		listBytes, err := json.Marshal(listCreated)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(listBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiPutList() http.HandlerFunc {
	return h.apiPostOrPutList(http.MethodPut)
}

func (h *HTTPServer) apiPostList() http.HandlerFunc {
	return h.apiPostOrPutList(http.MethodPost)
}

func (h *HTTPServer) apiDeleteList() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		list, err := h.todolistManager.Lists().Get(context.TODO(), id)
		if err != nil {
			log.Printf("error getting list by id '%v': %v\n", id, err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if list.ID == "" {
			log.Println("error: list ID not found")
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if err := h.todolistManager.Lists().Delete(context.TODO(), id); err != nil {
			log.Printf("error getting lists: %v\n", err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, "OK", http.StatusOK)
	})
}

func (h *HTTPServer) apiListItems() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		listid, ok := vars["listid"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		ctx := context.TODO()
		lists, err := h.todolistManager.Items(listid).List(ctx)
		if err != nil {
			log.Printf("error getting lists: %v\n", err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		listsBytes, err := json.Marshal(lists)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(listsBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiGetItem() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		listid, ok := vars["listid"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		ctx := context.TODO()
		list, err := h.todolistManager.Items(listid).Get(ctx, id)
		if err != nil {
			log.Printf("error getting list by id '%v': %v\n", id, err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if list.ID == "" {
			log.Println("error: list ID not found")
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		listBytes, err := json.Marshal(list)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(listBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiPostOrPutItem(method string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		var id string
		vars := mux.Vars(r)
		listid, ok := vars["listid"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		if method == http.MethodPut {
			existingID, ok := vars["id"]
			if !ok {
				log.Println("Failed to find name param")
				response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
				return
			}
			existingList, err := h.todolistManager.Items(listid).Get(context.TODO(), existingID)
			if err != nil || existingList.ID == "" {
				response(w, "NOT_FOUND", http.StatusNotFound)
				return
			}
			id = existingID
		}
		var item *types.Item
		body, err := io.ReadAll(r.Body)
		if err != nil {
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &item); err != nil {
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		item.ListID = listid
		if method == http.MethodPut {
			item.ID = id
		}
		itemCreated, err := h.todolistManager.Items(listid).Put(context.TODO(), item)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		itemBytes, err := json.Marshal(itemCreated)
		if err != nil {
			log.Println(err)
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, string(itemBytes), http.StatusOK)
	})
}

func (h *HTTPServer) apiPutItem() http.HandlerFunc {
	return h.apiPostOrPutItem(http.MethodPut)
}

func (h *HTTPServer) apiPostItem() http.HandlerFunc {
	return h.apiPostOrPutItem(http.MethodPost)
}

func (h *HTTPServer) apiDeleteItem() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		listid, ok := vars["listid"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		list, err := h.todolistManager.Items(listid).Get(context.TODO(), id)
		if err != nil {
			log.Printf("error getting list by id '%v': %v\n", id, err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if list.ID == "" {
			log.Println("error: list ID not found")
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		if err := h.todolistManager.Items(listid).Delete(context.TODO(), id); err != nil {
			log.Printf("error getting lists: %v\n", err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, "OK", http.StatusOK)
	})
}

func (h *HTTPServer) apiDeleteItemAll() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// claims, ok := r.Context().Value("jwt").(*types.JWTclaim)
		// if !ok {
		// 	log.Println("error: failed to get JWT claim from request")
		// 	response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
		// 	return
		// }
		vars := mux.Vars(r)
		listid, ok := vars["listid"]
		if !ok {
			log.Println("Failed to find name param")
			response(w, "INTERNAL_ERROR", http.StatusInternalServerError)
			return
		}
		if err := h.todolistManager.Items(listid).DeleteAll(context.TODO()); err != nil {
			log.Printf("error getting lists: %v\n", err)
			response(w, "NOT_FOUND", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		response(w, "OK", http.StatusOK)
	})
}
