**This is my first program in Go language. I created it for my use case - to backup files from pods running on my home k8s cluster and then restore them using initContainer - so I don't need to use persistent storage for the cluster and can easy delete/create the cluster back if required. For example, this tool make backups of my personal HomeAssistant, Vaultwarden and NextCloud**
```diff
**This tool is manupulating with files on S3 bucket. Do not use this tool without specifying proper settings or you may loose your data!**
```
# Tool to backup and restore local directories or MySQL databaseto AWS S3 storage

Simple tool to archive all files and subdirectories of desired local directory as one ZIP or TAR archive file and upload this file to S3 bucket, to desired folder.
Tool can also make mysql dump and then archive and upload it to S3 bucket/folder.

Later this files anb mysql database dump can be restored by the same tool to any other directory.

The tool can also work as a daemon and run periodical backups based on desired interval. In this mode the tool also do automatic pruning and keep only number of archives in S3 bucket (deleting old files from desired folder).

Set the following environment variables for this tool to work:

Mandatory variables:

* AWS_BUCKET - S3 bucket to use.
* AWS_KEY - key to use to access to the bucket
* AWS_SECRET_KEY - secret key to use to access to the bucket
* DIR_TO_BACKUP - absolute path for directory to backup (tool will backup all files and subdirectories inside it)
* DIR_TO_RESTORE - absolute path for directory to restore into
* MYSQL_PASSWORD - Password to connect to MySQL database to make mysql dump (if executing tool to make MySQL dump).

Optionally, set the following variables:
* S3_BUCKET_FOLDER - folder where to store ZIP archive. "podbackup" by default
* S3_FILE_PREFIX - ZIP archive name prefix. "podbackup" by default. Full filename will be <prefix>-<timestamp>.zip
* ENCRYPT_PASSWORD - encrypt/decrypt ZIP archives using this password. 
* BACKUP_INTERVAL - interval in seconds (if number like `3000`) or in minutes/hours (like `2m` or `24h`) to run periodical backups (if running as daemon). 1h by default.
* PRUNE_INTERVAL - interval in seconds (if number like `3000`) or in minutes/hours (like `2m` or `24h`) to run periodical pruning (if running as daemon). 2h by default.
* COPIES_TO_KEEP - number of copies to keep in S3 folder when executing pruning.
* FORCE_RESTORE - set to True if needed the tool to fail (exit with code 1) if it cannot restore files from backup. Useful if tool is using as IniContainer and you don't want main containers to run without restoring actual data.
* MYSQL_USER - user to connect to MySQL database when making mysql dump. Default value as root.
* MYSQL_HOST - IP address or hosname to use to connect to MySQL database. Default value is 127.0.0.1. Process will wait for connection to restore the database.
* MYSQL_PORT - port to use to connect to MySQL database. Default value is 3306.
*	ARCHIVE_TYPE - 	by default set to `zip` - tool will create ZIP archive (and encrypt it if ENCRYPT_PASSWORD is set). 
	              	Set to `tarzip` - to archive all files as tar archive and then zip it (encrypted if ENCRYPT_PASSWORD is set).
			            Use it if you need to save original ownership and mode of the files.
			            Set to `targz` - to archive all files as tar compressed archive. File mode and ownership persist during unpacking, however encryption is not supported.
	
## Usage

Run 
`podbackup <command>`

where commands are:

`backup` - run one time backup and backup all files from folder `DIR_TO_BACKUP` to S3 object storage.

`backup-sql` - run one time backup to make MySQL dump based on `MYSQL*` variables.

`backup-daemon` - work as daemon and run periodical backups according to BACKUP_INTERVAL environment variable (`1h` by default). In this mode daemon will do automatic pruning (default PRUNE_INTERVAL is `2h`) and keep only # of copies based on COPIES_TO_KEEP environment variable (3 by default). **All other files with prefix `S3_FILE_PREFIX` in the folder `S3_BUCKET_FOLDER` will be destroyed.**

`backup-sql-daemon` - work as daemon and run periodical MySQL database dumps according to BACKUP_INTERVAL environment variable (`1h` by default). In this mode daemon will do automatic pruning (default PRUNE_INTERVAL is `2h`) and keep only # of copies based on COPIES_TO_KEEP environment variable (3 by default). **All other files with prefix `S3_FILE_PREFIX` in the folder `S3_BUCKET_FOLDER` will be destroyed.**

