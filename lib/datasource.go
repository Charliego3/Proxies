package lib

import (
	"github.com/progrium/macdriver/macos/appkit"
	"github.com/progrium/macdriver/macos/foundation"
	"github.com/progrium/macdriver/objc"
)

type OutlineViewDatasource struct {
	_OutlineViewSortDescriptorsDidChange                   func(outlineView appkit.OutlineView, oldDescriptors []foundation.SortDescriptor)
	_OutlineViewChildOfItem                                func(outlineView appkit.OutlineView, index int, item objc.Object) objc.Object
	_OutlineViewPersistentObjectForItem                    func(outlineView appkit.OutlineView, item objc.Object) objc.Object
	_OutlineViewValidateDropProposedItemProposedChildIndex func(outlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) appkit.DragOperation
	_OutlineViewDraggingSessionWillBeginAtPointForItems    func(outlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, draggedItems []objc.Object)
	_OutlineViewAcceptDropItemChildIndex                   func(outlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) bool
	_OutlineViewDraggingSessionEndedAtPointOperation       func(outlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation)
	_OutlineViewSetObjectValueForTableColumnByItem         func(outlineView appkit.OutlineView, object objc.Object, tableColumn appkit.TableColumn, item objc.Object)
	_OutlineViewPasteboardWriterForItem                    func(outlineView appkit.OutlineView, item objc.Object) appkit.PasteboardWritingObject
	_OutlineViewNumberOfChildrenOfItem                     func(outlineView appkit.OutlineView, item objc.Object) int
	_OutlineViewObjectValueForTableColumnByItem            func(outlineView appkit.OutlineView, tableColumn appkit.TableColumn, item objc.Object) objc.Object
	_OutlineViewUpdateDraggingItemsForDrag                 func(outlineView appkit.OutlineView, draggingInfo appkit.DraggingInfoObject)
	_OutlineViewItemForPersistentObject                    func(outlineView appkit.OutlineView, object objc.Object) objc.Object
	_OutlineViewIsItemExpandable                           func(outlineView appkit.OutlineView, item objc.Object) bool
}

