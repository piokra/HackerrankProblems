package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

var MODULO int64 = 1e9 + 7

/*
 * Complete the countStrings function below.
 */

type RegexType uint8

type Alphabet uint8

const (
	Epsilon     Alphabet = 2
	A           Alphabet = 0
	B           Alphabet = 1
)

func runeToAlphabet(r rune) Alphabet {
	if r == 'a' {
		return A
	}

	if r == 'b' {
		return B
	}

	panic("Invalid rune")
}

type NodeType int

const (
	Default         = 0
	Start           = 1
	End             = 2
	Unreachable		= 3
)

type Node struct {
	nodeNumber  int
	edges           [2]*Node
	epsilonEdges    []*Node
	nodeType        NodeType
}

const MAX_NFA_GRAPH_SIZE = 400

type Graph struct {
	nodes       []*Node
	nodesLenght    int
}

func newGraph() *Graph {
	graph := new(Graph)
	return graph
}

func newNode(graph *Graph) *Node {
	node := new(Node)
	node.nodeNumber = graph.nodesLenght
	graph.nodes = append(graph.nodes, node)

	graph.nodesLenght ++
	return node
}



const (
	Rune        RegexType = 0
	And         RegexType = 1
	Or          RegexType = 2
	Any         RegexType = 4
)

var AlphabetStrings = map[Alphabet]string {
	Epsilon: "E1",
	Epsilon + 1: "E2",
	A: "A",
	B: "B",
}

type Regex struct {
	letter rune
	regexType RegexType
	left, right *Regex
}

func appendEpsilonEdge(nodeToAppendTo, nodeToAppend *Node) {
	nodeToAppendTo.epsilonEdges = append(nodeToAppendTo.epsilonEdges, nodeToAppend)
}

func _regexToNFA(regex *Regex, graph *Graph, prevStart, nextEnd *Node) (*Node, *Node) {
	var start, end *Node
	// prevStart, nextEnd = nil, nil
	if prevStart == nil {
		start = newNode(graph)
	} else {
		start = prevStart
	}

	if nextEnd == nil {
		end = newNode(graph)
	} else {
		end = nextEnd
	}
	switch regex.regexType {
	case Rune:
		// fmt.Println("Rune")
		letter := runeToAlphabet(regex.letter)
		if start.edges[letter] != nil {
			panic("Replacing edge")
		} else {
			start.edges[letter] = end
		}
		return start, end
	case Or:
		// fmt.Println("Or")
		// start := newNode(graph)
		// end := newNode(graph)

		leftStart, leftEnd := _regexToNFA(regex.left, graph, nil, nil)
		rightStart, rightEnd := _regexToNFA(regex.right, graph, nil, nil)

		appendEpsilonEdge(start, leftStart)
		appendEpsilonEdge(start, rightStart)

		appendEpsilonEdge(leftEnd, end)
		appendEpsilonEdge(rightEnd, end)
		return start, end
	case And:
		// fmt.Println("And")
		// start := newNode(graph)
		// end := newNode(graph)

		_, firstEnd := _regexToNFA(regex.left, graph, start, nil)
		_regexToNFA(regex.right, graph, firstEnd, end)

		// appendEpsilonEdge(start, firstStart)
		// appendEpsilonEdge(firstEnd, secondStart)
		// appendEpsilonEdge(secondEnd, end)

		return start, end
	case Any:
		// fmt.Println("Any")
		// start := newNode(graph)
		second := newNode(graph)
		// end := newNode(graph)

		appendEpsilonEdge(start, second)
		appendEpsilonEdge(start, end)

		_, leftEnd := _regexToNFA(regex.left, graph, second, nil)

		// appendEpsilonEdge(second, leftStart)
		appendEpsilonEdge(leftEnd, end)
		appendEpsilonEdge(end, second)
		return start, end

	}
	panic("Invalid node type")
}

func regexToNFA(regex *Regex) *Graph {
	graph := newGraph()
	_regexToNFA(regex, graph, nil, nil)
	return graph
}

func _nodeToEpsilonStarTransitions(node* Node, visited map[int]bool, nodeSet []int) []int {
	if node == nil {
		return nodeSet
	}
	if visited[node.nodeNumber] {
		return nodeSet
	}

	visited[node.nodeNumber] = true

	nodeSet = append(nodeSet, node.nodeNumber)


	// nodeSet = _nodeToEpsilonStarTransitions(node.edges[Epsilon+1], visited, nodeSet)

	for _, epsilonNode := range node.epsilonEdges {
		nodeSet = _nodeToEpsilonStarTransitions(epsilonNode, visited, nodeSet)
	}

	return nodeSet
}

