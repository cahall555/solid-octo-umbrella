package tools

const ToolGetActiveAlerts = `{
  "type": "function",
  "function": {
    "name": "get_active_alerts",
    "description": "Get active NWS alerts impacting a marine region.",
    "parameters": {
      "type": "object",
      "properties": {
        "region": {
          "type": "string",
          "description": "Marine region code (AL, AT, GL, GM, PA, PI)"
        }
      },
      "required": ["region"]
    }
  }
}`
