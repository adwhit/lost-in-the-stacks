set -e

echo "Running with Go 1.1\n"

echo "Stack:"
GOROOT=/home/alex/sources/go1.1/ go1.1 run -gcflags -l stack.go

echo "\nSnippet:"
GOROOT=/home/alex/sources/go1.1/ go1.1 run -gcflags -l snippet.go

echo "\nBench:"
GOROOT=/home/alex/sources/go1.1/ go1.1 run -gcflags -l split.go

echo "\nRunning with Go 1.3"
echo "\nBench:"
GOROOT=/home/alex/sources/go1.3/ go1.3 run -gcflags -l split.go