`prune` - manually run pruning (delete all old archives). **All other files with prefix `S3_FILE_PREFIX` in the folder `S3_BUCKET_FOLDER` will be destroyed.**

`list` - list files in S3 folder (based on `S3_BUCKET_FOLDER` environment variable).
	
`restore` - download file from S3 and restore files to directory (DIR_TO_RESTORE environment variable). Most recent archive will be used. To restore from another file provide archive name based on 'podbackup list' output (like podbackup/podbackup-20210802213807.zip) as an argument for `restore` command.

`restore-sql` - download file from S3 and restore MySQL database. Most recent archive will be used. To restore from another file provide archive name based on 'podbackup list' output (like podbackup/podbackup-20210802213807.zip) as an argument for `restore-sql` command.

## Use-case
I'm using this tool to backup/restore my Home Assistant, Vaultwarden and NextCloud data - so I don't need to use localstorage with my home k8s cluster.

The tool is working as a sidecar container for home-assistant pod and runs periodical backups (once per hour).

To restore the data the Home Assistant k8s deployment has initContainer with the same tool performing restore process. 

Example of home-assistant k8s deployment:
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: home-assistant
  namespace: home-assistant
  labels:
    app.kubernetes.io/name: home-assistant
spec:
  revisionHistoryLimit: 3
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      app.kubernetes.io/name: home-assistant
      app.kubernetes.io/instance: home-assistant
  template:
    metadata:
      labels:
        app.kubernetes.io/name: home-assistant
        app.kubernetes.io/instance: home-assistant
    spec:
      serviceAccountName: default
      terminationGracePeriodSeconds: 90
      dnsPolicy: ClusterFirst
      enableServiceLinks: true
      volumes:
        - name: hass-workingdir
          emptyDir: {}
      containers:
        - name: home-assistant
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 512Mi        
          image: "homeassistant/home-assistant:2021.6.3"
          imagePullPolicy: IfNotPresent
          securityContext:
            privileged: false
          env:
            - name: "TZ"
              value: "UTC"
          ports:
            - name: http
              containerPort: 8123
              protocol: TCP
          volumeMounts:
            - name: hass-workingdir
              mountPath: /config          
          livenessProbe:
            tcpSocket:
              port: 8123
            initialDelaySeconds: 30
            failureThreshold: 10
            timeoutSeconds: 1
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 8123
            initialDelaySeconds: 30
            failureThreshold: 10
            timeoutSeconds: 1
            periodSeconds: 10
          startupProbe:
            tcpSocket:
              port: 8123
            initialDelaySeconds: 300
            failureThreshold: 30
            timeoutSeconds: 2
            periodSeconds: 5
        - name: hass-backup-daemon
          image: quay.io/vadimzharov/podbackup:latest
          lifecycle:
            preStop:
              exec:
                command: ["/usr/local/bin/podbackup", "backup"]
          command: ["podbackup"]
          args: ["backup-daemon"]
          resources:
            requests:
              cpu: 50m
              memory: 64Mi
            limits:
              cpu: 50m
              memory: 64Mi
          volumeMounts:
            - name: hass-workingdir
              mountPath: /hass-workdir
          env:
          - name: AWS_BUCKET
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awsbucket
          - name: AWS_KEY
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awskey
          - name: AWS_SECRET_KEY
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awssecretkey
          - name: S3_BUCKET_FOLDER
            value: "homehass"
          - name: S3_FILE_PREFIX
            value: "hass-backup"
          - name: DIR_TO_BACKUP
            value: "/hass-workdir"
          - name: BACKUP_INTERVAL
            value: "1h"
      initContainers:
        - name: hass-restore
          image: quay.io/vadimzharov/podbackup:latest
          command: ["podbackup"]
          args: ["restore"]
          volumeMounts:
            - name: hass-workingdir
              mountPath: /hass-workdir
          env:
          - name: AWS_BUCKET
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awsbucket
          - name: AWS_KEY
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awskey
          - name: AWS_SECRET_KEY
            valueFrom:
              secretKeyRef:
                name: backup-restore-creds
                key: awssecretkey
          - name: S3_BUCKET_FOLDER
            value: "homehass"
          - name: S3_FILE_PREFIX
            value: "hass-backup"
          - name: DIR_TO_RESTORE
            value: "/hass-workdir"
```