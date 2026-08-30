package dldruntime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/tree-1508/dlm-butler-runtime/dldruntime/internal/buildinfo"
)

type Config struct {
	Store       Store
	APIBaseURL  string
	StorageRoot string
	StagingRoot string
	HTTPClient  *http.Client
}

type Runtime struct {
	store       Store
	apiBaseURL  string
	storageRoot string
	stagingRoot string
	httpClient  *http.Client

	mu             sync.Mutex
	gameProfiles   map[int64]int64
	uploadProfiles map[int64]int64
	uploads        map[int64]*Upload
	queues         map[string]*queuedInstall
	active         map[string]context.CancelFunc
}

type queuedInstall struct{ Params InstallQueueParams }

func New(cfg Config) (*Runtime, error) {
	if cfg.Store == nil {
		return nil, errors.New("store is required")
	}
	if cfg.APIBaseURL == "" {
		cfg.APIBaseURL = "https://api.itch.io"
	}
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = http.DefaultClient
	}
	return &Runtime{
		store: cfg.Store, apiBaseURL: strings.TrimRight(cfg.APIBaseURL, "/"), storageRoot: cfg.StorageRoot, stagingRoot: cfg.StagingRoot, httpClient: cfg.HTTPClient,
		gameProfiles: make(map[int64]int64), uploadProfiles: make(map[int64]int64), uploads: make(map[int64]*Upload), queues: make(map[string]*queuedInstall), active: make(map[string]context.CancelFunc),
	}, nil
}

func (r *Runtime) Handle(ctx context.Context, method string, raw json.RawMessage) (any, bool, error) {
	switch method {
	case "Version.Get":
		return map[string]any{"version": buildinfo.Version, "versionString": buildinfo.VersionString()}, false, nil
	case "Meta.Shutdown":
		return map[string]any{}, true, nil
	case "Profile.LoginWithOAuthCode":
		res, err := r.loginWithOAuthCode(ctx, raw)
		return res, false, err
	case "Fetch.ProfileGames":
		res, err := r.fetchProfileGames(ctx, raw)
		return res, false, err
	case "Fetch.GameUploads":
		res, err := r.fetchGameUploads(ctx, raw)
		return res, false, err
	case "Install.GetUploads":
		res, err := r.installGetUploads(ctx, raw)
		return res, false, err
	case "Install.PlanUpload":
		res, err := r.installPlanUpload(ctx, raw)
		return res, false, err
	case "Install.Queue":
		res, err := r.installQueue(raw)
		return res, false, err
	case "Install.Perform":
		res, err := r.installPerform(ctx, raw)
		return res, false, err
	case "Install.Cancel":
		res, err := r.installCancel(raw)
		return res, false, err
	default:
		return nil, false, methodNotFound()
	}
}

type oauthTokenResponse struct {
	Key *struct {
		ID     int64  `json:"id"`
		UserID int64  `json:"userId"`
		Key    string `json:"key"`
	} `json:"key"`
	Cookie map[string]string `json:"cookie"`
}

type profileResponse struct {
	User *User `json:"user"`
}
type profileGamesResponse struct {
	Games []*Game `json:"games"`
}
type gameResponse struct {
	Game *Game `json:"game"`
}
type uploadsResponse struct {
	Uploads []*Upload `json:"uploads"`
}
type uploadResponse struct {
	Upload *Upload `json:"upload"`
}

func (r *Runtime) apiRequest(ctx context.Context, method, path, apiKey string, form url.Values, dst any) error {
	var body *strings.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	} else {
		body = strings.NewReader("")
	}
	req, err := http.NewRequestWithContext(ctx, method, r.apiBaseURL+path, body)
	if err != nil {
		return err
	}
	if form != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	if apiKey != "" {
		req.Header.Set("Authorization", apiKey)
	}
	req.Header.Set("Accept", "application/vnd.itch.v2")
	req.Header.Set("User-Agent", "DLD-Butler-Runtime")
	res, err := r.httpClient.Do(req)
	if err != nil {
		return redactError(err)
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("remote operation failed with HTTP status %d", res.StatusCode)
	}
	if err := json.NewDecoder(res.Body).Decode(dst); err != nil {
		return errors.New("remote response decode failed")
	}
	return nil
}

