package cli

import (
	"fmt"
	"os"
)

// osExit is a variable that can be mocked in tests
var osExit = os.Exit

const helpText = `neo4j-aura-mcp - A Model Context Protocol Server to manage Neo4j Aura DB infrastructure

Usage:
  neo4j-aura-mcp [OPTIONS]

Options:
  -h, --help      Show this help message
  -v, --version   Show version information

Environment Variables:

	AURA_API_URI          The URI of the Aura API  - OPTIONAL - ( Default: https://api.neo4j.io/v1 )
	AURA_API_ID           Client Id used to obtain an token that is then used to acccess the Aura API - REQUIRED - ( Default: N/A ) 
	AURA_API_SECRET       Client Secret used to obtain an token that is then used to acccess the Aura API - REQUIRED - ( Default: N/A ) 
    AURA_API_READ_ONLY    Prevents updates and deletions to Aura infrastructure.  Set to False to allow that. Aura API - OPTIONAL - ( Default: True ) 
  

Examples:
  AURA_API_ID="asdkfjsdhjkgsdf-2nsndfgskdfsg"
  AURA_API_SECRET="dfkbdjdf98er598ds-98vnasns-12312309smn-bht"



`

func HandleArgs(version string) {
	if len(os.Args) <= 1 {
		return
	}

	flags := make(map[string]bool)
	var err error

	for _, arg := range os.Args[1:] {
		switch arg {
		case "-h", "--help":
			flags["help"] = true
		case "-v", "--version":
			flags["version"] = true
		default:
			err = fmt.Errorf("unknown flag or argument: %s", arg)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		osExit(1)
	}

	if flags["help"] {
		fmt.Print(helpText)
		osExit(0)
	}

	if flags["version"] {
		fmt.Printf("neo4j-aura-mcp version: %s\n", version)
		osExit(0)
	}
}
