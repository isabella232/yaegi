package interp

import (
	"fmt"
	"time"
)

// Function to run at CFG execution
type Builtin func(n *Node, f *Frame)

var builtin = [...]Builtin{
	Nop:          nop,
	ArrayLit:     arrayLit,
	Assign:       assign,
	AssignX:      assignX,
	Assign0:      assign0,
	Add:          add,
	And:          and,
	Call:         call,
	Case:         _case,
	CompositeLit: arrayLit,
	Dec:          nop,
	Equal:        equal,
	Greater:      greater,
	GetIndex:     getIndex,
	Inc:          inc,
	Land:         land,
	Lor:          lor,
	Lower:        lower,
	Mul:          mul,
	Not:          not,
	Range:        _range,
	Recv:         recv,
	Return:       _return,
	Send:         send,
	Sub:          sub,
}

var goBuiltin map[string]Builtin

func initGoBuiltin() {
	goBuiltin = make(map[string]Builtin)
	goBuiltin["make"] = _make
	goBuiltin["println"] = _println
	goBuiltin["sleep"] = sleep
}

// Run a Go function
func Run(def *Node, cf *Frame, recv *Node, rseq []int, args []*Node, rets []int, fork bool) {
	//println("run", def.Child[1].ident, "allocate", def.fsize)
	// Allocate a new Frame to store local variables
	anc := cf.anc
	if fork {
		anc = cf
	}
	f := Frame{anc: anc, data: make([]interface{}, def.fsize)}

	// Assign receiver value, if defined (for methods)
	if recv != nil {
		if rseq != nil {
			f.data[def.Child[0].findex] = valueSeq(recv, rseq, cf) // Promoted method
		} else {
			f.data[def.Child[0].findex] = value(recv, cf)
		}
	}

	// Pass func parameters by value: copy each parameter from caller frame
	// Get list of param indices built by FuncType at CFG
	paramIndex := def.Child[2].Child[0].val.([]int)
	i := 0
	for _, arg := range args {
		f.data[paramIndex[i]] = value(arg, cf)
		i++
		// Handle multiple results of a function call argmument
		for j := 1; j < arg.fsize; j++ {
			f.data[paramIndex[i]] = cf.data[arg.findex+j]
			i++
		}
	}

	// Execute the function body
	runCfg(def.Child[3].Start, &f)

	// Propagate return values to caller frame
	for i, ret := range rets {
		cf.data[ret] = f.data[i]
	}
}

// Functions set to run during execution of CFG

func value(n *Node, f *Frame) interface{} {
	switch n.kind {
	case BasicLit, FuncDecl:
		return n.val
	default:
		for level := n.level; level > 0; level-- {
			f = f.anc
		}
		//println(n.index, "val(", n.findex, "):", n.level, f.data[n.findex])
		if n.findex < 0 {
			//fmt.Println(n.index, "ident node value", n.val)
			return n.val
		}
		return f.data[n.findex]
	}
}

// Run by walking the CFG and running node builtin at each step
func runCfg(n *Node, f *Frame) {
	for n != nil {
		//fmt.Println(n.index, "runCfg run")
		n.run(n, f)
		if n.fnext == nil || value(n, f).(bool) {
			n = n.tnext
		} else {
			n = n.fnext
		}
	}
}

// assignX(n, f) implements assignement for a single call which returns multiple values
func assignX(n *Node, f *Frame) {
	//println(n.index, "in assignX")
	l := len(n.Child) - 1
	b := n.Child[l].findex
	for i, c := range n.Child[:l] {
		f.data[c.findex] = f.data[b+i]
	}
}

// assign(n, f) implements assignement with the same number of left and right values
func assign(n *Node, f *Frame) {
	l := len(n.Child) / 2
	for i, c := range n.Child[:l] {
		f.data[c.findex] = value(n.Child[l+i], f)
	}
}

// assign0(n, f) implements assignement of zero value
func assign0(n *Node, f *Frame) {
	l := len(n.Child) - 1
	z := n.typ.zero()
	for _, c := range n.Child[:l] {
		f.data[c.findex] = z
	}
}

