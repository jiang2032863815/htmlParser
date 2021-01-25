package htmlParser

type NodeType int
const(
	NodeTypeText=iota
	NodeTypeElement
)
type Node struct{
	/*元素类型*/
	Type NodeType
	/*标签名称*/
	NodeName string
	/*是否是双标签*/
	IsDouble bool
	/*标签属性*/
	Attributes map[string]string
	/*子元素*/
	ChildNodes []*Node
	/*元素内容*/
	Content string
}