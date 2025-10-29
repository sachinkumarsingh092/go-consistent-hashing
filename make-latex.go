package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

var tikzTemplate = template.Must(template.New("t").Parse(`
\documentclass[tikz]{standalone}
\usepackage{tikz}
\usetikzlibrary{arrows.meta}

\definecolor{Oc}{RGB}{0, 0, 0}

\begin{document}
\begin{tikzpicture}[>=Stealth, line cap=round, line join=round, scale=2.8]
  \pgfextra{\path[use as bounding box] (-1.3,-1.3) rectangle (1.3,1.3);}
  \def\R{1.0}
  \draw[thick] (0,0) circle (\R);

  % Angles array
  \def\angles#1{\ifcase#1 {{.Angles}}\fi}

  % Labels array
  \def\labels#1{\ifcase#1 {{.Labels}}\fi}

  \foreach \i in {0,...,{{.RangeEnd}}}{
    \pgfmathsetmacro{\ang}{\angles{\i}}
    \fill[Oc] (\ang:\R) circle (1.1pt);
    \node[red, font=\small, anchor=south east] at (\ang:\R) {\labels{\i}};
  }

\end{tikzpicture}
\end{document}
`))

func makeNodesCircleTikz() {
	ch := NewConsistentHasher(1024 * 1024)

	nodes := []string{"lubna", "luboo", "booboo", "honeybee", "sugar", "kashmiri-apple", "lubooboo", "honey"}
	for _, node := range nodes {
		if err := ch.AddNode(node); err != nil {
			panic(err)
		}
	}

	var angles []string

	// For each node index, compute the angle it represents. ringSize is 360
	// degrees.
	for _, s := range ch.slots {
		fraction := float32(s) / float32(ch.nslots)
		angle := fraction * 360.0
		angles = append(angles, fmt.Sprintf("%d", int(angle)))
	}
	f, err := os.Create("labels_circle.tex")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tikzTemplate.Execute(f, map[string]string{
		"Angles":   strings.Join(angles, `\or `),
		"Labels":   strings.Join(nodes, `\or `),
		"RangeEnd": fmt.Sprintf("%d", len(angles)-1),
	})
}

func main() {
	makeNodesCircleTikz()
}