func (r *Runtime) loginWithOAuthCode(ctx context.Context, raw json.RawMessage) (*ProfileLoginWithOAuthCodeResult, error) {
	var p ProfileLoginWithOAuthCodeParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.Code == "" || p.CodeVerifier == "" || p.RedirectURI == "" || p.ClientID == "" {
		return nil, invalidParams("code, codeVerifier, redirectUri and clientId are required")
	}
	form := url.Values{"grant_type": {"authorization_code"}, "code": {p.Code}, "code_verifier": {p.CodeVerifier}, "redirect_uri": {p.RedirectURI}, "client_id": {p.ClientID}}
	var token oauthTokenResponse
	if err := r.apiRequest(ctx, http.MethodPost, "/oauth/token", "", form, &token); err != nil || token.Key == nil || token.Key.Key == "" {
		return nil, operationFailed("OAuth code exchange failed")
	}
	var profile profileResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/profile", token.Key.Key, nil, &profile); err != nil || profile.User == nil {
		return nil, operationFailed("profile retrieval failed")
	}
	stored := StoredProfile{ID: profile.User.ID, APIKey: token.Key.Key, User: profile.User}
	if err := r.store.SaveProfile(stored); err != nil {
		return nil, operationFailed("secure profile persistence failed")
	}
	return &ProfileLoginWithOAuthCodeResult{Profile: &Profile{ID: stored.ID, LastConnected: time.Now().UTC(), User: stored.User}, Cookie: token.Cookie}, nil
}

func (r *Runtime) fetchProfileGames(ctx context.Context, raw json.RawMessage) (*FetchProfileGamesResult, error) {
	var p FetchProfileGamesParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.ProfileID == 0 {
		return nil, invalidParams("profileId is required")
	}
	profile, err := r.store.LoadProfile(p.ProfileID)
	if err != nil {
		return nil, operationFailed("profile is not available")
	}
	var apiRes profileGamesResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/profile/games", profile.APIKey, nil, &apiRes); err != nil {
		return nil, operationFailed("profile game request failed")
	}
	items := make([]*ProfileGame, 0, len(apiRes.Games))
	for _, g := range apiRes.Games {
		if g == nil || !matchProfileGame(g, p) {
			continue
		}
		items = append(items, &ProfileGame{Game: g, ViewsCount: g.ViewsCount, DownloadsCount: g.DownloadsCount, PurchasesCount: g.PurchasesCount, Published: g.Published})
		r.mu.Lock()
		r.gameProfiles[g.ID] = p.ProfileID
		r.mu.Unlock()
	}
	sortProfileGames(items, p.SortBy, p.Reverse)
	if p.Limit > 0 && int64(len(items)) > p.Limit {
		items = items[:p.Limit]
	}
	return &FetchProfileGamesResult{Items: items}, nil
}

func matchProfileGame(g *Game, p FetchProfileGamesParams) bool {
	if p.Search != "" && !strings.Contains(strings.ToLower(g.Title), strings.ToLower(p.Search)) {
		return false
	}
	if p.Filters.Visibility == "draft" && g.Published {
		return false
	}
	if p.Filters.Visibility == "published" && !g.Published {
		return false
	}
	if p.Filters.PaidStatus == "free" && g.MinPrice != 0 {
		return false
	}
	if p.Filters.PaidStatus == "paid" && g.MinPrice == 0 {
		return false
	}
	return true
}

func sortProfileGames(items []*ProfileGame, by string, reverse bool) {
	var less func(i, j int) bool
	switch by {
	case "title":
		less = func(i, j int) bool {
			return strings.ToLower(items[i].Game.Title) < strings.ToLower(items[j].Game.Title)
		}
	case "views":
		less = func(i, j int) bool { return items[i].ViewsCount > items[j].ViewsCount }
	case "downloads":
		less = func(i, j int) bool { return items[i].DownloadsCount > items[j].DownloadsCount }
	case "purchases":
		less = func(i, j int) bool { return items[i].PurchasesCount > items[j].PurchasesCount }
	default:
		return
	}
	sort.SliceStable(items, func(i, j int) bool {
		if reverse {
			return less(j, i)
		}
		return less(i, j)
	})
}

func (r *Runtime) resolveProfileID(gameID, uploadID, explicit int64) (int64, error) {
	if explicit != 0 {
		return explicit, nil
	}
	r.mu.Lock()
	if uploadID != 0 && r.uploadProfiles[uploadID] != 0 {
		id := r.uploadProfiles[uploadID]
		r.mu.Unlock()
		return id, nil
	}
	if gameID != 0 && r.gameProfiles[gameID] != 0 {
		id := r.gameProfiles[gameID]
		r.mu.Unlock()
		return id, nil
	}
	r.mu.Unlock()
	profiles, err := r.store.ListProfiles()
	if err != nil || len(profiles) != 1 {
		return 0, errors.New("profile cannot be resolved unambiguously")
	}
	return profiles[0].ID, nil
}

