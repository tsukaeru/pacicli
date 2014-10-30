package lib

import (
	"encoding/xml"
)

type PasswordResponse struct {
	XMLName  xml.Name `xml:"pwd-response" json:"-" toml:"-"`
	Message  string   `xml:"message"`
	Password string   `xml:"password"`
}

type VeList struct {
	XMLName xml.Name `xml:"ve-list" json:"-" toml:"-"`
	VeInfo  []struct {
		ID             int    `xml:"id,attr"`
		Name           string `xml:"name,attr"`
		Hostname       string `xml:"hostname,attr"`
		State          string `xml:"state,attr"`
		SubscriptionID int    `xml:"subscription-id,attr"`
	} `xml:"ve-info"`
}

type CPU struct {
	Number int `xml:"number,attr"`
	Power  int `xml:"power,attr"`
}

type VeDisk struct {
	StorageID string `xml:"storage-id,attr"`
	Created   bool   `xml:"created,attr"`
	GlobalID  int    `xml:"global-id,attr"`
	ID        int    `xml:"id,attr"`
	Type      string `xml:"type,attr"`
	Size      int    `xml:"size,attr"`
}

type Platform struct {
	TemplateInfo struct {
		Name string `xml:"name,attr"`
	} `xml:"template-info"`
	OSInfo struct {
		Type       string `xml:"type,attr"`
		Technology string `xml:"technology,attr"`
		Family     string `xml:"family,attr,omitempty"`
	} `xml:"os-info"`
}

type Network struct {
	PrivateIP IPAddr `xml:"private-ip,attr"`
	PublicIP  []struct {
		ChunkRef int    `xml:"chunk-ref,attr"`
		ID       int    `xml:"id,attr"`
		Address  IPAddr `xml:"address,attr"`
		Gateway  IPAddr `xml:"gateway,attr"`
	} `xml:"public-ip"`
	PublicIP6 []struct {
		ID      int    `xml:"id,attr"`
		Address IPAddr `xml:"address,attr"`
		Gateway IPAddr `xml:"gateway,attr"`
	} `xml:"public-ip6"`
}

type BackupSchedule struct {
	Name string `xml:"name,attr"`
}

type Console struct {
	Address IPAddr `xml:"address"`
	Port    int    `xml:"port"`
}

type Admin struct {
	Login    string `xml:"login,attr"`
	Password string `xml:"password"`
}

type Traffic struct {
	Sent     int `xml:"sent,attr"`
	Received int `xml:"received,attr"`
}

type Ve struct {
	XMLName         xml.Name       `xml:"ve" json:"-" toml:"-"`
	ID              int            `xml:"id"`
	UUID            string         `xml:"uuid"`
	Hnid            int            `xml:"hnId"`
	CustomerID      int            `xml:"customer-id"`
	Name            string         `xml:"name"`
	Hostname        string         `xml:"hostname"`
	Description     string         `xml:"description"`
	SubscriptionID  int            `xml:"subscription-id"`
	CPU             CPU            `xml:"cpu"`
	RAMSize         int            `xml:"ram-size"`
	Bandwidth       int            `xml:"bandwidth"`
	VeDisk          VeDisk         `xml:"ve-disk"`
	Platform        Platform       `xml:"platform"`
	Network         Network        `xml:"network"`
	BackupSchedule  BackupSchedule `xml:"backup-schedule"`
	Console         Console        `xml:"console"`
	State           string         `xml:"state"`
	PrimaryDiskID   int            `xml:"primary-disk-id"`
	TemplateID      int            `xml:"template-id"`
	Admin           Admin          `xml:"admin"`
	LastOperationRc int            `xml:"last-operation-rc"`
	AppInfo         []struct {
		AppTemplate   string `xml:"app-template,attr"`
		ForOS         string `xml:"for-os,attr"`
		InstalledAt   string `xml:"installed-at,attr"`
		InstalledOk   bool   `xml:"installed-ok,attr"`
		UninstalledAt string `xml:"uninstalled-at,attr"`
		UninstalledOk bool   `xml:"uninstalled-ok,attr"`
		AppTemplateID int    `xml:"app-template-id,attr"`
	} `xml:"app-info"`
	LoadBalancer               string     `xml:"load-balancer"`
	SteadyState                string     `xml:"steady-state"`
	Autoscale                  *Autoscale `xml:"autoscale"`
	CurrentResourceConsumption struct {
		CPU            int     `xml:"cpu,attr"`
		RAM            int     `xml:"ram,attr"`
		PrivateTraffic Traffic `xml:"private-traffic"`
		PublicTraffic  Traffic `xml:"public-traffic"`
	} `xml:"current-resource-consumption"`
}

