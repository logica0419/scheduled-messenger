package event

// JOINED, LEFT イベントリクエストスキーマ
type SystemEvent struct {
	EventTime string  `json:"eventTime,omitempty"`
	Channel   Channel `json:"channel,omitempty"`
}

// チャンネルスキーマ
type Channel struct {
	ID        string  `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Path      string  `json:"path,omitempty"`
	ParentID  string  `json:"parentId,omitempty"`
	Creator   Creator `json:"creator,omitempty"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
}

// チャンネルクリエイタースキーマ
type Creator struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	IconID      string `json:"iconId,omitempty"`
	Bot         bool   `json:"bot,omitempty"`
}

// チャンネルID 取得メソッド
func (r *SystemEvent) GetChannelID() string {
	return r.Channel.ID
}

// チャンネルパス 取得メソッド
func (r *SystemEvent) GetChannelPath() string {
	return r.Channel.Path
}
