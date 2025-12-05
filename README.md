# A MCP Server that allows a LLM / Agent to managed Neo4j Aura DB infrastructure

Official Model Context Protocol (MCP) server for Neo4j.

## Status

BETA - Active development; not yet suitable for production.

## Prerequisites

- Registered to use Neo4j Aura DB
- Client Id and secret to use the Aura API


## Configure VSCode (MCP)

Create / edit `mcp.json` (docs: https://code.visualstudio.com/docs/copilot/customization/mcp-servers):

```json
{
  "servers": {
    "neo4j": {
      "type": "stdio",
      "command": "neo4j-aura-mcp",
      "env": {
        "AURA_API_URI": "https://api.neo4j.io/v1",
        "AURA_API_ID":  "kjsdghkvp8w3rvaser",
        "AURA_API_SECRET": "asdkjhfdk-adsfhshdf-1238sd8-afjkajdfkjsaf",
        "AURA_API_READ_ONLY":  "True",
      }
    }
  }
}
```

Restart VSCode; open Copilot Chat and ask: "List Neo4j MCP tools" to confirm.

## Configure Claude Desktop

First, make sure you have Claude for Desktop installed. [You can install the latest version here](https://claude.ai/download).

We’ll need to configure Claude for Desktop for whichever MCP servers you want to use. To do this, open your Claude for Desktop App configuration at:

- (MacOS/Linux) `~/Library/Application Support/Claude/claude_desktop_config.json`
- (Windows) `$env:AppData\Claude\claude_desktop_config.json`

in a text editor. Make sure to create the file if it doesn’t exist.

You’ll then add the `neo4j-mcp` MCP in the mcpServers key:

```json
{
  "mcpServers": {
    "neo4j-mcp": {
      "type": "stdio",
      "command": "neo4j-aura-mcp",
      "args": [],
      "env": {
        "AURA_API_URI": "https://api.neo4j.io/v1",
        "AURA_API_ID":  "kjsdghkvp8w3rvaser",
        "AURA_API_SECRET": "asdkjhfdk-adsfhshdf-1238sd8-afjkajdfkjsaf",
        "AURA_API_READ_ONLY":  "True",
      }
    }
  }
}
```

__Notes:__
- Adjust env vars for your setup.  You must supply values for AURA_API_ID and AURA_API_SECRET
- Set `AURA_API_READ_ONLY=false` to enable changes to be made to existing Aura infrastructure

## Tools & Usage

Provided tools:

| Tool                  | ReadOnly | Purpose                                              | Notes                                                                                                                          |
| --------------------- | -------- | ---------------------------------------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
| `list-instances`      | `true`   | Lists all instances in Aura                          |                                                                                 |


### Readonly mode flag

To avoid a LLM / Agent having a 'moment' the MCP server cannot make changes to existing infrastructure.  If you want to do that, set  AURA_API_READ_ONLY to false

** WARNING - This is at your own risk.  LLM / Agents can be unpredictable. 

  

## Example Natural Language Prompts

Below are some example prompts you can try in Copilot or any other MCP client:

- "List all instances in Aura"


## Disclaimer
Although based on a fork of Neo4j official MCP server, this has no association with Neo4j.  It's a pet project of mine. 

** WARNING - Use at your own risk.  LLM / Agents can be unpredictable. 


