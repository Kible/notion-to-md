package requests

type GetBlockChildrenRequest struct {
	BlockID     string `json:"block_id"`
	PageSize    *int   `json:"page_size"`
	StartCursor string `json:"start_cursor"`
}
