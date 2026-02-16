package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DocumentStatus string

const (
	StatusPending  DocumentStatus = "PENDING"
	StatusApproved DocumentStatus = "APPROVED"
	StatusRejected DocumentStatus = "REJECTED"
)

type Document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DocNo     string             `bson:"doc_no" json:"doc_no"`
	Title     string             `bson:"title" json:"title"`
	Status    DocumentStatus     `bson:"status" json:"status"`
	Reason    string             `bson:"reason,omitempty" json:"reason,omitempty"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type DocumentItemDTO struct {
	ID        string         `json:"id" example:"6992a40c3fb0b6dcb1857948"`
	DocNo     string         `json:"doc_no" example:"IT03-0001"`
	Title     string         `json:"title" example:"Access Request"`
	Status    DocumentStatus `json:"status" example:"PENDING"`
	Reason    string         `json:"reason,omitempty" example:"Approved by IT manager"`
	CreatedAt string         `json:"created_at" example:"2026-02-16T04:58:52Z"`
	UpdatedAt string         `json:"updated_at" example:"2026-02-16T04:58:52Z"`
}

type DocumentListResponse struct {
	StatusCode int              `json:"statusCode" example:"200"`
	Data       []DocumentItemDTO `json:"data"`
}

type BulkDecisionRequest struct {
	DocumentIDs []string `json:"document_ids" example:"6992a066b7e55e55d9d1dd8f"`
	Reason      string   `json:"reason" example:"Approved by IT manager"`
}

type BulkApprovalResponse struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Message    string `json:"message" example:"approved"`
	Requested  int    `json:"requested" example:"2"`
	Approved   int64  `json:"approved" example:"2"`
}

type BulkRejectionResponse struct {
	StatusCode int    `json:"statusCode" example:"200"`
	Message    string `json:"message" example:"rejected"`
	Requested  int    `json:"requested" example:"2"`
	Rejected   int64  `json:"rejected" example:"2"`
}
