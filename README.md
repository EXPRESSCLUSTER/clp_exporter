# clp_exporter
- This is sample exporter for EXPRESSCLUSTER.
- This exporter can get the following data.
  - Performance data of each mirror disk resource
  - Elapse time of each monitor resource


# How to Build
1. Clone this repository.
   ```sh
   git clone https://github.com/EXPRESSCLUSTER/clp_exporter.git
   ```
1. Move to `src` directory.
   ```sh
   cd clp_exporter/src
   ```
1. Initialize and build clp_exporter.
   ```sh
   go mod init clp_exporter
   ```
   ```sh
   go mod tidy
   ```
   ```sh
   go build
   ```

# How to Use
```
+---------------------+
| Ubuntu Server 24.04 |
| Prometheus          |
+-+-------------------+
  |
  |  +----------------------------------+
  |  | Node 1                           |
  |  | - AlmaLinux 9.6                  |
  +--+ - EXPRESSCLUSTER X for Linux 5.3 |
  |  | - clp_exporter                   |
  |  +----------------------------------+
  |
  |  +----------------------------------+
  |  | Node 2                           |
  |  | - AlmaLinux 9.6                  |
  +--+ - EXPRESSCLUSTER X for Linux 5.3 |
     | - clp_exporter                   |
     +----------------------------------+

```

1. Install Prometheus.
1. Install EXPRESSCLUSTER and create a cluster.
1. Save clp_exporter on both Node 1 and 2.
1. Run clp_exporter.
   ```sh
   ./clp_exporter
   ```
1. Open `prometheus.yml` and add IP address of Node 1 and 2 as below.
   ```sh
   vim /etc/prometheus/prometheus.yml
   ```
   ```
     - job_name: clp
       static_configs:
         - targets: ['192.168.122.11:29090']
           labels:
             name: node1
         - targets: ['192.168.122.12:29090']
           labels:
             name: node2
   ```
1. Open web browser and access to Prometheus.
   ```sh
   http://<IP address>:9090
   ```
1. Add the metrics of EXPRESSCLUSTER.
   - clp_mirror_\<mirror disk resource name\>_\<performance data\>
   - clp_monitor_\<monitor resoruce name\>