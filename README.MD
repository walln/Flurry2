# Flurry A simple HTTP Proxy for Microservices

## Description

Flurry is a L7 application proxy designed to connect, monitor, and secure http requests. All requests to a flurry server are proxied upstream to the specified server. This enables highly concurrent traffic routing and middleware for your applications.

### Configuring

To configure a flurry server, run `flurry.exe` and pass the `-f` flag with the path to the config file. To create a config file see [configuring a flurry server](#creating-a-flurry-config-file). Flurry can also be added to any existing projects simply by including the flurry package.

```go

    import "github.com/walln/flurry2"

    func main() {
        server := flurry.Initialize()
	    server.ListenAndServe()
    }
```

Flurry listens on the port environment variable and automatically intializes a proxy server based on the specifcation in the config file. A config file is passed to the binary through the `-f` flag.

```bash
    ./flurry -f flurry-config.yaml
```

---

## Features

Flurry is a competent proxy gatway with many features for managing microservices.

- Traffic Monitoring and Logging
- Dynamic configuring of different routes
- Supports authentication plugins
- Packaged as a single self contained binary

---

### Authentication

Authentication schemas are entirely optional with flurry but can enable developers to control the access of their services granularly.
Flurry supports JWT, Firebase, and AWS Cognito...(soon). To Enable any of these authentication schemas, set the flurry `authentication` mode in the config file to either `JWT`, `FIREBASE`, `COGNITO`, or `NONE`.

#### JWT

Authenticating requests with JSON Web Tokens requires that the config file has the `signingMethod` filed set with a valid schema. Flurry supports:

- HMAC
- ECDSA
- RSA
- RSAPSS

In addtion to setting the signing method, the `JWT_SECRET` environment variable must be set.

#### Firebase Authentication

Flurry supports Firebase Authentication out of the box. If the `authentication` field in the config file is set to `FIREBASE`, the only other necessary component is the path to the GCP service acccount used to validate the token. This can be specified using the `GOOGLE_APPLICATION_CREDENTIALS` environment variable. If you happen to be deploying flurry on Google Cloud Platform, this should already be set for you.

---

## Creating a Flurry Config File

Flurry generates proxy routes and policies based off of a YAML file passed in at runtime.

A flurry config file is composed of two main parts, the Server configuration, and the proxy configuration.

#### Server Configuration

Flurry configurations are primarily used for enabling plugins and specifiyng plugin-specific information. Additionally server configurations are used for logging purposes and monitoring.

This exmple server configuration outlines simple flurry server with no authentication.

```yaml
name: 'Flurry Server'
authentication: NONE
```

Most server configuration values are optional and are not required for flurry to function. However, it is recommended to fill in these values as it can be useful for documentation and logging purposes. For instance, the `authentication` field is not required. The `NONE` value is implied if no value is supplied.

If using JWT authentication an additional field is required, the `signingMethod` field specifies the cryptographic method used to sign the token.

```yaml
---
name: Flurry Server
authentication: JWT
signingMethod: HMAC
```

#### Proxy Configuration

The primary content of a Flurry Config file is in the proxy settings. A flurry proxy is a predfined route that forwards all traffic to another host. Proxy configs are defined as an array of routes.

```yaml
routes:
  - endpoint: '/users/user/{id}'
    host: 'https://example.com'
    authenticated: true
    methods:
      - GET
      - POST
      - PUT
      - DELETE
  - endpoint: /users
    host: 'https://example.com'
    authenticated: false
    methods:
      - GET
```

When defining proxy routes a route on the flurry host and the declared proxy are linked. For example:

```bash
    https://my-flurry-server.io/users
        ->
    https://example.com/user
```

A flurry server cannot have multiple routes to the same endpoint. Design routes as if you are composing an API from multiple hosts. Flurry was created to enable APIs to be composed across distributed architectures.

Each proxy route is defined by four fields, the `endpoint`, the `host`, `authentication`, and the enabled `methods`.

##### Endpoint

The endpoint value specifies the routes to link on both the proxy and the flurry server.

```yaml
endpoint: /route/{param}
```

Endpoints can have request parameters embeded, variables in the request URL are resolved by Flurry and passed to the upstream server.

```bash
    https://flurry-example-server.io/catalog/{id}
```

Consider this request, any requests to this endpoint require an id value to be passed through.

```bash
    curl https://flurry-example-server.io/catalog/1
```

Requests need only a parameter name in the config file to forward these values. The specify a request parameter in a Flurry Config, surrond the variable name with braces. - `{foo}`

##### Host

Every proxy route requires a host to be defined, including the host schema. Upstream requests format the requests from flurry according to the HTTP specs to it is important to include this value.

```yaml
host: "https://google.com" #Valid

host: "http://example.com" #Valid

host: "my-example-site.com" #Not Valiid
```

Flurry proxy servers can support any number of hosts including other proxy servers. This can be useful for load balancing requests through successive gateways.

##### Methods

Flurry locks down forwarded requests based on the request methods that are accepted on the upstream server. To specify acceptable methods for an endpoint, set the `methods` field to an array of methods. Acceptable methods are outlined in RFC 2616 section 9.

```yaml
methods:
  - GET
  - DELETE
  - HEAD
```

Attempting to access a proxy route with a disallowed method will result in a 405 Method Not Allowed response code from the flurry server, blocking upstream traffic.

#### Authenticated

Every proxy route on a Flurry server has granular control of authentication and access policy. Proxy routes inherit their authentication methods from the Server Config. Authentication is not necessary but it is best-practice to set the field value for every proxy-route.

```yaml
authenticated: false
```
