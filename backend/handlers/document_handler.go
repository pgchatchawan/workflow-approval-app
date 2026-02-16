package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"backend/db"
	"backend/models"
)

func colDocuments() *mongo.Collection {
	return db.Database.Collection("documents")
}

func parseObjectIDs(hexIDs []string) ([]primitive.ObjectID, error) {
	oids := make([]primitive.ObjectID, 0, len(hexIDs))
	for _, id := range hexIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		oids = append(oids, oid)
	}
	return oids, nil
}

func isValidStatus(s string) bool {
	return s == string(models.StatusPending) ||
		s == string(models.StatusApproved) ||
		s == string(models.StatusRejected)
}

func uniqueTrimmedStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	for _, v := range in {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}


// ListDocuments godoc
// @Summary      List documents
// @Description  Get documents (optional filter by status)
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        status  query  string  false  "Filter status: PENDING | APPROVED | REJECTED"
// @Success      200  {object}  models.DocumentListResponse
// @Failure      500  {object}  map[string]any
// @Router       /api/documents [get]
func ListDocuments(c *fiber.Ctx) error {
	status := c.Query("status", "")
	filter := bson.M{}

	if status != "" && !isValidStatus(status) {
		return c.Status(400).JSON(fiber.Map{
			"statusCode": 400,
			"message":    "invalid status (PENDING, APPROVED, REJECTED)",
		})
	}

	if status != "" {
		filter["status"] = status
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cur, err := colDocuments().Find(ctx, filter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"message":    err.Error(),
		})
	}
	defer cur.Close(ctx)

	var docs []models.Document
	if err := cur.All(ctx, &docs); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"message":    err.Error(),
		})
	}

	items := make([]models.DocumentItemDTO, 0, len(docs))
	for _, d := range docs {
		items = append(items, models.DocumentItemDTO{
			ID:        d.ID.Hex(),
			DocNo:     d.DocNo,
			Title:     d.Title,
			Status:    d.Status,
			Reason:    d.Reason,
			CreatedAt: d.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt: d.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(models.DocumentListResponse{
		StatusCode: 200,
		Data:       items,
	})
}


// BulkApproveDocuments godoc
// @Summary      Approve documents (bulk)
// @Description  Approve only documents with status PENDING (cannot approve twice)
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        body  body      models.BulkDecisionRequest  true  "Approval payload"
// @Success      200   {object}  models.BulkApprovalResponse
// @Failure      400   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /api/documents/approval [post]
func ApproveDocuments(c *fiber.Ctx) error {
	var req models.BulkDecisionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "invalid body"})
	}

	req.DocumentIDs = uniqueTrimmedStrings(req.DocumentIDs)
	req.Reason = strings.TrimSpace(req.Reason)

	if len(req.DocumentIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "document_ids is required"})
	}
	if req.Reason == "" {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason is required"})
	}

	if len(req.Reason) < 3 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason must be at least 3 characters"})
	}
	if len(req.Reason) > 500 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason must be <= 500 characters"})
	}

	oids, err := parseObjectIDs(req.DocumentIDs)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "invalid document_ids"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":    bson.M{"$in": oids},
		"status": models.StatusPending,
	}
	update := bson.M{"$set": bson.M{
		"status":     models.StatusApproved,
		"reason":     req.Reason,
		"updated_at": time.Now(),
	}}

	res, err := colDocuments().UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}

	return c.JSON(models.BulkApprovalResponse{
		StatusCode: 200,
		Message:    "approved",
		Requested:  len(req.DocumentIDs),
		Approved:   res.ModifiedCount,
	})
}

// BulkRejectDocuments godoc
// @Summary      Reject documents (bulk)
// @Description  Reject only documents with status PENDING
// @Tags         Documents
// @Accept       json
// @Produce      json
// @Param        body  body      models.BulkDecisionRequest  true  "Rejection payload"
// @Success      200   {object}  models.BulkRejectionResponse
// @Failure      400   {object}  map[string]any
// @Failure      500   {object}  map[string]any
// @Router       /api/documents/rejection [post]
func RejectDocuments(c *fiber.Ctx) error {
	var req models.BulkDecisionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "invalid body"})
	}

	req.DocumentIDs = uniqueTrimmedStrings(req.DocumentIDs)
	req.Reason = strings.TrimSpace(req.Reason)

	if len(req.DocumentIDs) == 0 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "document_ids is required"})
	}
	if req.Reason == "" {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason is required"})
	}

	if len(req.Reason) < 3 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason must be at least 3 characters"})
	}
	if len(req.Reason) > 500 {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "reason must be <= 500 characters"})
	}

	oids, err := parseObjectIDs(req.DocumentIDs)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"statusCode": 400, "message": "invalid document_ids"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":    bson.M{"$in": oids},
		"status": models.StatusPending,
	}
	update := bson.M{"$set": bson.M{
		"status":     models.StatusRejected,
		"reason":     req.Reason,
		"updated_at": time.Now(),
	}}

	res, err := colDocuments().UpdateMany(ctx, filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"statusCode": 500, "message": err.Error()})
	}

	return c.JSON(models.BulkRejectionResponse{
		StatusCode: 200,
		Message:    "rejected",
		Requested:  len(req.DocumentIDs),
		Rejected:   res.ModifiedCount,
	})
}

// SeedMockDocuments godoc
// @Summary      Seed mock documents (dev only)
// @Description  Insert mock documents for testing (use ?reset=true to clear old data)
// @Tags         Documents
// @Produce      json
// @Param        reset  query     string  false  "Reset collection before seed (true/false)"
// @Success      200    {object}  map[string]any
// @Failure      500    {object}  map[string]any
// @Router       /api/documents/seed [post]
func SeedMockDocuments(c *fiber.Ctx) error {
	var mock []any
	now := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// จะลบข้อมูลเดิมก่อน
	if c.Query("reset") == "true" {
		_, err := colDocuments().DeleteMany(ctx, bson.M{})
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"statusCode": 500,
				"message":    "failed to reset collection",
			})
		}
	}

	// ===== PENDING 15 =====
	for i := 1; i <= 15; i++ {
		mock = append(mock, models.Document{
			DocNo:     fmt.Sprintf("IT03-%04d", i),
			Title:     fmt.Sprintf("IT Equipment Request %d", i),
			Status:    models.StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// ===== APPROVED 5 =====
	for i := 16; i <= 20; i++ {
		mock = append(mock, models.Document{
			DocNo:     fmt.Sprintf("IT03-%04d", i),
			Title:     fmt.Sprintf("Approved Request %d", i),
			Status:    models.StatusApproved,
			Reason:    "Approved by IT manager",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	// ===== REJECTED 5 =====
	for i := 21; i <= 25; i++ {
		mock = append(mock, models.Document{
			DocNo:     fmt.Sprintf("IT03-%04d", i),
			Title:     fmt.Sprintf("Rejected Request %d", i),
			Status:    models.StatusRejected,
			Reason:    "Not compliant with IT policy",
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	res, err := colDocuments().InsertMany(ctx, mock)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"statusCode": 500,
			"message":    err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"statusCode": 200,
		"message":    "mock data seeded successfully",
		"inserted":   len(res.InsertedIDs),
	})
}