func (datasource *OutlineViewDatasource) OutlineViewSortDescriptorsDidChange(outlineView appkit.OutlineView, oldDescriptors []foundation.SortDescriptor) {
	datasource._OutlineViewSortDescriptorsDidChange(outlineView, oldDescriptors)
}
func (datasource *OutlineViewDatasource) SetOutlineViewSortDescriptorsDidChange(f func(foutlineView appkit.OutlineView, oldDescriptors []foundation.SortDescriptor)) {
	datasource._OutlineViewSortDescriptorsDidChange = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewSortDescriptorsDidChange() bool {
	return datasource._OutlineViewSortDescriptorsDidChange != nil
}

func (datasource *OutlineViewDatasource) OutlineViewChildOfItem(outlineView appkit.OutlineView, index int, item objc.Object) objc.Object {
	return datasource._OutlineViewChildOfItem(outlineView, index, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewChildOfItem(f func(foutlineView appkit.OutlineView, index int, item objc.Object) objc.Object) {
	datasource._OutlineViewChildOfItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewChildOfItem() bool {
	return datasource._OutlineViewChildOfItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewPersistentObjectForItem(outlineView appkit.OutlineView, item objc.Object) objc.Object {
	return datasource._OutlineViewPersistentObjectForItem(outlineView, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewPersistentObjectForItem(f func(foutlineView appkit.OutlineView, item objc.Object) objc.Object) {
	datasource._OutlineViewPersistentObjectForItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewPersistentObjectForItem() bool {
	return datasource._OutlineViewPersistentObjectForItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewValidateDropProposedItemProposedChildIndex(outlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) appkit.DragOperation {
	return datasource._OutlineViewValidateDropProposedItemProposedChildIndex(outlineView, info, item, index)
}
func (datasource *OutlineViewDatasource) SetOutlineViewValidateDropProposedItemProposedChildIndex(f func(foutlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) appkit.DragOperation) {
	datasource._OutlineViewValidateDropProposedItemProposedChildIndex = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewValidateDropProposedItemProposedChildIndex() bool {
	return datasource._OutlineViewValidateDropProposedItemProposedChildIndex != nil
}

func (datasource *OutlineViewDatasource) OutlineViewDraggingSessionWillBeginAtPointForItems(outlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, draggedItems []objc.Object) {
	datasource._OutlineViewDraggingSessionWillBeginAtPointForItems(outlineView, session, screenPoint, draggedItems)
}
func (datasource *OutlineViewDatasource) SetOutlineViewDraggingSessionWillBeginAtPointForItems(f func(foutlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, draggedItems []objc.Object)) {
	datasource._OutlineViewDraggingSessionWillBeginAtPointForItems = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewDraggingSessionWillBeginAtPointForItems() bool {
	return datasource._OutlineViewDraggingSessionWillBeginAtPointForItems != nil
}

func (datasource *OutlineViewDatasource) OutlineViewAcceptDropItemChildIndex(outlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) bool {
	return datasource._OutlineViewAcceptDropItemChildIndex(outlineView, info, item, index)
}
func (datasource *OutlineViewDatasource) SetOutlineViewAcceptDropItemChildIndex(f func(foutlineView appkit.OutlineView, info appkit.DraggingInfoObject, item objc.Object, index int) bool) {
	datasource._OutlineViewAcceptDropItemChildIndex = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewAcceptDropItemChildIndex() bool {
	return datasource._OutlineViewAcceptDropItemChildIndex != nil
}

func (datasource *OutlineViewDatasource) OutlineViewDraggingSessionEndedAtPointOperation(outlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation) {
	datasource._OutlineViewDraggingSessionEndedAtPointOperation(outlineView, session, screenPoint, operation)
}
func (datasource *OutlineViewDatasource) SetOutlineViewDraggingSessionEndedAtPointOperation(f func(foutlineView appkit.OutlineView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation)) {
	datasource._OutlineViewDraggingSessionEndedAtPointOperation = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewDraggingSessionEndedAtPointOperation() bool {
	return datasource._OutlineViewDraggingSessionEndedAtPointOperation != nil
}

func (datasource *OutlineViewDatasource) OutlineViewSetObjectValueForTableColumnByItem(outlineView appkit.OutlineView, object objc.Object, tableColumn appkit.TableColumn, item objc.Object) {
	datasource._OutlineViewSetObjectValueForTableColumnByItem(outlineView, object, tableColumn, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewSetObjectValueForTableColumnByItem(f func(foutlineView appkit.OutlineView, object objc.Object, tableColumn appkit.TableColumn, item objc.Object)) {
	datasource._OutlineViewSetObjectValueForTableColumnByItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewSetObjectValueForTableColumnByItem() bool {
	return datasource._OutlineViewSetObjectValueForTableColumnByItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewPasteboardWriterForItem(outlineView appkit.OutlineView, item objc.Object) appkit.PasteboardWritingObject {
	return datasource._OutlineViewPasteboardWriterForItem(outlineView, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewPasteboardWriterForItem(f func(foutlineView appkit.OutlineView, item objc.Object) appkit.PasteboardWritingObject) {
	datasource._OutlineViewPasteboardWriterForItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewPasteboardWriterForItem() bool {
	return datasource._OutlineViewPasteboardWriterForItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewNumberOfChildrenOfItem(outlineView appkit.OutlineView, item objc.Object) int {
	return datasource._OutlineViewNumberOfChildrenOfItem(outlineView, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewNumberOfChildrenOfItem(f func(foutlineView appkit.OutlineView, item objc.Object) int) {
	datasource._OutlineViewNumberOfChildrenOfItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewNumberOfChildrenOfItem() bool {
	return datasource._OutlineViewNumberOfChildrenOfItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewObjectValueForTableColumnByItem(outlineView appkit.OutlineView, tableColumn appkit.TableColumn, item objc.Object) objc.Object {
	return datasource._OutlineViewObjectValueForTableColumnByItem(outlineView, tableColumn, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewObjectValueForTableColumnByItem(f func(foutlineView appkit.OutlineView, tableColumn appkit.TableColumn, item objc.Object) objc.Object) {
	datasource._OutlineViewObjectValueForTableColumnByItem = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewObjectValueForTableColumnByItem() bool {
	return datasource._OutlineViewObjectValueForTableColumnByItem != nil
}

func (datasource *OutlineViewDatasource) OutlineViewUpdateDraggingItemsForDrag(outlineView appkit.OutlineView, draggingInfo appkit.DraggingInfoObject) {
	datasource._OutlineViewUpdateDraggingItemsForDrag(outlineView, draggingInfo)
}
func (datasource *OutlineViewDatasource) SetOutlineViewUpdateDraggingItemsForDrag(f func(foutlineView appkit.OutlineView, draggingInfo appkit.DraggingInfoObject)) {
	datasource._OutlineViewUpdateDraggingItemsForDrag = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewUpdateDraggingItemsForDrag() bool {
	return datasource._OutlineViewUpdateDraggingItemsForDrag != nil
}

func (datasource *OutlineViewDatasource) OutlineViewItemForPersistentObject(outlineView appkit.OutlineView, object objc.Object) objc.Object {
	return datasource._OutlineViewItemForPersistentObject(outlineView, object)
}
func (datasource *OutlineViewDatasource) SetOutlineViewItemForPersistentObject(f func(foutlineView appkit.OutlineView, object objc.Object) objc.Object) {
	datasource._OutlineViewItemForPersistentObject = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewItemForPersistentObject() bool {
	return datasource._OutlineViewItemForPersistentObject != nil
}

func (datasource *OutlineViewDatasource) OutlineViewIsItemExpandable(outlineView appkit.OutlineView, item objc.Object) bool {
	return datasource._OutlineViewIsItemExpandable(outlineView, item)
}
func (datasource *OutlineViewDatasource) SetOutlineViewIsItemExpandable(f func(foutlineView appkit.OutlineView, item objc.Object) bool) {
	datasource._OutlineViewIsItemExpandable = f
}
func (datasource *OutlineViewDatasource) HasOutlineViewIsItemExpandable() bool {
	return datasource._OutlineViewIsItemExpandable != nil
}

type TableViewDataSourceDelegate struct {
	_TableViewSetObjectValueForTableColumnRow              func(tableView appkit.TableView, object objc.Object, tableColumn appkit.TableColumn, row int)
	_NumberOfRowsInTableView                               func(tableView appkit.TableView) int
	_TableViewSortDescriptorsDidChange                     func(tableView appkit.TableView, oldDescriptors []foundation.SortDescriptor)
	_TableViewDraggingSessionEndedAtPointOperation         func(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation)
	_TableViewDraggingSessionWillBeginAtPointForRowIndexes func(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, rowIndexes foundation.IndexSet)
	_TableViewAcceptDropRowDropOperation                   func(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) bool
	_TableViewObjectValueForTableColumnRow                 func(tableView appkit.TableView, tableColumn appkit.TableColumn, row int) objc.Object
	_TableViewPasteboardWriterForRow                       func(tableView appkit.TableView, row int) appkit.PasteboardWritingObject
	_TableViewUpdateDraggingItemsForDrag                   func(tableView appkit.TableView, draggingInfo appkit.DraggingInfoObject)
	_TableViewValidateDropProposedRowProposedDropOperation func(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) appkit.DragOperation
}

func (t *TableViewDataSourceDelegate) TableViewSetObjectValueForTableColumnRow(tableView appkit.TableView, object objc.Object, tableColumn appkit.TableColumn, row int) {
	t._TableViewObjectValueForTableColumnRow(tableView, tableColumn, row)
}

func (t *TableViewDataSourceDelegate) HasTableViewSetObjectValueForTableColumnRow() bool {
	return t._TableViewSetObjectValueForTableColumnRow != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewSetObjectValueForTableColumnRow(f func(tableView appkit.TableView, object objc.Object, tableColumn appkit.TableColumn, row int)) {
	t._TableViewSetObjectValueForTableColumnRow = f
}

func (t *TableViewDataSourceDelegate) NumberOfRowsInTableView(tableView appkit.TableView) int {
	return t._NumberOfRowsInTableView(tableView)
}

func (t *TableViewDataSourceDelegate) HasNumberOfRowsInTableView() bool {
	return t._NumberOfRowsInTableView != nil
}

func (t *TableViewDataSourceDelegate) SetNumberOfRowsInTableView(f func(tableView appkit.TableView) int) {
	t._NumberOfRowsInTableView = f
}

func (t *TableViewDataSourceDelegate) TableViewSortDescriptorsDidChange(tableView appkit.TableView, oldDescriptors []foundation.SortDescriptor) {
	t._TableViewSortDescriptorsDidChange(tableView, oldDescriptors)
}

func (t *TableViewDataSourceDelegate) HasTableViewSortDescriptorsDidChange() bool {
	return t._TableViewSortDescriptorsDidChange != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewSortDescriptorsDidChange(f func(tableView appkit.TableView, oldDescriptors []foundation.SortDescriptor)) {
	t._TableViewSortDescriptorsDidChange = f
}

func (t *TableViewDataSourceDelegate) TableViewDraggingSessionEndedAtPointOperation(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation) {
	t._TableViewDraggingSessionEndedAtPointOperation(tableView, session, screenPoint, operation)
}

func (t *TableViewDataSourceDelegate) HasTableViewDraggingSessionEndedAtPointOperation() bool {
	return t._TableViewDraggingSessionEndedAtPointOperation != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewDraggingSessionEndedAtPointOperation(f func(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, operation appkit.DragOperation)) {
	t._TableViewDraggingSessionEndedAtPointOperation = f
}

func (t *TableViewDataSourceDelegate) TableViewDraggingSessionWillBeginAtPointForRowIndexes(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, rowIndexes foundation.IndexSet) {
	t._TableViewDraggingSessionWillBeginAtPointForRowIndexes(tableView, session, screenPoint, rowIndexes)
}

func (t *TableViewDataSourceDelegate) HasTableViewDraggingSessionWillBeginAtPointForRowIndexes() bool {
	return t._TableViewDraggingSessionWillBeginAtPointForRowIndexes != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewDraggingSessionWillBeginAtPointForRowIndexes(f func(tableView appkit.TableView, session appkit.DraggingSession, screenPoint foundation.Point, rowIndexes foundation.IndexSet)) {
	t._TableViewDraggingSessionWillBeginAtPointForRowIndexes = f
}

func (t *TableViewDataSourceDelegate) TableViewAcceptDropRowDropOperation(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) bool {
	return t._TableViewAcceptDropRowDropOperation(tableView, info, row, dropOperation)
}

func (t *TableViewDataSourceDelegate) HasTableViewAcceptDropRowDropOperation() bool {
	return t._TableViewAcceptDropRowDropOperation != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewAcceptDropRowDropOperation(f func(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) bool) {
	t._TableViewAcceptDropRowDropOperation = f
}

func (t *TableViewDataSourceDelegate) TableViewObjectValueForTableColumnRow(tableView appkit.TableView, tableColumn appkit.TableColumn, row int) objc.Object {
	return t._TableViewObjectValueForTableColumnRow(tableView, tableColumn, row)
}

func (t *TableViewDataSourceDelegate) HasTableViewObjectValueForTableColumnRow() bool {
	return t._TableViewObjectValueForTableColumnRow != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewObjectValueForTableColumnRow(f func(tableView appkit.TableView, tableColumn appkit.TableColumn, row int) objc.Object) {
	t._TableViewObjectValueForTableColumnRow = f
}

func (t *TableViewDataSourceDelegate) TableViewPasteboardWriterForRow(tableView appkit.TableView, row int) appkit.PasteboardWritingObject {
	return t._TableViewPasteboardWriterForRow(tableView, row)
}

func (t *TableViewDataSourceDelegate) HasTableViewPasteboardWriterForRow() bool {
	return t._TableViewPasteboardWriterForRow != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewPasteboardWriterForRow(f func(tableView appkit.TableView, row int) appkit.PasteboardWritingObject) {
	t._TableViewPasteboardWriterForRow = f
}

func (t *TableViewDataSourceDelegate) TableViewUpdateDraggingItemsForDrag(tableView appkit.TableView, draggingInfo appkit.DraggingInfoObject) {
	t._TableViewUpdateDraggingItemsForDrag(tableView, draggingInfo)
}

func (t *TableViewDataSourceDelegate) HasTableViewUpdateDraggingItemsForDrag() bool {
	return t._TableViewUpdateDraggingItemsForDrag != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewUpdateDraggingItemsForDrag(f func(tableView appkit.TableView, draggingInfo appkit.DraggingInfoObject)) {
	t._TableViewUpdateDraggingItemsForDrag = f
}

func (t *TableViewDataSourceDelegate) TableViewValidateDropProposedRowProposedDropOperation(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) appkit.DragOperation {
	return t._TableViewValidateDropProposedRowProposedDropOperation(tableView, info, row, dropOperation)
}

func (t *TableViewDataSourceDelegate) HasTableViewValidateDropProposedRowProposedDropOperation() bool {
	return t._TableViewValidateDropProposedRowProposedDropOperation != nil
}

func (t *TableViewDataSourceDelegate) SetTableViewValidateDropProposedRowProposedDropOperation(f func(tableView appkit.TableView, info appkit.DraggingInfoObject, row int, dropOperation appkit.TableViewDropOperation) appkit.DragOperation) {
	t._TableViewValidateDropProposedRowProposedDropOperation = f
}
