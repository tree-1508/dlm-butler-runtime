//go:build dld_runtime

package daemon

import (
	"crawshaw.io/sqlite/sqlitex"
	"github.com/itchio/butler/buildinfo"
	"github.com/itchio/butler/butlerd"
	"github.com/itchio/butler/butlerd/messages"
	"github.com/itchio/butler/endpoints/fetch"
	"github.com/itchio/butler/endpoints/install"
	"github.com/itchio/butler/endpoints/profile"
	"github.com/itchio/butler/mansion"
)

var mainRouter *butlerd.Router

// GetRouter returns the DLD-051 least-privilege RPC router when built with
// -tags dld_runtime. It intentionally omits Launch, password/API-key login,
// uninstall, publishing, search, update and diagnostic/test endpoints.
func GetRouter(dbPool *sqlitex.Pool, mansionContext *mansion.Context) *butlerd.Router {
	if mainRouter != nil {
		return mainRouter
	}

	mainRouter = butlerd.NewRouter(dbPool, mansionContext.NewClient, mansionContext.HTTPClient, mansionContext.HTTPTransport)

	messages.VersionGet.Register(mainRouter, func(rc *butlerd.RequestContext, params butlerd.VersionGetParams) (*butlerd.VersionGetResult, error) {
		return &butlerd.VersionGetResult{
			Version:       buildinfo.Version,
			VersionString: buildinfo.VersionString,
		}, nil
	})
	messages.MetaShutdown.Register(mainRouter, func(rc *butlerd.RequestContext, params butlerd.MetaShutdownParams) (*butlerd.MetaShutdownResult, error) {
		rc.Shutdown()
		return &butlerd.MetaShutdownResult{}, nil
	})

	// OAuth-only profile lifecycle. Password and API-key login are deliberately absent.
	messages.ProfileLoginWithOAuthCode.Register(mainRouter, profile.LoginWithOAuthCode)
	messages.ProfileList.Register(mainRouter, profile.List)
	messages.ProfileForget.Register(mainRouter, profile.Forget)

	// Narrow owned-library read surface required by the future pilot contract.
	messages.FetchProfileGames.Register(mainRouter, fetch.FetchProfileGames)

	// Narrow acquisition/install surface. No uninstall, launch, prereq execution or shortcuts.
	messages.GameFindUploads.Register(mainRouter, install.GameFindUploads)
	messages.InstallGetUploads.Register(mainRouter, install.InstallGetUploads)
	messages.InstallPlanUpload.Register(mainRouter, install.InstallPlanUpload)
	messages.InstallQueue.Register(mainRouter, install.InstallQueue)
	messages.InstallPerform.Register(mainRouter, install.InstallPerform)
	messages.InstallCancel.Register(mainRouter, install.InstallCancel)

	return mainRouter
}