func nodeToEpsilonStarTransitions(node* Node) []int {
	var visited map[int]bool = make(map[int]bool)
	var nodeSet []int
	transitions := _nodeToEpsilonStarTransitions(node, visited, nodeSet)
	sort.Sort(sort.IntSlice(transitions))
	return transitions

}

func stringify(nodes []int) string {
	var runes [MAX_NFA_GRAPH_SIZE]rune
	for i, node := range nodes {
		runes[i] = rune(node)
	}
	return string(runes[:len(nodes)])
}

func letterStarStep(graph *Graph, nodes []int, letter Alphabet, letterTransitions [MAX_NFA_GRAPH_SIZE][2]int, epsilonTransitions [MAX_NFA_GRAPH_SIZE][]int) []int {
	skip := 0
	nodesCopy := make([]int, len(nodes))
	for i, node := range nodes {
		nodesCopy[i] = letterTransitions[node][letter]
		if nodesCopy[i] == -1 {
			skip ++
		}
	}

	sort.Sort(sort.IntSlice(nodesCopy))
	nodesCopy = nodesCopy[skip:]
	if len(nodes) == 0 {
		return nodes
	}

	nodesSet := make(map[int]bool)

	for _, node := range nodesCopy {
		for _, nextNode := range epsilonTransitions[node] {
			nodesSet[nextNode] = true
		}
	}

	nodesCopy = []int{}

	for key := range nodesSet {
		nodesCopy = append(nodesCopy, key)
	}

	sort.Sort(sort.IntSlice(nodesCopy))

	return nodesCopy
}

func isAccNode(node []int) bool {
	for _, n := range node {
		if n == 1 {
			return true
		}
	}
	return false
}

func isStartNode(node []int) bool {
	for _, n := range node {
		if n == 0 {
			return true
		}
	}
	return false
}

func nfaToDFA(graph* Graph) (*Graph, []int, []int) {
	var epsilonTransitions [MAX_NFA_GRAPH_SIZE][]int
	var letterTransitions [MAX_NFA_GRAPH_SIZE][2]int

	var newNodes [][]int = make([][]int, MAX_NFA_GRAPH_SIZE)

	for i, node := range graph.nodes {
		epsilonTransitions[i] = nodeToEpsilonStarTransitions(node)
		if node.edges[A] == nil {
			letterTransitions[i][A] = -1
		} else {
			letterTransitions[i][A] = node.edges[A].nodeNumber
		}

		if node.edges[B] == nil {
			letterTransitions[i][B] = -1
		} else {
			letterTransitions[i][B] = node.edges[B].nodeNumber
		}
		copiedNodes := make([]int, len(epsilonTransitions[i]))
		copy(copiedNodes, epsilonTransitions[i])
		newNodes[i] = copiedNodes

		//fmt.Println(epsilonTransitions[i])
		// fmt.Println(letterTransitions[i])
	}
	newNodes = newNodes[:graph.nodesLenght]
	var dfaVisitedMap map[string]bool = make(map[string]bool)
	var dfaNodesMap map[string]int = make(map[string]int)
	dfaGraph := newGraph()
	var accNodes []int
	var startNodes []int

	for len(newNodes) > 0 {
		popped := newNodes[len(newNodes)-1]
		newNodes = newNodes[:len(newNodes)-1]
		stringified := stringify(popped)
		// fmt.Println("Popped:", popped)
		if len(popped) == 0 || dfaVisitedMap[stringified] {
			continue
		}
		// fmt.Println(popped)


		var currentNode *Node
		nodeNb, ok := dfaNodesMap[stringified]
		if ok {
			currentNode = dfaGraph.nodes[nodeNb]
		} else {
			currentNode = newNode(dfaGraph)
			dfaNodesMap[stringified] = currentNode.nodeNumber
		}
		if isStartNode(popped) {
			startNodes = append(startNodes, currentNode.nodeNumber)
		}
		if isAccNode(popped) {
			accNodes = append(accNodes, currentNode.nodeNumber)
			currentNode.nodeType = End
		}



		dfaVisitedMap[stringified] = true

		afterA := letterStarStep(graph, popped, A, letterTransitions, epsilonTransitions)
		afterB := letterStarStep(graph, popped, B, letterTransitions, epsilonTransitions)
		// fmt.Println(afterA)
		// fmt.Println(afterB)
		stringifiedA := stringify(afterA)
		stringifiedB := stringify(afterB)

		var nextNode *Node
		if len(afterA) > 0 {
			nodeNb, ok = dfaNodesMap[stringifiedA]
			if ok {
				nextNode = dfaGraph.nodes[nodeNb]
			} else {
				nextNode = newNode(dfaGraph)
				dfaNodesMap[stringifiedA] = nextNode.nodeNumber
			}
			currentNode.edges[A] = nextNode

		}
		if len(afterB) > 0 {
			nodeNb, ok = dfaNodesMap[stringifiedB]
			if ok {
				nextNode = dfaGraph.nodes[nodeNb]
			} else {
				nextNode = newNode(dfaGraph)
				dfaNodesMap[stringifiedB] = nextNode.nodeNumber
			}
			currentNode.edges[B] = nextNode
		}

		if !dfaVisitedMap[stringifiedA] && len(afterA) > 0 {
			newNodes = append(newNodes, afterA)
		}

		if !dfaVisitedMap[stringifiedB] && len(afterB) > 0 {
			newNodes = append(newNodes, afterB)
		}


	}

	// for key, value := range dfaNodesMap {
	// fmt.Printf("%v : %d\n", []rune(key), value)
	// }
	return dfaGraph, startNodes, accNodes
}


