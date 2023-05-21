# track-cbf

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square)
![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square)
![AppVersion: 0.1.0](https://img.shields.io/badge/AppVersion-0.1.0-informational?style=flat-square)

track-cbf CronJob Helm chart

* [Values](#values)

---

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| event | string | `"cbf2023"` | The Cambridge Beer Festival event to pull data from. |
| extraArgs | list | `[]` | Additional arguments to pass to the container at run-time. e.g. to filter for specific drink categories: <br /><br /><pre> extraArgs: <br /> - -c <br /> - beer <br /> - -c <br /> - wine </pre> |
| image.pullPolicy | string | `"IfNotPresent"` | The image pull policy. |
| image.repository | string | `"track-cbf"` | The Docker image repository (and registry) to pull from. |
| image.tag | string | `""` | The Docker image tag, defaults to the `appVersion` property of the chart. |
| localPath | string | `"/var/lib/cbf2023/"` | Local path defines where on disk to create the Persistent Volume and mount it in the container. |
| nodeAffinity | string | `"ns3093148"` | Node Affinity defines which node to create the PV and run the CronJob on. |