func (r *Runtime) fetchGameUploads(ctx context.Context, raw json.RawMessage) (*FetchGameUploadsResult, error) {
	var p FetchGameUploadsParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.GameID == 0 {
		return nil, invalidParams("gameId is required")
	}
	profileID, err := r.resolveProfileID(p.GameID, 0, 0)
	if err != nil {
		return nil, operationFailed("profile cannot be resolved")
	}
	profile, err := r.store.LoadProfile(profileID)
	if err != nil {
		return nil, operationFailed("profile is not available")
	}
	var apiRes uploadsResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/games/"+strconv.FormatInt(p.GameID, 10)+"/uploads", profile.APIKey, nil, &apiRes); err != nil {
		return nil, operationFailed("game uploads request failed")
	}
	uploads := apiRes.Uploads
	if p.OnlyCompatible {
		uploads, _ = splitCompatibleUploads(uploads)
	}
	r.rememberUploads(profileID, uploads)
	return &FetchGameUploadsResult{Uploads: uploads}, nil
}

func (r *Runtime) installGetUploads(ctx context.Context, raw json.RawMessage) (*InstallGetUploadsResult, error) {
	var p InstallGetUploadsParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.GameID == 0 {
		return nil, invalidParams("gameId is required")
	}
	profileID, err := r.resolveProfileID(p.GameID, 0, p.ProfileID)
	if err != nil {
		return nil, operationFailed("profile cannot be resolved")
	}
	profile, err := r.store.LoadProfile(profileID)
	if err != nil {
		return nil, operationFailed("profile is not available")
	}
	var gr gameResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/games/"+strconv.FormatInt(p.GameID, 10), profile.APIKey, nil, &gr); err != nil || gr.Game == nil {
		return nil, operationFailed("game request failed")
	}
	var ur uploadsResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/games/"+strconv.FormatInt(p.GameID, 10)+"/uploads", profile.APIKey, nil, &ur); err != nil {
		return nil, operationFailed("game uploads request failed")
	}
	compatible, incompatible := splitCompatibleUploads(ur.Uploads)
	r.mu.Lock()
	r.gameProfiles[p.GameID] = profileID
	r.mu.Unlock()
	r.rememberUploads(profileID, ur.Uploads)
	return &InstallGetUploadsResult{Game: gr.Game, Uploads: compatible, IncompatibleUploads: incompatible}, nil
}

func (r *Runtime) rememberUploads(profileID int64, uploads []*Upload) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, u := range uploads {
		if u != nil {
			r.uploadProfiles[u.ID] = profileID
			r.uploads[u.ID] = u
		}
	}
}

func splitCompatibleUploads(uploads []*Upload) ([]*Upload, []*Upload) {
	var compatible, incompatible []*Upload
	for _, u := range uploads {
		if u == nil {
			continue
		}
		name := strings.ToLower(u.Filename)
		platformOK := u.Type != "default" || u.Platforms.Windows == ArchitecturesAll || u.Platforms.Windows == ArchitecturesAmd64
		formatOK := !strings.HasSuffix(name, ".7z") && !strings.HasSuffix(name, ".rar") && !strings.HasSuffix(name, ".xz") && !strings.HasSuffix(name, ".tar.xz")
		if platformOK && formatOK {
			compatible = append(compatible, u)
		} else {
			incompatible = append(incompatible, u)
		}
	}
	return compatible, incompatible
}

func newID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}

func (r *Runtime) installPlanUpload(ctx context.Context, raw json.RawMessage) (*InstallPlanUploadResult, error) {
	var p InstallPlanUploadParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.UploadID == 0 {
		return nil, invalidParams("uploadId is required")
	}
	profileID, err := r.resolveProfileID(0, p.UploadID, 0)
	if err != nil {
		return nil, operationFailed("profile cannot be resolved")
	}
	profile, err := r.store.LoadProfile(profileID)
	if err != nil {
		return nil, operationFailed("profile is not available")
	}
	var ur uploadResponse
	if err := r.apiRequest(ctx, http.MethodGet, "/uploads/"+strconv.FormatInt(p.UploadID, 10), profile.APIKey, nil, &ur); err != nil || ur.Upload == nil {
		return nil, operationFailed("upload request failed")
	}
	u := ur.Upload
	r.rememberUploads(profileID, []*Upload{u})
	typeName, ok := supportedInstallerType(u.Filename)
	info := &InstallPlanInfo{Upload: u, Build: u.Build, Type: typeName}
	if !ok {
		info.Error = "unsupported-format"
		info.ErrorMessage = "upload format is outside the DLD-051 narrow runtime"
		info.ErrorCode = 1
		return &InstallPlanUploadResult{Info: info}, nil
	}
	info.DiskUsage = &DiskUsageInfo{FinalDiskUsage: u.Size, NeededFreeSpace: u.Size * 2, Accuracy: "estimate"}
	return &InstallPlanUploadResult{Info: info}, nil
}

