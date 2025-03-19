#Clone the repo
#Execute the following commands 
1. go get github.com/Dineshkumarramaraj/mcp-kubernetes/kubernetes-listpods
2. go install github.com/Dineshkumarramaraj/mcp-kubernetes/kubernetes-listpods

#Add the below content in "claude_desktop_config.json"

{
  "mcpServers": {
    "kubernetes_listpods": {
      "command": "kubernetes-listpods",
      "args": []
    }
  }
}

Note: In command, try to give pull path of go bin eg: ~/go/bin/kubernetes-listpods
![Image](https://github.com/user-attachments/assets/3be7f196-b042-4941-a7ce-060405e00359)

#If you dont have claude desktop installed on your laptop, then try running below command in developer mode

1. npx @modelcontextprotocol/inspector go run main.go
2. Try "localhost:5173" in your browser. 

