package param

type GetPresenceRequest struct {
	UserID []uint
}

type GetPresenceResponse struct {
	Items []GetPresenceItem
}

type GetPresenceItem struct {
	UserID    uint
	Timestamp int64
}
