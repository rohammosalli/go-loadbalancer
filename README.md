## Loadbalancer Overview

Simple load balancer / reverse proxy written in GoLang.

You can configure your LB backend in `config.json` file

```json
{
    "URLS": ["http://192.168.10.1","http://192.168.10.2","http://192.168.10.3"]
}
```
***How to run the Loadbalancer***

Run the LB
```bash
make start_lb
```
Stop the LB
```bash
make stop_lb
```
clean the build 
```bash 
make clean 
```
