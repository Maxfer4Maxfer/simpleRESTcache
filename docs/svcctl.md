# svcctl

**svcclt** is a simpleRESTcache service management tool.
**svcclt** is command line interface for control a simpleRESTcache instance.

## Build 
```bash
go built ./srcctl/main.go
```

## CLI arguments
```bash
  -h string
    	A address of a simpleRestCache instance (default "srcsvc")
  -p int
    	A control port of a simpleRestCache instance (default 8081)
```

## Usage
    srcctl -h <host> -p <port> COMMAND
    |
    |- stat 	        # Display statistic of cache usage
    |   |- all		    # Display all from cache
    |   |- top <N>	    # Display top <N> popular requests to cache. <N> number 
    |   |- last <N>	    # Display last <N> unpopular requests to cache. <N> number
    |
    |- cache	Manage cache
    |   |- all		    # Display all from cache
    |   |- clean		Delete all cache records
    |   |- refresh		Refresh all cache records
    |
    |- settings	Display settings of a cache system