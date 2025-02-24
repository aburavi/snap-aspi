# ratelimiter

1. fix windows
key per menit, increment it and set expire
2. Sliding Log
use sorted set instead of set
3. sliding windows
every second in unix increment,count 1 menit (max time) and compare with max limit


source:
https://redis.com/redis-best-practices/basic-rate-limiting/
https://medium.com/@NlognTeam/design-a-scalable-rate-limiting-algorithm-system-design-nlogn-895abba44b77
