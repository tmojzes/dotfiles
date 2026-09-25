package main

// armadaCommandHelp is the armada-xo bot's command reference, embedded in
// the tool description so LLM callers know the available commands and their
// arguments. Source: the bot's own help output, lightly cleaned (chat
// timestamps and duplicated usage headers removed).
const armadaCommandHelp = `Usage information: parameters shown in italics in the bot's help (typically region=<region> and carrier=<carrier>) are optional, and will restrict results to a particular region (e.g. us-south) and/or a particular carrier (e.g. carrier5).

accountVPCZoneFlavors bssAccount=<bssAccount> zone=<zone,...> flavor=<flavor,...> tag=<tag,...> groupBy=<zone|flavor> region=<region>: Displays flavors available in VPC infrastructure zones with optional flavor filtering and grouping
accountVPCZones bssAccount=<bssAccount> region=<region>: Displays VPC infrastructure zone map for a user account
bf.setCBR region=<region> train=<train>: Set DesiredContextBasedRestrictionEventID if required
bssAccountState bssAccount=<bssAccount,...>: Will return state information about account in BSS from BSS account id.
buildSyncGhostClusterCmd <cluster>: builds command to synchronized a single cluster with ghost
buildSyncIGWorkerCountCmd <id>: will build a command for a conductor to run to synchronize the actual worker count for an instance group
cloudflare-rayid <rayID>: Will return information about a cloudflare request from a ray id (5b7767e3febbcf58 or 5b7767e3febbcf58-IAD)
clusterInfoReport region=<region>: Will generate a report with the counts of different, deployed cluster types in a region. See SQL WIKI with help sql for more.
coeSQL region=<region> name=<classicMigrationReport>: Will perform one of the predefined SQL queries.
  classicMigrationReport: Get classic cluster and related worker information
cse.cluster <cluster>: validates CSE settings for a single cluster
cse.region <region>: validates CSE settings for a region
debug.GoroutineCount region=<region> service=<armada-api|armada-xo|armada-reaper>: Will retrieve go routine count from an xo instance.
debug.Stacktraces region=<region> service=<armada-api|armada-xo|armada-reaper>: Will retrieve go stacktraces from an xo instance to a file.
deleteFailedWorkersReport region=<region> issue=<issue>: Will generate reports with details of all delete_failed workers, which can be used to generate input to the armada-billing-reconcile job. Optional issue parameter for use by @Phillips bot (owned by Troutbridge squad)
eek.summary account=<account,...> region=<region>: shows envelope encryption key summary
etcdcerts region=<region>: Will show etcd cert info.
fluentdFilterConfig <clusterID>: Retrieves the filter configuration for a cluster
fluentdLogConfig <clusterID>: Retrieves the logging configuration for a cluster
getAccountQuota <accountid>
getCargoCounts region=<region>: get the instance counts for every cargo field
getPrunedClusterFields region=<region> cluster=<cluster>: retrieves all remaining fields for a pruned cluster
globalAccountClusters <accountID> resourceGroupID=<resourceGroupID>: Calls global api to get the account's clusters
globalAdminCluster <accountID> clusterID=<clusterID>: Calls aramda-api admin get cluster url to get the account's cluster as API knows it
health.issueDetails cluster=<cluster> issueID=<issueID,...> region=<region>: Retrieves detailed information for specific health issues
health.issues cluster=<cluster> region=<region>: Retrieves health issues for a cluster
ingress.status cluster=<cluster> region=<region>: Retrieves Ingress Status report for a cluster
ledger.AccountResources crnPrefix=<crnPrefix> startTime=<YYYY-MM-DDTHH:MM:SS±0000> endTime=<YYYY-MM-DDTHH:MM:SS±0000> limit=<limit> offset=<offset> resourceType=<cluster|worker|satconnector|dhost|nlb|ccontract|qnassignment|odflicense|meteredresource> region=<region>: retrieves ledger account resources
ledger.Cluster id=<id> show=<show> region=<region>: cluster ledger record
ledger.CommittedContract id=<id> show=<show> region=<region>: committed contract ledger record
ledger.CountRows resourceType=<cluster|worker|satconnector|dhost|nlb|ccontract|qnassignment|odflicense|meteredresource> entryType=<Use one of the following values: Registered/Deregistered/MeteringStartRequired/MeteringStarted/MeteringAdjustmentRequired/MeteringAdjusted/MeteringRenewed/MeteringStopRequired/MeteringStopped/DeleteRequested/CommitmentUnfulfilled> region=<region>: count rows in specified resource's ledger table, optionally counts specified entryTypes in the resource's table
ledger.DatastoreReport region=<region>: generate report on ledger datastore usage (slow)
ledger.DedicatedHost id=<id> show=<show> region=<region>: dedicated host ledger record
ledger.EntryRecords crnPrefix=<crnPrefix> startTime=<YYYY-MM-DDTHH:MM:SS±0000> endTime=<YYYY-MM-DDTHH:MM:SS±0000> limit=<limit> offset=<offset> resourceType=<cluster|worker|satconnector|dhost|nlb|ccontract|qnassignment|odflicense|meteredresource> payloadTypes=<Registered/MeteringStarted/.../all> region=<region>: retrieves ledger entry records
ledger.MeteredResource id=<id> show=<show> region=<region>: Generic Metered Resource ledger record
ledger.NLB id=<id> show=<show> region=<region>: nlb ledger record
ledger.ODFLicense id=<id> show=<show> region=<region>: OpenShift Data Foundation License ledger record
ledger.Ping region=<region>: ledger ping
ledger.QueueNodeAssignment id=<id> show=<show> region=<region>: queue node assignment ledger record
ledger.SatConnector id=<id> show=<show> region=<region>: satconnector ledger record
ledger.Worker id=<id> show=<show> region=<region>: worker ledger record
metrics region=<region>: Will retrieve metrics from an xo instance as a file.
nodepool id=<id>: Check the node pool CR for an instance group
nps-link email=<email> firstName=<firstName> lastName=<lastName> lang=<lang> test=<true|false>: Generate a Medallia link for an NPS survey
ping region=<region>: Will show the region and environment of armada-xo, by default all regions are shown
providerClusters provider=<provider> region=<region>: Will list clusters and accounts that have a specified default provider
queryGhostClusters <account> cluster=<cluster> region=<region>: queries ghost for IKS clusters in a given account
queryGhostControllers <account> controller=<controller> region=<region>: queries ghost for Satellite locations in a given account
reaper.summary region=<region>: reaper job summary for the designated region
reservationReport region=<region> date=<YYYY-MM>: Request report of previous month reservation utilisation (slow)
sat.acl.ip-check cluster=<cluster> ip=<ip>: check if an IP is allowed against to a location cruiser
sat.acl.nodes-check cluster=<cluster>: check if all the nodes ip addresses are in the configured satellite acl list
sqlCmd region=<region> name=<accountClusterVersions|activeClusterBoms|activeClustersDataCenter|addonVersions|clusterSupportInfo|clusterOSInfo|clusterHealthByAccount> isIBM=<|true|false|True|False> cluster=<cluster> account=<account> os=<os> health_state=<health_state>: Will perform one of the predefined SQL queries.
  accountClusterVersions: Get account_id, cluster_id, name, ansible_bom_version, and is_ibmer for active clusters under a specific account, REQUIRED account=<accountID>
  activeClusterBoms: Get number of active clusters for each BOMVersion
  activeClustersDataCenter: Get count of active clusters by datacenter
  addonVersions: Get accountID, clusterID, addonID, addonHealthState, and addonVersion of active clusters, cluster=<clusterID>, account=<accountID>
  clusterHealthByAccount: Get clusterID, name, and health state for active clusters with workers, account=<accountID>, health_state=<normal/critical/unsupported>
  clusterOSInfo: Get accountID, clusterID, name, OS, and isIBMer for active clusters with workers, os=<operating system>, isIBM=<true/false>
  clusterSupportInfo: Get accountID, clusterID, name, BOMVersion, EOS date, and isIBMer bool of unsupported clusters, isIBM=<true/false>, account=<accountID>
sqlQuery region=<region> headers=true SELECT...: Will perform an SQL query.
storage.getVolumes region=<region> volumeType=<volumeType> status=<status,...> count=<count>: returns details or counts of volume resources
suspendedClustersReport region=<region> date=<YYYY-MM-DD> issue=<issue>: Will generate a report with details of all suspended clusters, which can be used as an input to the armada-billing-reconcile job. Default behaviour returns all clusters beyond threshold which are now candidates for deletion. Optional date parameter instead retrieves all clusters suspended since a particular date. Optional issue parameter for use by @Phillips bot (owned by Troutbridge squad)
unsatisfiedMeteringStartRequiredReport region=<region> startDate=<YYYY-MM-DD> endDate=<YYYY-MM-DD> issue=<issue>: Will generate a report with details of all unsatisfied metering start entries and associated plan payloads. Default behaviour returns all unsatisfied metering start required entries for clusters and workers over the past month where start was required more than 1 day ago. Optional start and end date parameters instead retrieve unsatisfied metering start required entries for the given time period. Optional issue parameter for use by @Phillips bot (owned by Troutbridge squad)
viewNodeQueueExpiry region=<region> account=<account>: View expiration dates of pre-existing node-queue tokens

The following commands list resources per account, cluster, location, provider, reservation, satellite, or zone. Each has two forms: with fields=<fields> to show the fields specified in the /-separated fields list for all matching resources, or without it to show a summary of all matching resources. Optional region=<region> and carrier=<carrier> parameters restrict results to a particular region and/or carrier.

accountClusters <accountid> fields=<fields> region=<region> carrier=<carrier>
accountDedicatedHostPools <accountid> fields=<fields> region=<region> carrier=<carrier>
accountMeteredResources <accountid> fields=<fields> region=<region> carrier=<carrier>
accountReservations <accountid> fields=<fields> region=<region> carrier=<carrier>
accountSatelliteConnectors <accountid> fields=<fields> region=<region> carrier=<carrier>
accountSatelliteLocations <accountid> fields=<fields> region=<region> carrier=<carrier>
clusterAlbs <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterCleanupTasks <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterIngressInstances <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterInstanceGroups <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterMasters <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterMeteredResources <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterPortableSubnets <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterSecrets <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterSubdomains <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterTaggedResources <tagclusterid> fields=<fields> region=<region> carrier=<carrier>
clusterWorkerPools <clusterid> fields=<fields> region=<region> carrier=<carrier>
clusterWorkers <clusterid> fields=<fields> region=<region> carrier=<carrier>
dedicatedHostPoolGroups <[accountid/]dhpoolid> fields=<fields> region=<region> carrier=<carrier>
dedicatedHostPoolHosts <[accountid/]dhpoolid> fields=<fields> region=<region> carrier=<carrier>
dedicatedHostPoolWorkerPools <dhpoolid> fields=<fields> region=<region> carrier=<carrier>
locationClusters <controllerid> fields=<fields> region=<region> carrier=<carrier>
providerDedicatedHostFlavors <provider> fields=<fields> region=<region> carrier=<carrier>
providerFlavors <provider> fields=<fields> region=<region> carrier=<carrier>
reservationContracts <[accountid/]ccomputeid> fields=<fields> region=<region> carrier=<carrier>
reservationWorkerPools <ccomputeid> fields=<fields> region=<region> carrier=<carrier>
satelliteHostWorkers <[accountid/controllerid/]nodeid> fields=<fields> region=<region> carrier=<carrier>
satelliteLocationHosts <[accountid/]controllerid> fields=<fields> region=<region> carrier=<carrier>
zoneDedicatedHostFlavors <provider/zone> fields=<fields> region=<region> carrier=<carrier>
zoneFlavors <provider/zone> fields=<fields> region=<region> carrier=<carrier>`
