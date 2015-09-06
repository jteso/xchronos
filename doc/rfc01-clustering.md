# RFC01 - Clustering

## Scheduler election
key: `/zeitd/lock/scheduler`

value: `{agentIp}`

## Job offers
key:   `/zeitd/offers/<job_id>`

value: `job{}` 

## Job takers
key:  `/zeitd/lock/<job_id>`

value: `{executorIp}`

## Executors List (**)
key:   `/zeitd/executors/<agent_id>`

value: `agent{}`

## Job List
key:   `/zeitd/jobs/<job_id>`

value: `job{}`








