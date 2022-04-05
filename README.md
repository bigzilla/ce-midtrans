# CloudEvent x Midtrans

CloudEvent adapter for Midtrans notifications.

## Usage

Specify below environment variable.

| Variable     | Description                             | Example                 |
| ------------ | --------------------------------------- | ----------------------- |
| `K_SINK`     | Sink URI.                               | `http://localhost:8080` |
| `SERVER_KEY` | Midtrans server key.                    | `SB-Mid-server-abcdefg` |
| `PORT`       | *Optional*. Server port. Default `8080` | `8080`                  |

### With Knative Eventing

Using ContainerSource (K_SINK is injected by ContainerSource).

```yaml
apiVersion: sources.knative.dev/v1
kind: ContainerSource
metadata:
  name: <name>
spec:
  template:
    spec:
      containers:
        - image: ghcr.io/injustease/ce-midtrans
          name: ce-midtrans
          env:
            - name: SERVER_KEY
              valueFrom:
                secretKeyRef:
                  name: <secret-name>
                  key: <secret-key>
  sink:
    ref:
      apiVersion: <apiVersion>
      kind: <kind>
      name: <sink>
```

Using Knative serving

```yaml
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: <name>
spec:
  template:
    spec:
      containers:
        - image: ghcr.io/injustease/ce-midtrans
          name: ce-midtrans
          env:
            - name: K_SINK
              value: <sink-uri>
            - name: SERVER_KEY
              valueFrom:
                secretKeyRef:
                  name: <secret-name>
                  key: <secret-key>
```

## TODO

- [x] Payment Notification
- [ ] Recurring Notification
- [ ] Pay Account Notification
- [ ] Test
- [x] Usage docs
- [ ] Event Registry
