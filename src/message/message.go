package message

import "sync"

var subID = 0
var lock = sync.RWMutex{}
var msg = make(map[int]map[int]MsgProcFun)

type MsgProcFun func(int, ...any)

func Sub(msgID int, fun MsgProcFun) int {
	lock.Lock()
	defer lock.Unlock()

	subID += 1

	funs, ok := msg[msgID]
	if !ok {
		msg[msgID] = map[int]MsgProcFun{subID: fun}
		return subID
	}

	funs[subID] = fun
	return subID
}

func Unsub(msgID int, subID int) {
	lock.Lock()
	defer lock.Unlock()

	funs, ok := msg[msgID]
	if !ok {
		return
	}

	delete(funs, subID)
}

func Pub(id int, args ...any) {
	lock.RLock()
	defer lock.RUnlock()

	funs, ok := msg[id]
	if !ok {
		return
	}

	for _, fun := range funs {
		fun(id, args)
	}
}
