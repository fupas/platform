# fupas/platform

This library is a lightweight golang runtime abstaction to build cloud-native services independent from the infrastructure the service is deployed on.

 The interfaces are inspired by the Google AppEngine and Google Cloud Platform APIs and first platform provider targets GCP only. Deployment in containers independent from a particular Cloud Provider will be next (TBD!).

In addition to provide basic service abstractions (e.g. logging) `fupas/platform` also implements "higher-functions" e.g. opinionated service authentication and authorization or a simple key/value store.

## Packages

| Packages               | Description                                     |
|------------------------|-------------------------------------------------|
| platform               | not sure ...                                    |
| observer               | Provides observability, i.e. logging and error reporting abstractions |
| analytics              | Reporting of 'events' e.g. page visits or e-commerce transactions     |
| services (deprecated)  | will be remove or replaced, don't bother |

## Dependencies

While `fupas/platform` SHOULD have as less dependencies to other libraries as possible, it does not re-invent the wheel. It uses, or "borrows" code from the following modules/repos/projects:

| Dependency                | Repo                               |
| --------------------------|------------------------------------|
| labstack/echo             | https://github.com/labstack/echo   |
| dapr/dapr                 | https://github.com/dapr/dapr       |
| dapr/components-contrib   | https://github.com/dapr/components-contrib |

## References

A list of useful posts, resources and references:

* https://brunoscheufler.com/blog/2019-04-26-choosing-the-right-go-web-framework