func assignField(n *Node, f *Frame) {
	(*f.data[n.findex].(*interface{})) = value(n.Child[1], f)
}

func assignMap(n *Node, f *Frame) {
	f.data[n.findex].(map[interface{}]interface{})[value(n.Child[0].Child[1], f)] = value(n.Child[1], f)
}

func and(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) & value(n.Child[1], f).(int)
}

func not(n *Node, f *Frame) {
	f.data[n.findex] = !value(n.Child[0], f).(bool)
}

func _println(n *Node, f *Frame) {
	for i, m := range n.Child[1:] {
		if i > 0 {
			fmt.Printf(" ")
		}
		fmt.Printf("%v", value(m, f))

		// Handle multiple results of a function call argmument
		for j := 1; j < m.fsize; j++ {
			fmt.Printf(" %v", f.data[m.findex+j])
		}
	}
	fmt.Println("")
}

func call(n *Node, f *Frame) {
	//println(n.index, "call", n.Child[0].ident)
	// TODO: method detection should be done at CFG, and handled in a separate callMethod()
	var recv *Node
	var rseq []int
	if n.Child[0].kind == SelectorExpr {
		recv = n.Child[0].Child[0]
		rseq = n.Child[0].Child[1].val.([]int)
	}
	//fn := n.val.(*Node)
	fn := value(n.Child[0], f).(*Node)
	//fmt.Println("fn:", fn, f.data[n.Child[0].findex])
	var ret []int
	if len(fn.Child[2].Child) > 1 {
		if fieldList := fn.Child[2].Child[1]; fieldList != nil {
			ret = make([]int, len(fieldList.Child))
			for i, _ := range fieldList.Child {
				ret[i] = n.findex + i
			}
		}
	}
	Run(fn, f, recv, rseq, n.Child[1:], ret, false)
}

// Same as call(), but execute function in a goroutine
func callGoRoutine(n *Node, f *Frame) {
	//println(n.index, "call", n.Child[0].ident)
	// TODO: method detection should be done at CFG, and handled in a separate callMethod()
	var recv *Node
	var rseq []int
	if n.Child[0].kind == SelectorExpr {
		recv = n.Child[0].Child[0]
		rseq = n.Child[0].Child[1].val.([]int)
	}
	fn := n.val.(*Node)
	var ret []int
	if len(fn.Child[2].Child) > 1 {
		if fieldList := fn.Child[2].Child[1]; fieldList != nil {
			ret = make([]int, len(fieldList.Child))
			for i, _ := range fieldList.Child {
				ret[i] = n.findex + i
			}
		}
	}
	go Run(fn, f, recv, rseq, n.Child[1:], ret, false)
}

func getIndexAddr(n *Node, f *Frame) {
	a := value(n.Child[0], f).(*[]interface{})
	f.data[n.findex] = &(*a)[value(n.Child[1], f).(int)]
}

func getIndex(n *Node, f *Frame) {
	a := value(n.Child[0], f).(*[]interface{})
	f.data[n.findex] = (*a)[value(n.Child[1], f).(int)]
}

func getIndexMap(n *Node, f *Frame) {
	m := value(n.Child[0], f).(map[interface{}]interface{})
	f.data[n.findex] = m[value(n.Child[1], f)]
}

func getMap(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(map[interface{}]interface{})
}

func getIndexSeq(n *Node, f *Frame) {
	a := value(n.Child[0], f).(*[]interface{})
	seq := value(n.Child[1], f).([]int)
	l := len(seq) - 1
	for _, i := range seq[:l] {
		a = (*a)[i].(*[]interface{})
	}
	f.data[n.findex] = (*a)[seq[l]]
}

func valueSeq(n *Node, seq []int, f *Frame) interface{} {
	a := f.data[n.findex].(*[]interface{})
	l := len(seq) - 1
	for _, i := range seq[:l] {
		a = (*a)[i].(*[]interface{})
	}
	return (*a)[seq[l]]
}