type CreateVe struct {
	XMLName        xml.Name `xml:"ve" json:"-" toml:"-"`
	CustomNs       bool     `xml:"custom-ns,attr,omitempty"`
	Name           string   `xml:"name"`
	Hostname       string   `xml:"hostname"`
	Description    string   `xml:"description"`
	SubscriptionID int      `xml:"subscription-id,omitempty"`
	CPU            CPU      `xml:"cpu"`
	RAMSize        int      `xml:"ram-size"`
	Bandwidth      int      `xml:"bandwidth"`
	NoOfPublicIP   int      `xml:"no-of-public-ip,omitempty"`
	NoOfPublicIPv6 int      `xml:"no-of-public-ipv6,omitempty"`
	VeDisk         struct {
		Local   bool `xml:"local,attr"`
		Primary bool `xml:"primary,attr,omitempty"`
		Size    int  `xml:"size,attr"`
	} `xml:"ve-disk"`
	Platform       Platform `xml:"platform"`
	BackupSchedule *struct {
		Name string `xml:"name,attr"`
	} `xml:"backup-schedule"`
}

type ChangeCPU struct {
	Number int `xml:"number,attr,omitempty"`
	Power  int `xml:"power,attr,omitempty"`
}

type AddIP struct {
	Number int `xml:"number,attr"`
}

type DropIP struct {
	IP IPAddrList `xml:"ip,attr"`
}

type ReconfigureIP struct {
	AddIP  *AddIP  `xml:"add-ip"`
	DropIP *DropIP `xml:"drop-ip"`
}

type ReconfigureVe struct {
	XMLName         xml.Name       `xml:"reconfigure-ve" json:"-" toml:"-"`
	Description     string         `xml:"description,omitempty"`
	ChangeCPU       *ChangeCPU     `xml:"change-cpu"`
	RAMSize         int            `xml:"ram-size,omitempty"`
	Bandwidth       int            `xml:"bandwidth,omitempty"`
	ReconfigureIPv4 *ReconfigureIP `xml:"reconfigure-ipv4"`
	ReconfigureIPv6 *ReconfigureIP `xml:"reconfigure-ipv6"`
	PrimaryDiskSize int            `xml:"primary-disk-size,omitempty"`
	CustomNs        *int           `xml:"custom-ns"`
}

type VeHistory struct {
	XMLName    xml.Name `xml:"ve-history" json:"-" toml:"-"`
	VeSnapshot []struct {
		CPU                    int       `xml:"cpu,attr"`
		RAM                    int       `xml:"ram,attr"`
		LocalDisk              int       `xml:"local-disk,attr"`
		Nbd                    int       `xml:"nbd,attr"`
		Bandwidth              int       `xml:"bandwidth,attr"`
		LastTouchedFrom        string    `xml:"last-touched-from,attr"`
		State                  string    `xml:"state,attr"`
		SteadyState            string    `xml:"steady-state,attr"`
		LastChangedBy          string    `xml:"last-changed-by,attr"`
		EventTimestamp         Timestamp `xml:"event-timestamp,attr"`
		NoOfPublicIP           int       `xml:"no-of-public-ip,attr"`
		NoOfPublicIPv6         int       `xml:"no-of-public-ipv6,attr"`
		IsLb                   bool      `xml:"is-lb,attr"`
		PrivateIncomingTraffic int       `xml:"private-incoming-traffic,attr"`
		PrivateOutgoingTraffic int       `xml:"private-outgoing-traffic,attr"`
		PublicIncomingTraffic  int       `xml:"public-incoming-traffic,attr"`
		PublicOutgoingTraffic  int       `xml:"public-outgoing-traffic,attr"`
	} `xml:"ve-snapshot"`
}

