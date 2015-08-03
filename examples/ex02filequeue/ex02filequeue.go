package main

import (
	sci "github.com/samuell/scipipe"
)

func main() {
	sci.InitLogAudit()

	fq := sci.FQ("hej1.txt", "hej2.txt", "hej3.txt")
	fw := sci.Sh("echo {i:in} > {o:out}")
	fw.OutPathFuncs["out"] = func(t *sci.ShellTask) string { return t.GetInPath("in") }
	sn := sci.NewSink()

	fw.InPorts["in"] = fq.Out
	sn.In = fw.OutPorts["out"]

	pl := sci.NewPipeline()
	pl.AddProcs(fq, fw, sn)

	pl.Run()
}
