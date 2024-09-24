package model

type SessionData struct {
	UserId        string `json:"user_id" redis:"user_id"`
	Authenticated bool   `json:"authenticated" redis:"authenticated"`
}

// UnmarshalBinary implement encoding.BinaryUnmarshaler interface.
// func (sd *SessionData) UnmarshalBinary(b []byte) error {
// 	return json.Unmarshal(b, sd)
// }
