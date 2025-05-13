module main

go 1.24.3

require (
	github.com/gorilla/mux v1.8.1 
	github.com/rs/cors v1.11.1 

	bfs v0.0.0
	dfs v0.0.0
	bidirectional v0.0.0
)

replace (
	bfs => ../bfs
	dfs => ../dfs
	bidirectional => ../bidirectional
)
