package htmlParser

import (
	"fmt"
	"testing"
)
func printSpace(cnt int){
	for cnt>0{
		fmt.Print(" ")
		cnt--
	}
}
func printNode(nd *Node,ex int){
	if nd.Type==NodeTypeElement{
		printSpace(ex)
		fmt.Println(nd.NodeName)
		if nd.NodeName=="title"{
			fmt.Println("title="+nd.ChildNodes[0].Content)
		}
		for _,v:=range nd.ChildNodes{
			printNode(v,ex+3)
		}
	}
}
func TestNewParser(t *testing.T) {
	var document=NewParser().Parse(`
		<!DOCTYPE html>
		<html>
			<head>
				<title>AAA</title>
				<meta charset="UTF-8"/>
			</head>
			<body>
				<div>
					<img src="aaa">
					<div>
						<ul>
							<li></li>
							<li></li>
						</ul>
					</div>
				</div>
			</body>
		</html>
	`)
	printNode(document,0)
}