package jdb

import (
	"testing"
	"time"

	"github.com/silas/jdb/internal/ptr"
	"github.com/stretchr/testify/require"
)

type rowFieldsSet struct {
	Kind            bool
	ID              bool
	ParentKind      bool
	ParentID        bool
	Data            bool
	UniqueStringKey bool
	StringKey       bool
	NumericKey      bool
	CreateTime      bool
	UpdateTime      bool
}

func TestRowScanMeta(t *testing.T) {
	none := struct{}{}
	r, err := rowScanMeta(none, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{})

	kind := struct {
		Kind string `jdb:"-kind"`
	}{"kind"}
	r, err = rowScanMeta(kind, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{Kind: true})
	require.Equal(t, "kind", r.Kind)

	kindPtr := struct {
		Kind *string `jdb:"-kind"`
	}{ptr.String("kind")}
	r, err = rowScanMeta(kindPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{Kind: true})
	require.Equal(t, "kind", r.Kind)

	id := struct {
		ID string `jdb:"-id"`
	}{"id"}
	r, err = rowScanMeta(id, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ID: true})
	require.Equal(t, "id", r.ID)

	idPtr := struct {
		ID *string `jdb:"-id"`
	}{ptr.String("id")}
	r, err = rowScanMeta(idPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ID: true})
	require.Equal(t, "id", r.ID)

	parentKind := struct {
		ParentKind string `jdb:"-parentkind"`
	}{"parentKind"}
	r, err = rowScanMeta(parentKind, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ParentKind: true})
	require.Equal(t, ptr.String("parentKind"), r.ParentKind)

	parentKindPtr := struct {
		ParentKind *string `jdb:"-parentkind"`
	}{ptr.String("parentKind")}
	r, err = rowScanMeta(parentKindPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ParentKind: true})
	require.Equal(t, ptr.String("parentKind"), r.ParentKind)

	parentID := struct {
		ParentID string `jdb:"-parentid"`
	}{"parentID"}
	r, err = rowScanMeta(parentID, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ParentID: true})
	require.Equal(t, ptr.String("parentID"), r.ParentID)

	parentIDPtr := struct {
		ParentID *string `jdb:"-parentid"`
	}{ptr.String("parentID")}
	r, err = rowScanMeta(parentIDPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{ParentID: true})
	require.Equal(t, ptr.String("parentID"), r.ParentID)

	uniqueStringKey := struct {
		UniqueStringKey string `jdb:",uniquestringkey"`
	}{"uniqueStringKey"}
	r, err = rowScanMeta(uniqueStringKey, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{UniqueStringKey: true})
	require.Equal(t, ptr.String("uniqueStringKey"), r.UniqueStringKey)

	uniqueStringKeyPtr := struct {
		UniqueStringKey *string `jdb:",uniquestringkey"`
	}{ptr.String("uniqueStringKey")}
	r, err = rowScanMeta(uniqueStringKeyPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{UniqueStringKey: true})
	require.Equal(t, ptr.String("uniqueStringKey"), r.UniqueStringKey)

	stringKey := struct {
		StringKey string `jdb:",stringkey"`
	}{"stringKey"}
	r, err = rowScanMeta(stringKey, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{StringKey: true})
	require.Equal(t, ptr.String("stringKey"), r.StringKey)

	stringKeyPtr := struct {
		StringKey *string `jdb:",stringkey"`
	}{ptr.String("stringKey")}
	r, err = rowScanMeta(stringKeyPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{StringKey: true})
	require.Equal(t, ptr.String("stringKey"), r.StringKey)

	numericKey := struct {
		NumericKey float64 `jdb:",numerickey"`
	}{12.5}
	r, err = rowScanMeta(numericKey, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{NumericKey: true})
	require.Equal(t, ptr.Float64(12.5), r.NumericKey)

	numericKeyPtr := struct {
		NumericKey *float64 `jdb:",numerickey"`
	}{ptr.Float64(12.5)}
	r, err = rowScanMeta(numericKeyPtr, false)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{NumericKey: true})
	require.Equal(t, ptr.Float64(12.5), r.NumericKey)

	now := time.Now()

	createTime := struct {
		CreateTime time.Time `jdb:"-createtime"`
	}{now}
	r, err = rowScanMeta(createTime, true)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{CreateTime: true})
	require.Equal(t, &now, r.CreateTime)

	createTimePtr := struct {
		CreateTime *time.Time `jdb:"-createtime"`
	}{&now}
	r, err = rowScanMeta(createTimePtr, true)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{CreateTime: true})
	require.Equal(t, &now, r.CreateTime)

	createTimeNil := struct {
		CreateTime *time.Time `jdb:"-createtime"`
	}{nil}
	r, err = rowScanMeta(createTimeNil, true)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{})

	updateTime := struct {
		UpdateTime time.Time `jdb:"-updatetime"`
	}{now}
	r, err = rowScanMeta(updateTime, true)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{UpdateTime: true})
	require.Equal(t, &now, r.UpdateTime)

	updateTimePtr := struct {
		UpdateTime *time.Time `jdb:"-updatetime"`
	}{&now}
	r, err = rowScanMeta(updateTimePtr, true)
	require.NoError(t, err)
	requireFieldsSet(t, r, rowFieldsSet{UpdateTime: true})
	require.Equal(t, &now, r.UpdateTime)
}

func requireFieldsSet(t *testing.T, r *row, set rowFieldsSet) {
	require.Equal(t, set.Kind, r.Kind != "")
	require.Equal(t, set.ID, r.ID != "")
	require.Equal(t, set.ParentKind, r.ParentKind != nil)
	require.Equal(t, set.ParentID, r.ParentID != nil)
	require.Equal(t, set.Data, r.Data != nil)
	require.Equal(t, set.UniqueStringKey, r.UniqueStringKey != nil)
	require.Equal(t, set.StringKey, r.StringKey != nil)
	require.Equal(t, set.NumericKey, r.NumericKey != nil)
	require.Equal(t, set.CreateTime, r.CreateTime != nil)
	require.Equal(t, set.UpdateTime, r.UpdateTime != nil)
}