func wipe(runes []rune, l, r int) {
	for i := l; i<r; i++ {
		runes[i] = 'x'
	}
}

func has(r []rune, c rune) bool {
	for i := range r {
		if r[i] == c {
			return true
		}
	}
	return false
}

func where(r []rune, c rune) int {
	for i := range r {
		if r[i] == c {
			return i
		}
	}
	return len(r)
}

func withoutSpecialCharacters(r []rune) string {
	var buffer [100]rune
	n := 0
	for i := range r {
		buffer[n] = r[i]
		n++
	}
	return string(buffer[:n])
}

func parsePartRegex(r []rune, regexStack []*Regex, regexIterator *int) *Regex {
	// fmt.Println("PPR", string(r), *regexIterator)
	var regexType RegexType

	if len(r) == 1 {
		regexType = Rune
		reg := Regex{r[0], regexType, nil, nil}
		preg := new(Regex)
		*preg = reg
		wipe(r, 0, 1)
		return preg
	}

	if has(r, '*') {
		regexType = Any
	} else if has(r, '|') {
		regexType = Or
	} else {
		regexType = And
	}
	ret := new(Regex)
	count := 0
	swap := false
	for i := range(r) {
		if r[i] == 'a' || r[i] == 'b' {
			if i > 0 && r[i-1] == '(' {
				swap = true
			}
			count ++

			regexStack[*regexIterator] = parsePartRegex(r[i:i+1], regexStack,
				regexIterator)
			*regexIterator++
		}
	}

	if count == 1 && swap && regexType != Any {
		regexStack[*regexIterator-2], regexStack[*regexIterator - 1] = regexStack[*regexIterator-1], regexStack[*regexIterator - 2]
	}

	wipe(r, 0, len(r))
	switch regexType {

	case Any:
		*regexIterator --
		*ret = Regex{rune(0), regexType, regexStack[*regexIterator], nil}
		return ret
	default:
		*regexIterator -= 2
		// fmt.Println("RSI", *regexIterator)
		*ret = Regex{rune(0), regexType,
			regexStack[*regexIterator], regexStack[*regexIterator+1]}
		return ret
	}
}
func parseRegex(r []rune) *Regex {

	stack := make([]int, len(r))
	stack_it := 0

	var regexStack []*Regex = make([]*Regex, 100)
	regexStackIt := 0

	for i, c := range r {
		if c == '(' {
			stack[stack_it] = i
			stack_it ++
		} else if c == ')' {
			stack_it --
			// fmt.Println(string(r[stack[stack_it]+1:i]))
			regex := parsePartRegex(r[stack[stack_it]:i+1], regexStack, &regexStackIt)
			regexStack[regexStackIt] = regex
			regexStackIt++
			//wipe(r, stack[stack_it], i+1)
		}
	}

	return regexStack[0]
}