type VeResourceUsageReport struct {
	XMLName           xml.Name `xml:"ve-resource-usage-report" json:"-" toml:"-"`
	VeName            string   `xml:"ve-name,attr"`
	VeID              int      `xml:"ve-id,attr"`
	OS                string   `xml:"os,attr"`
	Technology        string   `xml:"technology,attr"`
	LifeTimeInMinutes int      `xml:"life-time-in-minutes,attr"`
	IsLoadBalancer    bool     `xml:"is-load-balancer,attr"`
	ResourceUsage     []struct {
		Value             int    `xml:"value,attr"`
		ResourceUsageType string `xml:"resource-usage-type,attr"`
		ResourceType      string `xml:"resource-type,attr"`
	} `xml:"resource-usage"`
	VeTraffic []struct {
		TrafficType string `xml:"traffic-type,attr"`
		Used        int    `xml:"used,attr"`
	} `xml:"ve-traffic"`
	ActiveBackupSchedule []struct {
		ScheduleName string `xml:"schedule-name,attr"`
	} `xml:"active-backup-schedule"`
}

type Firewall struct {
	XMLName xml.Name `xml:"firewall" json:"-" toml:"-"`
	Rule    []struct {
		ID         int      `xml:"id,attr,omitempty" json:"-" toml:"-"`
		Name       string   `xml:"name,attr"`
		Protocol   string   `xml:"protocol,attr"`
		LocalPort  int      `xml:"local-port,attr"`
		RemotePort int      `xml:"remote-port,attr"`
		RemoteNet  []IPAddr `xml:"remote-net"`
	} `xml:"rule"`
}

type Backup struct {
	XMLName        xml.Name  `xml:"backup" json:"-" toml:"-"`
	ImBackupID     int       `xml:"im-backup-id,attr"`
	CloudBackupID  string    `xml:"cloud-backup-id,attr"`
	ScheduleName   string    `xml:"schedule-name,attr"`
	Started        Timestamp `xml:"started,attr"`
	Ended          Timestamp `xml:"ended,attr"`
	Successful     bool      `xml:"successful,attr"`
	BackupSize     int       `xml:"backup-size,attr"`
	BackupNodeName string    `xml:"backup-node-name,attr"`
	Description    string    `xml:"description"`
}

type VeBackups struct {
	XMLName xml.Name `xml:"ve-backups" json:"-" toml:"-"`
	Backup  []Backup `xml:"backup"`
}

type Threshold struct {
	Threshold *int `xml:"threshold,attr"`
	Period    int  `xml:"period,attr"`
}

type AutoscaleRule struct {
	XMLName           xml.Name   `xml:"autoscale-rule" json:"-" toml:"-"`
	Enabled           *bool      `xml:"enabled,attr"`
	Deleted           *bool      `xml:"deleted,attr"`
	Metric            string     `xml:"metric,attr"`
	Version           *int       `xml:"version,attr"`
	Updated           *Timestamp `xml:"updated,attr,omitempty"`
	UpdateDeliveredOk *bool      `xml:"update-delivered-ok,attr"`
	UpdateDelivered   *Timestamp `xml:"update-delivered,attr,omitempty"`
	AllowMigration    *bool      `xml:"allow-migration,attr"`
	AllowRestart      *bool      `xml:"allow-restart,attr"`
	Limits            *struct {
		Min  int `xml:"min,attr"`
		Max  int `xml:"max,attr"`
		Step int `xml:"step,attr"`
	} `xml:"limits"`
	Thresholds *struct {
		Up   *Threshold `xml:"up"`
		Down *Threshold `xml:"down"`
	} `xml:"thresholds"`
}

type Autoscale struct {
	XMLName xml.Name `xml:"autoscale" json:"-" toml:"-"`
	Current *struct {
		AutoscaleRule []AutoscaleRule `xml:"autoscale-rule"`
	} `xml:"current"`
	Ongoing *struct {
		AutoscaleRule []AutoscaleRule `xml:"autoscale-rule"`
	} `xml:"ongoing"`
}

