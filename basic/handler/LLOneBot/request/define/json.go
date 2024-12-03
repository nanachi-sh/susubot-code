package define

type MessageChain struct {
	Data map[string]any `json:"data"`
	Type string         `json:"type"`
}

type (
	Request struct {
		Action string          `json:"action"`
		Params *Request_Params `json:"params,omitempty"`
		Echo   string          `json:"echo"`
	}
	Request_Params struct {
		GroupId   string           `json:"group_id,omitempty"`
		UserId    string           `json:"user_id,omitempty"`
		MessageId string           `json:"message_id,omitempty"`
		Message   []map[string]any `json:"message,omitempty"`
	}
)
