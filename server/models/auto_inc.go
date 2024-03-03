package models

import "sync"

type autoInc struct {
	sync.Mutex
	id int
}

var autoIncRoomId = autoInc{
	id: 0,
}
