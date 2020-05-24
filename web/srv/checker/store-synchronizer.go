package checker

import (
	"log"
	"sort"
	"time"

	"github.com/aaaasmile/live-client/web/idl"
)

type StoreSynchronizer struct {
	stickyObjs    map[string]*idl.ObjectInfo
	updatedObjs   map[string]*idl.ObjectInfo
	insertedObjs  map[string]*idl.ObjectInfo
	chAbort       chan struct{}
	ObjTypeInProv idl.ObjTypeInProv
	tag           string
	Debug         bool
}

func (s *StoreSynchronizer) Initialize(tag string) {
	s.stickyObjs = make(map[string]*idl.ObjectInfo)
	s.updatedObjs = make(map[string]*idl.ObjectInfo)
	s.insertedObjs = make(map[string]*idl.ObjectInfo)
	s.chAbort = make(chan struct{})
	s.tag = tag
}

func (s *StoreSynchronizer) ListenChanges(chChanges chan idl.ObjInfoChange, prov idl.ObjProvider, readyPrev chan idl.ResErr) chan idl.ResErr {
	if chChanges == nil {
		panic("Unitialized delete channel")
	}
	readyMe := make(chan idl.ResErr)

	go func() {
		log.Println("Start listening for changes", s.ObjTypeInProv.String())
		haschanges := false
	loop:
		for {
			select {
			case oic, more := <-chChanges:
				if s.Debug {
					log.Println("change arrived ", oic, more)
				}
				if more {
					haschanges = true
					switch oic.ChangeType {
					case idl.OOCTupdate:
						s.addForUpdate(oic.Obj)
					case idl.OOCTconfirm:
						s.unchanged(oic.Obj)
					case idl.OOCTdelete:
						s.addDeleted(oic.Obj)
					case idl.OOCTinsert:
						s.addForInsert(oic.Obj)
					}
				} else {
					break loop
				}
			case <-s.chAbort:
				log.Printf("ABORT: Storesynchronizer %s listener terminated", s.ObjTypeInProv.String())
				return
			}
		}
		if readyPrev != nil {
			log.Println("Waiting for previous")
			itemRes := <-readyPrev
			if itemRes.Err != nil {
				log.Println(" --+- Error --+- inside previous process ", itemRes.Err)
				readyMe <- itemRes // rethrow error
				return
			}
		}
		var err error
		if haschanges {
			log.Println("start flush...")
			err = s.flushSynchInProvider(prov)
			if err != nil {
				log.Println("-+- Flush sync error ", err)
			}
			log.Println("end flush")
		}
		readyMe <- idl.ResErr{Err: err}
		close(readyMe)
	}()

	return readyMe
}

func (s *StoreSynchronizer) addDeleted(obj *idl.ObjectInfo) {
	if s.insertedObjs[obj.Key] != nil {
		delete(s.insertedObjs, obj.Key)
	}
	if s.updatedObjs[obj.Key] != nil {
		delete(s.updatedObjs, obj.Key)
	}
	s.stickyObjs[obj.Key] = obj
}

func (s *StoreSynchronizer) addForInsert(obj *idl.ObjectInfo) {
	delete(s.stickyObjs, obj.Key)
	s.insertedObjs[obj.Key] = obj
}

func (s *StoreSynchronizer) addForUpdate(obj *idl.ObjectInfo) {
	delete(s.stickyObjs, obj.Key)
	s.updatedObjs[obj.Key] = obj
}

func (s *StoreSynchronizer) unchanged(obj *idl.ObjectInfo) {
	delete(s.stickyObjs, obj.Key)
}

func (sync *StoreSynchronizer) Abort() {
	select {
	case <-sync.chAbort:
	default:
		log.Printf("[%s]Synchronizer will be closed", sync.ObjTypeInProv.String())
		sync.chAbort <- struct{}{}
		close(sync.chAbort)
	}
}

func (sync *StoreSynchronizer) flushSynchInProvider(prov idl.ObjProvider) error {
	// sync with provider (sql lite) in background
	log.Printf("[%s]flushSynchInProvider. Insert/Update/Delete: %d/%d/%d", sync.ObjTypeInProv.String(),
		len(sync.insertedObjs), len(sync.updatedObjs), len(sync.stickyObjs))

	start := time.Now()
	trx, err := prov.GetNewTransaction()
	if err != nil {
		log.Println("Sync with provider failed: ", err)
		return err
	}

	if len(sync.insertedObjs) > 0 {
		keys := make([]string, 0, len(sync.insertedObjs))
		for k := range sync.insertedObjs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if sync.Debug {
				log.Println("Insert Key ", k)
			}
			err := prov.DoInsertObject(trx, sync.insertedObjs[k], sync.ObjTypeInProv)
			if err != nil {
				return err
			}
		}
	}

	for _, obj := range sync.updatedObjs {
		err := prov.DoUpdateObject(trx, obj, sync.ObjTypeInProv)
		if err != nil {
			return err
		}
	}

	for _, obj := range sync.stickyObjs {
		err := prov.DoDeleteObject(trx, obj, sync.ObjTypeInProv)
		if err != nil {
			return err
		}
	}

	trx.Commit()
	log.Printf("[%s] Provider is now in synch with the store %s, elapsed %v",
		sync.tag, sync.ObjTypeInProv.String(), time.Now().Sub(start))
	return nil
}
