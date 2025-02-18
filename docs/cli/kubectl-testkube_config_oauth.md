## kubectl-testkube config oauth

Set oauth credentials for api uri in testkube client

```
kubectl-testkube config oauth <value> [flags]
```

### Options

```
      --client-id string       client id for authentication provider
      --client-secret string   client secret for authentication provider
  -h, --help                   help for oauth
      --provider string        authentication provider, currently available: github (default "github")
      --scope stringArray      scope for authentication provider
```

### Options inherited from parent commands

```
  -a, --api-uri string      api uri, default value read from config if set
  -c, --client string       client used for connecting to Testkube API one of proxy|direct (default "proxy")
      --namespace string    Kubernetes namespace, default value read from config if set (default "testkube")
      --oauth-enabled       enable oauth
      --telemetry-enabled   enable collection of anonumous telemetry data
      --verbose             show additional debug messages
```

### SEE ALSO

* [kubectl-testkube config](kubectl-testkube_config.md)	 - Set feature configuration value

