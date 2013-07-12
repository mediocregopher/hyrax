package types

import (
    "fmt"
    "encoding/binary"
)

// ConnId is a unique value that's given to every connection on this hyrax node
type ConnId uint64
func (cid *ConnId) Serialize() []byte {
    size := binary.Size(*cid)
    buf := make([]byte,size)
    binary.PutUvarint(buf,uint64(*cid))
    return buf
}

func ConnIdDeserialize(s []byte) (ConnId,error) {
    ui,br := binary.Uvarint(s)
    if br <= 0 { return 0,fmt.Errorf("Error deserializing %02X\n",s) }
    return ConnId(ui),nil
}
