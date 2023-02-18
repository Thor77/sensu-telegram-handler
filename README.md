# Sensu Go Telegram Handler

**Looking for new maintainer(s)**. As I'm not using Sensu for my monitoring needs anymore I'll only provide security updates for this project.

## Installation

```
go build -o /usr/bin/sensu-telegram-handler .
```

## Configuration

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "telegram"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-telegram-handler --api-token <bot api token> --chatid <your chat id>",
        "timeout": 10,
        "filters": [
            "is_incident",
            "not_silenced"
        ]
    }
}
```

## Usage

```
Sensu Go handler for sending telegram notifications

Usage:
  sensu-telegram-handler [flags]

Flags:
  -a, --api-token string
  -c, --chatid uint
  -h, --help               help for sensu-telegram-handler
  -t, --template string     (default "**{{.Entity.Name}}/{{.Check.Name}}**: {{.Check.State}}\n`{{.Check.Output}}`")
```

