This package manages a pool of working proxies to keep your requests from getting their 30 min cooldown limit on requessting songs

* **Refresh:** Pulls fresh list of ips from github and tests them to see if they are actually alive
* **Validation:** Uses 50 workers at once to check up to 300 proxies quickly so you don't wait forever (this will be changed to using only na proxies)
* **GetRandomClient:** Picks proxy at random and gives it a fresh cookiejar so every request looks like a new user
* **Proxy Check:** Every proxy is tested against a google endpoint if it's too slow or dead, it wont be used

Some current problems is that it relies on global state with heavy locking, which causes unnecessary blocking and makes it hard to test. It also basically recreates new HTTP clients and transports for every proxy check instead of reusing them, wasting resources and slowing everything down. Error handling is mostly ignored, hiding problems and making debugging difficult. Also the google proxy validation is kinda useless. The fixed number of workers is not adaptable so it'll either overload or not use ur cores properly.
