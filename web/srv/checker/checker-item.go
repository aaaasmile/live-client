package checker

import (
	"strings"

	"github.com/aaaasmile/live-client/util"
	"github.com/aaaasmile/live-client/web/idl"
)

type DateTimeEqual int
type ContentEqual int
type PresenceType int
type OtherDiffType int

const (
	ODTNone OtherDiffType = iota
	ODTVersionListDiff
	ODTContent
	ODTCName
)

const (
	DTCCompareNotAvailable DateTimeEqual = iota
	DTCServerItemNewer
	DTCFileSourceItemRecent
	DTCDateEqual
)

const (
	ContentCheckNotAvailable ContentEqual = iota
	ContentDiffer
	ContentSame
)

const (
	PresTypeBoth PresenceType = iota
	PresTypeServerOnly
	PresTypeFileSourceOnly
)

func (g OtherDiffType) MarshalJSON() ([]byte, error) {
	res := ""
	switch g {
	case ODTNone:
		res = "None"
	case ODTVersionListDiff:
		res = "Diff Version"
	case ODTContent:
		res = "Diff Content"
	case ODTCName:
		res = "Diff Name"
	}
	return util.WrapJsonToString(res)
}

func (g DateTimeEqual) MarshalJSON() ([]byte, error) {
	res := ""
	switch g {
	case DTCCompareNotAvailable:
		res = "Not available"
	case DTCServerItemNewer:
		res = "Server is newer"
	case DTCFileSourceItemRecent:
		res = "File Source is newer"
	case DTCDateEqual:
		res = "Same date/time"
	}
	return util.WrapJsonToString(res)
}

func (g ContentEqual) MarshalJSON() ([]byte, error) {
	res := ""
	switch g {
	case ContentCheckNotAvailable:
		res = "Not available"
	case ContentDiffer:
		res = "Different"
	case ContentSame:
		res = "Equal"
	}
	return util.WrapJsonToString(res)
}

func (g PresenceType) MarshalJSON() ([]byte, error) {
	res := ""
	switch g {
	case PresTypeBoth:
		res = "Both"
	case PresTypeServerOnly:
		res = "Server only"
	case PresTypeFileSourceOnly:
		res = "File Source only"
	}
	return util.WrapJsonToString(res)
}

type CheckerItem struct {
	KeyStore         string
	Name             string
	AreEqual         bool
	VersionListEqual bool
	ContentEqual     ContentEqual
	PresenceType     PresenceType
	OtherDiffType    OtherDiffType
	DateTimeEqual    DateTimeEqual
	DateTimeExp      string
	VersionExp       string
	VersionExp2      string
	DateTimeExp2     string
}

func NewCheckerItem(pt PresenceType, oi *idl.ObjectInfo) *CheckerItem {
	ci := CheckerItem{
		KeyStore:     oi.Key,
		PresenceType: pt,
		Name:         oi.Name,
		VersionExp:   oi.VersionList,
		DateTimeExp:  oi.Timestamp.Format("02-01-2006 15:01:02"),
	}
	return &ci
}

func normalizeVersionList(vl string) string {
	arr := strings.Split(vl, ",")
	res := ""
	for ix, item := range arr {
		spaced := strings.Split(item, " ")
		if ix > 0 {
			res += ","
		}
		res += spaced[0]
	}
	return res
}

func (ci *CheckerItem) HasDiff(serverItem *idl.ObjectInfo, sourceItem *idl.ObjectInfo) bool {
	// if serverItem.ObjectID == 90900 {
	// 	fmt.Println("**Diff ", serverItem, sourceItem)
	// }
	switch ts := serverItem.Timestamp; {
	case ts.After(sourceItem.Timestamp):
		ci.DateTimeEqual = DTCServerItemNewer
	case ts.Before(sourceItem.Timestamp):
		ci.DateTimeEqual = DTCFileSourceItemRecent
	case ts == sourceItem.Timestamp:
		ci.DateTimeEqual = DTCDateEqual
	default:
		ci.DateTimeEqual = DTCCompareNotAvailable
	}
	if ci.DateTimeEqual != DTCDateEqual && ci.DateTimeEqual != DTCCompareNotAvailable {
		//fmt.Printf("*** diff time %s %d %d %s\n", sourceItem.Name, sourceItem.ObjectID, sourceItem.Type, sourceItem.VersionList)
	}

	srcVerList := normalizeVersionList(sourceItem.VersionList)
	navVerList := normalizeVersionList(serverItem.VersionList)
	ci.VersionListEqual = (navVerList == srcVerList)
	// if !ci.VersionListEqual {
	// 	fmt.Printf("*** diff VL: %s --- %s\n", navVerList, srcVerList)
	// }

	ci.PresenceType = PresTypeBoth

	if serverItem.Checksum != "" && sourceItem.Checksum != "" {
		if serverItem.Checksum != sourceItem.Checksum {
			ci.ContentEqual = ContentDiffer
		} else {
			ci.ContentEqual = ContentSame
		}
	} else {
		ci.ContentEqual = ContentCheckNotAvailable
	}

	switch {
	//case n.Name != sourceItem.Name:
	//	ci.OtherDiffType = ODTCName // TODO
	case ci.ContentEqual == ContentDiffer:
		ci.OtherDiffType = ODTContent
	case !ci.VersionListEqual:
		ci.OtherDiffType = ODTVersionListDiff
	default:
		ci.OtherDiffType = ODTNone
	}

	ci.AreEqual = ci.VersionListEqual
	if ci.DateTimeEqual != DTCCompareNotAvailable {
		ci.AreEqual = ci.AreEqual && (ci.DateTimeEqual == DTCDateEqual)
	}

	if ci.ContentEqual != ContentCheckNotAvailable {
		ci.AreEqual = ci.AreEqual && (ci.ContentEqual == ContentSame)
	}

	ci.DateTimeExp2 = sourceItem.Timestamp.Format("02-01-2006 15:01:02")
	ci.VersionExp2 = sourceItem.VersionList

	return !ci.AreEqual
}

type CheckerItemColl []*CheckerItem

func (a CheckerItemColl) Len() int      { return len(a) }
func (a CheckerItemColl) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CheckerItemColl) Less(i, j int) bool {
	return a[i].KeyStore < a[j].KeyStore
}
