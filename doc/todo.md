# MVP - Backlog

[x] Finish end-to-end the integration of the scheduler within the agent
[x] Integrate the config file into the bootstrap process
[x] Develop a simple CLI:
    [x] Run node as agent 
    [x] Schedule all jobs into cluster
    [-] See basic stats

# BUGS

[ ] bug1: no more than one job can be scheduled at the time
[ ] bug2: only wait-index=0 works when listening for new jobs, i.e. no past jobs can be re-scheduled
[ ] bug3: cron parser does not look to work very accurately
[ ] bug4: jobs can be overwriden easily. Need to compare job ids in the enqueue function of the scheduler
[x] bug5: nice to have a cli command to purge all registered jobs:

 ```
    xchronos rm --id=payment
    xchronos rm --all
 ```   