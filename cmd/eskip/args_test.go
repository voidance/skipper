package main

import (
	"net/url"
	"os"
	"testing"
)

func preserveArgs(args []string, f func()) {
	os.Args, args = append([]string{"eskip", "cmd"}, args...), os.Args
	defer func() {
		os.Args = args
	}()

	etcdUrls = ""
	etcdStorageRoot = ""
	inlineRoutes = ""
	inlineRouteIds = ""

	initFlags()

	f()
}

func TestProcessArgs(t *testing.T) {
	for i, item := range []struct {
		args  []string
		fail  bool
		err   error
		media []*medium
	}{{

		// invalid flag
		[]string{"-invalid"},
		true,
		nil,
		nil,
	}, {

		// missing etcd urls
		[]string{"-etcd-urls"},
		true,
		nil,
		nil,
	}, {

		// invalid url
		[]string{"-etcd-urls", "::"},
		true,
		nil,
		nil,
	}, {

		// empty args
		nil,
		false,
		nil,
		nil,
	}, {

		// etcd-urls
		[]string{"-etcd-urls", "https://etcd1.example.org:4242,https://etcd2.example.org:4545"},
		false,
		nil,
		[]*medium{{
			typ: etcd,
			urls: []*url.URL{
				{Scheme: "https", Host: "etcd1.example.org:4242"},
				{Scheme: "https", Host: "etcd2.example.org:4545"}},
			path: "/skipper"}},
	}, {

		// inline routes
		[]string{"-routes", `Method("POST") -> "https://www.example.org"`},
		false,
		nil,
		[]*medium{{
			typ:   inline,
			eskip: `Method("POST") -> "https://www.example.org"`}},
	}, {

		// inline ids
		[]string{"-ids", "route1,route2"},
		false,
		nil,
		[]*medium{{
			typ: inlineIds,
			ids: []string{"route1", "route2"}}},
	}, {

		// etcd storage root
		[]string{"-etcd-storage-root", "/skipper-test"},
		false,
		nil,
		[]*medium{{
			typ: etcd,
			urls: []*url.URL{
				{Scheme: "http", Host: "127.0.0.1:2379"},
				{Scheme: "http", Host: "127.0.0.1:4001"}},
			path: "/skipper-test"}},
	}, {

		// too many files
		[]string{"file1", "file2"},
		true,
		invalidNumberOfArgs,
		nil,
	}, {

		// file
		[]string{"file1"},
		false,
		nil,
		[]*medium{{
			typ:  file,
			path: "file1"}},
	}} {
		preserveArgs(item.args, func() {
			media, err := processArgs()
			if item.fail {
				if err == nil {
					t.Error("failed to fail", i)
				}

				if item.err != nil && err != item.err {
					t.Error("invalid error", i)
				}
			} else {
				if err != nil {
					t.Error(err)
				}

				if len(media) == len(item.media) {
					for j, m := range item.media {
						checkMedium(t, m, media[j], i, j)
					}
				} else {
					t.Error("invalid number of parsed media", i)
				}
			}
		})
	}
}
