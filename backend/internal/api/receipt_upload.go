// receipt_upload.go implements the two endpoints around a milestone's
// spending-receipt file (see db/schema.sql's milestone_receipts comment for
// why this exists — the on-chain receiptHash is only a fingerprint, never
// the file itself):
//   - POST /api/campaigns/{id}/milestones/{idx}/receipt — the charity
//     uploads the actual document after their submitReceipt transaction
//     confirms on-chain.
//   - GET  /api/campaigns/{id}/milestones/{idx}/receipt/file — serves that
//     document back, so a donor can open and read it.
//
// A fifth exception to server.go's "mostly read-only" doc comment, and like
// POST /api/chat, a different kind again: this one writes to local disk, not
// the database (the DB row is just a pointer to the file). No ownership
// check, same trust level as the API's other write endpoints — see
// server.go's package doc for why that's an accepted, documented risk here.
package api

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/zoujunhui1/trust-chain/backend/internal/store"
)

const maxReceiptUploadBytes = 8 << 20 // 8MB — receipts are photos/PDFs, not video

var allowedReceiptContentTypes = map[string]string{
	"image/png":       ".png",
	"image/jpeg":      ".jpg",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

// handleUploadMilestoneReceipt: POST /api/campaigns/{id}/milestones/{idx}/receipt
// — multipart/form-data with a "file" field (the document) and an optional
// "note" field (a short expense description). The hash is computed here,
// server-side, from the bytes actually received — never trusted from the
// client — so "this file's hash matches the on-chain receipt" is something
// anyone can verify, not something the backend merely asserts.
func (s *Server) handleUploadMilestoneReceipt(w http.ResponseWriter, r *http.Request) {
	campaignID, ok := parseID(w, r)
	if !ok {
		return
	}
	idx, ok := parseMilestoneIdx(w, r)
	if !ok {
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxReceiptUploadBytes+1<<20) // small slack for the rest of the multipart form
	if err := r.ParseMultipartForm(maxReceiptUploadBytes); err != nil {
		writeError(w, http.StatusBadRequest, "文件太大或表单格式不对 / file too large or malformed form")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, "缺少文件 / missing file field")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedReceiptContentTypes[contentType]
	if !ok {
		writeError(w, http.StatusBadRequest, "只支持 PNG/JPEG/WEBP/PDF / only PNG, JPEG, WEBP or PDF are accepted")
		return
	}

	note := strings.TrimSpace(r.FormValue("note"))

	if err := os.MkdirAll(s.receiptsDir(), 0o755); err != nil {
		writeError(w, http.StatusInternalServerError, "存储初始化失败 / storage init failed")
		return
	}

	// 用 campaignId+idx 做确定性文件名：重新上传直接替换旧文件，跟
	// SetMilestoneReceiptFile 的 ON DUPLICATE KEY UPDATE 语义一致。
	// A deterministic filename from campaignId+idx: re-uploading replaces the
	// old file outright, matching SetMilestoneReceiptFile's ON DUPLICATE KEY
	// UPDATE semantics.
	relPath := fmt.Sprintf("receipts/%d-%d%s", campaignID, idx, ext)
	absPath := filepath.Join(s.uploadsDir, relPath)

	out, err := os.Create(absPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "写入文件失败 / failed to write file")
		return
	}
	defer out.Close()

	hasher := sha256.New()
	size, err := io.Copy(out, io.TeeReader(file, hasher))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "写入文件失败 / failed to write file")
		return
	}

	sum := hex.EncodeToString(hasher.Sum(nil))
	receipt := store.MilestoneReceiptFile{
		FilePath:    relPath,
		FileName:    header.Filename,
		ContentType: contentType,
		FileSize:    size,
		SHA256:      sum,
	}
	if note != "" {
		receipt.Note = &note
	}
	if err := s.store.SetMilestoneReceiptFile(r.Context(), campaignID, idx, receipt); err != nil {
		writeError(w, http.StatusInternalServerError, "写入失败 / write failed")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"sha256":      sum,
		"fileName":    header.Filename,
		"contentType": contentType,
		"fileSize":    size,
	})
}

// handleGetMilestoneReceiptFile: GET /api/campaigns/{id}/milestones/{idx}/receipt/file
// — serves the uploaded document back so it can be viewed inline (an <img>
// or a PDF preview) or downloaded.
func (s *Server) handleGetMilestoneReceiptFile(w http.ResponseWriter, r *http.Request) {
	campaignID, ok := parseID(w, r)
	if !ok {
		return
	}
	idx, ok := parseMilestoneIdx(w, r)
	if !ok {
		return
	}

	f, found, err := s.store.GetMilestoneReceiptFile(r.Context(), campaignID, idx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "查询失败 / query failed")
		return
	}
	if !found {
		writeError(w, http.StatusNotFound, "还没有上传凭证 / no receipt uploaded yet")
		return
	}

	file, err := os.Open(filepath.Join(s.uploadsDir, f.FilePath))
	if err != nil {
		writeError(w, http.StatusNotFound, "文件不存在 / file not found on disk")
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "读取文件失败 / failed to stat file")
		return
	}

	w.Header().Set("Content-Type", f.ContentType)
	http.ServeContent(w, r, f.FileName, stat.ModTime(), file)
}

func (s *Server) receiptsDir() string {
	return filepath.Join(s.uploadsDir, "receipts")
}

// parseMilestoneIdx 从路径取 {idx} 并转成 uint32。
// parseMilestoneIdx reads {idx} from the path and parses it as uint32.
func parseMilestoneIdx(w http.ResponseWriter, r *http.Request) (uint32, bool) {
	idx, err := strconv.ParseUint(r.PathValue("idx"), 10, 32)
	if err != nil {
		writeError(w, http.StatusBadRequest, "idx 非法 / invalid milestone index")
		return 0, false
	}
	return uint32(idx), true
}