func mul(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) * value(n.Child[1], f).(int)
}

func add(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) + value(n.Child[1], f).(int)
}

func sub(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) - value(n.Child[1], f).(int)
}

func equal(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f) == value(n.Child[1], f)
}

func inc(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) + 1
}

func greater(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) > value(n.Child[1], f).(int)
}

func land(n *Node, f *Frame) {
	if v := value(n.Child[0], f).(bool); v {
		f.data[n.findex] = value(n.Child[1], f).(bool)
	} else {
		f.data[n.findex] = v
	}
}

func lor(n *Node, f *Frame) {
	if v := value(n.Child[0], f).(bool); v {
		f.data[n.findex] = v
	} else {
		f.data[n.findex] = value(n.Child[1], f).(bool)
	}
}

func lower(n *Node, f *Frame) {
	f.data[n.findex] = value(n.Child[0], f).(int) < value(n.Child[1], f).(int)
}

func nop(n *Node, f *Frame) {}

func _return(n *Node, f *Frame) {
	for i, c := range n.Child {
		f.data[i] = value(c, f)
	}
}

// create an array of litteral values
func arrayLit(n *Node, f *Frame) {
	a := make([]interface{}, len(n.Child)-1)
	for i, c := range n.Child[1:] {
		a[i] = value(c, f)
	}
	f.data[n.findex] = &a
}

// Create a map of litteral values
func mapLit(n *Node, f *Frame) {
	m := make(map[interface{}]interface{})
	for _, c := range n.Child[1:] {
		m[value(c.Child[0], f)] = value(c.Child[1], f)
	}
	f.data[n.findex] = m
}

// Create a struct object
func compositeLit(n *Node, f *Frame) {
	l := len(n.typ.field)
	a := n.typ.zero().(*[]interface{})
	for i := 0; i < l; i++ {
		if i < len(n.Child[1:]) {
			c := n.Child[i+1]
			(*a)[i] = value(c, f)
			//println(n.index, "compositeLit, set field", i, value(c, f))
		} else {
			(*a)[i] = n.typ.field[i].typ.zero()
		}
	}
	f.data[n.findex] = a
}

// Create a struct Object, filling fields from sparse key-values
func compositeSparse(n *Node, f *Frame) {
	a := n.typ.zero().(*[]interface{})
	for _, c := range n.Child[1:] {
		// index from key was pre-computed during CFG
		(*a)[c.findex] = value(c.Child[1], f)
	}
	f.data[n.findex] = a
}

func _range(n *Node, f *Frame) {
	i, index := 0, n.Child[0].findex
	if f.data[index] != nil {
		i = f.data[index].(int)
	}
	a := value(n.Child[2], f).(*[]interface{})
	if i >= len(*a) {
		f.data[n.findex] = false
		return
	}
	f.data[index] = i + 1
	f.data[n.Child[1].findex] = (*a)[i]
	f.data[n.findex] = true
}

func _case(n *Node, f *Frame) {
	f.data[n.findex] = value(n.anc.anc.Child[0], f) == value(n.Child[0], f)
}

// Allocates and initializes a slice, a map or a chan.
func _make(n *Node, f *Frame) {
	typ := value(n.Child[1], f).(*Type)
	switch typ.cat {
	case ArrayT:
		f.data[n.findex] = make([]interface{}, value(n.Child[2], f).(int))
	case ChanT:
		f.data[n.findex] = make(chan interface{})
	case MapT:
		f.data[n.findex] = make(map[interface{}]interface{})
	}
}

// Read from a channel
func recv(n *Node, f *Frame) {
	f.data[n.findex] = <-value(n.Child[0], f).(chan interface{})
}

// Write to a channel
func send(n *Node, f *Frame) {
	value(n.Child[0], f).(chan interface{}) <- value(n.Child[1], f)
}

// Temporary, for debugging purppose
func sleep(n *Node, f *Frame) {
	duration := time.Duration(value(n.Child[1], f).(int))
	time.Sleep(duration * time.Millisecond)
}