func _printRegexTree(regex *Regex, shift int) {
	if regex == nil {
		fmt.Println("Whoops!")
		return
	}
	switch regex.regexType {
	case Any:
		fmt.Println(strings.Repeat(" ", shift), "Any:")
		_printRegexTree(regex.left, shift + 2)
	case And:
		fmt.Println(strings.Repeat(" ", shift), "And:")
		_printRegexTree(regex.left, shift + 2)
		_printRegexTree(regex.right, shift + 2)
	case Or:
		fmt.Println(strings.Repeat(" ", shift), "Or:")
		_printRegexTree(regex.left, shift + 2)
		_printRegexTree(regex.right, shift + 2)
	case Rune:
		fmt.Println(strings.Repeat(" ", shift), regex.letter)
	default:
		fmt.Println("Whoops invalid type!")
	}
}

func printRegexTree(regex *Regex) {
	_printRegexTree(regex, 0)
}

func printGraph(graph *Graph) {
	for i:=0; i<graph.nodesLenght; i++ {
		fmt.Printf("%d: \n", i)
		for j, edge := range graph.nodes[i].edges {
			if edge != nil {
				fmt.Println(" ", AlphabetStrings[Alphabet(j)], edge.nodeNumber)
			}
		}
	}
}

func dfaToCountingMatrix(graph* Graph) [][]int64 {
	n := len(graph.nodes)
	ret := make([][]int64, n)
	for i := range ret {
		ret[i] = make([]int64, n)
	}

	for i, node := range graph.nodes {
		if node.edges[A] != nil {
			j := node.edges[A].nodeNumber
			ret[i][j] += 1
		}

		if node.edges[B] != nil {
			j := node.edges[B].nodeNumber
			ret[i][j] += 1
		}
	}
	return ret
}

func copyMatrix(toCopy [][]int64) [][]int64 {
	n := len(toCopy)
	ret := make([][]int64, n)
	for i := range ret {
		ret[i] = make([]int64, n)
	}

	for i:=0; i<n; i++ {
		for j:=0; j<n; j++ {
			ret[i][j] = toCopy[i][j]
		}
	}
	return ret
}

func idMatrix(n int) [][]int64 {
	ret := make([][]int64, n)
	for i := range ret {
		ret[i] = make([]int64, n)
	}

	for i:=0; i<n; i++ {
		ret[i][i] = 1
	}
	return ret
}

func isStringAccepted(start *Node, graph *Graph, str []Alphabet) bool {

	for _, c := range str {
		if start.edges[c] == nil {
			return false
		}
		start = start.edges[c]
	}
	if start.nodeType == End {
		return true
	}
	return false
}

func fastPower(A [][]int64, n int64) [][]int64 {
	t := n
	log := int64(0)
	for t > 0 {
		log += 1
		t /= 2
	}

	bitfield := make([]bool, log)
	matrices := make([][][]int64, log)

	ans := idMatrix(len(A))

	t = n
	i := 1
	matrices[0] = A
	bitfield[0] = (t % 2 == 1)
	t /= 2
	cur := A
	for t > 0 {
		if t % 2 == 1 {
			bitfield[i] = true
		}
		matrices[i] = matMul(cur, cur)
		cur = matrices[i]
		i += 1
		t /= 2
	}
	// fmt.Println(bitfield)
	for i, bit := range bitfield {
		if bit {
			ans = matMul(ans, matrices[i])
		}
	}
	return ans
}

func removeUnreachableNodes(graph* Graph, q0 *Node) (*Graph, map[int]*Node) {
	reachableStates := make(map[int]bool)
	newStates := make(map[int]bool)

	newStates[q0.nodeNumber] = true
	reachableStates[q0.nodeNumber] = true

	for len(newStates) > 0 {
		temp := make(map[int]bool)
		for i := range newStates {
			if graph.nodes[i].edges[A] != nil {
				temp[graph.nodes[i].edges[A].nodeNumber] = true
			}

			if graph.nodes[i].edges[B] != nil {
				temp[graph.nodes[i].edges[B].nodeNumber] = true
			}

		}

		for i := range reachableStates {
			delete(temp, i)
		}
		newStates = temp
		for i := range newStates {
			reachableStates[i] = true
		}
	}

	newGraphMap := make(map[int]*Node)
	newGraph := newGraph()

	for _, node := range graph.nodes {
		if reachableStates[node.nodeNumber] {
			newNode := newNode(newGraph)
			newGraphMap[node.nodeNumber] = newNode
		}
	}

	for _, node := range graph.nodes {
		if reachableStates[node.nodeNumber] {
			newNode := newGraphMap[node.nodeNumber]
			if node.edges[A] != nil {
				newNode.edges[A] = newGraphMap[node.edges[A].nodeNumber]
			}

			if node.edges[B] != nil {
				newNode.edges[B] = newGraphMap[node.edges[B].nodeNumber]
			}
		}
	}

	return newGraph, newGraphMap
}

