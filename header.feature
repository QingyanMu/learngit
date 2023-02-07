@annotations @router.cookie @release-1.22
Feature: Header
  An Ingress may define header rules in its annotations.

  If a Request matches a rule defined in a Ingress which have rewrite-header annotations,
  BFE should modify the name and value of the header field defined
  Scenario: An Ingress with annotation `bfe.ingress.kubernetes.io/rewrite-header.actions`
    Given an Ingress resource with rewrite annotation
    """
    apiVersion: networking.k8s.io/v1
    kind: Ingress
    metadata:
      name: rewrite-header-field
      namespace: ingress-bfe
      annotations:
        bfe.ingress.kubernetes.io/rewrite-header.actions: >-
          [
            {
              "cmd": "REQ_HEADER_SET",
              "params": ["key1", "value1"]
            },
            {
              "cmd": "REQ_HEADER_SET",
              "params": ["key2", "value2"]
            },
            {
              "cmd": "REQ_HEADER_SET",
              "params": ["key3", "value3"]
            },
            {
              "cmd": "REQ_HEADER_ADD",
              "params": ["key2", "_addvalue"]
            },
            {
              "cmd": "REQ_HEADER_DEL",
              "params": ["key3"]
            }
          ]
    spec:
      rules:
        - host: "foo.com"
          http:
            paths:
              - path: /bar
                pathType: Prefix
                backend:
                  service:
                    name: foo-exact
                    port:
                      number: 3000
    """
    And The Ingress status shows the IP address or FQDN where it is exposed
    When I send a "GET" request to "http://rewrite-url.com/bar"
    Then the response status code must be 200
    And the value of the "key1" field in the header must be "value1"
    And the value of the "key2" field in the header must be "value2"
    And the value of the "key3" field in the header must be "value3"
    And the value of the "key2" field in the header must be "value2_addvalue" 
    And header's "key3" field and its value must be deleted
