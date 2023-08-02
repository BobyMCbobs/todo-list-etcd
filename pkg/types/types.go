package types

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type List struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	AlertOn               time.Time `json:"alertOn"`
	AuthorID              string    `json:"authorID"`
	Items                 []*Item   `json:"items"`
	CreationTimestamp     string    `json:"creationTimestamp"`
	ModificationTimestamp string    `json:"modificationTimestamp"`
	DeletionTimestamp     string    `json:"deletionTimestamp"`
	Revision              int64     `json:"revision"`
}

type Item struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	AlertOn               time.Time `json:"alertOn"`
	AuthorID              string    `json:"authorID"`
	ListID                string    `json:"listID"`
	Complete              string    `json:"complete"`
	CreationTimestamp     string    `json:"creationTimestamp"`
	ModificationTimestamp string    `json:"modificationTimestamp"`
	DeletionTimestamp     string    `json:"deletionTimestamp"`
	Revision              int64     `json:"revision"`
}

type JWTclaim struct {
	ID         string   `json:"id"`
	InstanceID string   `json:"instanceID"`
	Scopes     []string `json:"scopes"`
	Type       string   `json:"type"`
	AuthNonce  string   `json:"authNonce"`
	jwt.StandardClaims
}