type AutoscaleData struct {
	XMLName       xml.Name `xml:"autoscale-data" json:"-" toml:"-"`
	AutoscaleRule []AutoscaleRule
}

type ResourceConsumptionAndAutoscaleHistory struct {
	XMLName                   xml.Name `xml:"resource-consumption-and-autoscale-history" json:"-" toml:"-"`
	ResourceConsumptionSample []struct {
		RAMUsage               int       `xml:"ram-usage,attr"`
		CPUUsage               int       `xml:"cpu-usage,attr"`
		PrivateIncomingTraffic int       `xml:"private-incoming-traffic,attr"`
		PrivateOutgoingTraffic int       `xml:"private-outgoing-traffic,attr"`
		PublicIncomingTraffic  int       `xml:"public-incoming-traffic,attr"`
		PublicOutgoingTraffic  int       `xml:"public-outgoing-traffic,attr"`
		NodeSeqNo              int       `xml:"node-seq-no,attr"`
		NodeTimestamp          Timestamp `xml:"node-timestamp,attr"`
		PaciTimestamp          Timestamp `xml:"paci-timestamp,attr"`
		CPU                    int       `xml:"cpu,attr"`
		RAM                    int       `xml:"ram,attr"`
		Bandwidth              int       `xml:"bandwidth,attr"`
	} `xml:"resource-consumption-sample"`
	AutoscaleEvent struct {
		Direction     string    `xml:"direction,attr"`
		RuleVersion   int       `xml:"rule-version,attr"`
		NodeSeqNo     int       `xml:"node-seq-no,attr"`
		NodeTimestamp Timestamp `xml:"node-timestamp,attr"`
		Metric        string    `xml:"metric,attr"`
		NewValue      int       `xml:"new-value,attr"`
		NodeUUID      string    `xml:"node-uuid,attr"`
		Started       string    `xml:"started,attr"`
		Ended         string    `xml:"ended,attr"`
		EndedOk       bool      `xml:"ended-ok,attr"`
	} `xml:"autoscale-event"`
	AutoscaleRule []AutoscaleRule `xml:"autoscale-rule"`
}

type ApplicationList struct {
	XMLName             xml.Name              `xml:"application-list" json:"-" toml:"-"`
	ApplicationTemplate []ApplicationTemplate `xml:"application-template"`
}

type ApplicationTemplate struct {
	XMLName xml.Name `xml:"application-template" json:"-" toml:"-"`
	ID      int      `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	// Active      bool   `xml:"active,attr"`
	ForOS       string `xml:"for-os,attr"`
	Description string `xml:"description"`
}

type ImageList struct {
	XMLName   xml.Name `xml:"image-list" json:"-" toml:"-"`
	ImageInfo []struct {
		Name           string    `xml:"name,attr"`
		Description    string    `xml:"description,attr"`
		Size           int       `xml:"size,attr"`
		Created        Timestamp `xml:"created,attr"`
		SubscriptionID int       `xml:"subscription-id,attr"`
		ImageOf        string    `xml:"image-of,attr"`
		Location       string    `xml:location,attr"`
	} `xml:"image-info"`
}

