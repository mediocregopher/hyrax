package custom

import (
	"bytes"
	"errors"
	. "hyrax-server/storage"
	stypes "hyrax-server/types"
	types "hyrax/types"
)

// EAdd adds the connection's id (and name) to an ekg's set of things it's
// watching, and adds the ekg's information to the connection's set of ekgs its
// hooked up to
func EAdd(cid stypes.ConnId, pay *types.Payload) (interface{}, error) {
	connekgkey := ConnEkgKey(cid)
	connekgval := ConnEkgVal(pay.Domain, pay.Id, pay.Name)
	_, err := CmdPretty(SADD, connekgkey, connekgval)
	if err != nil {
		return nil, err
	}

	ekgkey := EkgKey(pay.Domain, pay.Id)
	ekgval := EkgVal(cid, pay.Name)
	_, err = CmdPretty(SADD, ekgkey, ekgval)
	return OK, err
}

// ERem removes the connection's id (and name) from an ekg's set of things it's
// watching, and removes the ekg's information from the connection's set of ekgs
// its hooked up to
func ERem(cid stypes.ConnId, pay *types.Payload) (interface{}, error) {
	ekgkey := EkgKey(pay.Domain, pay.Id)
	ekgval := EkgVal(cid, pay.Name)
	_, err := CmdPretty(SREM, ekgkey, ekgval)
	if err != nil {
		return nil, err
	}

	connekgkey := ConnEkgKey(cid)
	connekgval := ConnEkgVal(pay.Domain, pay.Id, pay.Name)
	_, err = CmdPretty(SREM, connekgkey, connekgval)
	return OK, err
}

// CleanConnEkg takes in a connection id and cleans up all of its ekgs, and the
// set which keeps track of those ekgs. It also sends out alerts for all the
// ekgs it's hooked up to, since this only gets called on a disconnect.
func CleanConnEkg(cid stypes.ConnId) error {
	connekgkey := ConnEkgKey(cid)
	r, err := CmdPretty(SMEMBERS, connekgkey)
	if err != nil {
		return err
	}

	ekgs := r.([][]byte)

	for i := range ekgs {
		domain, id, name := DeconstructConnEkgVal(ekgs[i])
		ekgkey := EkgKey(domain, id)
		ekgval := EkgVal(cid, name)
		_, err = CmdPretty(SREM, ekgkey, ekgval)
		if err != nil {
			return err
		}

		cmd := types.Command{
			Command: DISCONNECT,
			Payload: types.Payload{
				Domain: domain,
				Id:     id,
				Name:   name,
			},
		}
		MonMakeAlert(&cmd)
	}

	_, err = CmdPretty(DEL, connekgkey)
	return err

}

// EMembers returns the list of names being monitored by an ekg
func EMembers(cid stypes.ConnId, pay *types.Payload) (interface{}, error) {
	ekgkey := EkgKey(pay.Domain, pay.Id)
	r, err := CmdPretty(SMEMBERS, ekgkey)
	if err != nil {
		return nil, err
	}

	members := r.([][]byte)
	for i := range members {
		_, name := DeconstructEkgVal(members[i])
		members[i] = name
	}

	return members, nil
}

// ECard returns the number of connection/name combinations being monitored
func ECard(cid stypes.ConnId, pay *types.Payload) (interface{}, error) {
	ekgkey := EkgKey(pay.Domain, pay.Id)
	return CmdPretty(SCARD, ekgkey)
}

// EIsMember returns whether or not the given name is being monitored by the ekg
func EIsMember(cid stypes.ConnId, pay *types.Payload) (interface{}, error) {

	if !(len(pay.Values) > 0) {
		return nil, errors.New("ERR wrong number of arguments for 'eismember' command")
	}

	ekgkey := EkgKey(pay.Domain, pay.Id)
	r, err := CmdPretty(SMEMBERS, ekgkey)
	if err != nil {
		return nil, err
	}

	members := r.([][]byte)
	for i := range members {
		_, name := DeconstructEkgVal(members[i])
		if bytes.Equal(name, pay.Values[0]) {
			return 1, nil
		}
	}

	return 0, nil
}
