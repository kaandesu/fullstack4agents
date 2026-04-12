# Agent Coordination

## How it works

There is always exactly **one file** in this folder (besides this README), named with an ISO timestamp: `20260412T120000Z.md`

The filename IS the version. Agents keep track of the last timestamp they read.

**At the start of every session:**
1. `ls .coordination/` — read the filename
2. If the filename differs from your last remembered timestamp → read the file
3. If it matches → skip reading, nothing changed

**After completing work:**
1. Write a new file with the current timestamp as the name (compact ISO: `YYYYMMDDTHHmmssZ`)
2. Delete the old timestamped file

## File format (keep it short — max ~15 lines)

```
by: <agent-name>
at: <ISO timestamp>

did:
- short bullet of what was completed

next:
- what this agent plans to do next (optional)

note:
- anything the other agent must know (API changes, blockers, decisions)
```

Only write what the other agent actually needs to know. Omit empty sections entirely.
