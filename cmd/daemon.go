/*
 * Copyright (C) 2024 Delusoire
 * SPDX-License-Identifier: GPL-3.0-or-later
 */

package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/Delusoire/bespoke-cli/v3/paths"
	"github.com/charmbracelet/log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/net/publicsuffix"

	"github.com/gorilla/websocket"
)

var (
	DaemonAddr    = "localhost:7967"
	AllowedOrigin = "https://xpui.app.spotify.com"
	daemon        bool
)

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run daemon",
	Run: func(cmd *cobra.Command, args []string) {
		if daemon {
			rootLogger.Info("Starting daemon")
			startDaemon(rootLogger)
		}
	},
}

var daemonStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start daemon",
	Run: func(cmd *cobra.Command, args []string) {
		rootLogger.Info("Starting daemon")
		startDaemon(rootLogger)
	},
}

var daemonEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable daemon",
	Run: func(cmd *cobra.Command, args []string) {
		rootLogger.Info("Enabling daemon")
		daemon = true
		viper.Set("daemon", daemon)
		viper.WriteConfig()
	},
}

var daemonDisableCmd = &cobra.Command{
	Use:   "disable",
	Short: "Disable daemon",
	Run: func(cmd *cobra.Command, args []string) {
		rootLogger.Info("Disabling daemon")
		daemon = false
		viper.Set("daemon", daemon)
		viper.WriteConfig()
	},
}

func init() {
	cobra.OnInitialize(func() {
		viper.SetDefault("daemon", true)
		daemon = viper.GetBool("daemon")
	})

	rootCmd.AddCommand(daemonCmd)

	daemonCmd.AddCommand(daemonStartCmd, daemonEnableCmd, daemonDisableCmd)

	viper.SetDefault("daemon", false)
}

func startDaemon(logger *log.Logger) {
	c := make(chan struct{})
	var (
		watcherCtx    context.Context
		watcherCancel context.CancelFunc
	)

	startWatcher := func() {
		watcherCtx, watcherCancel = context.WithCancel(context.Background())
		go watchSpotifyApps(watcherCtx, spotifyDataPath, logger.WithPrefix("Watcher"))
	}

	viper.OnConfigChange(func(in fsnotify.Event) {
		restartWatcher := false

		_daemon := viper.GetBool("daemon")
		_mirror := viper.GetBool("mirror")
		_spotifyDataPath := viper.GetString("spotify-data-path")
		_spotifyExecPath := viper.GetString("spotify-exec-path")
		_spotifyConfigPath := viper.GetString("spotify-config-path")

		if _spotifyDataPath != spotifyDataPath {
			if watcherCancel != nil {
				// TODO: wait for watcher to stop
				watcherCancel()
			}
			restartWatcher = true
		}

		daemon = _daemon
		mirror = _mirror
		spotifyDataPath = _spotifyDataPath
		spotifyExecPath = _spotifyExecPath
		spotifyConfigPath = _spotifyConfigPath

		if !daemon {
			close(c)
		}

		if restartWatcher {
			startWatcher()
		}
	})

	go viper.WatchConfig()
	startWatcher()
	go func() {
		setupProxy(logger.WithPrefix("Proxy"))
		setupWebSocket(logger.WithPrefix("WebSocket"))
		err := http.ListenAndServe(DaemonAddr, nil)
		logger.Fatalf("failed to start server: %s", err)
	}()

	<-c

	os.Exit(0)
}

func watchSpotifyApps(ctx context.Context, spotifyDataPath string, logger *log.Logger) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logger.Fatal(err)
	}
	defer watcher.Close()

	logger.Infof("watching: %s", paths.GetSpotifyAppsPath(spotifyDataPath))
	if err := watcher.Add(paths.GetSpotifyAppsPath(spotifyDataPath)); err != nil {
		logger.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Info("stopping")
			return
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}
			logger.Infof("event: %s", event)
			if event.Has(fsnotify.Create) {
				if strings.HasSuffix(event.Name, "xpui.spa") {
					if err := execApply(logger); err != nil {
						logger.Warn(err)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			logger.Warn(err)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: improve security
		return true
	},
}

func setupWebSocket(logger *log.Logger) {
	http.HandleFunc("/rpc", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Infof("failed to upgrade: %s", err)
			return
		}
		defer c.Close()

		for {
			_, p, err := c.ReadMessage()
			if err != nil {
				logger.Warnf("failed to read: %s", err)
				break
			}

			incoming := string(p)
			logger.Infof("recv: %s", incoming)
			res, err := HandleProtocol(incoming)
			if err != nil {
				logger.Warnf("protocol error: %s", err)
			}
			if res != "" {
				c.WriteMessage(websocket.TextMessage, []byte(res))
			}
		}
	})
}

func setupProxy(logger *log.Logger) {
	proxy := (&httputil.ReverseProxy{
		Transport: &CustomTransport{Transport: http.DefaultTransport},
		Rewrite: func(r *httputil.ProxyRequest) {
			p, ok := strings.CutPrefix(r.In.URL.Path, "/proxy/")
			if !ok {
				logger.Fatal(errors.New("proxy received invalid path"))
			}
			u, err := url.Parse(p)
			if err != nil {
				logger.Fatal(fmt.Errorf("proxy received invalid path: %w", err))
			}

			r.Out.URL = u
			r.Out.Host = ""

			xSetHeaders := r.In.Header.Get("X-Set-Headers")
			r.Out.Header.Del("X-Set-Headers")
			var headers map[string]string
			if err := json.Unmarshal([]byte(xSetHeaders), &headers); err == nil {
				for k, v := range headers {
					if v == "undefined" {
						r.Out.Header.Del(k)
					} else {
						r.Out.Header.Set(k, v)
					}
				}
			}
		},
	})

	http.HandleFunc("/proxy/{url}", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Origin", AllowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Set-Headers")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		proxy.ServeHTTP(w, r)
	})
}

var jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

type CustomTransport struct {
	Transport http.RoundTripper
}

func (t *CustomTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if jar != nil {
		for _, cookie := range jar.Cookies(req.URL) {
			req.AddCookie(cookie)
		}
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if jar != nil {
		if rc := resp.Cookies(); len(rc) > 0 {
			jar.SetCookies(req.URL, rc)
		}
	}

	resp.Header.Set("Access-Control-Allow-Origin", AllowedOrigin)
	resp.Header.Set("Access-Control-Allow-Credentials", "true")

	if loc, err := resp.Location(); err == nil {
		proxyUrl := "http://" + DaemonAddr + "/proxy/"
		resp.Header.Set("Location", proxyUrl+url.PathEscape(loc.String()))
	}

	return resp, nil
}