type VeImage struct {
	XMLName        xml.Name  `xml:"ve-image" json:"-" toml:"-"`
	ID             int       `xml:"id,attr"`
	BnodeUUID      string    `xml:"bnode-uuid,attr"`
	CustomerID     int       `xml:"customer-id,attr"`
	SubscriptionID int       `xml:"subscription-id,attr"`
	Name           string    `xml:"name,attr"`
	Hostname       string    `xml:"hostname,attr"`
	Description    string    `xml:"description"`
	CPUNumber      int       `xml:"cpu-number,attr"`
	CPUPower       int       `xml:"cpu-power,attr"`
	RAMSize        int       `xml:"ram-size,attr"`
	Bandwidth      int       `xml:"bandwidth,attr"`
	Login          string    `xml:"login,attr"`
	PrimaryDiskID  int       `xml:"primary-disk-id,attr"`
	ImageSize      int       `xml:"image-size,attr"`
	Created        Timestamp `xml:"created,attr"`
	ImageOf        string    `xml:"image-of,attr"`
	NoOfPublicIP   int       `xml:"no-of-public-ip,attr"`
	NoOfPublicIPv6 int       `xml:"no-of-public-ipv6,attr"`
	CustomNs       bool      `xml:"custom-ns,attr"`
	Disks          []struct {
		ID      int    `xml:"id,attr"`
		Type    string `xml:"type,attr"`
		Primary bool   `xml:"primary,attr"`
		Size    int    `xml:"size,attr"`
	} `xml:"disks"`
	Platform Platform `xml:"platform"`
}

type LbList struct {
	XMLName      xml.Name `xml:"lb-list" json:"-" toml:"-"`
	LoadBalancer []struct {
		Name           string `xml:"name,attr"`
		State          string `xml:"state,attr"`
		SubscriptionID int    `xml:"subscription-id,attr"`
	} `xml:"load-balancer"`
}

type LoadBalancer struct {
	XMLName         xml.Name       `xml:"load-balancer" json:"-" toml:"-"`
	ID              int            `xml:"id"`
	UUID            string         `xml:"uuid"`
	Hnid            int            `xml:"hnId"`
	CustomerID      int            `xml:"customer-id"`
	Name            string         `xml:"name"`
	Hostname        string         `xml:"hostname"`
	Description     string         `xml:"description"`
	SubscriptionID  int            `xml:"subscription-id"`
	CPU             CPU            `xml:"cpu"`
	RAMSize         int            `xml:"ram-size"`
	Bandwidth       int            `xml:"bandwidth"`
	VeDisk          VeDisk         `xml:"ve-disk"`
	Platform        Platform       `xml:"platform"`
	Network         Network        `xml:"network"`
	BackupSchedule  BackupSchedule `xml:"backup-schedule"`
	Console         Console        `xml:"console"`
	State           string         `xml:"state"`
	PrimaryDiskID   int            `xml:"primary-disk-id"`
	TemplateID      int            `xml:"template-id"`
	Admin           Admin          `xml:"admin"`
	LastOperationRc int            `xml:"last-operation-rc"`
	UsedBy          []struct {
		VeName string `xml:"ve-name,attr"`
		IP     IPAddr `xml:"ip,attr"`
	} `xml:"used-by"`
}

type Template struct {
	XMLName                  xml.Name `xml:"template" json:"-" toml:"-"`
	ID                       int      `xml:"id,attr"`
	Name                     string   `xml:"name,attr"`
	OSType                   string   `xml:"osType,attr"`
	Technology               string   `xml:"technology,attr"`
	Active                   bool     `xml:"active,attr"`
	Default                  bool     `xml:"default,attr"`
	RootLogin                string   `xml:"root-login,attr"`
	MinHddSize               int      `xml:"min-hdd-size,attr"`
	PwdRegex                 string   `xml:"pwd-regex,attr"`
	HighWaterMarkForDelivery int      `xml:"high-watermark-for-delivery,attr"`
	LowWaterMarkForDelivery  int      `xml:"low-watermark-for-delivery,attr"`
	Option                   []struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	} `xml:"option"`
}

type TemplateList struct {
	XMLName  xml.Name   `xml:"template-list" json:"-" toml:"-"`
	Template []Template `xml:"template"`
}

type BackupScheduleList struct {
	XMLName        xml.Name `xml:"backup-schedule-list" json:"-" toml:"-"`
	BackupSchedule []struct {
		ID   int    `xml:"id,attr"`
		Name string `xml:"name,attr"`
		//		CronExpression  string `xml:"cron-expression,attr"`
		Description     string `xml:"description"`
		Enabled         bool   `xml:"enabled,attr"`
		BackupsToKeep   int    `xml:"backups-to-keep,attr"`
		NoOfIncremental int    `xml:"no-of-incremental,attr"`
	} `xml:"backup-schedule"`
}
