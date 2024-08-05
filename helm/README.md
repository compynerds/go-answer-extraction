# helm-chart-template
A templatized, one-size-fits-all Helm chart.

### Explanation of Configs in the `values.yaml` file.
##### Global
- `name:` Primarily used to name the pods and label resources.
- `namespace:` Namespace to deploy resources to.

##### Deployment
- `deployment:` This block houses all the configurations for the primary pod. Currently only supports single-container deployments.
- `image.name:` Image to be pulled. IE when using `docker pull`
- `image.tag:` Tag of the image to be pulled. Typically `latest`
- `restartPolicy:` Usually `Always`. Can also be `Never` or `OnFailure`. Read more about pod lifecycles [HERE](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/).
- `kind:` Default is `Deployment`. Could also be `DaemonSet`
- `deployAnnotations:` Annotations to apply to the deployment. The will apply to all containers in a deployment. This is typically where annotations will be set.
- `podAnnotations:` Annotations to apply to a specific container within the pod. Not typically used.
- `labels:` Additional labels to add to the deployment. Defaults are `app`, `chart`, `release`, and `heritage`.
- `livenessProbe:` The Healthcheck config to ensure the pod stays in a ready state. Should typically be configured to hit a similar endpoint as incoming traffic. Read more about the healthchecks [HERE](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
- `readinessProbe:` The Healthcheck config to check that the pod is ready to receive traffic when it is first created. Read more about the healthchecks [HERE](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/).
- `ports:` Configuration of ports that can be used to access this pod. Read more about the possible port configurations [HERE](https://matthewpalmer.net/kubernetes-app-developer/articles/kubernetes-ports-targetport-nodeport-service.html).
- `volumes:` Creating volumes that can be made available to the pods via VolumeMounts. Learn more about using volumes in pods [HERE](https://kubernetes.io/docs/concepts/storage/volumes/).
- `volumeMounts:` Location inside of the pod to mount the Volume and make it available to the pod.
- `updateStrategy:` Behavior used to update the pods when multiples have been deployed. Learn more about Update Strategies [HERE](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#update-strategies).
- `resources:` The requests and limits are used extensively by HPAs. These values should be set for every deployment but finding the correct settings requires knowing the needs of the application which isn't always possible until after the application has been deployed. Learn more about HPAs [HERE](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/).

#### Service
- `annotations:` Annotations to be applied to the service. Sometimes helpful when using tools that automatically create Route53 records or auto-provision ALBs
- `type:` The type determines the service's behavior. Learn more about the different service types [HERE](https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types).
- `ports:` the configuration of the ports that are used for the service.

#### Secrets
- `secrets:` Should contain a list of `key:value` pairs that will be converted in to environment variables for the pod.

#### Configmaps
- `configmaps:` Each item here will be translated in to its own configmap that can then be mapped to volumes to be made available to the pods. For instance, when creating a config file for a sub-process on the pod.
- `configmaps.data` should contain `key:value` pairs and the value should _always_ be a string. Follow the example for a multi-line string.
