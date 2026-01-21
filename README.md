# r-a-d.io-infinite

Barebones implementation for infinite [r-a-d.io](https://r-a-d.io/search) song requests bypassing the 30 minute request limit (please use responsibly), this cli simply [gets proxies](https://raw.githubusercontent.com/monosans/proxy-list/main/proxies/http.txt) > reads the site source code and scrape the amazing one time [gorilla/csrf](https://github.com/gorilla/csrf) token and tricks the server to accept your song requests. **qol features may be worked on soon.**

- utilizes go's awesome concurrency to verify proxies
- stateless networking net/http which disposes each request making it fresh
- regex to get csrf token and store it in cookiejar
- inject POST request with stolen token, song id and spoofed Referer headers

## Disclaimer

This tool is for educational and research purposes only. By using this software, you acknowledge that:

sethispr does not condone or support the malicious spamming of community-run services.

You assume all risks associated with bypassing rate limits or automated scraping, sethispr is not responsible for any IP bans, blacklisting, or legal consequences.
