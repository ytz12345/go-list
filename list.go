package golist

type IntList interface {
	// 检查一个元素是否存在，如果存在则返回 true，否则返回 false
	Contains(value int) bool

	// 插入一个元素，如果此操作成功插入一个元素，则返回 true，否则返回 false
	Insert(value int) bool

	// 删除一个元素，如果此操作成功删除一个元素，则返回 true，否则返回 false
	Delete(value int) bool

	// 遍历此有序链表的所有元素，如果 f 返回 false，则停止遍历
	Range(f func(value int) bool)

	// 返回有序链表的元素个数
	Len() int
}
