package main

import (
	"bytes"
	sci "github.com/scipipe/scipipe"
)

func main() {
	sci.InitLogAudit()

	foo := NewFooer()
	f2b := NewFoo2Barer()
	snk := sci.NewSink()

	foo.OutFoo.Connect(f2b.InFoo)
	snk.Connect(f2b.OutBar)

	pl := sci.NewPipelineRunner()
	pl.AddProcesses(foo, f2b, snk)
	pl.Run()
}

// ------------------------------------------------------------------------
// Components
// ------------------------------------------------------------------------

// Fooer

type Fooer struct {
	InnerProcess *sci.SciProcess
	OutFoo       *sci.FilePort
}

func NewFooer() *Fooer {
	// Initiate task from a "shell like" pattern, though here we
	// just specify the out-port foo
	innerFoo := sci.NewFromShell("fooer", "{o:foo}")
	// Set the output formatter to a static string
	innerFoo.SetPathStatic("foo", "foo.txt")
	// Create the custom execute function, with pure Go code
	innerFoo.CustomExecute = func(task *sci.SciTask) {
		task.OutTargets["foo"].WriteTempFile([]byte("foo\n"))
	}
	// Connect the ports of the outer task to the inner, generic one
	fooer := &Fooer{
		InnerProcess: innerFoo,
		OutFoo:       sci.NewFilePort(),
	}
	return fooer
}

func (p *Fooer) Run() {
	// Connect inner ports to outer ones again, in order to update
	// connectivity after the workflow wiring has taken place.
	p.InnerProcess.Out["foo"] = p.OutFoo
	// Run the inner process
	p.InnerProcess.Run()
}

func (p *Fooer) IsConnected() bool {
	return p.OutFoo.IsConnected()
}

// Foo2Barer

type Foo2Barer struct {
	InnerProcess *sci.SciProcess
	InFoo        *sci.FilePort
	OutBar       *sci.FilePort
}

func NewFoo2Barer() *Foo2Barer {
	// Initiate task from a "shell like" pattern, though here we
	// just specify the in-port foo and the out-port bar
	innerProc := sci.NewFromShell("foo2bar", "{i:foo}{o:bar}")
	// Set the output formatter to extend the path on the "bar"" in-port
	innerProc.SetPathExtend("foo", "bar", ".bar.txt")
	// Create the custom execute function, with pure Go code
	innerProc.CustomExecute = func(task *sci.SciTask) {
		task.OutTargets["bar"].WriteTempFile(bytes.Replace(task.InTargets["foo"].Read(), []byte("foo"), []byte("bar"), 1))
	}

	// Connect the ports of the outer task to the inner, generic one
	return &Foo2Barer{
		InnerProcess: innerProc,
		InFoo:        sci.NewFilePort(),
		OutBar:       sci.NewFilePort(),
	}
}

func (p *Foo2Barer) Run() {
	// Connect inner ports to outer ones again, in order to update
	// connectivity after the workflow wiring has taken place.
	p.InnerProcess.In["foo"] = p.InFoo
	p.InnerProcess.Out["bar"] = p.OutBar
	// Run the inner process
	p.InnerProcess.Run()
}

func (p *Foo2Barer) IsConnected() bool {
	return p.InFoo.IsConnected() &&
		p.OutBar.IsConnected()
}
