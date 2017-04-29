package pgproto3

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

type FunctionCallResponse struct {
	Result []byte
}

func (*FunctionCallResponse) Backend() {}

func (dst *FunctionCallResponse) Decode(src []byte) error {
	if len(src) < 4 {
		return &invalidMessageFormatErr{messageType: "FunctionCallResponse"}
	}
	rp := 0
	resultSize := int(binary.BigEndian.Uint32(src[rp:]))
	rp += 4

	if resultSize == -1 {
		dst.Result = nil
		return nil
	}

	if len(src[rp:]) != resultSize {
		return &invalidMessageFormatErr{messageType: "FunctionCallResponse"}
	}

	dst.Result = src[rp:]
	return nil
}

func (src *FunctionCallResponse) MarshalBinary() ([]byte, error) {
	var bigEndian BigEndianBuf
	buf := &bytes.Buffer{}

	buf.WriteByte('V')
	buf.Write(bigEndian.Uint32(uint32(4 + 4 + len(src.Result))))

	if src.Result == nil {
		buf.Write(bigEndian.Int32(-1))
	} else {
		buf.Write(bigEndian.Int32(int32(len(src.Result))))
		buf.Write(src.Result)
	}

	return buf.Bytes(), nil
}

func (src *FunctionCallResponse) MarshalJSON() ([]byte, error) {
	var formattedValue map[string]string
	var hasNonPrintable bool
	for _, b := range src.Result {
		if b < 32 {
			hasNonPrintable = true
			break
		}
	}

	if hasNonPrintable {
		formattedValue = map[string]string{"binary": hex.EncodeToString(src.Result)}
	} else {
		formattedValue = map[string]string{"text": string(src.Result)}
	}

	return json.Marshal(struct {
		Type   string
		Result map[string]string
	}{
		Type:   "FunctionCallResponse",
		Result: formattedValue,
	})
}