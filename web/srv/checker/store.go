package checker

import (
	"log"

	"github.com/aaaasmile/live-client/web/idl"
)

type Store struct {
	InfoObjects   map[string]*idl.ObjectInfo
	Debug         bool
	ObjTypeInProv idl.ObjTypeInProv
	synchro       *StoreSynchronizer
	chChanges     chan idl.ObjInfoChange
}

func NewStore(otp idl.ObjTypeInProv) Store {
	log.Println("New store ", otp.String())
	res := Store{
		InfoObjects:   make(map[string]*idl.ObjectInfo),
		ObjTypeInProv: otp,
	}

	return res
}

func (st *Store) InsertOrUpdateSingleObj(oi *idl.ObjectInfo) idl.ObjOpChangeType {
	// Keep in mind that this function should be called always inside the same routine.
	// Changes will be propagate to the provider (sqlite) only if a StoreSynchronizer was activated before with StartSyncWithProv.
	key := oi.Key
	if key == "" {
		panic("Key could not be empty")
	}
	ntfy := idl.ObjInfoChange{Obj: oi}
	if st.InfoObjects[key] == nil {
		ntfy.ChangeType = idl.OOCTinsert
		st.chChanges <- ntfy
		st.InfoObjects[key] = oi
		if st.Debug {
			log.Println("Insert item in store ", oi)
		}
	} else {
		if !oi.IsEqual(st.InfoObjects[key]) {
			oldOi := st.InfoObjects[key]
			ntfy.Obj.SourceFile.DbLiteID = oldOi.SourceFile.DbLiteID

			ntfy.ChangeType = idl.OOCTupdate
			st.chChanges <- ntfy
			st.InfoObjects[key] = oi
			if st.Debug {
				log.Println("Update  item in store ", oi)
			}
		} else {
			ntfy.ChangeType = idl.OOCTconfirm
			st.chChanges <- ntfy
			if st.Debug {
				log.Println("Unchanged item", oi.Key)
			}
		}
	}
	return ntfy.ChangeType
}

func (st *Store) UpdateSingleObj(oi *idl.ObjectInfo) {
	key := oi.Key
	ntfy := idl.ObjInfoChange{Obj: oi}
	ntfy.ChangeType = idl.OOCTupdate
	st.chChanges <- ntfy
	st.InfoObjects[key] = oi
	//fmt.Println("*** store update oi ", oi)
	if st.Debug {
		log.Println("Update  item in store ", oi)
	}
}

func (st *Store) ResetInfoWithTrack() {
	log.Println("Reset store with notification", st.ObjTypeInProv.String())
	for _, v := range st.InfoObjects {
		ntfy := idl.ObjInfoChange{ChangeType: idl.OOCTdelete, Obj: v}
		st.chChanges <- ntfy
	}
	st.InfoObjects = make(map[string]*idl.ObjectInfo)
}

func (st *Store) InsertStoreKeys(keys []string) {
	log.Println("Insert keys in store for", st.ObjTypeInProv.String())
	count := 0
	for _, key := range keys {
		if st.InfoObjects[key] == nil {
			count++
			oi := &idl.ObjectInfo{
				Key: key,
			}
			ntfy := idl.ObjInfoChange{ChangeType: idl.OOCTinsert, Obj: oi}
			st.chChanges <- ntfy
			st.InfoObjects[key] = oi
		}
	}
	log.Println("Inserted items ", count)
}

func (st *Store) DeleteStoreKeys(keys []string) {
	log.Println("Delete keys in store for", st.ObjTypeInProv.String())
	for _, key := range keys {
		st.DeleteStoreKey(key)
	}
}

func (st *Store) DeleteStoreKey(key string) {
	if st.InfoObjects[key] != nil {
		ntfy := idl.ObjInfoChange{ChangeType: idl.OOCTdelete, Obj: st.InfoObjects[key]}
		st.chChanges <- ntfy
		delete(st.InfoObjects, key)
	}
}

func (st *Store) EndTrackingChanges() {
	close(st.chChanges) // trigger end of changes
	log.Println("End tracking changes in store ", st.ObjTypeInProv.String())
}

func (st *Store) Abort() {
	log.Println("Abort tracking")
	st.synchro.Abort()
}

func (st *Store) PopulateObjs(prov idl.ObjProvider) error {
	log.Println("Populate store  objects from provider ", st.ObjTypeInProv.String())
	st.InfoObjects = make(map[string]*idl.ObjectInfo)

	objs, err := prov.DoReadAllObj(st.ObjTypeInProv)
	if err != nil {
		return err
	}

	for _, obj := range objs {
		if obj.Key == "" {
			panic("Key could not be empty")
		}
		st.InfoObjects[obj.Key] = obj // changes are coming from provider, so no need to track those changes.
		if st.Debug {
			log.Println("Set in store obj ", obj)
		}
	}
	log.Println("Store populated with items: ", len(objs), st.ObjTypeInProv.String())
	return nil
}

func (st *Store) StartSyncWithProv(prov idl.ObjProvider, tag string, readyPrev chan idl.ResErr) chan idl.ResErr {
	log.Println("Start sync for delete with provider ", st.ObjTypeInProv.String())
	// Start collecting changes in store (InfoObjects) to be persisted into a provider (sqlite db).
	// Channels (e.g chInsert) are used to collect those changes. Changes are written in the provider
	// usinf the method sync.EndSyncWithProv.

	sync := StoreSynchronizer{
		ObjTypeInProv: st.ObjTypeInProv,
	}
	st.synchro = &sync

	sync.Initialize(tag)

	bufLen := 10
	if len(st.InfoObjects) > 10 {
		bufLen = len(st.InfoObjects) / 2
	}
	st.chChanges = make(chan idl.ObjInfoChange, bufLen)
	readyCh := sync.ListenChanges(st.chChanges, prov, readyPrev)

	return readyCh
}
