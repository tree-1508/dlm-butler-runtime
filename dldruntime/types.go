package dldruntime

import "time"

type User struct {
	ID            int64  `json:"id"`
	Username      string `json:"username"`
	DisplayName   string `json:"displayName"`
	Developer     bool   `json:"developer"`
	PressUser     bool   `json:"pressUser"`
	URL           string `json:"url"`
	CoverURL      string `json:"coverUrl"`
	StillCoverURL string `json:"stillCoverUrl"`
}

type Architectures string

const (
	ArchitecturesAll   Architectures = "all"
	Architectures386   Architectures = "386"
	ArchitecturesAmd64 Architectures = "amd64"
)

type Platforms struct {
	Windows Architectures `json:"windows,omitempty"`
	Linux   Architectures `json:"linux,omitempty"`
	OSX     Architectures `json:"osx,omitempty"`
}

type Game struct {
	ID             int64      `json:"id"`
	URL            string     `json:"url"`
	Title          string     `json:"title"`
	ShortText      string     `json:"shortText,omitempty"`
	Type           string     `json:"type"`
	Classification string     `json:"classification"`
	CoverURL       string     `json:"coverUrl,omitempty"`
	StillCoverURL  string     `json:"stillCoverUrl,omitempty"`
	MinPrice       int64      `json:"minPrice,omitempty"`
	CanBeBought    bool       `json:"canBeBought,omitempty"`
	HasDemo        bool       `json:"hasDemo,omitempty"`
	Platforms      Platforms  `json:"platforms"`
	User           *User      `json:"user,omitempty"`
	UserID         int64      `json:"userId,omitempty"`
	ViewsCount     int64      `json:"viewsCount,omitempty"`
	DownloadsCount int64      `json:"downloadsCount,omitempty"`
	PurchasesCount int64      `json:"purchasesCount,omitempty"`
	Published      bool       `json:"published,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	PublishedAt    *time.Time `json:"publishedAt,omitempty"`
}

type Build struct {
	ID            int64  `json:"id"`
	ParentBuildID int64  `json:"parentBuildId"`
	State         string `json:"state"`
	UploadID      int64  `json:"uploadId"`
	GameID        int64  `json:"gameId"`
	UserID        int64  `json:"userId"`
	Version       int64  `json:"version"`
	UserVersion   string `json:"userVersion"`
}

type Upload struct {
	ID          int64      `json:"id"`
	Storage     string     `json:"storage"`
	Host        string     `json:"host,omitempty"`
	Filename    string     `json:"filename"`
	DisplayName string     `json:"displayName"`
	Size        int64      `json:"size"`
	ChannelName string     `json:"channelName"`
	Build       *Build     `json:"build,omitempty"`
	BuildID     int64      `json:"buildId,omitempty"`
	Type        string     `json:"type"`
	Preorder    bool       `json:"preorder"`
	Demo        bool       `json:"demo"`
	Platforms   Platforms  `json:"platforms"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

type Profile struct {
	ID            int64     `json:"id"`
	LastConnected time.Time `json:"lastConnected"`
	User          *User     `json:"user"`
}

type ProfileLoginWithOAuthCodeParams struct {
	Code         string `json:"code"`
	CodeVerifier string `json:"codeVerifier"`
	RedirectURI  string `json:"redirectUri"`
	ClientID     string `json:"clientId"`
}

type ProfileLoginWithOAuthCodeResult struct {
	Profile *Profile          `json:"profile"`
	Cookie  map[string]string `json:"cookie"`
}

type FetchProfileGamesParams struct {
	ProfileID int64  `json:"profileId"`
	Limit     int64  `json:"limit"`
	Search    string `json:"search"`
	SortBy    string `json:"sortBy"`
	Filters   struct {
		Visibility string `json:"visibility"`
		PaidStatus string `json:"paidStatus"`
	} `json:"filters"`
	Reverse bool   `json:"reverse"`
	Cursor  string `json:"cursor"`
	Fresh   bool   `json:"fresh"`
}

type ProfileGame struct {
	Game           *Game `json:"game"`
	ViewsCount     int64 `json:"viewsCount"`
	DownloadsCount int64 `json:"downloadsCount"`
	PurchasesCount int64 `json:"purchasesCount"`
	Published      bool  `json:"published"`
}

type FetchProfileGamesResult struct {
	Items      []*ProfileGame `json:"items"`
	NextCursor string         `json:"nextCursor,omitempty"`
	Stale      bool           `json:"stale,omitempty"`
}

type FetchGameUploadsParams struct {
	GameID         int64 `json:"gameId"`
	OnlyCompatible bool  `json:"compatible"`
	Fresh          bool  `json:"fresh"`
}

type FetchGameUploadsResult struct {
	Uploads []*Upload `json:"uploads"`
	Stale   bool      `json:"stale,omitempty"`
}

type InstallGetUploadsParams struct {
	GameID    int64 `json:"gameId"`
	ProfileID int64 `json:"profileId,omitempty"`
}

type InstallGetUploadsResult struct {
	Game                *Game     `json:"game"`
	Uploads             []*Upload `json:"uploads"`
	IncompatibleUploads []*Upload `json:"incompatibleUploads,omitempty"`
}

type InstallPlanUploadParams struct {
	ID                string `json:"id"`
	UploadID          int64  `json:"uploadId"`
	DownloadSessionID string `json:"downloadSessionId"`
}

type InstallPlanUploadResult struct {
	Info *InstallPlanInfo `json:"info"`
}

type InstallPlanInfo struct {
	Upload       *Upload        `json:"upload"`
	Build        *Build         `json:"build,omitempty"`
	Type         string         `json:"type"`
	DiskUsage    *DiskUsageInfo `json:"diskUsage,omitempty"`
	Error        string         `json:"error,omitempty"`
	ErrorMessage string         `json:"errorMessage,omitempty"`
	ErrorCode    int64          `json:"errorCode,omitempty"`
}

type DiskUsageInfo struct {
	FinalDiskUsage  int64  `json:"finalDiskUsage"`
	NeededFreeSpace int64  `json:"neededFreeSpace"`
	Accuracy        string `json:"accuracy"`
}

type InstallQueueParams struct {
	CaveID            string  `json:"caveId"`
	Reason            string  `json:"reason"`
	InstallLocationID string  `json:"installLocationId"`
	NoCave            bool    `json:"noCave"`
	InstallFolder     string  `json:"installFolder"`
	Game              *Game   `json:"game,omitempty"`
	Upload            *Upload `json:"upload,omitempty"`
	Build             *Build  `json:"build,omitempty"`
	IgnoreInstallers  bool    `json:"ignoreInstallers,omitempty"`
	StagingFolder     string  `json:"stagingFolder"`
	QueueDownload     bool    `json:"queueDownload"`
	FastQueue         bool    `json:"fastQueue"`
	ProfileID         int64   `json:"profileId,omitempty"`
}

type InstallQueueResult struct {
	ID                string  `json:"id"`
	Reason            string  `json:"reason"`
	CaveID            string  `json:"caveId"`
	Game              *Game   `json:"game"`
	Upload            *Upload `json:"upload"`
	Build             *Build  `json:"build,omitempty"`
	InstallFolder     string  `json:"installFolder"`
	StagingFolder     string  `json:"stagingFolder"`
	InstallLocationID string  `json:"installLocationId"`
}

type InstallPerformParams struct {
	ID            string `json:"id"`
	StagingFolder string `json:"stagingFolder"`
}

type InstallEvent struct {
	Type      string               `json:"type"`
	Timestamp time.Time            `json:"timestamp"`
	Install   *InstallInstallEvent `json:"install,omitempty"`
	Problem   *ProblemInstallEvent `json:"problem,omitempty"`
}

type InstallInstallEvent struct {
	Manager string `json:"manager"`
}

type ProblemInstallEvent struct {
	Error      string `json:"error"`
	ErrorStack string `json:"errorStack"`
}

type InstallPerformResult struct {
	CaveID string         `json:"caveId"`
	Events []InstallEvent `json:"events"`
}

type InstallCancelParams struct {
	ID string `json:"id"`
}

type InstallCancelResult struct {
	DidCancel bool `json:"didCancel"`
}