func countStrings(r string, l int32) int32 {
	regexTree := parseRegex([]rune(r))
	nfa := regexToNFA(regexTree)
	// printGraph(nfa)
	dfa, startNodes, accNodes := nfaToDFA(nfa)

	//fmt.Println("Pre:", len(dfa.nodes))

	newDFA, newGraph := removeUnreachableNodes(dfa, dfa.nodes[startNodes[0]])
	var newStartNodes, newAccNodes []int

	for _, node := range startNodes {
		newStartNodes = append(newStartNodes, newGraph[node].nodeNumber)
	}

	for _, node := range accNodes {
		newNodeRef, ok := newGraph[node]
		if ok {
			newAccNodes = append(newAccNodes, newNodeRef.nodeNumber)
		}

	}
	dfa = newDFA
	startNodes, accNodes = newStartNodes, newAccNodes

	//fmt.Println("Post:", len(dfa.nodes))

	matrix := dfaToCountingMatrix(dfa)
	powMatrix := fastPower(matrix, int64(l))
	// fmt.Println(dfa.nodesLenght)
	// powMatrix := idMatrix(len(matrix))
	// for i := 0; i<int(l); i++ {
	//     powMatrix = matMul(powMatrix, matrix)
	// }

	// fmt.Println(powMatrix[startNodes[0]])

	ans := int64(0)
	for _, start := range startNodes {
		for _, end := range accNodes {
			ans += powMatrix[start][end] % MODULO
			ans %= MODULO
		}
	}

	// trialString := make([]Alphabet, l)
	// printableString := make([]rune, l)
	// test := []Alphabet{0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1, 0, 1}
	// fmt.Println(1 << l)
	// for i := 0; i< (1 << l); i++ {
	//     t := i
	//     j := 0
	//     for t > 0 {
	//         trialString[j] = Alphabet(t % 2)
	//         t /= 2
	//         j++
	//     }
	//     if isStringAccepted(dfa.nodes[startNodes[0]], dfa, trialString) {
	//         for i := range trialString {
	//             printableString[i] = rune(trialString[i] + 97)
	//         }
	//         fmt.Println(string(printableString))
	//     }
	// }

	// if !isStringAccepted(dfa.nodes[startNodes[0]], dfa, test) {
	//     // panic("string not accepted!!")
	// }

	//fmt.Println(ans)
	// printRegexTree(regexTree)
	// fmt.Println("~~~~~~~")
	// printGraph(nfaToDFA(nfa))
	// fmt.Println("-------")
	// printGraph(nfa)

	// fmt.Println("=======")
	return int32(ans)

}


func matMul(A, B [][]int64) [][]int64 {
	n:=len(A)
	C := make([][]int64, n)
	for i := range C {
		C[i] = make([]int64, n)
	}

	for i := 0; i<n; i++ {
		for j := 0; j<n; j++ {
			for k := 0; k<n; k++ {
				C[i][j] += (A[i][k]*B[k][j]) % MODULO
			}
			C[i][j] %= MODULO
		}
	}
	return C
}

func testMultiplication() {
	A := [][]int64 {{1, 2}, {3, 4}}
	B := [][]int64 {{7, 4}, {3, 11}}
	fmt.Println(matMul(A, B))
}

func makeTestGraph() *Graph {
	ret := newGraph()
	v1, v2, v3, v4 := newNode(ret), newNode(ret), newNode(ret), newNode(ret)
	v1.edges[A] = v2
	v2.epsilonEdges = []*Node{v1}
	v2.edges[B] = v3
	v4.edges[A] = v3
	v4.epsilonEdges = []*Node{v1}
	v1.edges[B] = v4
	return ret
}

func makeUnreachableTestGraph() *Graph {
	ans := makeTestGraph()
	newNode(ans)
	newNode(ans)
	return ans
}



func main() {
	reader := bufio.NewReaderSize(os.Stdin, 1024 * 1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 1024 * 1024)

	tTemp, err := strconv.ParseInt(readLine(reader), 10, 64)
	checkError(err)
	t := int32(tTemp)

	for tItr := 0; tItr < int(t); tItr++ {
		rl := strings.Split(readLine(reader), " ")

		r := rl[0]

		lTemp, err := strconv.ParseInt(rl[1], 10, 64)
		checkError(err)
		l := int32(lTemp)

		result := countStrings(r, l)

		fmt.Fprintf(writer, "%d\n", result)
	}

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
