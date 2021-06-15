package main

import (
	"fmt"
	"os/exec"
	"sync"
)

var (
	lpfs = []string{
		"go-addr-util",
		"go-conn-security-multistream",
		"go-eventbus",
		"go-flow-metrics",
		"go-libp2p",
		"go-libp2p-autonat",
		"go-libp2p-blankhost",
		"go-libp2p-circuit",
		"go-libp2p-connmgr",
		"go-libp2p-core",
		"go-libp2p-crypto",
		"go-libp2p-discovery",
		"go-libp2p-kad-dht",
		"go-libp2p-kbucket",
		"go-libp2p-loggables",
		"go-libp2p-mplex",
		"go-libp2p-netutil",
		"go-libp2p-noise",
		"go-libp2p-peer",
		"go-libp2p-peerstore",
		"go-libp2p-pnet",
		"go-libp2p-pubsub-router",
		"go-libp2p-quic-transport",
		"go-libp2p-record",
		"go-libp2p-routing-helpers",
		"go-libp2p-secio",
		"go-libp2p-swarm",
		"go-libp2p-testing",
		"go-libp2p-tls",
		"go-libp2p-transport-upgrader",
		"go-libp2p-yamux",
		"go-mplex",
		"go-msgio",
		"go-nat",
		"go-netroute",
		"go-openssl",
		"go-reuseport",
		"go-reuseport-transport",
		"go-sockaddr",
		"go-stream-muxer-multistream",
		"go-tcp-transport",
		"go-ws-transport",
		"go-yamux",
	}
	libp2p = []string{
		"bbloom",
		"go-bitswap",
		"go-block-format",
		"go-blockservice",
		"go-cid",
		"go-cidutil",
		"go-datastore",
		"go-ds-badger",
		"go-ds-flatfs",
		"go-ds-leveldb",
		"go-ds-measure",
		"go-filestore",
		"go-fs-lock",
		"go-graphsync",
		"go-ipfs",
		"go-ipfs-blockstore",
		"go-ipfs-chunker",
		"go-ipfs-cmds",
		"go-ipfs-config",
		"go-ipfs-delay",
		"go-ipfs-ds-help",
		"go-ipfs-exchange-interface",
		"go-ipfs-exchange-offline",
		"go-ipfs-files",
		"go-ipfs-pinner",
		"go-ipfs-posinfo",
		"go-ipfs-pq",
		"go-ipfs-provider",
		"go-ipfs-routing",
		"go-ipfs-util",
		"go-ipld-cbor",
		"go-ipld-format",
		"go-ipld-git",
		"go-ipns",
		"go-log",
		"go-merkledag",
		"go-metrics-interface",
		"go-mfs",
		"go-path",
		"go-peertaskqueue",
		"go-unixfs",
		"go-verifcid",
		"interface-go-ipfs-core",
	}
)

func main() {
	var wg sync.WaitGroup
	for _, g := range lpfs {
		wg.Add(1)
		go func(g string) {
			defer wg.Done()
			path := "git@github.com:ipfs/" + g + ".git"
			fmt.Println(path)
			cmd := exec.Command("git", "clone", path)
			cmd.Dir = "/Users/shuxian/Projects/stack-labs/debug-libs/ipfs"
			stdout, err := cmd.CombinedOutput()
			fmt.Printf("start %s-%s\n", stdout, err)
		}(g)
	}

	wg.Wait()
}
