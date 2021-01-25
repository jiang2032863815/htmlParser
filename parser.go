package htmlParser

type parseNode struct {
	tp NodeType
	l,r int
}
type Parser struct {
	ele []parseNode
	nodeNameMp map[int]string
	nodeAttributesMp map[int]map[string]string
	eleMatch map[int]int
	eleIsEnd map[int]bool
	content string
}
func isAvailableNameChar(c uint8)bool{
	return c=='!'||c>='0'&&c<='9'||c>='a'&&c<='z'||c>='A'&&c<='Z'
}
func checkIsEnd(content string)bool{
	for i:=0;i<len(content);i++{
		if content[i]=='/'{
			return true
		}
	}
	return false
}
func readNodeName(content string)string{
	var ret string
	var L=len(content)
	var idx=0
	for idx<L&&!isAvailableNameChar(content[idx]){
		idx++
	}
	for idx<L&&isAvailableNameChar(content[idx]){
		ret+=string(content[idx])
		idx++
	}
	return ret
}
func readAttributes(content string)map[string]string{
	var ret=map[string]string{}
	var L=len(content)
	var idx=0
	for idx<L&&!isAvailableNameChar(content[idx]){
		idx++
	}
	for idx<L&&isAvailableNameChar(content[idx]){
		idx++
	}
	for idx<L{
		for idx<L&&!isAvailableNameChar(content[idx]){
			idx++
		}
		if idx<L{
			var nameStr string
			for idx<L&&isAvailableNameChar(content[idx])&&content[idx]!='='{
				nameStr+=string(content[idx])
				idx++
			}
			var valueStr string
			if idx<L&&content[idx]=='='{
				idx++
			}
			for idx<L&&isAvailableNameChar(content[idx]){
				if content[idx]!='"'&&content[idx]!='\''{
					valueStr+=string(content[idx])
				}
				idx++
			}
			ret[nameStr]=valueStr
		}
	}
	return ret
}
func(p *Parser)build(l,r int)[]*Node{
	if l>r{
		return nil
	}
	var nowNodes []*Node
	for i:=l;i<=r;{
		if p.ele[i].tp==NodeTypeText{
			nowNodes=append(nowNodes,&Node{NodeTypeText,"",false,nil,nil,p.content[p.ele[i].l:p.ele[i].r+1]})
			i++
		}else{
			if p.eleMatch[i]==i{
				nowNodes=append(nowNodes,&Node{NodeTypeElement,p.nodeNameMp[i],false,p.nodeAttributesMp[i],nil,""})
			}else{
				nowNodes=append(nowNodes,&Node{NodeTypeElement,p.nodeNameMp[i],true,p.nodeAttributesMp[i],p.build(i+1,p.eleMatch[i]-1),""})
			}
			i=p.eleMatch[i]+1
		}
	}
	return nowNodes
}
func(p *Parser)init(){
	var L= len(p.content)
	var idx=0
	var xMatch=map[int]int{}
	var xSt []int
	for i:=0;i<L;i++{
		if p.content[i]=='<'{
			xSt=append(xSt,i)
		}else if p.content[i]=='>'{
			if len(xSt)>0{
				xMatch[xSt[len(xSt)-1]]=i
				xMatch[i]=xSt[len(xSt)-1]
				xSt=xSt[:len(xSt)-1]
			}
		}
	}
	for idx<L{
		if xMatch[idx]==0{
			var pNd=parseNode{NodeTypeText,idx,idx-1}
			for idx<L&&xMatch[idx]==0{
				pNd.r++
				idx++
			}
			p.ele=append(p.ele,pNd)
		}else{
			if idx+1<=xMatch[idx]-1 {
				var pNd = parseNode{NodeTypeElement, idx + 1, xMatch[idx] - 1}
				p.ele=append(p.ele,pNd)
			}
			idx=xMatch[idx]+1
		}
	}
	for i:=0;i<len(p.ele);i++{
		if p.ele[i].tp==NodeTypeElement{
			p.nodeNameMp[i]=readNodeName(p.content[p.ele[i].l:p.ele[i].r+1])
			p.eleIsEnd[i]=checkIsEnd(p.content[p.ele[i].l:p.ele[i].r+1])
			p.nodeAttributesMp[i]=readAttributes(p.content[p.ele[i].l:p.ele[i].r+1])
			p.eleMatch[i]=i
		}
	}
	idx=0
	L=len(p.ele)
	var nameMpSt=map[string][]int{}
	for idx<L{
		if p.ele[idx].tp==NodeTypeText{
			idx++
		}else{
			if p.eleIsEnd[idx]{
				var st=nameMpSt[p.nodeNameMp[idx]]
				if len(st)>0{
					p.eleMatch[st[len(st)-1]]=idx
					st=st[:len(st)-1]
				}
				nameMpSt[p.nodeNameMp[idx]]=st
				idx++
			}else{
				var st=nameMpSt[p.nodeNameMp[idx]]
				st=append(st,idx)
				nameMpSt[p.nodeNameMp[idx]]=st
				idx++
			}
		}
	}
	for _,v:=range nameMpSt{
		for _,x:=range v{
			p.eleMatch[x]=x
		}
	}
}
func(p *Parser)Parse(content string)*Node{
	p.content=content
	p.eleMatch=map[int]int{}
	p.nodeNameMp=map[int]string{}
	p.eleIsEnd=map[int]bool{}
	p.nodeAttributesMp=map[int]map[string]string{}
	p.init()
	var document=&Node{NodeTypeElement,"Document",true,nil,p.build(0,len(p.ele)-1),""}
	return document
}
func NewParser()*Parser{
	return &Parser{}
}