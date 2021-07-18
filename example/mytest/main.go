package main

import (
	// "context"
	// 	// "fmt"
	// 	// "time"

	"log"
	"os"

	// 	"strconv"
	// 	"syscall"

	// 	// "time"
	// 	// "fmt"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
	"github.com/hanwen/go-fuse/v2/fuse/nodefs"
	"github.com/hanwen/go-fuse/v2/fuse/pathfs"
)

type numberNode struct {
	// Must embed an Inode for the struct to work as a node.
	fs.Inode

	// num is the integer represented in this file/directory
	num int
}

func main() {
	// This is where we'll mount the FS
	mntDir := "/tmp/x"
	os.Mkdir(mntDir, 0755)
	root := &numberNode{num: 10}

	
	pfs := pathfs.NewPathNodeFs(pathfs.NewLoopbackFileSystem("/"),
		&pathfs.PathNodeFsOptions{Debug: true })

	svr, errs := 
	server, err := fs.Mount(mntDir, root, &fs.Options{
		MountOptions: fuse.MountOptions{
			// Set to true to see how the file system works.
			Debug: true,
		},
	})
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Mounted on %s", mntDir)
	log.Printf("Unmount by calling 'fusermount -u %s'", mntDir)

	// Wait until unmount before exiting
	server.Wait()
}
