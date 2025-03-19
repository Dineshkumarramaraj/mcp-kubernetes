npx @modelcontextprotocol/inspector go run main.go

go install .

{
  "mcpServers": {
    "mcp_kubernetes_listpods": {
      "command": "mcp-kubernetes-listpods",
      "args": []
    }
  }
}

Note: In command, try to give pull path of go bin eg: ~/go/bin/mcp-kubernetes-listpods

go get github.com/Dineshkumarramaraj/mymcp/mcp-kubernetes-listpods
go install github.com/Dineshkumarramaraj/mymcp/mcp-kubernetes-listpods
