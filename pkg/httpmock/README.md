# HTTPMock

HTTPMock is composed of:

- a Go HTTP Mock library for writing tests,
- an HTTP scenario recording tool.

This library is used to record scenario that can be replayed locally, exactly as recorded.
It is helpful to prepare tests for your Go code when there is underlying HTTP requests. You can thus be sure
the test will replay real traffic consistently. Those recorded tests are also useful to document the behaviour 
of your code against a specific version of API or content. When the target HTTP endpoint changes and breaks your
code, you can thus now easily generate a diff of the HTTP content to understand the change in behaviour and
adapt your code accordingly.
