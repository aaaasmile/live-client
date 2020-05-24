package checker

import (
	"fmt"
	"log"
	"sort"

	"github.com/aaaasmile/live-client/web/idl"
)

type AskFetchSourceFn func() (*idl.SourceFile, error)

type Checker struct {
	ResultView          []*CheckerItem
	ServerOnlyCount     int
	FileSourceOnlyCount int
	EqualCount          int
	DiffCount           int
	Debug               bool
}

func (ck *Checker) String() string {
	res := ""
	// for _, item := range ck.ResultView {
	// 	res += fmt.Sprintf("[ %v ]", item)
	// }
	res += "\n"
	res += fmt.Sprintf("Diff: %d\n", ck.DiffCount)
	res += fmt.Sprintf("Equal: %d\n", ck.EqualCount)
	res += fmt.Sprintf("Server only: %d\n", ck.ServerOnlyCount)
	res += fmt.Sprintf("FileSource only: %d\n", ck.FileSourceOnlyCount)
	return res
}

func (ck *Checker) CreateResultView(storeServer, storeSourceFile *Store) {
	log.Printf("Compare %d server items with %d source file items\n", len(storeServer.InfoObjects), len(storeSourceFile.InfoObjects))
	res := []*CheckerItem{}
	onlyInNav := CheckerItemColl{}
	onlyInSourceFile := CheckerItemColl{}
	diffNavFile := CheckerItemColl{}
	ck.DiffCount, ck.EqualCount, ck.ServerOnlyCount, ck.FileSourceOnlyCount = 0, 0, 0, 0

	for k1, item1 := range storeServer.InfoObjects {
		ci := NewCheckerItem(PresTypeServerOnly, item1)
		if item2, ok := storeSourceFile.InfoObjects[k1]; ok {
			if ci.HasDiff(item1, item2) {
				diffNavFile = append(diffNavFile, ci)
			} else {
				if ck.Debug {
					//log.Println("same item: ", item2)
				}
				ck.EqualCount++
			}
		} else {
			onlyInNav = append(onlyInNav, ci)
		}
	}

	for k2, item2 := range storeSourceFile.InfoObjects {
		ci := NewCheckerItem(PresTypeFileSourceOnly, item2)
		if _, ok := storeServer.InfoObjects[k2]; !ok {
			onlyInSourceFile = append(onlyInSourceFile, ci)
		}
	}
	ll1 := len(onlyInNav)
	ll2 := len(onlyInSourceFile)
	ldif := len(diffNavFile)
	ck.DiffCount, ck.ServerOnlyCount, ck.FileSourceOnlyCount = ldif, ll1, ll2

	if ll1 > 0 {
		if ck.Debug {
			log.Println("Elements on Server only: ", ll1)
		}
		sort.Sort(onlyInNav)
		for _, item := range onlyInNav {
			res = append(res, item)
		}
	}
	if ll2 > 0 {
		if ck.Debug {
			log.Println("Elements in Source File only: ", ll2)
		}
		sort.Sort(onlyInSourceFile)
		for _, item := range onlyInSourceFile {
			res = append(res, item)
		}
	}

	if ldif > 0 {
		if ck.Debug {
			log.Println("Elements different: ", ldif)
		}
		sort.Sort(diffNavFile)
		for _, item := range diffNavFile {
			res = append(res, item)
		}
	}

	ck.ResultView = res

	log.Println("Result compare view:", ck.String())
}