func (r *Runtime) installQueue(raw json.RawMessage) (*InstallQueueResult, error) {
	var p InstallQueueParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if !p.NoCave || p.Upload == nil || p.Game == nil || p.ProfileID == 0 || p.InstallFolder == "" || p.StagingFolder == "" {
		return nil, invalidParams("DLD runtime requires noCave, profileId, game, upload, installFolder and stagingFolder")
	}
	if p.QueueDownload {
		return nil, invalidParams("queueDownload is not supported by the DLD runtime")
	}
	if _, ok := supportedInstallerType(p.Upload.Filename); !ok {
		return nil, invalidParams("upload format is not supported by the DLD runtime")
	}
	if err := requireContained(r.storageRoot, p.InstallFolder); err != nil {
		return nil, invalidParams("installFolder is outside the configured Storage Root")
	}
	if err := requireContained(r.stagingRoot, p.StagingFolder); err != nil {
		return nil, invalidParams("stagingFolder is outside the configured Staging Root")
	}
	if _, err := r.store.LoadProfile(p.ProfileID); err != nil {
		return nil, operationFailed("profile is not available")
	}
	id, err := newID()
	if err != nil {
		return nil, operationFailed("could not allocate operation id")
	}
	if p.Reason == "" {
		p.Reason = "install"
	}
	r.mu.Lock()
	r.queues[id] = &queuedInstall{Params: p}
	r.uploadProfiles[p.Upload.ID] = p.ProfileID
	r.uploads[p.Upload.ID] = p.Upload
	r.mu.Unlock()
	return &InstallQueueResult{ID: id, Reason: p.Reason, Game: p.Game, Upload: p.Upload, Build: p.Build, InstallFolder: p.InstallFolder, StagingFolder: p.StagingFolder}, nil
}

func (r *Runtime) installCancel(raw json.RawMessage) (*InstallCancelResult, error) {
	var p InstallCancelParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.ID == "" {
		return nil, invalidParams("id is required")
	}
	r.mu.Lock()
	cancel := r.active[p.ID]
	r.mu.Unlock()
	if cancel == nil {
		return &InstallCancelResult{DidCancel: false}, nil
	}
	cancel()
	return &InstallCancelResult{DidCancel: true}, nil
}

func (r *Runtime) installPerform(ctx context.Context, raw json.RawMessage) (*InstallPerformResult, error) {
	var p InstallPerformParams
	if err := decodeParams(raw, &p); err != nil {
		return nil, err
	}
	if p.ID == "" || p.StagingFolder == "" {
		return nil, invalidParams("id and stagingFolder are required")
	}
	r.mu.Lock()
	q := r.queues[p.ID]
	r.mu.Unlock()
	if q == nil {
		return nil, invalidParams("unknown install id")
	}
	if p.StagingFolder != q.Params.StagingFolder {
		return nil, invalidParams("stagingFolder does not match queued operation")
	}
	opCtx, cancel := context.WithCancel(ctx)
	r.mu.Lock()
	if _, exists := r.active[p.ID]; exists {
		r.mu.Unlock()
		cancel()
		return nil, operationFailed("operation is already active")
	}
	r.active[p.ID] = cancel
	r.mu.Unlock()
	defer func() { cancel(); r.mu.Lock(); delete(r.active, p.ID); delete(r.queues, p.ID); r.mu.Unlock() }()
	manager, err := r.performInstall(opCtx, q.Params)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, operationFailed("operation canceled")
		}
		return nil, operationFailed("install operation failed")
	}
	return &InstallPerformResult{Events: []InstallEvent{{Type: "install", Timestamp: time.Now().UTC(), Install: &InstallInstallEvent{Manager: manager}}}}, nil
}

func redactError(err error) error {
	if err == nil {
		return nil
	}
	return errors.New("remote operation failed")
}
