package types

import (
    "strconv"
)

// ConnId is a unique value that's given to every connection on this hyrax node
type ConnId uint64
func (cid *ConnId) Serialize() string {
    return strconv.Itoa(int(*cid))
}

func ConnIdDeserialize(s string) (ConnId,error) {
    i,err := strconv.Atoi(s)
    if err != nil { return 0,err }
    return ConnId(i),nil
}

