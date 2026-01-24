This package manages a pool of working proxies to keep your requests from getting their 30 min cooldown limit on requessting songs

* **Refresh:** Pulls fresh list of ips from github and tests them to see if they are actually alive
* **Validation:** Uses 50 workers at once to check up to 300 proxies quickly so you don't wait forever (this will be changed to using only na proxies)
* **GetRandomClient:** Picks proxy at random and gives it a fresh cookiejar so every request looks like a new user
* **Proxy Check:** Every proxy is tested against a google endpoint if it's too slow or dead, it wont be used
